package soggy

import (
)

type Context struct {
  Req *Request
  Res *Response
  Env *Env
  Next func(error)
}
