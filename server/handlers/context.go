package handlers

import (
    "github.com/dnp1/conversa/server/errors"
    "net/http"
    "io"
)

type Context interface {
    GetString(name string) (string, bool)
    ShouldGetString(name string) (string, errors.Error)
    Param(key string) string
    GetQuery(string) (string, bool)
    SetCookie(cookie *http.Cookie)
    DeleteCookie(name string)
    BindJSON(data interface{}) errors.Error
    Cookie(string) (string, error)
    Next()
    Abort()
    Set(name string, value interface{})
}

type ChannelContext interface {
    Param(string) string
    SSEvent(name string, message interface{})
    Stream(func(w io.Writer) bool)
}
