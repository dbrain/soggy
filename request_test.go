package soggy

import (
  "testing"
  "net/http"
)

func TestNewRequestTrivial(t *testing.T) {
  request := NewRequest(&http.Request{})
  if request == nil {
    t.Error("request is nil")
  }
}

func TestSetRelativePath(t *testing.T) {
  request := NewRequest(&http.Request{})
  if request == nil {
    t.Error("request is nil")
  }
  request.SetRelativePath("/test/", "/test/flowers/yes/")
  if request.RelativePath != "/flowers/yes/" {
    t.Error("expected RelativePath to be /flowers/yes/ but was " + request.RelativePath)
  }
}
