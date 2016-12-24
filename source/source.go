package source

import (
  "../client"
)

type Parser func(string, *client.Response) (interface{}, error)

type Source interface {
  GetName() string
  Iterator() chan string
  Parser() Parser
}
