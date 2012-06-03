package soggy

import (
  "log"
  "time"
)

type RequestLoggerMiddleware struct {}

func (requestLogger *RequestLoggerMiddleware) Execute(context *Context) {
  startTime := time.Now()
  req := context.Req
  log.Println(req.Method, "request for", context.Req.URL)
  context.Next(nil)
  log.Println("Request took", time.Since(startTime))
}
