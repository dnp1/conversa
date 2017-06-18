package session

import (
    "net/http"
    "time"
    "github.com/dnp1/conversa/server/errors"
    "github.com/dnp1/conversa/server/handlers"
)


func New(model Model) *handler {
    return &handler{model: model}
}

type Model interface {
    Create(username, password string) (string, errors.Error)
    Delete(token string) errors.Error
}

type handler struct {
    cookieName string
    model      Model
}


func (session *handler) Login(context handlers.Context, resp handlers.JsonResponse) {
    var body struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := context.BindJSON(&body); err != nil {
        resp.SetError(err)
    } else if key, err := session.model.Create(body.Username, body.Password); err != nil {
        if err.Validation() || err.Empty() {
            err = errors.Authentication(err)
        }
        resp.SetError(err)
    } else {
        const msg = "sucessful sign-in"
        cookie := http.Cookie{Name: session.cookieName, Value: key, Expires: time.Now().Add(1 * 24 * time.Hour)}
        context.SetCookie(&cookie)
        resp.SetMessage(msg)
        resp.SetStatus(http.StatusCreated)
    }
}

func (session *handler) Logout(context handlers.Context, resp handlers.JsonResponse) {
    if token, err := context.Cookie(session.cookieName); err == http.ErrNoCookie  {
        resp.SetStatus(http.StatusNoContent)
        resp.SetMessage("You was not signed-in!")
    } else if err != nil {
        resp.SetError(errors.Internal(err))
    } else {
        if err := session.model.Delete(token); err != nil {
            if err.Empty() {
                const msg = "looks like your cookie is outdated!"
                resp.SetMessage(msg)
                resp.SetStatus(http.StatusResetContent)
                context.DeleteCookie(session.cookieName)
            } else {
                resp.SetError(err)
                return
            }
        } else {
            const msg = "sucessful sign-out"
            resp.SetStatus(http.StatusOK)
            resp.SetMessage(msg)
            context.DeleteCookie(session.cookieName)
        }
    }
}

