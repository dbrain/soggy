package main

import (
  "github.com/dbrain/soggy"
  "net/http"
  "path/filepath"
)

type APIGetExample struct {
  Name string
  Age int
}

func WebServer() *soggy.Server {
  server := soggy.NewServer("/")
  server.Get("/", func (ctx *soggy.Context) (template string, opts interface{}) {
    return "multiserver_root.html", map[string]string{ "ip": ctx.Req.RemoteAddr }
  })
  server.Get("/about", func () string {
    return "How can you ask questions, when you have no mouth."
  })
  // The static server middleware will call next if the file is not found (to hit 404 or in this case)
  server.Get("/static/oldorange.jpg", func (ctx *soggy.Context) error {
    path, err := filepath.Abs("public/orange.jpg")
    if err != nil {
      return err
    }
    http.ServeFile(ctx.Res, ctx.Req.OriginalRequest, path)
    return nil
  })
  server.Use(&soggy.RequestLoggerMiddleware{}, soggy.NewStaticServerMiddleware("/static"), server.Router)
  return server
}

func APIServer() *soggy.Server {
  server := soggy.NewServer("/api")
  server.Get("/whoami", func () *APIGetExample {
    return &APIGetExample{"Daniel Brain", 27}
  })
  server.Post("/echo/(.*)/(.*)", func (key string, value string) map[string]string {
    result := make(map[string]string)
    result[key] = value
    return result
  })
  server.Use(&soggy.RequestLoggerMiddleware{}, soggy.NewStaticServerMiddleware("/static"), server.Router)
  return server
}

func AdminServer() *soggy.Server {
  server := soggy.NewServer("/admin")
  server.Get("/", func (ctx *soggy.Context) (template string, opts interface{}) {
    return "multiserver_admin.html", map[string]string{ "ip": ctx.Req.RemoteAddr }
  })
  server.Get("/hacks", func () string {
    return "You gosh darn hacked my admin section."
  })
  server.Use(&soggy.RequestLoggerMiddleware{}, soggy.NewStaticServerMiddleware("/static"), server.Router)
  return server
}

func main() {
  app := soggy.NewApp()
  app.AddServers(WebServer(), APIServer(), AdminServer())
  app.Listen("0.0.0.0:9999")
}
