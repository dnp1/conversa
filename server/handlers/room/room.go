package room

import (
    "net/http"
    "github.com/dnp1/conversa/server/errors"
    "github.com/dnp1/conversa/server/data/room"
    "github.com/dnp1/conversa/server/handlers"
)

func New(model Model) *handler {
    return &handler{
        model: model,
    }
}

type Model interface {
    All() ([]room.Data, errors.Error)
    Create(username, body string) errors.Error
    Delete(username, roomName string) errors.Error
    Rename(username, roomName, newRoomName string) errors.Error
}

type handler struct {
    model Model
}

func (room *handler) List(context handlers.Context, resp handlers.JsonResponse) {
    if data, err := room.model.All(); err != nil {
        resp.SetError(err)
    } else {
        const msg = "list of rooms"
        resp.SetMessage(msg)
        resp.SetData(data)
        resp.SetStatus(http.StatusOK)
    }
}

func (room *handler) Create(context handlers.Context, resp handlers.JsonResponse) {
    var body struct {
        Name string `json:"name"`
    }
    var user = context.Param("user")

    if username, err := context.ShouldGetString("username"); err != nil {
        resp.SetError(err)
    } else if user != username {
        const msg = "permission denied"
        resp.SetError(errors.Authorization(errors.FromString(msg)))
    } else if err := context.BindJSON(&body); err != nil {
        resp.SetError(err)
    } else if err := room.model.Create(user, body.Name); err != nil {
        resp.SetError(err)
    } else {
        const msg = "room created with success!"
        resp.SetMessage(msg)
        resp.SetStatus(http.StatusCreated)
    }
}

func (room *handler) Delete(context handlers.Context, resp handlers.JsonResponse) {
    var user = context.Param("user")
    var name = context.Param("room")
    if username, err := context.ShouldGetString("username"); err != nil {
        resp.SetError(err)
    } else if user != username {
        const msg = "permission denied"
        resp.SetError(errors.Authorization(errors.FromString(msg)))
    } else if err := room.model.Delete(user,name); err != nil {
        if err.Empty() {
            resp.SetStatus(http.StatusNoContent)
            resp.SetMessage("no room found for deletion.")
        } else {
            resp.SetError(err)
        }
    } else {
        const msg = "room deleted with success!"
        resp.SetMessage(msg)
        resp.SetStatus(http.StatusOK)
    }
}

func (room *handler) Edit(context handlers.Context, resp handlers.JsonResponse) {
    var body struct {
        Name string `json:"name"`
    }
    var user = context.Param("user")

    if username, err := context.ShouldGetString("username"); err != nil {
        resp.SetError(err)
    } else if user != username {
        const msg = "permission denied"
        resp.SetError(errors.Authorization(errors.FromString(msg)))
    } else if err := context.BindJSON(&body); err != nil {
        resp.SetError(err)
    } else if err := room.model.Rename(user, context.Param("room"), body.Name); err != nil {
        resp.SetError(err)
    } else {
        const msg = "room edited with success!"
        resp.SetMessage(msg)
        resp.SetStatus(http.StatusOK)
    }
}
