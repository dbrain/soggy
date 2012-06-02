package soggy

import (
  "strings"
  "net/http"
  "log"
)

type Servers []*Server

type ErrorHandler func(interface{}, *Context)

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
  TemplateEngines map[string]TemplateEngineFunc
}

func (server *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  var next func(interface{})
  var context *Context

  env := NewEnv()
  wrappedReq := NewRequest(req)
  wrappedReq.SetRelativePath(server.Mountpoint, SaneURLPath(req.URL.Path))
  wrappedRes := NewResponse(res, server)

  middlewares := server.middleware
  nextIndex := 0
  next = func (err interface{}) {
    if err != nil {
      server.ErrorHandler(err, context)
    } else if nextIndex < len(middlewares) {
      currentIndex := nextIndex
      nextIndex++
      middlewares[currentIndex].Execute(context)
    }
  }

  context = &Context{ wrappedReq, wrappedRes, server, env, next }
  next(nil)
}

func (server *Server) SetMountpoint(mountpoint string) {
  server.Mountpoint = SaneURLPath(mountpoint)
}

func (server *Server) IsValidForPath(path string) bool {
  return strings.HasPrefix(path, server.Mountpoint)
}

func (server *Server) Engine(ext string, engine TemplateEngine) {
  server.EngineFunc(ext, func(filename string, options interface{}) ([]byte, error) {
    return engine.SoggyEngine(filename, options);
  })
}

func (server *Server) EngineFunc(ext string, engineFunc TemplateEngineFunc) {
  server.TemplateEngines[ext] = engineFunc
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

func DefaultErrorHandler(err interface{}, ctx *Context) {
  res := ctx.Res
  res.WriteHeader(http.StatusInternalServerError)
  switch err.(type) {
    case error:
      res.Write([]byte(err.(error).Error()))
    case string:
      res.Write([]byte(err.(string)))
    default:
      res.Write([]byte("An error occured processing your request"))
  }

  log.Println("An error occured for", ctx.Req.RelativePath, err)
}

func NewServer(mountpoint string) *Server {
  server := &Server{ Router: NewRouter(), Config: make(ServerConfig),
    ErrorHandler: DefaultErrorHandler, TemplateEngines: make(map[string]TemplateEngineFunc) }
  server.SetMountpoint(mountpoint)
  return server
}
