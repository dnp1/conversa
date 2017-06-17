package controller

import (
    "gopkg.in/gin-gonic/gin.v1"
    "github.com/dnp1/conversa/server/model/session"
    "github.com/dnp1/conversa/server/model/user"
    "github.com/dnp1/conversa/server/model/room"
    "github.com/dnp1/conversa/server/model/message"
)

type RouterBuilder struct {
    Session session.Session
    User    user.User
    Room    room.Room
    Message message.Message
}

func (rb *RouterBuilder) Build() *gin.Engine {
    sessionCtrl := SessionController{
        Session: rb.Session,
    }
    usersCtrl := UsersController{
        User: rb.User,
    }
    roomCtrl := RoomController{
        Room: rb.Room,
    }
    messageCtrl := MessageController{
        Message: rb.Message,
    }

    r := gin.Default()
    r.POST("/sessions", sessionCtrl.Login)
    r.POST("/users", usersCtrl.CreateUser)
    r.DELETE("/sessions", sessionCtrl.Logout)


    authenticated := r.Group("")
    authentication := Authentication{Session:rb.Session}
    authenticated.Use(authentication.Middleware)
    authenticated.GET("/rooms", roomCtrl.ListRooms)
    authenticated.GET("/users/:user/rooms", roomCtrl.ListUserRooms)
    authenticated.GET("/users/:user/rooms/:room/messages", messageCtrl.ListMessages)
    authenticated.POST("/users/:user/rooms", roomCtrl.CreateRoom)
    authenticated.DELETE("/users/:user/rooms/:room", roomCtrl.DeleteRoom)
    authenticated.PATCH("/users/:user/rooms/:room", roomCtrl.EditRoom)
    authenticated.POST("/users/:user/rooms/:room/messages", messageCtrl.CreateMessage)
    authenticated.PATCH("/users/:user/rooms/:room/messages/:message", messageCtrl.EditMessage)
    authenticated.DELETE("/users/:user/rooms/:room/messages/:message", messageCtrl.DeleteMessage)


    return r
}


