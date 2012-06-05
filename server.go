package soggy

import (
  "strings"
  "net/http"
  "log"
  "io"
  "path/filepath"
)

const (
  CONFIG_VIEW_PATH = "viewPath"
  CONFIG_STATIC_PATH = "staticPath"
  DEFAULT_VIEW_PATH = "views"
  DEFAULT_STATIC_PATH = "public"
)

type Servers []*Server

type ErrorHandler func(*Context, interface{})

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

func (config ServerConfig) SetViewPath(viewPath string) error {
  viewPath, err := filepath.Abs(viewPath)
  if err != nil {
    return err
  }
  config[CONFIG_VIEW_PATH] = viewPath
  return nil
}

func (config ServerConfig) SetStaticPath(staticPath string) error {
  staticPath, err := filepath.Abs(staticPath)
  if err != nil {
    return err
  }
  config[CONFIG_STATIC_PATH] = staticPath
  return nil
}

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
      server.ErrorHandler(context, err)
    } else if nextIndex < len(middlewares) {
      currentIndex := nextIndex
      nextIndex++
      middlewares[currentIndex].Execute(context)
    }
  }

  context = &Context{ wrappedReq, wrappedRes, server, env, next }
  next(nil)
}

func (server *Server) TemplatePath(filename string) (ext, path string) {
  filePath := filepath.Join(server.Config[CONFIG_VIEW_PATH].(string), filename)
  return filepath.Ext(filePath), filePath
}

func (server *Server) SetMountpoint(mountpoint string) {
  server.Mountpoint = SaneURLPath(mountpoint)
}

func (server *Server) IsValidForPath(path string) bool {
  return strings.HasPrefix(path, server.Mountpoint)
}

func (server *Server) Engine(ext string, engine TemplateEngine) {
  server.EngineFunc(ext, func(writer io.Writer, filename string, options interface{}) error {
    return engine.SoggyEngine(writer, filename, options);
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

func DefaultErrorHandler(ctx *Context, err interface{}) {
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
  server := &Server{ Router: NewRouter(),
    Config: NewServerConfig(),
    ErrorHandler: DefaultErrorHandler,
    TemplateEngines: make(map[string]TemplateEngineFunc) }
  server.SetMountpoint(mountpoint)
  server.Engine("html", &HTMLTemplateEngine{})
  return server
}

func NewServerConfig() ServerConfig {
  config := make(ServerConfig)
  config.SetViewPath(DEFAULT_VIEW_PATH)
  config.SetStaticPath(DEFAULT_STATIC_PATH)
  return config
}
