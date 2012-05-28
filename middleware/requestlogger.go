package middleware

import (
  "log"
  "time"
  ".."
)

var RequestLogger = &LoggerMiddleware{}

type LoggerMiddleware struct {
}

func (requestLogger *LoggerMiddleware) Execute(req *soggy.Request, res *soggy.Response, env *soggy.Env, nextMiddleware func(error)) {
  startTime := time.Now()
  log.Println("Request for", req.URL)
  nextMiddleware(nil)
  log.Println("Request took", time.Since(startTime))
}
