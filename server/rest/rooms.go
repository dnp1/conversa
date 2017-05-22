package rest

import (
    "gopkg.in/gin-gonic/gin.v1"
    "github.com/dnp1/conversa/server/room"
    "net/http"
)

type RoomController struct {
    Room room.Room
}

func (rc *RoomController) ListRooms(c *gin.Context) {
    var resp ResponseBody
    defer resp.WriteJSON(c)

    if data, err := rc.Room.All(); err != nil {
        resp.FillWithUnexpected(err)
    } else {
        const msg = "list of rooms"
        resp.FillWithData(http.StatusOK, msg, data)
    }
}

func (rc *RoomController) ListUserRooms(c *gin.Context) {
    var resp ResponseBody
    defer resp.WriteJSON(c)

    if data, err := rc.Room.AllByUser(c.Param("user")); err != nil {
        resp.FillWithUnexpected(err)
    } else {
        const msg = "list of user's rooms"
        resp.FillWithData(http.StatusOK, msg, data)
    }
}

type CreateRoom struct {
    Name string `json:"name"`
}

func (rc *RoomController) CreateRoom(c *gin.Context) {
    var body CreateRoom
    var resp ResponseBody
    defer resp.WriteJSON(c)

    var user = c.Param("user")
    if username, ok := GetString(c, "username"); !ok {
        resp.FillWithUnexpected(ErrContextSetAssertion)
    } else if user != username {
        const msg = "permission denied"
        resp.Fill(http.StatusBadRequest, msg)
    } else if err := c.BindJSON(&body); err != nil {
        const msg = "body sent is not a valid json"
        resp.Fill(http.StatusBadRequest, msg)
    } else if err := rc.Room.Create(user, body.Name); err == room.ErrRoomNameAlreadyExists {
        resp.Fill(http.StatusConflict, err.Error())
    } else if err == room.ErrRoomNameHasInvalidCharacters || err == room.ErrRoomNameWrongLength {
        resp.Fill(http.StatusBadRequest, err.Error())
    } else if err != nil {
        resp.FillWithUnexpected(err)
    } else {
        const msg = "room created with success!"
        resp.Fill(http.StatusCreated, msg)
    }
}

func (rc *RoomController) DeleteRoom(c *gin.Context) {
    var resp ResponseBody
    defer resp.WriteJSON(c)

    var user = c.Param("user")
    if username, ok := GetString(c, "username"); !ok {
        resp.FillWithUnexpected(ErrContextSetAssertion)
    } else if user != username {
        const msg = "permission denied"
        resp.Fill(http.StatusBadRequest, msg)
    } else if err := rc.Room.Delete(c.Param("user"), c.Param("room")); err != nil {
        resp.Fill(http.StatusNoContent, err.Error())
    } else {
        const msg = "room deleted with success!"
        resp.Fill(http.StatusOK, msg)
    }
}

func (rc *RoomController) EditRoom(c *gin.Context) {
    var body CreateRoom
    var resp ResponseBody
    defer resp.WriteJSON(c)

    var user = c.Param("user")
    if username, ok := GetString(c, "username"); !ok {
        resp.FillWithUnexpected(ErrContextSetAssertion)
    } else if user != username {
        const msg = "permission denied"
        resp.Fill(http.StatusBadRequest, msg)
    } else if err := c.BindJSON(&body); err != nil {
        const msg = "body sent is not a valid json"
        resp.Fill(http.StatusBadRequest, msg)
    } else if err := rc.Room.Rename(c.Param("user"), c.Param("room"), body.Name); err != nil {
        resp.Fill(http.StatusConflict, err.Error()) //TODO:improve it
    } else {
        const msg = "room edited with success!"
        resp.Fill(http.StatusOK, msg)
    }
}