package grabster

import (
  "net/http"
  "./grabber"
)

type Response struct{
  Url string
  Status int
  Headers http.Header
  Body []byte
  Err error
}

func HandleSync(iterator chan string, cachePath string) chan Response {
  handler := make(chan Response)
  g := grabber.New(cachePath)
  go func() {
    for url := range iterator {
      r := Response{}
      r.Url = url
      r.Status, r.Headers, r.Body, r.Err = g.Get(url)
      handler <- r
    }
    close(handler)
  }()
  return handler
}

func HandleAll(iterator chan string, cachePath string) chan Response {
  handler := make(chan Response)
  g := grabber.New(cachePath)
  go func() {
    received := 0
    sent := 0
    ready := false
    for url := range iterator {
      received++
      go func(url string) {
        r := Response{}
        r.Url = url
        r.Status, r.Headers, r.Body, r.Err = g.Get(url)
        handler <- r
        sent++
        if (ready && received == sent) {
          close(handler)
        }
      }(url)
    }
    ready = true
  }()
  return handler
}
