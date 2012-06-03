package soggy

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
    app.AddServers(NewServer("/muffin"))
    if len(app.servers) != 1 {
        t.Error("expected one server")
    }
    app.AddServers(NewServer("/longerthanmuffin"))
    if (len(app.servers) != 2) {
        t.Error("expected 2 servers")
    }
    if (app.servers[0].Mountpoint != "/longerthanmuffin/") {
        t.Error("expected longest mountpoint server to be first")
    }
    if (app.servers[1].Mountpoint != "/muffin/") {
        t.Error("expected smallest mountpoint server to be last")
    }
}

