package soggy

import (
)

type Env map[string]interface{}

func NewEnv() *Env {
  return &Env{}
}
