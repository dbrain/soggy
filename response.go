package express

import (
  "net/http"
)

const (
  POWERED_BY_HEADER = "X-Powered-By"
  POWERED_BY = "express.go"
)

type Response struct {
  http.ResponseWriter
}

func (self *Response) WriteString(s string) (int, error) {
    return self.Write([]byte(s))
}

func (self *Response) SetHeader(header, value string) {
    self.Header().Set(header, value)
}

func NewResponse(res http.ResponseWriter) *Response {
  wrappedResponse := &Response{res}
  wrappedResponse.Header().Set(POWERED_BY_HEADER, POWERED_BY)
  return wrappedResponse;
}
