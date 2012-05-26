package express

import (
  "log"
  "time"
)

var RequestLogger = &LoggerMiddleware{}

type LoggerMiddleware struct {
}

func (requestLogger *LoggerMiddleware) Execute(req *Request, res *Response, env map[string]interface{}, nextMiddleware func(error)) {
  startTime := time.Now()
  log.Println("Request for", req.URL)
  nextMiddleware(nil)
  log.Println("Request took", time.Since(startTime))
}
