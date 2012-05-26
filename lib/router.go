package express

import (
  "regexp"
)

const (
  GET_METHOD = "GET"
  POST_METHOD = "POST"
  DELETE_METHOD = "DELETE"
  PUT_METHOD = "PUT"
  HEAD_METHOD = "HEAD"
  ALL_METHODS = "*"

  ANY_PATH = "(.*)"
)

type Router struct {
  Routes []Route
}

type Route struct {
  method string
  path *regexp.Regexp
  handler RouteHandler
}

type RouteHandler func(*Request, *Response, map[string]interface{}, func(error))

func (router *Router) AddRoute(method string, path string, routeHandler RouteHandler) {
  pathRegexp, err := regexp.Compile("^" + path + "$")
  if err != nil { panic(err) }
  router.Routes = append(router.Routes, Route{ method: method, path: pathRegexp, handler: routeHandler })
}

func (router *Router) Execute(req *Request, res *Response, env map[string]interface{}, nextMiddleware func(error)) {
  var next func(error)
  routes := router.Routes
  maxIndex := len(routes)
  nextIndex := 0
  next = func (err error) {
    if err != nil {
      nextMiddleware(err)
    } else if nextIndex < maxIndex {
      currentIndex := nextIndex
      nextIndex++
      route := routes[currentIndex]
      if (route.method == req.Method || route.method == ALL_METHODS) && (route.path.MatchString(req.URL.Path)) {
        route.handler(req, res, env, next)
      } else {
        next(nil)
      }
    } else {
      nextMiddleware(nil)
    }
  }
  next(nil)
}

func NewRouter() *Router {
  return &Router{ Routes: make([]Route, 0, 5) }
}
