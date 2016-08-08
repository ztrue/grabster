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
      data, err := process(grabber, url, parser)
      handler <- &Result{url, data, err}
    }
  }()
  return handler
}

func process(grabber *grab.Grabber, url string, parser source.Parser) (interface{}, error) {
  response, grabberErr := grabber.Get(url)
  if grabberErr != nil {
    return nil, grabberErr
  }
  return parser(response)
}

// TODO HandleAsync
