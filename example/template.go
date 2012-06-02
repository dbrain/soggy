package main

import (
  "soggy"
  "errors"
)

func main() {
  app, server := soggy.NewDefaultApp()
  server.Get("/echo/(.*)", func (echo string) (string, map[string]string) {
    return "template_example.html", map[string]string{
      "echo": echo }
  })
  server.Get("/break", func () error {
    return errors.New("This be broked.")
  })
  server.Use(server.Router)
  app.Listen("0.0.0.0:9999")
}
