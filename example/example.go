package main

import (
  ".."
  "../middleware"
)

func main() {
  server := express.NewServer("/web");
  server.Get("/i/like/cheese", func (req *express.Request, res *express.Response, env *express.Env, next func(error)) {
    res.Header().Set("Content-Type", "text/plain")
    res.Write([]byte("This is an example server. Hell yeah. It likes cheese"))
  })
  server.Get("/", func (req *express.Request, res *express.Response, env *express.Env, next func(error)) {
    res.Header().Set("Content-Type", "text/plain")
    res.Write([]byte("This is an example server. Hell yeah."))
  })
  server.All(express.ANY_PATH, func (req *express.Request, res *express.Response, env *express.Env, next func(error)) {
    res.Header().Set("Content-Type", "text/plain")
    res.Write([]byte("404 Page would go here"))
  })
  server.Use(middleware.RequestLogger, server.Router)

  app := express.NewApp()
  app.AddServer(server)
  app.AddServer(express.NewServer("/abc"))
  app.AddServer(express.NewServer("/abc123"))

  app.Listen("0.0.0.0:9999")
}
