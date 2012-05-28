package soggy

import (
  "log"
  "net/http"
  "sort"
)

type Middleware interface {
  Execute(*Context)
}

type App struct {
  servers Servers
}

func (app *App) AddServer(server *Server) {
  app.servers = append(app.servers, server)
  sort.Sort(app.servers)
}

func (app *App) RequestHandler() http.HandlerFunc  {
  return http.HandlerFunc(func (res http.ResponseWriter, req *http.Request) {
    for _, server := range app.servers {
      if path := SaneURLPath(req.URL.Path); server.IsValidForPath(path) {
        var next func(error)
        var context *Context

        env := NewEnv()
        wrappedReq := NewRequest(req, server)
        wrappedReq.SetRelativePath(server.Mountpoint, path)
        wrappedRes := NewResponse(res)

        middlewares := server.middleware
        maxIndex := len(middlewares)
        nextIndex := 0
        next = func (err error) {
          if err != nil {
            panic(err) // TODO This should passed to an error handler
          } else if nextIndex < maxIndex {
            currentIndex := nextIndex
            nextIndex++
            middlewares[currentIndex].Execute(context)
          }
        }

        context = &Context{ wrappedReq, wrappedRes, env, next }
        next(nil)
        break
      }
    }
  })
}

func (app *App) Listen(address string) {
  httpServer := &http.Server{
    Addr: address,
    Handler: app.RequestHandler() }
  log.Println("Listening on", address)
  err := httpServer.ListenAndServe()
  if err != nil { panic(err) }
}

func NewApp() *App {
  return &App{}
}

func NewDefaultApp() (*App, *Server) {
  app := NewApp()
  server := NewServer("/")
  app.AddServer(server)
  return app, server
}
