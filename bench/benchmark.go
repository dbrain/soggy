package main

import (
  ".."
  "net/http"
  "fmt"
)

type basicHandler struct {
}
func (handler *basicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/plain")
  fmt.Fprintf(w, "Hello, World\r\n")
}

func adminOnly(ctx *soggy.Context, name string) map[string]string {
  if name == "admin" {
    ctx.Next(nil)
    return nil
  }
  return map[string]string{ "error": "Not admin" }
}

func StartSingleServer() {
  app, server := soggy.NewDefaultApp()

  server.Get("/handler", &basicHandler{})

  server.Get("/html/(.*)/(.*)", func (param1 string, param2 string) string {
    return param1 + " " + param2
  })

  server.Get("/htmlNoParams/(.*)/(.*)", func (ctx *soggy.Context) string {
    req := ctx.Req
    return req.URLParams[0] + " " + req.URLParams[1]
  })

  server.Get("/json/(.*)/(.*)", func (param1 string, param2 string) map[string]string {
    return map[string]string{ param1: param2 }
  })

  server.Get("/template/(.*)/(.*)", func (param1 string, param2 string) (string, interface{}) {
    return "kitchensink.html", map[string]string{
      "name": "Bob Barker",
      "age": "27" }
  })

  server.Get("/self", func (ctx *soggy.Context) {
    res := ctx.Res
    res.WriteHeader(http.StatusOK)
    res.Set("Content-Type", "text/plain")
    res.WriteString("Cannn dooo.")
  })

  server.Get("/bundle/(.*)", adminOnly, func () map[string]string {
    return map[string]string{ "ok": "Hey admin!" }
  })

  server.Use(server.Router)
  app.Listen(":9991")
}

func WebServer() *soggy.Server {
  server := soggy.NewServer("/")
  server.Get("/", func () string {
    return "root"
  })
  server.Use(server.Router)
  return server
}

func APIServer() *soggy.Server {
  server := soggy.NewServer("/api")
  server.Get("/", func () string {
    return "api"
  })
  server.Use(server.Router)
  return server
}

func AdminServer() *soggy.Server {
  server := soggy.NewServer("/admin")
  server.Get("/", func () string {
    return "admin"
  })
  server.Use(server.Router)
  return server
}

func StartMultipleServer() {
  app := soggy.NewApp()
  app.AddServers(WebServer(), APIServer(), AdminServer())
  app.Listen(":9992")
}

type MiddlewareA struct {}
func (middleware *MiddlewareA) Execute(ctx *soggy.Context) {
  ctx.Next(nil)
}
type MiddlewareB struct {}
func (middleware *MiddlewareB) Execute(ctx *soggy.Context) {
  ctx.Next(nil)
}
type MiddlewareC struct {}
func (middleware *MiddlewareC) Execute(ctx *soggy.Context) {
  ctx.Next(nil)
}

func StartMultipleMiddleware() {
  app, server := soggy.NewDefaultApp()
  server.Get("/", func () string {
    return "oh hi"
  })
  server.Use(&MiddlewareA{}, &MiddlewareB{}, &MiddlewareC{}, server.Router)
  app.Listen(":9993")
}

func StartBasicServer() {
  server := &http.Server{
    Addr: ":9994",
    Handler: &basicHandler{} }
  server.ListenAndServe()
}

func main() {
  go StartBasicServer()
  go StartSingleServer()
  go StartMultipleServer()
  StartMultipleMiddleware()
}
