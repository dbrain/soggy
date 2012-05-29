package soggy

import (
  "log"
  "time"
)

type RequestLoggerMiddleware struct {}

func (requestLogger *RequestLoggerMiddleware) Execute(context *Context) {
  startTime := time.Now()
  log.Println("Request for", context.Req.URL)
  context.Next(nil)
  log.Println("Request took", time.Since(startTime))
}
