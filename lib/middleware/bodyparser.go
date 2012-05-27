package middleware

import (
  "log"
  ".."
)

var BodyParser = &BodyParserMiddleware{}

type BodyParserMiddleware struct {
}

func (bodyParser *BodyParserMiddleware) Execute(req *express.Request, res *express.Response, env map[string]interface{}, nextMiddleware func(error)) {
  log.Println("Body parser currently not implemented")
}