# soggy
  Fast, simple web framework for [go](http://golang.org). Based on [express](https://github.com/visionmedia/express/) and [web.go](https://github.com/hoisie/web).
  
```go
app, server := soggy.NewDefaultApp()
server.Get("/echo/(.*)", func (echo string) (string, map[string]string) {
  return "template_example.html", map[string]string{ "echo": echo }
})
server.Use(server.Router)
app.Listen("0.0.0.0:9999")
```

## Features
  * Routing
  * Middleware
  * HTTP helpers
  * Server mounting
  * Easily pluggable view system
  * Speed of Go

## Installation

### Requirements
  * Go 1.0.1 (tested), instructions [here](http://golang.org/doc/install.html).

### Steps

    $ go get github.com/dbrain/soggy
    
## Examples
  See [examples](https://github.com/dbrain/soggy/tree/master/example).