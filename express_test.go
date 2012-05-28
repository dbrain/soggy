package express

import (
    "testing"
)

func TestNewDefaultAppTrivial(t *testing.T) {
    app, server := NewDefaultApp()
    if app == nil {
        t.Error("app is nil")
    }
    if server == nil {
        t.Error("server is nil")
    }
}

func TestNewAppTrivial(t *testing.T) {
    app := NewApp()
    if app == nil {
        t.Error("app is nil")
    }
}

func TestAddServerToApp(t *testing.T) {
    app := NewApp()
    if len(app.servers) != 0 {
        t.Error("expected zero servers")
    }
    app.AddServer(NewServer("/muffin"))
    if len(app.servers) != 1 {
        t.Error("expected one server")
    }
}

