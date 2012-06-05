package main

import (
  "github.com/dbrain/soggy"
  "log"
  "errors"
  "fmt"
  "net/http"
  "io"
)

type HandlerExample struct {}
func (handlerEx *HandlerExample) ServeHTTP(res http.ResponseWriter, req *http.Request) {
  res.Write([]byte("Oh what a lovely handler"))
}

type MiddlewareExample struct {}
func (middleware *MiddlewareExample) Execute(ctx *soggy.Context) {
  log.Println("I've hit the custom middleware")
  // Call the next middleware, if this isn't called the request will end
  ctx.Next(nil)
}

type TemplateExample struct {
  Name string
  Age int
}

func main() {
  app, server := soggy.NewDefaultApp()

  // You can use normal handlers as routes
  server.Get("/handler", &HandlerExample{})

  // Or handler funcs
  server.Get("/handlerFunc", func(res http.ResponseWriter, req *http.Request) {
    res.Write([]byte("Oh what a lovely handler func"))
  })

  // URLParams can be read from the Req
  server.Get("/reqParams/(.*)/(.*)", func (ctx *soggy.Context) string {
    return fmt.Sprintln("Req params are", ctx.Req.URLParams[0], "and", ctx.Req.URLParams[1])
  })

  // URLParams will also be passed in to the function
  server.Get("/params/(.*)/(.*)", func (ctx *soggy.Context, param1 string, param2 string) string {
    return fmt.Sprintln("Params are", param1, "and", param2)
  })

  // You don't need Context, and If theres more params than URLParams theyll be blank
  server.Get("/moreParams/(.*)/(.*)", func (param1 string, param2 string, param3 string) string {
    return fmt.Sprintln("moreParams are", param1, "and", param2, "missing param3", param3)
  })

  // If theres less params than URLParams soggy will put in as much as it can
  server.Get("/lessParams/(.*)/(.*)", func (ctx *soggy.Context, param1 string) string {
    return fmt.Sprintln("lessParams are", param1, "param2 on URLParams is", ctx.Req.URLParams[1])
  })

  // Marshal non string or render returns to JSON
  server.Get("/returnJSON/(.*)", func (name string) *TemplateExample {
    return &TemplateExample{name, 27}
  })

  // Output string returns as HTML
  server.Get("/returnHTML", func () string {
    return `<html><body>I'm html</body></html>`
  })

  // Render templates (relative to config.ViewPath)
  server.Get("/render/(.*)", func (name string) (string, interface{}) {
    return "kitchensink.html", &TemplateExample{name, 256}
  })

  // No return expects you to write to ctx.Res yourself
  server.Get("/writeYourself", func (ctx *soggy.Context) {
    res := ctx.Res
    res.WriteHeader(http.StatusOK)
    res.Set("Content-Type", "text/plain")
    res.WriteString("Cannn dooo.")
  })

  // The last return can be an error which will be the same as calling ctx.Next(error)
  server.Get("/returnNext", func () error {
    return errors.New("This was supposed to happen")
  })

  // Next with error should hit error handler
  server.Get("/nextError", func (ctx *soggy.Context) {
    ctx.Next("Ice ice baby")
  })

  // Calling next will continue through the routes looking for another match
  server.Get("/keepOnNexting", func (ctx *soggy.Context) {
    ctx.Next(nil)
  })

  // Panics will also hit the error handler
  server.Get("/panic", func () {
    panic("This panic was supposed to hit an error handler")
  })

  // Override the default error handler
  server.ErrorHandler = func(ctx *soggy.Context, err interface{}) {
    res := ctx.Res
    res.WriteHeader(http.StatusInternalServerError)
    switch err.(type) {
      case error:
        res.Write([]byte("Overriden ErrorHandler, err: " + err.(error).Error()))
      case string:
        res.Write([]byte("Overriden ErrorHandler, string: " + err.(string)))
      default:
        res.Write([]byte("Overriden ErrorHandler, default: An error occured processing your request"))
    }
  }

  // Handle any requests that haven't matched the above.
  server.All(soggy.ANY_PATH, func (context *soggy.Context) {
    res := context.Res
    res.Header().Set("Content-Type", "text/plain")
    res.WriteString("404 Page for " + context.Req.RelativePath)
  })

  // Add some custom middleware
  server.Use(&MiddlewareExample{})

  // Add a new view handler (will be used with a template return with a filename ending in json)
  server.EngineFunc("json", func(writer io.Writer, filename string, options interface{}) error {
    writer.Write([]byte("Nothing happens here"))
    return nil
  })

  server.Use(&soggy.RequestLoggerMiddleware{}, server.Router)
  app.Listen("0.0.0.0:9999")
}
