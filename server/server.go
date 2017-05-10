package main

import (
    "fmt"
    "time"
    "net/http"
)

func newServer() *http.Server {
    host := env("HOST", "0.0.0.0")
    port := env("PORT", "5001")
    srv := &http.Server{
        Addr: fmt.Sprintf("%s:%s", host, port),
        ReadTimeout: 60 * time.Second,
        WriteTimeout: 60 * time.Second,
        ReadHeaderTimeout: 10 * time.Second,
        MaxHeaderBytes: 1 << 11,
        Handler: NewRouter(),
    }
    return srv
}