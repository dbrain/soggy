package soggy

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

type RouteHandler func(*Context)

func (router *Router) AddRoute(method string, path string, routeHandler RouteHandler) {
  pathRegexp, err := regexp.Compile("^" + SaneURLPath(path) + "$")
  if err != nil { panic(err) }
  router.Routes = append(router.Routes, Route{ method: method, path: pathRegexp, handler: routeHandler })
}

func (router *Router) Execute(middlewareCtx *Context) {
  var next func(error)
  var context *Context

  method := middlewareCtx.Req.Method
  relativePath := middlewareCtx.Req.RelativePath

  routes := router.Routes
  maxIndex := len(routes)
  nextIndex := 0
  next = func (err error) {
    if err != nil {
      middlewareCtx.Next(err)
    } else if nextIndex < maxIndex {
      currentIndex := nextIndex
      nextIndex++
      route := routes[currentIndex]
      if (route.method == method || route.method == ALL_METHODS) && (route.path.MatchString(relativePath)) {
        route.handler(context)
      } else {
        next(nil)
      }
    } else {
      middlewareCtx.Next(nil)
    }
  }

  context = &Context{ middlewareCtx.Req, middlewareCtx.Res, middlewareCtx.Env, next }
  next(nil)
}

func NewRouter() *Router {
  return &Router{ Routes: make([]Route, 0, 5) }
}
