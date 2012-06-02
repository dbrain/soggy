package main

import (
  "soggy"
  "html/template"
  "io"
)

func main() {
  app, server := soggy.NewDefaultApp()
  server.Get("/(.*)", func (echo string) (string, map[string]string) {
    return "template_example.html", map[string]string{
      "echo": echo }
  })
  server.EngineFunc("html", func(writer io.Writer, filename string, options interface{}) error {
    template, err := template.ParseFiles(filename)
    if err != nil {
      return err
    }
    return template.Execute(writer, options)
  })
  server.Use(server.Router)
  app.Listen("0.0.0.0:9999")
}
