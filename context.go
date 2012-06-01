package soggy

import (
)

type Context struct {
  Req *Request
  Res *Response
  Server *Server
  Env Env
  Next func(interface{})
}

func NewContext(server *Server) *Context {
  return &Context{ Server: server }
}
