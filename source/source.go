package source

import (
  "../client"
)

type Parser func(*client.Response) (interface{}, error)

type Source interface {
  Name() string
  Iterator() chan string
  Parser() Parser
}
