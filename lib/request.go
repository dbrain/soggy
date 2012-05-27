package express

import (
  "net/http"
)

type Request struct {
  *http.Request
  Server *Server
  RelativePath string
}

func (req *Request) SetRelativePath(mountpoint string, path string) {
  if mountpoint == "/" {
    req.RelativePath = path
  } else {
    relativePath := path[len(mountpoint):len(path)]
    if len(relativePath) == 0 {
      relativePath = "/"
    }
    req.RelativePath = relativePath
  }
}

func NewRequest(req *http.Request, server *Server) *Request {
  return &Request{http.Request: req, Server: server}
}
