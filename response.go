package soggy

import (
  "net/http"
  "encoding/json"
  "bytes"
  "os"
  "io"
)

const (
  POWERED_BY_HEADER = "X-Powered-By"
  POWERED_BY = "sogginess"
)

const (
  HTML_CONTENT_TYPE = "text/html; charset=utf-8"
  JSON_CONTENT_TYPE = "application/json; charset=utf-8"
)

type Response struct {
  http.ResponseWriter
  server *Server
}

func (res *Response) Render(status int, file string, params interface{}) (err interface{}) {
  buf := new(bytes.Buffer)

  ext, template := res.server.TemplatePath(file)
  if _, err := os.Stat(template); err != nil {
    return err
  }

  engine := res.server.TemplateEngines[ext[1:]]
  if engine == nil {
    return "No engine defined for " + ext[1:]
  }

  err = engine(buf, template, params)
  if err != nil {
    return err
  }
  res.Set("Content-Type", HTML_CONTENT_TYPE)
  res.Set("Content-Length", string(buf.Len()))
  res.WriteHeader(status)
  _, err = io.Copy(res, buf)
  return err
}

func (res *Response) Html(status int, html string) (err interface{}) {
  res.Set("Content-Type", HTML_CONTENT_TYPE)
  res.WriteHeader(status)
  _, err = res.WriteString(html)
  return err
}

func (res *Response) Json(status int, jsonIn interface{}) (err interface{}) {
  res.Set("Content-Type", JSON_CONTENT_TYPE)
  res.WriteHeader(status)
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

func NewResponse(res http.ResponseWriter, server *Server) *Response {
  wrappedResponse := &Response{res, server}
  wrappedResponse.Set(POWERED_BY_HEADER, POWERED_BY)
  return wrappedResponse;
}
