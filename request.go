package soggy

import (
  "net/http"
  "net/url"
)

type Request struct {
  *http.Request
  RelativePath string
  Server *Server
}

func (req *Request) SetRelativePath(mountpoint string, path string) {
  if mountpoint == "/" {
    req.RelativePath = path
  } else {
    relativePath := path[len(mountpoint)-1:len(path)]
    req.RelativePath = relativePath
  }
}

func NewRequest(req *http.Request, server *Server) *Request {
  return &Request{http.Request: req, Server: server}
}

func newStubRequest(method, path string) *Request {
  url, err := url.ParseRequestURI(path)
  if err != nil {
    panic("invalid path")
  }
  req := &http.Request{ URL: url, Method: method }
  return &Request{http.Request: req, RelativePath: path, Server: nil}
}

