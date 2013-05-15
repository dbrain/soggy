package soggy

import (
  "net/http"
  "net/url"
  "strings"
  "encoding/json"
  "errors"
  "io/ioutil"
)

type URLParams []string

type Request struct {
  *http.Request
  ID string
  RelativePath string
  URLParams URLParams

  bodyParsed bool
  bodyType string
  parsedBody interface{}
  bodyParseError error
}

var BodyTypeJson = "json"

func (req *Request) SetRelativePath(mountpoint string, path string) {
  if mountpoint == "/" {
    req.RelativePath = path
  } else {
    relativePath := path[len(mountpoint)-1:len(path)]
    req.RelativePath = relativePath
  }
}

func (req *Request) GetBody(jsonStruct interface{}) (string, interface{}, error) {
  if req.bodyParsed {
    return req.bodyType, req.parsedBody, req.bodyParseError
  }

  defer func () { req.bodyParsed = true }()
  contentType := req.Header.Get("Content-Type")
  if contentType == "" {
    req.bodyParseError = errors.New("No content type specified")
    return "", nil, req.bodyParseError
  }

  contentTypeParts := strings.Split(contentType, ";")
  mimeType := contentTypeParts[0]
  switch {
  case strings.HasPrefix(mimeType, "application/") && strings.HasSuffix(mimeType, "json"):
    req.bodyType = BodyTypeJson
    return req.parseJSON(jsonStruct)
  }

  return "", nil, errors.New("Unsupported content type " + contentType)
}

func (req *Request) parseJSON(jsonStruct interface{}) (string, interface{}, error) {
  if jsonStruct == nil {
    jsonStruct = map[string]interface{}{}
  }
  body, err := ioutil.ReadAll(req.Body);
  if err != nil {
    req.bodyParseError = err
    return req.bodyType, nil, req.bodyParseError
  }
  if err := json.Unmarshal(body, &jsonStruct); err != nil {
    req.bodyParseError = err
    return req.bodyType, nil, req.bodyParseError
  }
  req.parsedBody = jsonStruct
  return BodyTypeJson, req.parsedBody, nil
}

func NewRequest(req *http.Request) *Request {
  return &Request{ ID: UIDString(), Request: req, URLParams: make(URLParams, 2) }
}

func newStubRequest(method, path string) *Request {
  url, err := url.ParseRequestURI(path)
  if err != nil {
    panic("Invalid path")
  }
  req := &http.Request{ URL: url, Method: method }
  return &Request{ ID: UIDString(), Request: req, RelativePath: path }
}
