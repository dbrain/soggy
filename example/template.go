package main

import (
  "soggy"
)

func main() {
  app, server := soggy.NewDefaultApp()
  server.Get("/(.*)", func (echo string) (string, map[string]string) {
    return "template_example.html", map[string]string{
      "echo": echo }
  })
  server.Use(server.Router)
  app.Listen("0.0.0.0:9999")
}
