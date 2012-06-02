package soggy

import (
  "html/template"
  "io"
)

type HTMLTemplateEngine struct {
}

func (engine *HTMLTemplateEngine) SoggyEngine(writer io.Writer, filename string, options interface{}) error {
  template, err := template.ParseFiles(filename)
  if err != nil {
    return err
  }
  return template.Execute(writer, options)
}
