package middleware

import (
  "log"
  ".."
)

var BodyParser = &BodyParserMiddleware{}

type BodyParserMiddleware struct {
}

func (bodyParser *BodyParserMiddleware) Execute(req *soggy.Request, res *soggy.Response, env *soggy.Env, nextMiddleware func(error)) {
  log.Println("Body parser currently not implemented")
}
