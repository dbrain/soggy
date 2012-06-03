package soggy

import (
  "io"
  "log"
  "net/http"
  "sort"
)

type TemplateEngine interface {
  SoggyEngine(writer io.Writer, filename string, options interface{}) error
}

type TemplateEngineFunc func(writer io.Writer, filename string, options interface{}) error

type Middleware interface {
  Execute(*Context)
}

type App struct {
  servers Servers
}

func (app *App) AddServers(servers ...*Server) {
  for _, server := range servers {
    app.servers = append(app.servers, server)
  }
  sort.Sort(app.servers)
}

func (app *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  for _, server := range app.servers {
    if server.IsValidForPath(SaneURLPath(req.URL.Path)) {
      server.ServeHTTP(res, req)
      break
    }
  }
}

func (app *App) Listen(address string) {
  httpServer := &http.Server{
    Addr: address,
    Handler: app }
  log.Println("Starting to listen on", address)
  err := httpServer.ListenAndServe()
  if err != nil { panic(err) }
}

func NewApp() *App {
  return &App{}
}

func NewDefaultApp() (*App, *Server) {
  app := NewApp()
  server := NewServer("/")
  app.AddServers(server)
  return app, server
}
