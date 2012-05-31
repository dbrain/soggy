package soggy

import (
  "net/http"
  "net/url"
)

type Request struct {
  *http.Request
  RelativePath string
}

func (req *Request) SetRelativePath(mountpoint string, path string) {
  if mountpoint == "/" {
    req.RelativePath = path
  } else {
    relativePath := path[len(mountpoint)-1:len(path)]
    req.RelativePath = relativePath
  }
}

func NewRequest(req *http.Request) *Request {
  return &Request{http.Request: req}
}

func newStubRequest(method, path string) *Request {
  url, err := url.ParseRequestURI(path)
  if err != nil {
    panic("invalid path")
  }
  req := &http.Request{ URL: url, Method: method }
  return &Request{http.Request: req, RelativePath: path}
}

