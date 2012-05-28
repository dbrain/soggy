package middleware

import (
  "log"
  ".."
)

var BodyParser = &BodyParserMiddleware{}

type BodyParserMiddleware struct {
}

func (bodyParser *BodyParserMiddleware) Execute(context *soggy.Context) {
  log.Println("Body parser currently not implemented")
}
