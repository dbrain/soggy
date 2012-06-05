package soggy

import (
  "regexp"
  "log"
  "reflect"
  "net/http"
)

const (
  CALL_TYPE_EMPTY = iota
  CALL_TYPE_CTX_ONLY
  CALL_TYPE_CTX_AND_PARAMS
  CALL_TYPE_PARAMS_ONLY
  CALL_TYPE_HANDLER_FUNC
)

const (
  RETURN_TYPE_EMPTY = iota
  RETURN_TYPE_STRING
  RETURN_TYPE_JSON
  RETURN_TYPE_RENDER
)

const (
  GET_METHOD = "GET"
  POST_METHOD = "POST"
  DELETE_METHOD = "DELETE"
  PUT_METHOD = "PUT"
  HEAD_METHOD = "HEAD"
  ALL_METHODS = "*"
)

const (
  ANY_PATH = "(.*)"
)

type Router struct {
  Routes []*Route
}

type Route struct {
  method string
  path *regexp.Regexp
  handler reflect.Value
  argCount int
  callType int
  returnType int
  returnHasError bool
}

var contextType = reflect.TypeOf(Context{})
var requestType = reflect.TypeOf(http.Request{})
var errorType = reflect.TypeOf((*error)(nil)).Elem()
var httpHandlerType = reflect.TypeOf((*http.Handler)(nil)).Elem()
var httpResponseWriterType = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()

func (route *Route) CacheCallType() {
  handlerType := route.handler.Type()

  if (handlerType.Kind() == reflect.Ptr && handlerType.Elem().Implements(httpHandlerType)) || handlerType.Implements(httpHandlerType) {
    httpHandler := route.handler.Interface().(http.Handler)
    route.handler = reflect.ValueOf(func (res http.ResponseWriter, req *http.Request) {
      httpHandler.ServeHTTP(res, req)
    })
    route.callType = CALL_TYPE_HANDLER_FUNC
    route.argCount = 2
    return
  }

  argCount := handlerType.NumIn()
  route.argCount = argCount
  if argCount == 0 {
    route.callType = CALL_TYPE_EMPTY
    return
  }

  firstArg := handlerType.In(0)
  if firstArg.Kind() == reflect.Ptr && firstArg.Elem() == contextType {
    if argCount > 1 {
      route.callType = CALL_TYPE_CTX_AND_PARAMS
    } else {
      route.callType = CALL_TYPE_CTX_ONLY
    }
    return
  }

  if argCount == 2 {
    secondArg := handlerType.In(1)
    if firstArg.Implements(httpResponseWriterType) && secondArg.Kind() == reflect.Ptr && secondArg.Elem() == requestType {
      route.callType = CALL_TYPE_HANDLER_FUNC
      return
    }
  }

  route.callType = CALL_TYPE_PARAMS_ONLY
}

func (route *Route) CacheReturnType() {
  handlerType := route.handler.Type()
  outCount := handlerType.NumOut();
  if outCount == 0 {
    route.returnType = RETURN_TYPE_EMPTY
    route.returnHasError = false
    return
  }

  hasError := handlerType.Out(outCount - 1) == errorType
  route.returnHasError = hasError
  if hasError { outCount-- }
  if outCount == 0 {
    route.returnType = RETURN_TYPE_EMPTY
    return
  }

  if outCount > 2 {
    panic("Handler has more return values than expected.")
  } else if outCount == 2 {
    route.returnType = RETURN_TYPE_RENDER
    return
  } else if handlerType.Out(0).Kind() == reflect.String {
    route.returnType = RETURN_TYPE_STRING
    return
  }

  route.returnType = RETURN_TYPE_JSON
  return
}

func (route *Route) CallHandler(ctx *Context, relativePath string) {
  var args []reflect.Value
  callType := route.callType

  urlParams := route.path.FindStringSubmatch(relativePath)[1:]
  ctx.Req.URLParams = urlParams

  switch callType {
  case CALL_TYPE_HANDLER_FUNC:
    args = []reflect.Value{ reflect.ValueOf(ctx.Res), reflect.ValueOf(ctx.Req.OriginalRequest) }
  case CALL_TYPE_CTX_ONLY, CALL_TYPE_CTX_AND_PARAMS:
    args = append(args, reflect.ValueOf(ctx))
  }

  if callType == CALL_TYPE_PARAMS_ONLY || callType == CALL_TYPE_CTX_AND_PARAMS {
    for _, param := range urlParams {
      args = append(args, reflect.ValueOf(param))
    }
  }

  if len(args) < route.argCount {
    log.Println("Route", route.path.String(), "expects", route.argCount, "arguments but only got", len(args), ". Padding.")
    for len(args) < route.argCount {
      args = append(args, reflect.ValueOf(""))
    }
  } else if len(args) > route.argCount {
    log.Println("Route", route.path.String(), "expects", route.argCount, "arguments but got", len(args), ". Trimming.")
    args = args[:route.argCount]
  }

  result, err := route.safelyCall(args)
  if err != nil {
    ctx.Next(err)
    return
  }

  err = route.renderResult(ctx, result)
  if err != nil {
    ctx.Next(err)
  }
}

func (route *Route) renderResult(ctx *Context, result []reflect.Value) interface{} {
  if route.returnHasError {
    err := result[len(result)-1]
    if !err.IsNil() {
      return err.Interface()
    }
  }

  switch route.returnType {
    case RETURN_TYPE_EMPTY:
      return nil
    case RETURN_TYPE_RENDER:
      return ctx.Res.Render(http.StatusOK, result[0].String(), result[1].Interface())
    case RETURN_TYPE_STRING:
      return ctx.Res.Html(http.StatusOK, result[0].String())
    case RETURN_TYPE_JSON:
      return ctx.Res.Json(http.StatusOK, result[0].Interface())
  }
  return nil
}

func (route *Route) safelyCall(args []reflect.Value) (result []reflect.Value, err interface{}) {
  defer func() {
    if err = recover(); err != nil {
      log.Println("Handler for route", route.path.String(), "paniced with", err)
    }
  }()
  return route.handler.Call(args), err
}

func (router *Router) AddRoute(method string, path string, handler interface{}) {
  rawRegex := "^" + SaneURLPath(path) + "$"
  routeRegex, err := regexp.Compile(rawRegex)
  if err != nil {
    log.Println("Could not compile route regex", rawRegex, ":", err)
    return
  }
  handlerValue := reflect.ValueOf(handler)
  route := &Route{ method: method, path: routeRegex, handler: handlerValue }
  route.CacheCallType()
  route.CacheReturnType()

  router.Routes = append(router.Routes, route)
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
  var next func(interface{})
  var context *Context

  method := middlewareCtx.Req.Method
  relativePath := middlewareCtx.Req.RelativePath

  routeIndex := 0
  next = func (err interface{}) {
    if err != nil {
      middlewareCtx.Next(err)
      return
    }

    var route *Route
    route, routeIndex = router.findRoute(method, relativePath, routeIndex)
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
