package main

import (
  "./lib"
  "./lib/middleware"
  "encoding/json"
  "io/ioutil"
  "os"
  "log"
)
var config map[string]interface{}

func main() {
  configContents, err := ioutil.ReadFile(os.Getenv("HOME") + "/.commander/config.json")
  if err != nil { panic(err) }

  var unmarshalledConfig interface{}
  err = json.Unmarshal(configContents, &unmarshalledConfig)
  if err != nil { panic(err) }

  config = unmarshalledConfig.(map[string]interface{})
  log.Println("Loaded config for " + config["name"].(string))

  server := express.NewServer("/web");
  server.Get("/", express.RouteHandler(func (req *express.Request, res *express.Response, env *express.Env, next func(error)) {
    res.Header().Set("Content-Type", "text/plain")
    res.Write([]byte("This is an example server. Hell yeah."))
  }))
  server.All(express.ANY_PATH, express.RouteHandler(func (req *express.Request, res *express.Response, env *express.Env, next func(error)) {
    res.Header().Set("Content-Type", "text/plain")
    res.Write([]byte("404 Page would go here"))
  }))
  server.Use(middleware.RequestLogger, server.Router)

  app := express.NewApp()
  app.AddServer(server)
  app.AddServer(express.NewServer("/abc"))
  app.AddServer(express.NewServer("/abc123"))

  app.Listen(config["host"].(string) + ":" + config["port"].(string))

  // routes := map[string]mango.App{
  //   "/health$": Health,
  //   "/commands(|/([^/]+))$": Commands }

  // stack := new(mango.Stack)
  // stack.Address = config["host"].(string) + ":" + config["port"].(string)
  // stack.Middleware(mango.Routing(routes))
  // stack.Run(Default)
}
