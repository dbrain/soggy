package soggy

import (
  "testing"
)

func TestSaneURLPath(t *testing.T) {
  path := SaneURLPath("/wee")
  if path != "/wee/" {
    t.Error("expected a trailing slash to be added")
  }
  path = SaneURLPath("/wee/")
  if path != "/wee/" {
    t.Error("expected the path to be unchanged")
  }
}
