package grabber

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
  }
  return &Grabber{storage}
}

func (g *Grabber) Get(url string) (int, http.Header, []byte, error) {
  fileName := g.getPath(url)
  if g.storage != nil && g.storage.Exists(fileName) {
    status, headers, body, err := g.getFromCache(fileName)
    if err == nil {
      return status, headers, body, nil
    }
    // TODO Log errors if any
  }
  return g.loadActual(url)
}

func (g *Grabber) ClearCache(url string) error {
  if g.storage == nil {
    return nil
  }
  // TODO Delete folders also if empty
  return g.storage.Delete(g.getPath(url))
}

func (g *Grabber) getFromCache(fileName string) (int, http.Header, []byte, error) {
  data, readErr := g.storage.Read(fileName)
  if readErr != nil {
    return 0, http.Header{}, []byte{}, readErr
  }
  lines := strings.Split(string(data), "\n")
  if len(lines) < 3 {
    return 0, http.Header{}, []byte{}, nil
  }
  status, convErr := strconv.Atoi(lines[0])
  if convErr != nil {
    return 0, http.Header{}, []byte{}, convErr
  }
  var headers http.Header
  jsonErr := json.Unmarshal([]byte(lines[1]), &headers)
  if jsonErr != nil {
    return 0, http.Header{}, []byte{}, jsonErr
  }
  body := []byte(strings.Join(lines[2:], "\n"))
  return status, headers, body, nil
}

func (g *Grabber) loadActual(url string) (int, http.Header, []byte, error) {
  status, headers, body, clientErr := client.Get(url)
  if g.storage != nil && clientErr == nil {
    cacheErr := g.cacheActual(url, status, headers, body)
    if cacheErr != nil {
      // TODO Log cache errors if any
    }
  }
  return status, headers, body, clientErr
}

func (g *Grabber) cacheActual(url string, status int, headers http.Header, body []byte) error {
  jsonHeaders, jsonErr := json.Marshal(headers)
  if jsonErr != nil {
    return jsonErr
  }
  parts := []string{strconv.Itoa(status), string(jsonHeaders), string(body)}
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
