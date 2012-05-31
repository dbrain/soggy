package soggy

import (
  "strings"
  "net/http"
  "log"
)

type Servers []*Server

type ErrorHandler func(error, *Context)

func (servers Servers) Len() int {
  return len(servers)
}

func (servers Servers) Less(i, j int) bool {
  return len([]rune(servers[i].Mountpoint)) > len([]rune(servers[j].Mountpoint))
}

func (servers Servers) Swap(i, j int) {
  servers[i], servers[j] = servers[j], servers[i]
}

type ServerConfig map[string]interface{}

type Server struct {
  Mountpoint string
  middleware []Middleware
  Router *Router
  Config ServerConfig
  ErrorHandler ErrorHandler
}

func (server *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  var next func(error)
  var context *Context

  env := NewEnv()
  wrappedReq := NewRequest(req, server)
  wrappedReq.SetRelativePath(server.Mountpoint, SaneURLPath(req.URL.Path))
  wrappedRes := NewResponse(res)

  middlewares := server.middleware
  nextIndex := 0
  next = func (err error) {
    if err != nil {
      server.ErrorHandler(err, context)
    } else if nextIndex < len(middlewares) {
      currentIndex := nextIndex
      nextIndex++
      middlewares[currentIndex].Execute(context)
    }
  }

  context = &Context{ wrappedReq, wrappedRes, env, next }
  next(nil)
}

func (server *Server) SetMountpoint(mountpoint string) {
  server.Mountpoint = SaneURLPath(mountpoint)
}

func (server *Server) IsValidForPath(path string) bool {
  return strings.HasPrefix(path, server.Mountpoint)
}

func (server *Server) Use(middleware ...Middleware) {
  server.middleware = append(server.middleware, middleware...)
}

func (server *Server) Get(path string, routeHandler interface{}) {
  server.Router.AddRoute(GET_METHOD, path, routeHandler);
}

func (server *Server) Post(path string, routeHandler interface{}) {
  server.Router.AddRoute(POST_METHOD, path, routeHandler);
}

func (server *Server) Put(path string, routeHandler interface{}) {
  server.Router.AddRoute(PUT_METHOD, path, routeHandler);
}

func (server *Server) Delete(path string, routeHandler interface{}) {
  server.Router.AddRoute(DELETE_METHOD, path, routeHandler);
}

func (server *Server) All(path string, routeHandler interface{}) {
  server.Router.AddRoute(ALL_METHODS, path, routeHandler);
}

func DefaultErrorHandler(err error, ctx *Context) {
  res := ctx.Res
  res.WriteHeader(http.StatusInternalServerError)
  res.Write([]byte(err.Error()))
  log.Println("An error occured for", ctx.Req.RelativePath, err)
}

func NewServer(mountpoint string) *Server {
  server := &Server{ Router: NewRouter(), Config: make(ServerConfig),
    ErrorHandler: DefaultErrorHandler }
  server.SetMountpoint(mountpoint)
  return server
}
