package client

import (
  "io/ioutil"
  "net/http"
)

func Get(url string) (int, http.Header, []byte, error) {
  res, err := http.Get(url)
  if err != nil {
    return 0, http.Header{}, []byte{}, err
  }
  defer res.Body.Close()
  body, err := ioutil.ReadAll(res.Body)
  return res.StatusCode, res.Header, body, err
}
