package soggy

import (
  "log"
  "net/http"
  "mime"
  "encoding/json"
)

const (
  POWERED_BY_HEADER = "X-Powered-By"
  POWERED_BY = "sogginess"
)

type Response struct {
  http.ResponseWriter
}

func (res *Response) Render(status int, file string, params interface{}) error {
  log.Println("Would have rendered", file, "with", params, "if this was implemented.")
  return nil
}

func (res *Response) Html(status int, html string) error {
  res.Set("Content-Type", mime.TypeByExtension(".html"))
  _, err := res.WriteString(html)
  return err
}

func (res *Response) Json(status int, jsonIn interface{}) error {
  res.Set("Content-Type", mime.TypeByExtension(".json"))
  jsonOut, err := json.Marshal(jsonIn)
  if err == nil {
    _, err = res.Write(jsonOut)
  }
  return err
}

func (res *Response) WriteString(s string) (int, error) {
    return res.Write([]byte(s))
}

func (res *Response) Set(header, value string) {
    res.Header().Set(header, value)
}

func NewResponse(res http.ResponseWriter) *Response {
  wrappedResponse := &Response{res}
  wrappedResponse.Set(POWERED_BY_HEADER, POWERED_BY)
  return wrappedResponse;
}
