package user

import (
    "net/http"
    "github.com/dnp1/conversa/server/handlers"
    "github.com/dnp1/conversa/server/errors"
)


func New(model Model) *handler {
    return &handler{model: model}
}

type Model interface {
    Create(username, password, passwordConfirmation string) errors.Error
}

type handler struct {
    model Model
}

func (user *handler) Create(context handlers.Context, resp handlers.JsonResponse) {
    var body struct {
        Username string `json:"username"`
        Password string `json:"password"`
        PasswordConfirmation string `json:"passwordConfirmation"`
    }

    if err := context.BindJSON(&body); err != nil {
        resp.SetError(err)
    } else if err:= user.model.Create(body.Username, body.Password, body.PasswordConfirmation); err != nil {
        resp.SetError(err)
    } else {
        const msg = "user created with success"
        resp.SetMessage(msg)
        resp.SetStatus(http.StatusCreated)
    }
}