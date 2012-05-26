package main

import (
  "./lib"
  "github.com/paulbellamy/mango"
  "encoding/json"
  "io/ioutil"
  "os"
  "fmt"
)
var config map[string]interface{}
var jsonHeaders = mango.Headers{"Content-Type": []string{"application/json; charset=utf-8"}}

func Default(env mango.Env) (mango.Status, mango.Headers, mango.Body) {
  return 404, mango.Headers{}, mango.Body("Not found.")
}

func Health(env mango.Env) (mango.Status, mango.Headers, mango.Body) {
  return 200, jsonHeaders, `{ "status": "ok" }`
}

func Commands(env mango.Env) (mango.Status, mango.Headers, mango.Body) {
  commands, err := json.Marshal(config["commands"])
  if err != nil { return 500, jsonHeaders, `{ "error": "Exploded trying to parse the commands JSON" }` }

  return 200, jsonHeaders, mango.Body(string(commands))
}

func ExecuteFunc(command map[string]interface{}) func(mango.Env) (mango.Status, mango.Headers, mango.Body) {
  return func (env mango.Env) (mango.Status, mango.Headers, mango.Body) {
    return 200, mango.Headers{}, mango.Body("unimplemented")
  }
}

func main() {
  configContents, err := ioutil.ReadFile(os.Getenv("HOME") + "/.commander/config.json")
  if err != nil { panic(err) }

  var unmarshalledConfig interface{}
  err = json.Unmarshal(configContents, &unmarshalledConfig)
  if err != nil { panic(err) }

  config = unmarshalledConfig.(map[string]interface{})
  fmt.Println("Loaded config for " + config["name"].(string))

  app := express.NewServer();
  app.Get("/", express.RouteHandler(func (req *express.Request, res *express.Response, env map[string]interface{}, next func(error)) {
    res.Header().Set("Content-Type", "text/plain")
    res.Write([]byte("This is an example server. Hell yeah."))
  }))
  app.All(express.ANY_PATH, express.RouteHandler(func (req *express.Request, res *express.Response, env map[string]interface{}, next func(error)) {
    res.Header().Set("Content-Type", "text/plain")
    res.Write([]byte("404 Page would go here"))
  }))
  app.Use(express.RequestLogger, app.Router)
  app.Listen(config["host"].(string) + ":" + config["port"].(string))

  // routes := map[string]mango.App{
  //   "/health$": Health,
  //   "/commands(|/([^/]+))$": Commands }

  // stack := new(mango.Stack)
  // stack.Address = config["host"].(string) + ":" + config["port"].(string)
  // stack.Middleware(mango.Routing(routes))
  // stack.Run(Default)
}
