package grabster

import (
  "net/http"
  "sync"
  "./grab"
)

type Response struct {
  Url string
  Status int
  Headers http.Header
  Body []byte
  Err error
}

func HandleSync(iterator chan string, cachePath string) chan *Response {
  handler := make(chan *Response)
  grabber := grab.New(cachePath)
  go func() {
    for url := range iterator {
      response, err := grabber.Get(url)
      handler <- &Response{url, response.Status, response.Headers, response.Body, err}
    }
    close(handler)
  }()
  return handler
}

// TODO Refactor, fix errors
func HandleAll(iterator chan string, cachePath string) chan *Response {
  mutex := &sync.Mutex{}
  handler := make(chan *Response)
  grabber := grab.New(cachePath)
  go func() {
    received := 0
    sent := 0
    ready := false
    for url := range iterator {
      received++
      go func(url string) {
        response, err := grabber.Get(url)
        handler <- &Response{url, response.Status, response.Headers, response.Body, err}
        mutex.Lock()
        sent++
        if (ready && received == sent) {
          close(handler)
        }
        mutex.Unlock()
      }(url)
    }
    ready = true
  }()
  return handler
}
