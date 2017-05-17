package server

import (
    "gopkg.in/gin-gonic/gin.v1"
    "github.com/dnp1/conversa/conversa-server/session"
    "github.com/dnp1/conversa/conversa-server/user"
    "github.com/dnp1/conversa/conversa-server/room"
    "database/sql"
    "github.com/dnp1/conversa/conversa-server/message"
)

type RouterBuilder struct {
    db      *sql.DB
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



    authorized := r.Group("")
    authorization := Authorization{Session:rb.Session}
    authorized.Use(authorization.Middleware)
    authorized.POST("/users/:user/rooms", roomCtrl.CreateRoom)
    authorized.DELETE("/users/:user/rooms/:room", roomCtrl.DeleteRoom)
    authorized.PATCH("/users/:user/rooms/:room", roomCtrl.EditRoom)
    authenticated.POST("/users/:user/rooms/:room/messages", messageCtrl.CreateMessage)
    authenticated.PATCH("/users/:user/rooms/:room/messages/:message", messageCtrl.EditMessage)
    authenticated.DELETE("/users/:user/rooms/:room/messages/:message", messageCtrl.DeleteMessage)


    return r
}

func NewRouter(db *sql.DB) *gin.Engine {
    builder := RouterBuilder{
        Session: session.Builder{DB:db}.Build(),
        User: user.Builder{DB:db}.Build(),
        Room: room.Builder{DB:db}.Build(),
        Message: message.Builder{DB:db}.Build(),
    }
    return builder.Build()
}

