package grabster

import (
  "sync"
  "time"
  "./grab"
  "./source"
)

type Result struct {
  Url string
  Data interface{}
  Err error
  Cached bool
}

func HandleSync(s source.Source, cachePath string, timeout time.Duration) chan *Result {
  handler := make(chan *Result)
  grabber := grab.New(cachePath + "/" + s.GetName())
  parser := s.Parser()
  go func() {
    defer close(handler)
    for url := range s.Iterator() {
      data, cached, err := func(url string) (interface{}, bool, error) {
        response, cached, err := grabber.Get(url)
        if err != nil {
          return response, cached, err
        }
        data, err := parser(response)
        return data, cached, err
      }(url)
      handler <- &Result{url, data, err, cached}
      time.Sleep(timeout)
    }
  }()
  return handler
}

func HandleAsync(s source.Source, cachePath string, timeout time.Duration) chan *Result {
  handler := make(chan *Result)
  grabber := grab.New(cachePath + "/" + s.GetName())
  parser := s.Parser()
  go func() {
    var wg sync.WaitGroup
    for url := range s.Iterator() {
      wg.Add(1)
      go func(url string) {
        defer wg.Done()
        data, cached, err := func(url string) (interface{}, bool, error) {
          response, cached, err := grabber.Get(url)
          if err != nil {
            return response, cached, err
          }
          data, err := parser(response)
          return data, cached, err
        }(url)
        handler <- &Result{url, data, err, cached}
      }(url)
      time.Sleep(timeout)
    }
    go func() {
      wg.Wait()
      close(handler)
    }()
  }()
  return handler
}
