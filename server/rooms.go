package server

import (
    "gopkg.in/gin-gonic/gin.v1"
    "github.com/dnp1/conversa/server/room"
    "net/http"
)

type RoomController struct {
    Room room.Room
}

func ListRooms(c *gin.Context) {
    notImplemented(c)
}

func ListUserRooms(c *gin.Context) {
    notImplemented(c)
}

func RetrieveRoom(c *gin.Context) {
    notImplemented(c)
}


type CreateRoom struct {
    Name string `json:"name"`
}
func (rc *RoomController) CreateRoom(c *gin.Context) {
    var body CreateRoom
    if err := c.BindJSON(&body); err != nil {
        c.AbortWithError(http.StatusBadRequest, err)
    } else if err := rc.Room.Create(c.Param("user"), body.Name); err != nil {
        c.AbortWithError(http.StatusConflict, err)
    } else {
        c.Status(http.StatusOK)
    }
}
func JoinRoom(c *gin.Context) {
    notImplemented(c)
}
func DeleteRoom(c *gin.Context) {
    notImplemented(c)
}
func EditRoom(c *gin.Context) {
    notImplemented(c)
}