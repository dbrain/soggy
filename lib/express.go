package express

import (
  "fmt"
  "net/http"
)

type Middleware interface {
  Execute(*Request, *Response, map[string]interface{}, func(error))
}

type Server struct {
  httpServer *http.Server
  middleware []Middleware
  Router *Router
}

func (server *Server) requestHandler() http.HandlerFunc  {
  return http.HandlerFunc(func (res http.ResponseWriter, req *http.Request) {
    var next func (error)
    env := make(map[string]interface{})
    wrappedReq := NewRequest(req)
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
        middlewares[currentIndex].Execute(wrappedReq, wrappedRes, env, next)
      }
    }
    next(nil)
  })
}

func (server *Server) Use(middleware ...Middleware) {
  server.middleware = append(server.middleware, middleware...)
}

func (server *Server) Get(path string, routeHandler RouteHandler) {
  server.Router.AddRoute(GET_METHOD, path, routeHandler);
}

func (server *Server) Post(path string, routeHandler RouteHandler) {
  server.Router.AddRoute(POST_METHOD, path, routeHandler);
}

func (server *Server) Put(path string, routeHandler RouteHandler) {
  server.Router.AddRoute(PUT_METHOD, path, routeHandler);
}

func (server *Server) Delete(path string, routeHandler RouteHandler) {
  server.Router.AddRoute(DELETE_METHOD, path, routeHandler);
}

func (server *Server) All(path string, routeHandler RouteHandler) {
  server.Router.AddRoute(ALL_METHODS, path, routeHandler);
}

func (server *Server) Listen(address string) {
  httpServer := &http.Server{
    Addr: address,
    Handler: server.requestHandler() }
  server.httpServer = httpServer
  fmt.Println("Listening on", address)
  err := httpServer.ListenAndServe()
  if err != nil { panic(err) }
}

func NewServer() *Server {
  return &Server{ Router: NewRouter() }
}
