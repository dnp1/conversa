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
    sc := sessionController {
        Session: rb.Session,
    }
    uc := usersController{
        User: rb.User,
    }
    r := gin.New()
    r.POST("/session", sc.Login)
    r.DELETE("/session", sc.Logout)

    r.POST("/users", uc.CreateUser)

    auth := r.Group("/")
    authMiddleware := AuthenticationMiddleware{Session:rb.Session}
    auth.Use(authMiddleware.AuthMiddleware)
    auth.GET("/users/:user/rooms", ListRooms)
    auth.GET("/users/:user/rooms/:room", RetrieveRoom)
    auth.POST("/users/:user/rooms", CreateRoom)
    auth.POST("/users/:user/rooms/:room", JoinRoom)
    auth.DELETE("/users/:user/rooms/:room", LeaveRoom)
    auth.PATCH("/users/:user/rooms/:room", EditRoom)
    auth.GET("/users/:user/rooms/:room/messages", ListMessages)
    auth.POST("/users/:user/rooms/:room/messages", CreateMessage)
    auth.PATCH("/users/:user/rooms/:room/messages/:message", EditMessage)
    auth.DELETE("/users/:user/rooms/:room/messages/:message", DeleteMessage)

    return r
}


func NewRouter() *gin.Engine {
    builder := RouterBuilder{
        Session: session.New(),
    }
    return builder.Build()
}

