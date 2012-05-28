package express

import (
  "strings"
)

type Servers []*Server

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

func NewServer(mountpoint string) *Server {
  server := &Server{ Router: NewRouter() }
  server.SetMountpoint(mountpoint)
  return server
}
