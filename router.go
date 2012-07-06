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

type SoggyRouter interface {
  Middleware
  AddRoute(method string, path string, handler ...interface{})
}

type Router struct {
  RouteBundles []*RouteBundle
}

type RouteBundle struct {
  method string
  path *regexp.Regexp
  Routes []*Route
}

type Route struct {
  handler reflect.Value
  argCount int
  callType int
  returnType int
  returnHasError bool
  returnHasStatusCode bool
}

var contextType = reflect.TypeOf(Context{})
var requestType = reflect.TypeOf(http.Request{})
var errorType = reflect.TypeOf((*error)(nil)).Elem()
var httpHandlerType = reflect.TypeOf((*http.Handler)(nil)).Elem()
var httpResponseWriterType = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()

func (route *Route) CacheCallType(routePath *regexp.Regexp) {
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

  if handlerType.Kind() != reflect.Func {
    panic("Route handlers must be a http.Handler or a func. Broken for " + routePath.String())
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
  outSkip := 0
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

  hasStatusCode := handlerType.Out(0).Kind() == reflect.Int
  route.returnHasStatusCode = hasStatusCode
  if hasStatusCode {
    outCount--
    outSkip++
  }
  if (outCount == 0) {
    route.returnType = RETURN_TYPE_EMPTY
    return
  }

  if outCount > 2 {
    panic("Handler has more return values than expected.")
  } else if outCount == 2 {
    route.returnType = RETURN_TYPE_RENDER
    return
  } else if handlerType.Out(outSkip).Kind() == reflect.String {
    route.returnType = RETURN_TYPE_STRING
    return
  }

  route.returnType = RETURN_TYPE_JSON
}

func (routeBundle *RouteBundle) CallBundle(ctx *Context, relativePath string) {
  routes := routeBundle.Routes
  if len(routes) == 1 {
    routes[0].CallHandler(ctx, routeBundle.path, relativePath)
  } else {
    var next func(interface{})
    var routeCtx *Context
    nextIndex := 0
    next = func (err interface{}) {
      if err != nil {
        ctx.Next(err)
      } else if nextIndex < len(routes) {
        currentIndex := nextIndex
        nextIndex++
        routes[currentIndex].CallHandler(routeCtx, routeBundle.path, relativePath)
      }
    }
    routeCtx = &Context{ ctx.Req, ctx.Res, ctx.Server, ctx.Env, next }
    next(nil)
  }
}

func (route *Route) CallHandler(ctx *Context, routePath *regexp.Regexp, relativePath string) {
  var args []reflect.Value
  callType := route.callType

  urlParams := routePath.FindStringSubmatch(relativePath)[1:]
  ctx.Req.URLParams = urlParams

  switch callType {
  case CALL_TYPE_HANDLER_FUNC:
    args = []reflect.Value{ reflect.ValueOf(ctx.Res), reflect.ValueOf(ctx.Req.Request) }
  case CALL_TYPE_CTX_ONLY, CALL_TYPE_CTX_AND_PARAMS:
    args = append(args, reflect.ValueOf(ctx))
  }

  if callType == CALL_TYPE_PARAMS_ONLY || callType == CALL_TYPE_CTX_AND_PARAMS {
    for _, param := range urlParams {
      args = append(args, reflect.ValueOf(param))
    }
  }

  if len(args) < route.argCount {
    log.Println("Route", routePath.String(), "expects", route.argCount, "arguments but only got", len(args), ". Padding.")
    for len(args) < route.argCount {
      args = append(args, reflect.ValueOf(""))
    }
  } else if len(args) > route.argCount {
    log.Println("Route", routePath.String(), "expects", route.argCount, "arguments but got", len(args), ". Trimming.")
    args = args[:route.argCount]
  }

  result, err := route.safelyCall(args, routePath)
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

  statusCode := http.StatusOK
  if route.returnHasStatusCode {
    statusCode = result[0].Interface().(int)
    result = result[1:len(result)]
  }

  // If the return is the zero value for the type its not rendered
  // This is to allow routes to pass control on without rendering when control comes back.
  // This may come back to bite me or cause weird issues for users. Document well.
  switch route.returnType {
    case RETURN_TYPE_EMPTY:
      return nil
    case RETURN_TYPE_RENDER:
      if template := result[0].String(); template != "" {
        return ctx.Res.Render(statusCode, template, result[1].Interface())
      }
    case RETURN_TYPE_STRING:
      if html := result[0].String(); html != "" {
        return ctx.Res.Html(statusCode, html)
      }
    case RETURN_TYPE_JSON:
      if !result[0].IsNil() {
        return ctx.Res.Json(statusCode, result[0].Interface())
      }
  }
  return nil
}

func (route *Route) safelyCall(args []reflect.Value, routePath *regexp.Regexp) (result []reflect.Value, err interface{}) {
  defer func() {
    if err = recover(); err != nil {
      log.Println("Handler for route", routePath.String(), "paniced with", err)
    }
  }()
  return route.handler.Call(args), err
}

func (router *Router) AddRoute(method string, path string, handlers ...interface{}) {
  rawRegex := "^" + SaneURLPath(path) + "$"
  routeRegex, err := regexp.Compile(rawRegex)
  if err != nil {
    log.Println("Could not compile route regex", rawRegex, ":", err)
    return
  }
  routeBundle := &RouteBundle{ method, routeRegex, make([]*Route, 0, 1) }
  for _, handler := range handlers {
    handlerValue := reflect.ValueOf(handler)
    route := &Route{ handler: handlerValue }
    route.CacheCallType(routeRegex)
    route.CacheReturnType()
    routeBundle.Routes = append(routeBundle.Routes, route)
  }

  router.RouteBundles = append(router.RouteBundles, routeBundle)
}

func (router *Router) findRoute(method, relativePath string, start int) (*RouteBundle, int) {
  routeBundles := router.RouteBundles
  for i := start; i < len(routeBundles); i++ {
    route := routeBundles[i];
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

    var routeBundle *RouteBundle
    routeBundle, routeIndex = router.findRoute(method, relativePath, routeIndex)
    if routeBundle != nil {
      routeBundle.CallBundle(context, relativePath)
    } else {
      middlewareCtx.Next(nil)
    }
  }

  context = &Context{ middlewareCtx.Req, middlewareCtx.Res, middlewareCtx.Server, middlewareCtx.Env, next }
  next(nil)
}

func NewRouter() *Router {
  return &Router{ RouteBundles: make([]*RouteBundle, 0, 5) }
}
