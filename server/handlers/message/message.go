package message

import (
    "net/http"
    "github.com/dnp1/conversa/server/errors"
    "github.com/dnp1/conversa/server/handlers"
)

func New(model Model) *handler {
    return &handler{
        model: model,
    }
}

type Model interface {
    Create(username, roomName, senderName, content string) errors.Error
    Edit(username, roomName, messageOwner, msg, content string) errors.Error
    Delete(username, roomName, messageOwner, msg string) errors.Error
}


type handler struct {
    model Model
}

type Body struct {
    Content string `json:"content"`
}

func (message *handler) Create(req handlers.Context, resp handlers.JsonResponse) {
    var (
        username  = req.Param("user")
        roomName  = req.Param("room")
        body     Body
    )
    if senderName, err := req.ShouldGetString("username"); err != nil {
        resp.SetError(err)
    } else if err := req.BindJSON(&body); err != nil {
        resp.SetError(err)
    } else if err := message.model.Create(username, roomName, senderName, body.Content);
        err != nil {
        resp.SetError(err)
    } else {
        const msg = "message created with success"
        resp.SetMessage(msg)
        resp.SetStatus(http.StatusCreated)
    }
}

func (message *handler) Edit(req handlers.Context, resp handlers.JsonResponse) {
    var (
        msg       = req.Param("message")
        username  = req.Param("user")
        roomName  = req.Param("room")
        body     Body
    )
    if owner, err := req.ShouldGetString("username"); err != nil {
        resp.SetError(err)
    } else if err := req.BindJSON(&body); err != nil {
        resp.SetError(err)
    } else if err := message.model.Edit(username, roomName, owner, msg, body.Content);
        err != nil {
        resp.SetError(err)
    } else {
        resp.SetMessage("message updated with success")
        resp.SetStatus(http.StatusOK)
    }
}

func (message *handler) Delete(req handlers.Context, resp handlers.JsonResponse) {
    var (
        msg = req.Param("message")
        username = req.Param("user")
        roomName = req.Param("room")
    )
    if owner, err := req.ShouldGetString("username"); err != nil {
        resp.SetError(err)
    } else if err := message.model.Delete(username, roomName, owner, msg); err != nil {
        if !err.Empty() {
            resp.SetStatus(http.StatusNoContent)
            resp.SetMessage("no content found for delete!")
        } else {
            resp.SetError(err)
        }
    } else {
        const msg = "message deleted with success"
        resp.SetMessage(msg)
        resp.SetStatus(http.StatusOK)
    }
}