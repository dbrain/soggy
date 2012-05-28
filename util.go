package express

import (
  "strings"
)

func SaneURLPath(path string) string {
  if !strings.HasSuffix(path, "/") {
    path = path + "/"
  }
  return path
}
