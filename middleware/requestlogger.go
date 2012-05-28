package middleware

import (
  "log"
  "time"
  ".."
)

var RequestLogger = &LoggerMiddleware{}

type LoggerMiddleware struct {
}

func (requestLogger *LoggerMiddleware) Execute(context *soggy.Context) {
  startTime := time.Now()
  log.Println("Request for", context.Req.URL)
  context.Next(nil)
  log.Println("Request took", time.Since(startTime))
}
