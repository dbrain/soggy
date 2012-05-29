package soggy

import (
  "testing"
)

func myHandler(ctx *Context) {
  ctx.Env["value"] = true
}

func TestAddRouteAppendsSlashToPath(t *testing.T) {
  router := NewRouter()
  router.AddRoute("GET", "/foo", myHandler)
  if len(router.Routes) == 0 {
    t.Error("where's my route?")
  }
  if router.Routes[0].path.String() != "^/foo/$" {
    t.Error("route does not match expected pattern")
  }
}

func TestFindRoute(t *testing.T) {
  router := NewRouter()
  router.AddRoute("GET", "/foo/", myHandler)
  handler, n := router.findRoute("GET", "/foo/", 0)
  if handler == nil {
    t.Error("expected to find a handler")
  }
  if n != 0 {
    t.Error("expected route index to be zero")
  }
}

//
// TODO currently too much of a PITA to stub out Request properly.
//      do this later.
//
/*
func TestExecuteRouteWithGoodRoute(t *testing.T) {
  router := NewRouter()
  router.AddRoute("GET", "/tehmuffin", myHandler)
  env := NewEnv()
  context := &Context{
    Req: NewStubRequest("GET", "/tehmuffin"), Res: nil, Env: env, Next: nil}
  router.Execute(context)
  if env["value"] == nil {
    t.Error("expected 'value' env variable to be set")
  }
}
*/

