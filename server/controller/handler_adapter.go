package controller

import (
    "fmt"
    "time"
    "gopkg.in/gin-gonic/gin.v1"
    "github.com/dnp1/conversa/server/errors"
    "net/http"
    "github.com/dnp1/conversa/server/handlers"
)

func WrapChannelContext(f func(context handlers.ChannelContext)) gin.HandlerFunc{
    return func(c *gin.Context) {
        f(c)
    }
}
func WrapMiddleware(f func (req handlers.Context, resp handlers.JsonResponse)) gin.HandlerFunc {
    return func(c *gin.Context) {
        var resp  = handlers.NewResponse()
        f(&contextAdapter{c}, resp)
        if c.IsAborted() {
            resp.WriteJSON(c)
        }
    }
}
func WrapContext(f func (req handlers.Context, resp handlers.JsonResponse)) gin.HandlerFunc {
    return func(c *gin.Context) {
        var resp  = handlers.NewResponse()
        f(&contextAdapter{c}, resp)
        resp.WriteJSON(c)
    }
}

type contextAdapter struct {
    *gin.Context
}

func (context *contextAdapter) GetString(name string) (string, bool) {
    if iValue, ok := context.Get(name); !ok {
        return "", false
    } else if value, ok := iValue.(string); !ok {
        return "", false
    } else {
        return value, true
    }
}

func  (context *contextAdapter)  ShouldGetString(name string) (string, errors.Error) {
    if str, ok := context.GetString(name); ok {
        return str, nil
    } else {
        var msg = fmt.Errorf("Context %q should be setted.", name)
        return "", errors.Internal(msg)
    }
}


func (context *contextAdapter) DeleteCookie(name string) {
    cookie := http.Cookie{Name:name, Expires:time.Now().Add(-1 * 24 * time.Hour), Value: "deleted"}
    http.SetCookie(context.Writer, &cookie)
}

func (context *contextAdapter) SetCookie(cookie *http.Cookie) {
    http.SetCookie(context.Writer, cookie)
}

func (context *contextAdapter) BindJSON(data interface{}) errors.Error {
    if err := context.Context.BindJSON(data); err != nil {
        return errors.Validation(err)
    }
    return nil
}