package dist

import (
  "hash/crc32"
  "fmt"
  "path/filepath"
  "strconv"
  "strings"
)

func ConvertPath(path string, distSteps, distRange int) string {
  dir := filepath.Dir(path)
  if dir != "" {
    dir += "/"
  }
  fileName := filepath.Base(path)
  parts := calculate(fileName, distSteps, distRange)
  distPath := strings.Join(parts, "/")
  if distPath != "" {
    distPath += "/"
  }
  return dir + distPath + fileName
}

func calculate(fileName string, distSteps, distRange int) []string {
  partLength := len(strconv.Itoa(distRange - 1))
  sum := int(crc32.ChecksumIEEE([]byte(fileName)))
  var parts []string
  for i := 0; i < distSteps; i++ {
    partFormat := "%0" + strconv.Itoa(partLength) + "d"
    part := fmt.Sprintf(partFormat, sum % distRange)
    parts = append([]string{part}, parts...)
    sum /= distRange
  }
  return parts[:]
}
