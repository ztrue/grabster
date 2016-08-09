package grabster

import (
  "./grab"
  "./source"
)

type Result struct {
  Url string
  Data interface{}
  Err error
}

func HandleSync(s source.Source, cachePath string) chan *Result {
  handler := make(chan *Result)
  grabber := grab.New(cachePath + "/" + s.Name())
  parser := s.Parser()
  go func() {
    defer close(handler)
    for url := range s.Iterator() {
      data, err := func(url string) (interface{}, error) {
        response, err := grabber.Get(url)
        if err != nil {
          return nil, err
        }
        return parser(response)
      }(url)
      handler <- &Result{url, data, err}
    }
  }()
  return handler
}

// TODO HandleAsync
