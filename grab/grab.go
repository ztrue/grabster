package grab

import (
  "crypto/md5"
  "encoding/hex"
  "encoding/json"
  "net/http"
  "strconv"
  "strings"
  "../client"
  "../filestorage"
)

type Grabber struct {
  storage *filestorage.Storage
}

func New(cachePath string) *Grabber {
  var storage *filestorage.Storage
  if cachePath != "" {
    storage = filestorage.New()
    storage.BasePath = cachePath
    storage.IsDist = true
    storage.Compressed = true
  }
  return &Grabber{storage}
}

func (g *Grabber) Get(url string) (*client.Response, bool, error) {
  fileName := g.getPath(url)
  if g.storage != nil && g.storage.Exists(fileName) {
    response, err := g.getFromCache(fileName)
    if err == nil {
      return response, true, nil
    }
    // TODO Log errors if any
  }
  response, err := g.loadActual(url)
  return response, false, err
}

func (g *Grabber) ClearCache(url string) error {
  if g.storage == nil {
    return nil
  }
  // TODO Delete folders also if empty
  return g.storage.Delete(g.getPath(url))
}

func (g *Grabber) getFromCache(fileName string) (*client.Response, error) {
  data, readErr := g.storage.Read(fileName)
  if readErr != nil {
    return &client.Response{0, http.Header{}, []byte{}}, readErr
  }
  lines := strings.Split(string(data), "\n")
  if len(lines) < 4 {
    // TODO Return error
    return &client.Response{0, http.Header{}, []byte{}}, nil
  }
  status, convErr := strconv.Atoi(lines[1])
  if convErr != nil {
    return &client.Response{0, http.Header{}, []byte{}}, convErr
  }
  var headers http.Header
  jsonErr := json.Unmarshal([]byte(lines[2]), &headers)
  if jsonErr != nil {
    return &client.Response{0, http.Header{}, []byte{}}, jsonErr
  }
  body := []byte(strings.Join(lines[3:], "\n"))
  return &client.Response{status, headers, body}, nil
}

func (g *Grabber) loadActual(url string) (*client.Response, error) {
  response, clientErr := client.Get(url)
  if g.storage != nil && clientErr == nil {
    cacheErr := g.cacheActual(url, response)
    if cacheErr != nil {
      // TODO Log cache errors if any
    }
  }
  return response, clientErr
}

func (g *Grabber) cacheActual(url string, response *client.Response) error {
  jsonHeaders, jsonErr := json.Marshal(response.Headers)
  if jsonErr != nil {
    return jsonErr
  }
  parts := []string{
    url,
    strconv.Itoa(response.Status),
    string(jsonHeaders),
    string(response.Body),
  }
  data := []byte(strings.Join(parts, "\n"))
  fileName := g.getPath(url)
  return g.storage.Write(fileName, data)
}

func (g *Grabber) getPath(url string) string {
  return g.urlToFileName(url)
}

func (g *Grabber) urlToFileName(url string) string {
  sum := md5.Sum([]byte(url))
  return hex.EncodeToString(sum[:])
}
