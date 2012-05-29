package soggy

import (
  "net/http"
  "net/url"
)

type Request struct {
  URL *url.URL
  Method string
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
  return &Request{URL: req.URL, Method: req.Method, Server: server}
}

func newStubRequest(method, path string) *Request {
  url, err := url.ParseRequestURI(path)
  if err != nil {
    panic("invalid path")
  }
  return &Request{URL: url, Method: method, RelativePath: path, Server: nil}
}

