package server

import (
    "gopkg.in/gin-gonic/gin.v1"
    "github.com/dnp1/conversa/server/session"

    "github.com/dnp1/conversa/server/user"
)

type RouterBuilder struct {
    Session session.Session
    User user.User
}

func (rb * RouterBuilder) Build() *gin.Engine {
    sessionCtrl := SessionController{
        Session: rb.Session,
    }
    usersController := UsersController{
        User: rb.User,
    }
    authentication := Authentication{Session:rb.Session}
    r := gin.New()
    r.POST("/session", sessionCtrl.Login)
    r.DELETE("/session", sessionCtrl.Logout)

    r.POST("/users", usersController.CreateUser)

    authenticated := r.Group("/")
    authenticated.Use(authentication.Middleware)
    //auth.GET("/users", usersController.)
    authenticated.GET("/rooms", ListRooms)

    authenticated.GET("/users/:user/rooms/:room/messages", ListMessages)
    authenticated.POST("/users/:user/rooms/:room/messages", CreateMessage)
    authenticated.PATCH("/users/:user/rooms/:room/messages/:message", EditMessage)
    authenticated.DELETE("/users/:user/rooms/:room/messages/:message", DeleteMessage)

    authorized := authenticated.Group("/")
    authorized.GET("/users/:user/rooms", ListUserRooms)
    authorized.POST("/users/:user/rooms", CreateRoom)
    authorized.DELETE("/users/:user/rooms/:room", DeleteRoom)
    authorized.PATCH("/users/:user/rooms/:room", EditRoom)

    return r
}


func NewRouter() *gin.Engine {
    builder := RouterBuilder{
        Session: session.New(),
    }
    return builder.Build()
}

