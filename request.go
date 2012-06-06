package soggy

import (
  "net/http"
  "net/url"
)

type URLParams []string

type Request struct {
  *http.Request
  RelativePath string
  URLParams URLParams
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
  return &Request{ Request: req, URLParams: make(URLParams, 2) }
}

func newStubRequest(method, path string) *Request {
  url, err := url.ParseRequestURI(path)
  if err != nil {
    panic("invalid path")
  }
  req := &http.Request{ URL: url, Method: method }
  return &Request{ Request: req, RelativePath: path }
}

