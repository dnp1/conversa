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
    )
    var offset = int64(10)
    var limit = int64(10)
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

func (mc *MessageController) CreateMessage(c *gin.Context) {
    notImplemented(c)
}
func (mc *MessageController) EditMessage(c *gin.Context) {
    notImplemented(c)
}
func (mc *MessageController) DeleteMessage(c *gin.Context) {
    notImplemented(c)
}