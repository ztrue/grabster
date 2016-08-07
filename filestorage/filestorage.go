package filestorage

import (
  "path/filepath"
  "./dist"
  "./fs"
)

type Storage struct {
  BasePath string
  IsDist bool
  DistSteps int
  DistRange int
}

func New() *Storage {
  return &Storage{"", false, 3, 1000}
}

func (s *Storage) Exists(path string) bool {
  return fs.Exists(s.convertPath(path))
}

func (s *Storage) Read(path string) ([]byte, error) {
  return fs.Read(s.convertPath(path))
}

func (s *Storage) Write(path string, body []byte) error {
  return fs.Write(s.convertPath(path), body)
}

func (s *Storage) Delete(path string) error {
  return fs.Delete(s.convertPath(path))
}

func (s *Storage) convertPath(path string) string {
  if s.BasePath != "" {
    path = filepath.Dir(s.BasePath + "/") + "/" + path
  }
  if s.IsDist {
    path = dist.ConvertPath(path, s.DistSteps, s.DistRange)
  }
  return path
}
