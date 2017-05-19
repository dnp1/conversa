package server

import (
    "gopkg.in/gin-gonic/gin.v1"
    "github.com/dnp1/conversa/conversa-server/message"
    "strconv"
    "net/http"
)

type MessageController struct {
    Message message.Message
}

func (mc *MessageController) ListMessages(c *gin.Context) {
    var (
        username = c.Param("user")
        roomName = c.Param("room")
        offset = int64(10)
        limit = int64(10)
    )
    if queryLimit, ok := c.GetQuery("limit"); ok {
        if i, err := strconv.ParseInt(queryLimit, 10, 64); err != nil {
            c.AbortWithError(http.StatusBadRequest, err)
        } else {
            limit = i
        }
    }
    if queryOffset, ok := c.GetQuery("limit"); ok {
        if i, err := strconv.ParseInt(queryOffset, 10, 64); err != nil {
            c.AbortWithError(http.StatusBadRequest, err)
        } else {
            offset = i
        }
    }
    if data, err := mc.Message.All(username, roomName, limit, offset); err != nil {
        c.AbortWithError(http.StatusInternalServerError, err)
    } else {
        c.JSON(http.StatusOK, data)
    }
}

type MessageBody struct {
    Content string `json:"content"`
}

func (mc *MessageController) CreateMessage(c *gin.Context) {
    var (
        username = c.Param("user")
        roomName = c.Param("room")
        body MessageBody
    )
    if senderName, ok := GetString(c, "username"); !ok {
        c.AbortWithStatus(http.StatusInternalServerError)
    } else if err := c.BindJSON(&body); err != nil {
        c.AbortWithError(http.StatusBadRequest, err)
    } else if err := mc.Message.Create(username, roomName, senderName, body.Content); err != nil {
        c.AbortWithError(http.StatusInternalServerError, err) // TODO:improve it
    } else {
        c.Status(http.StatusOK)
    }
}

func (mc *MessageController) EditMessage(c *gin.Context) {
    var (
        msg = c.Param("message")
        username = c.Param("user")
        roomName = c.Param("room")
        body MessageBody
    )
    if messageOwner, ok := GetString(c, "username"); !ok {
        c.AbortWithStatus(http.StatusInternalServerError)
    } else if err := c.BindJSON(&body); err != nil {
        c.AbortWithError(http.StatusBadRequest, err)
    } else if err := mc.Message.Edit(username, roomName, messageOwner, msg, body.Content); err != nil {
        c.AbortWithError(http.StatusInternalServerError, err) // TODO:improve it
    } else {
        c.Status(http.StatusOK)
    }
}
func (mc *MessageController) DeleteMessage(c *gin.Context) {
    var (
        msg = c.Param("message")
        username = c.Param("user")
        roomName = c.Param("room")
    )
    if messageOwner, ok := GetString(c, "username"); !ok {
        c.AbortWithStatus(http.StatusInternalServerError)
    } else if err := mc.Message.Delete(username, roomName, messageOwner, msg); err != nil {
        c.AbortWithError(http.StatusInternalServerError, err) // TODO:improve it
    } else {
        c.Status(http.StatusOK)
    }
}