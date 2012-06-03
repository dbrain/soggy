package main

import (
  "github.com/dbrain/soggy"
  "log"
  "errors"
)

func main() {
  app, server := soggy.NewDefaultApp()
  server.Get("/i/like/cheese", func (context *soggy.Context) {
    res := context.Res
    res.Header().Set("Content-Type", "text/plain")
    res.Write([]byte("This is an example server. Hell yeah. It likes cheese"))
  })
  server.Get("/i/ate/it/blah.html", func (context *soggy.Context) {
    res := context.Res
    res.Header().Set("Content-Type", "text/html")
    res.Write([]byte("<html><body>It ates you too</body></html>"))
  })
  server.Get("/", func (context *soggy.Context) {
    res := context.Res
    res.Header().Set("Content-Type", "text/plain")
    res.Write([]byte("This is an example server. Hell yeah."))
  })
  server.Get("/jebus", func (context *soggy.Context) {
    log.Println("In route for /jebus")
    // This should hit a 404 page for /jebus
    context.Next(nil)
  })
  server.Get("/error", func (context *soggy.Context) {
    log.Println("Im going to error")
    context.Next(errors.New("Uh oh spaghettios"))
  })
  server.Get("/echo/(.*)", func (echo string) string {
    return echo
  })
  server.Get("/json/(.*)", func (echo string) interface{} {
    echoMap := make(map[string]string)
    echoMap["echo"] = echo
    return echoMap
  })
  server.All(soggy.ANY_PATH, func (context *soggy.Context) {
    res := context.Res
    res.Header().Set("Content-Type", "text/plain")
    res.Write([]byte("404 Page would go here for: " + context.Req.RelativePath))
  })
  server.Use(&soggy.RequestLoggerMiddleware{}, server.Router)

  app.AddServer(soggy.NewServer("/abc"))
  app.AddServer(soggy.NewServer("/abc123"))

  app.Listen("0.0.0.0:9999")
}
