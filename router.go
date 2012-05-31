package soggy

import (
  "regexp"
  "log"
  "reflect"
)

const (
  GET_METHOD = "GET"
  POST_METHOD = "POST"
  DELETE_METHOD = "DELETE"
  PUT_METHOD = "PUT"
  HEAD_METHOD = "HEAD"
  ALL_METHODS = "*"

  ANY_PATH = "(.*)"

  CALL_TYPE_EMPTY = iota
  CALL_TYPE_CTX_ONLY
  CALL_TYPE_CTX_AND_PARAMS
  CALL_TYPE_PARAMS_ONLY

  RETURN_TYPE_EMPTY = iota
  RETURN_TYPE_STRING
)

type Router struct {
  Routes []*Route
}

type Route struct {
  method string
  path *regexp.Regexp
  handler reflect.Value
  callType int
  returnType int
}

var contextType = reflect.TypeOf(Context{})

func DetermineCallType(handlerType reflect.Type) int {
  argCount := handlerType.NumIn()
  if argCount == 0 {
    return CALL_TYPE_EMPTY
  }

  firstArg := handlerType.In(0)
  if firstArg.Kind() == reflect.Ptr && firstArg.Elem() == contextType {
    if argCount > 1 {
      return CALL_TYPE_CTX_AND_PARAMS
    } else {
      return CALL_TYPE_CTX_ONLY
    }
  }

  return CALL_TYPE_PARAMS_ONLY
}

func (route Route) CallHandler(ctx *Context, relativePath string) {
  // match := route.path.FindStringSubmatch(relativePath)
  var args []reflect.Value
  args = append(args, reflect.ValueOf(ctx))
  route.handler.Call(args)
}

func (router *Router) AddRoute(method string, path string, handler interface{}) {
  rawRegex := "^" + SaneURLPath(path) + "$"
  routeRegex, err := regexp.Compile(rawRegex)
  if err != nil {
    log.Println("Could not compile route regex", rawRegex, err)
    return
  }
  handlerValue := reflect.ValueOf(handler)
  router.Routes = append(router.Routes, &Route{
    method: method, path: routeRegex, handler: handlerValue, callType: DetermineCallType(handlerValue.Type()) })
}

func (router *Router) findRoute(method, relativePath string, start int) (*Route, int) {
  routes := router.Routes
  for i := start; i < len(routes); i++ {
    route := routes[i];
    if route.method == method || route.method == ALL_METHODS {
      if route.path.MatchString(relativePath) {
        return route, i + 1
      }
    }
  }
  return nil, -1
}

func (router *Router) Execute(middlewareCtx *Context) {
  var next func(error)
  var context *Context

  method := middlewareCtx.Req.Method
  relativePath := middlewareCtx.Req.RelativePath

  n := 0
  next = func (err error) {
    if err != nil {
      middlewareCtx.Next(err)
      return
    }

    var route *Route
    route, n = router.findRoute(method, relativePath, n)
    if route != nil {
      route.CallHandler(context, relativePath)
    } else {
      middlewareCtx.Next(nil)
    }
  }

  context = &Context{ middlewareCtx.Req, middlewareCtx.Res, middlewareCtx.Server, middlewareCtx.Env, next }
  next(nil)
}

func NewRouter() *Router {
  return &Router{ Routes: make([]*Route, 0, 5) }
}
