package fs

import (
  "io/ioutil"
  "os"
  "path/filepath"
)

func Exists(path string) bool {
  _, err := os.Stat(path)
  return !os.IsNotExist(err)
}

func Read(path string) ([]byte, error) {
  return ioutil.ReadFile(path)
}

func Write(path string, body []byte) error {
  dir := filepath.Dir(path)
  err := os.MkdirAll(dir, 0777)
  if err != nil {
    return err
  }
  return ioutil.WriteFile(path, body, 0777)
}

func Delete(path string) error {
  return os.RemoveAll(path)
}
