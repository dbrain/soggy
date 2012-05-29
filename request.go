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
  reqWrapper := &Request{}
  reqWrapper.URL = req.URL
  reqWrapper.Method = req.Method
  reqWrapper.Server = server
  return reqWrapper
}
