package soggy

import (
  "net/http"
  "mime"
  "encoding/json"
  "bytes"
  "os"
  "io"
)

const (
  POWERED_BY_HEADER = "X-Powered-By"
  POWERED_BY = "sogginess"
)

type Response struct {
  http.ResponseWriter
  server *Server
}

func (res *Response) Render(status int, file string, params interface{}) (err interface{}) {
  res.WriteHeader(status)
  res.Set("Content-Type", mime.TypeByExtension(".html"))
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
  res.Set("Content-Length", string(buf.Len()))
  _, err = io.Copy(res, buf)
  return err
}

func (res *Response) Html(status int, html string) (err interface{}) {
  res.WriteHeader(status)
  res.Set("Content-Type", mime.TypeByExtension(".html"))
  _, err = res.WriteString(html)
  return err
}

func (res *Response) Json(status int, jsonIn interface{}) (err interface{}) {
  res.WriteHeader(status)
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

func NewResponse(res http.ResponseWriter, server *Server) *Response {
  wrappedResponse := &Response{res, server}
  wrappedResponse.Set(POWERED_BY_HEADER, POWERED_BY)
  return wrappedResponse;
}
