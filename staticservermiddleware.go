package soggy

import (
  "net/http"
  "path/filepath"
  "strings"
  "os"
)

type StaticServerMiddleware struct {
  path string
}

func (staticServer *StaticServerMiddleware) GetRelativeFilePath(urlPath string) string {
  if staticServer.path == "/" {
    return urlPath
  }
  return urlPath[len(staticServer.path)-1:len(urlPath)]
}

func (staticServer *StaticServerMiddleware) IsValidForPath(path string) bool {
  return strings.HasPrefix(path, staticServer.path)
}

func (staticServer *StaticServerMiddleware) Execute(ctx *Context) {
  req := ctx.Req
  if (req.Method == GET_METHOD || req.Method == HEAD_METHOD) && staticServer.IsValidForPath(req.RelativePath) {
    staticPath := ctx.Server.Config[CONFIG_STATIC_PATH].(string)
    staticFile := filepath.Join(staticPath, staticServer.GetRelativeFilePath(req.RelativePath))
    if stat, err := os.Stat(staticFile); err == nil && !stat.IsDir() {
      http.ServeFile(ctx.Res, req.OriginalRequest, staticFile)
    } else {
      ctx.Next(nil)
    }
  } else {
    ctx.Next(nil)
  }
}

func NewStaticServerMiddleware(path string) *StaticServerMiddleware {
  return &StaticServerMiddleware{ SaneURLPath(path) }
}
