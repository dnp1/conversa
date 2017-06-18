package authentication

import (
    "github.com/dnp1/conversa/server/handlers"
    "github.com/dnp1/conversa/server/errors"
    "github.com/dnp1/conversa/server/data/session"
)

func New(cookieName string, model Model) *middleware {
    return &middleware{
        cookieName: cookieName,
        model:      model,
    }
}

type Model interface {
    Retrieve(string) (*session.Data, errors.Error)
}

type middleware struct {
    cookieName string
    model      Model
}

func (h *middleware) Middleware(context handlers.Context, resp handlers.JsonResponse) {
    if token, err := context.Cookie(h.cookieName); err != nil {
        resp.SetError(errors.Authentication(err))
        context.Abort()
    } else if data, err := h.model.Retrieve(token); err != nil {
        if err.Empty() {
            err = errors.Authentication(err)
        }
        resp.SetError(err)
        context.Abort()
    } else{
        context.Set("username", data.Username)
        context.Next()
    }
}
