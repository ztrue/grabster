package dist

import (
  "fmt"
  "hash/crc32"
  "strconv"
  "strings"
)

func ConvertPath(path string, distSteps, distRange int) string {
  if distRange <= 0 || distRange <= 0 {
    return path
  }
  delimiter := "/"
  chars := []rune(path)
  isDir := false
  isNotEmpty := len(chars) > 0
  if isNotEmpty {
    lastChar := chars[len(chars) - 1]
    isDir = string(lastChar) == delimiter
  }
  index := 1
  if isDir {
    index = 2
  }
  rawParts := strings.Split(path, delimiter)
  fileName := rawParts[len(rawParts) - index]
  distParts := calculate(fileName, distSteps, distRange)
  parts := rawParts[:len(rawParts) - index]
  parts = append(parts, distParts...)
  if isNotEmpty {
    parts = append(parts, fileName)
  }
  if isDir {
    parts = append(parts, "")
  }
  return strings.Join(parts, delimiter)
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
