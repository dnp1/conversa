package server

import (
    "gopkg.in/gin-gonic/gin.v1"
    "github.com/dnp1/conversa/server/room"
    "net/http"
)

type RoomController struct {
    Room room.Room
}

func (rc *RoomController) ListRooms(c *gin.Context) {
    if data, err := rc.Room.All(); err != nil{
        c.AbortWithError(http.StatusInternalServerError, err)
    } else {
        c.JSON(http.StatusOK, data)
    }
}

func ListUserRooms(c *gin.Context) {
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

func (rc *RoomController) DeleteRoom(c *gin.Context) {
    if err := rc.Room.Delete(c.Param("user"), c.Param("room")); err != nil {
        c.AbortWithError(http.StatusConflict, err)
    } else {
        c.Status(http.StatusOK)
    }
}

func (rc *RoomController) EditRoom(c *gin.Context) {
    var body CreateRoom
    if err := c.BindJSON(&body); err != nil {
        c.AbortWithError(http.StatusBadRequest, err)
    } else if err := rc.Room.Rename(c.Param("user"), c.Param("room"), body.Name); err != nil {
        c.AbortWithError(http.StatusConflict, err)
    } else {
        c.Status(http.StatusOK)
    }
}