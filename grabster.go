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
    for url := range iterator {
      response, grabberErr := grabber.Get(url)
      if grabberErr != nil {
        handler <- &Result{url, nil, grabberErr}
        continue
      }
      data, parserErr := parser(response)
      if parserErr != nil {
        handler <- &Result{url, nil, parserErr}
        continue
      }
      handler <- &Result{url, data, nil}
    }
    close(handler)
  }()
  return handler
}

// TODO HandleAsync
