package main

import (
  ".."
  "../middleware"
)

func main() {
  server := soggy.NewServer("/");
  server.Get("/i/like/cheese", func (req *soggy.Request, res *soggy.Response, env *soggy.Env, next func(error)) {
    res.Header().Set("Content-Type", "text/plain")
    res.Write([]byte("This is an example server. Hell yeah. It likes cheese"))
  })
  server.Get("/i/ate/it/blah.html", func (req *soggy.Request, res *soggy.Response, env *soggy.Env, next func(error)) {
    res.Header().Set("Content-Type", "text/html")
    res.Write([]byte("<html><body>It ates you too</body></html>"))
  })
  server.Get("/", func (req *soggy.Request, res *soggy.Response, env *soggy.Env, next func(error)) {
    res.Header().Set("Content-Type", "text/plain")
    res.Write([]byte("This is an example server. Hell yeah."))
  })
  server.All(soggy.ANY_PATH, func (req *soggy.Request, res *soggy.Response, env *soggy.Env, next func(error)) {
    res.Header().Set("Content-Type", "text/plain")
    res.Write([]byte("404 Page would go here"))
  })
  server.Use(middleware.RequestLogger, server.Router)

  app := soggy.NewApp()
  app.AddServer(server)
  app.AddServer(soggy.NewServer("/abc"))
  app.AddServer(soggy.NewServer("/abc123"))

  app.Listen("0.0.0.0:9999")
}
