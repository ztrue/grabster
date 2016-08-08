package grabster

import (
  "sync"
  "./client"
  "./grab"
)

type Result struct {
  Url string
  Response *client.Response
  Err error
}

func HandleSync(iterator chan string, cachePath string) chan *Result {
  handler := make(chan *Result)
  grabber := grab.New(cachePath)
  go func() {
    for url := range iterator {
      response, err := grabber.Get(url)
      handler <- &Result{url, response, err}
    }
    close(handler)
  }()
  return handler
}

// TODO Refactor, fix errors
func HandleAll(iterator chan string, cachePath string) chan *Result {
  mutex := &sync.Mutex{}
  handler := make(chan *Result)
  grabber := grab.New(cachePath)
  go func() {
    received := 0
    sent := 0
    ready := false
    for url := range iterator {
      received++
      go func(url string) {
        response, err := grabber.Get(url)
        handler <- &Result{url, response, err}
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
