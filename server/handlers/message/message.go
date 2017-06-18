package message

import (
    "strconv"
    "net/http"
    "github.com/dnp1/conversa/server/errors"
    "github.com/dnp1/conversa/server/data/message"
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
    All(username, roomName string, limit, offset int64) ([]message.Data, errors.Error)
    Delete(username, roomName, messageOwner, msg string) errors.Error
}


type handler struct {
    model Model
}

func (message *handler) List(req handlers.Context, resp handlers.JsonResponse) {
    var (
        username = req.Param("user")
        roomName = req.Param("room")
        offset = int64(10)
        limit = int64(10)
    )
    if queryLimit, ok := req.GetQuery("limit"); ok {
        if i, err := strconv.ParseInt(queryLimit, 10, 64); err != nil {
            const msg = "msgs limit should be a integer"
            resp.SetError(errors.Validation(errors.FromString(msg)))
        } else {
            limit = i
        }
    }
    if queryOffset, ok := req.GetQuery("offset"); ok {
        if i, err := strconv.ParseInt(queryOffset, 10, 64); err != nil {
            const msg = "msgs offset should be a integer"
            resp.SetError(errors.Validation(errors.FromString(msg)))
        } else {
            offset = i
        }
    }
    if data, err := message.model.All(username, roomName, limit, offset); err != nil {
        resp.SetError(err)
    } else {
        resp.SetMessage("message's list")
        resp.SetData(data)
        resp.SetStatus(http.StatusOK)
    }
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