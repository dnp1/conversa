package server

import (
    "github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
    r := gin.New()
    r.POST("/session", Login)
    r.DELETE("/session", Logout)
    r.POST("/sign-up", CreateUser)

    auth := r.Group("/")
    auth.Use(AuthMiddleware)
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

