package client

import (
  "io/ioutil"
  "net/http"
  "time"
)

type Response struct {
  Status int
  Headers http.Header
  Body []byte
}

func Get(url string) (*Response, error) {
  return GetWithTimeout(url, 0)
}

func GetWithTimeout(url string, timeout time.Duration) (*Response, error) {
  client := http.Client{
    Timeout: timeout,
  }
  res, err := client.Get(url)
  if err != nil {
    return &Response{0, http.Header{}, []byte{}}, err
  }
  defer res.Body.Close()
  body, err := ioutil.ReadAll(res.Body)
  return &Response{res.StatusCode, res.Header, body}, err
}
