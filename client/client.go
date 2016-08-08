package client

import (
  "io/ioutil"
  "net/http"
)

type Response struct {
  Status int
  Headers http.Header
  Body []byte
}

func Get(url string) (*Response, error) {
  res, err := http.Get(url)
  if err != nil {
    return &Response{0, http.Header{}, []byte{}}, err
  }
  defer res.Body.Close()
  body, err := ioutil.ReadAll(res.Body)
  return &Response{res.StatusCode, res.Header, body}, err
}
