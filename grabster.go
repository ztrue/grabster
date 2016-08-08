package grabster

import (
  "./client"
  "./grab"
)

type Parser func(*client.Response) (interface{}, error)

type Result struct {
  Url string
  Data interface{}
  Err error
}

func HandleSync(iterator chan string, parser Parser, cachePath string) chan *Result {
  handler := make(chan *Result)
  grabber := grab.New(cachePath)
  go func() {
    defer close(handler)
    for url := range iterator {
      data, err := process(grabber, url, parser)
      handler <- &Result{url, data, err}
    }
  }()
  return handler
}

func process(grabber *grab.Grabber, url string, parser Parser) (interface{}, error) {
  response, grabberErr := grabber.Get(url)
  if grabberErr != nil {
    return nil, grabberErr
  }
  return parser(response)
}

// TODO HandleAsync
