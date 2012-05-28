package soggy

import (
    "testing"
)

func TestNewEnvTrivial(t *testing.T) {
    env := NewEnv()
    if env == nil {
        t.Error("NewEnv() failed")
    }
}

