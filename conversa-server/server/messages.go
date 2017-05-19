package server

import (
    "gopkg.in/gin-gonic/gin.v1"
    "github.com/dnp1/conversa/conversa-server/message"
    "strconv"
    "net/http"
    "github.com/pkg/errors"
)

var ErrContextSetAssertion = errors.New("Server assertion wasn't true")

type MessageController struct {
    Message message.Message
}

func (mc *MessageController) ListMessages(c *gin.Context) {
    var resp ResponseBody
    defer resp.WriteJSON(c)

    var (
        username = c.Param("user")
        roomName = c.Param("room")
        offset = int64(10)
        limit = int64(10)
    )
    if queryLimit, ok := c.GetQuery("limit"); ok {
        if i, err := strconv.ParseInt(queryLimit, 10, 64); err != nil {
            const msg = "msgs limit should be a integer"
            resp.Fill(http.StatusBadRequest, msg)
        } else {
            limit = i
        }
    }
    if queryOffset, ok := c.GetQuery("offset"); ok {
        if i, err := strconv.ParseInt(queryOffset, 10, 64); err != nil {
            const msg = "msgs offset should be a integer"
            resp.Fill(http.StatusBadRequest, msg)
        } else {
            offset = i
        }
    }

    if data, err := mc.Message.All(username, roomName, limit, offset); err != nil {
        resp.FillWithUnexpected(err)
    } else {
        resp.FillWithData(http.StatusOK, "message's list", data)
    }
}

type MessageBody struct {
    Content string `json:"content"`
}

func (mc *MessageController) CreateMessage(c *gin.Context) {
    var resp ResponseBody
    defer resp.WriteJSON(c)

    var (
        username = c.Param("user")
        roomName = c.Param("room")
        body MessageBody
    )
    if senderName, ok := GetString(c, "username"); !ok {
        resp.FillWithUnexpected(ErrContextSetAssertion)
    } else if err := c.BindJSON(&body); err != nil {
        const msg = "body sent is not a valid json"
        resp.Fill(http.StatusBadRequest, msg)
    } else if err := mc.Message.Create(username, roomName, senderName, body.Content);
        err == message.ErrMessageIsEmpty || err == message.ErrCouldNotFound {
        resp.Fill(http.StatusBadRequest, err.Error())
    } else if err != nil {
        resp.FillWithUnexpected(err)
    } else {
        const msg = "message created with success"
        resp.Fill(http.StatusCreated, msg)
    }
}

func (mc *MessageController) EditMessage(c *gin.Context) {
    var resp ResponseBody
    defer resp.WriteJSON(c)

    var (
        msg = c.Param("message")
        username = c.Param("user")
        roomName = c.Param("room")
        body MessageBody
    )
    if messageOwner, ok := GetString(c, "username"); !ok {
        resp.FillWithUnexpected(ErrContextSetAssertion)
    } else if err := c.BindJSON(&body); err != nil {
        const msg = "body sent is not a valid json"
        resp.Fill(http.StatusBadRequest, msg)
    } else if err := mc.Message.Edit(username, roomName, messageOwner, msg, body.Content);
        err == message.ErrMessageIsEmpty || err == message.ErrCouldNotFound {
        resp.Fill(http.StatusBadRequest, err.Error())
    } else if err != nil {
        resp.FillWithUnexpected(err)
    } else {
        const msg = "message updated with success"
        resp.Fill(http.StatusOK, msg)
    }
}
func (mc *MessageController) DeleteMessage(c *gin.Context) {
    var resp ResponseBody
    defer resp.WriteJSON(c)

    var (
        msg = c.Param("message")
        username = c.Param("user")
        roomName = c.Param("room")
    )
    if messageOwner, ok := GetString(c, "username"); !ok {
        resp.FillWithUnexpected(ErrContextSetAssertion)
    } else if err := mc.Message.Delete(username, roomName, messageOwner, msg); err == message.ErrCouldNotFound {
        resp.Fill(http.StatusNoContent, err.Error())
    } else if err != nil {
        resp.FillWithUnexpected(err)
    } else {
        const msg = "message deleted with success"
        resp.Fill(http.StatusOK, msg)
    }
}