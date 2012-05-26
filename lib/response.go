package express

import (
  "net/http"
)

type Response struct {
  http.ResponseWriter
}

func NewResponse(res http.ResponseWriter) *Response {
  return &Response{res}
}
