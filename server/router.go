package server

import (
    "gopkg.in/gin-gonic/gin.v1"
    "github.com/dnp1/conversa/server/session"

    "github.com/dnp1/conversa/server/user"
    "github.com/dnp1/conversa/server/room"
)

type RouterBuilder struct {
    Session session.Session
    User user.User
    Room room.Room
}

func (rb * RouterBuilder) Build() *gin.Engine {
    sessionCtrl := SessionController{
        Session: rb.Session,
    }
    usersController := UsersController{
        User: rb.User,
    }


    r := gin.New()
    r.POST("/session", sessionCtrl.Login)
    r.DELETE("/session", sessionCtrl.Logout)

    r.POST("/users", usersController.CreateUser)

    authenticated := r.Group("")
    authentication := Authentication{Session:rb.Session}
    authenticated.Use(authentication.Middleware)
    //auth.GET("/users", usersController.)
    authenticated.GET("/rooms", ListRooms)

    authenticated.GET("/users/:user/rooms/:room/messages", ListMessages)
    authenticated.POST("/users/:user/rooms/:room/messages", CreateMessage)
    authenticated.PATCH("/users/:user/rooms/:room/messages/:message", EditMessage)
    authenticated.DELETE("/users/:user/rooms/:room/messages/:message", DeleteMessage)
    authenticated.GET("/users/:user/rooms", ListUserRooms)

    authorized := r.Group("")
    authorization := Authorization{Session:rb.Session}
    authorized.Use(authorization.Middleware)
    authorized.POST("/users/:user/rooms", CreateRoom)
    authorized.DELETE("/users/:user/rooms/:room", DeleteRoom)
    authorized.PATCH("/users/:user/rooms/:room", EditRoom)

    return r
}


func NewRouter() *gin.Engine {
    builder := RouterBuilder{
        Session: session.New(),
        User: user.New(),
        Room: room.New(),
    }
    return builder.Build()
}

