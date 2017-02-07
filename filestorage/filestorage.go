package filestorage

import (
  "path/filepath"
  "github.com/ztrue/grabster/filestorage/dist"
  "github.com/ztrue/grabster/filestorage/fs"
  "github.com/mholt/archiver"
)

const CompressedExt = ".zip"

type Storage struct {
  BasePath string
  IsDist bool
  DistSteps int
  DistRange int
  Compressed bool
}

func New() *Storage {
  return &Storage{
    BasePath: "",
    IsDist: false,
    DistSteps: 3,
    DistRange: 1000,
    Compressed: false,
  }
}

func (s *Storage) Exists(path string) bool {
  decompressed, compressed := s.convertPath(path)
  return fs.Exists(decompressed) || fs.Exists(compressed)
}

func (s *Storage) Read(path string) ([]byte, error) {
  decompressed, compressed := s.convertPath(path)
  decompressedExists := fs.Exists(decompressed)
  compressedExists := fs.Exists(compressed)
  if !s.Compressed {
    if !decompressedExists {
      if err := s.DecompressFile(compressed); err != nil {
        return []byte{}, err
      }
      defer fs.Delete(compressed)
    }
    data, err := fs.Read(decompressed)
    if err != nil {
      return data, err
    }
    if compressedExists {
      defer fs.Delete(compressed)
    }
    return data, nil
  }
  if compressedExists {
    if err := s.DecompressFile(compressed); err != nil {
      return []byte{}, err
    }
    defer fs.Delete(decompressed)
  }
  data, err := fs.Read(decompressed)
  if err != nil {
    return data, err
  }
  if !compressedExists {
    if err := s.CompressFile(decompressed); err != nil {
      return []byte{}, err
    }
    defer fs.Delete(decompressed)
  }
  return data, nil
}

func (s *Storage) Write(path string, body []byte) error {
  decompressed, compressed := s.convertPath(path)
  if err := fs.Write(decompressed, body); err != nil {
    return err
  }
  if s.Compressed {
    if err := s.CompressFile(decompressed); err != nil {
      return err
    }
    defer fs.Delete(decompressed)
  } else if fs.Exists(compressed) {
    defer fs.Delete(compressed)
  }
  return nil
}

func (s *Storage) Delete(path string) error {
  decompressed, compressed := s.convertPath(path)
  if fs.Exists(decompressed) {
    if err := fs.Delete(decompressed); err != nil {
      return err
    }
  }
  if fs.Exists(compressed) {
    return fs.Delete(compressed)
  }
  return nil
}

func (s *Storage) Compress(output string, inputs []string) error {
  return archiver.Zip.Make(output, inputs)
}

func (s *Storage) Decompress(input string, outputDir string) error {
  return archiver.Zip.Open(input, outputDir)
}

func (s *Storage) CompressFile(fileName string) error {
  return s.Compress(fileName + CompressedExt, []string{fileName})
}

func (s *Storage) DecompressFile(fileName string) error {
  return s.Decompress(fileName, filepath.Dir(fileName))
}

func (s *Storage) convertPath(path string) (string, string) {
  if s.BasePath != "" {
    path = filepath.Dir(s.BasePath + "/") + "/" + path
  }
  if s.IsDist {
    path = dist.ConvertPath(path, s.DistSteps, s.DistRange)
  }
  return path, path + CompressedExt
}
