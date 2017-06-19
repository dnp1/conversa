package controller

import (
    "gopkg.in/gin-gonic/gin.v1"
    "github.com/dnp1/conversa/server/handlers"
)

type ( // dependencies
    Authentication interface {
        Middleware(req handlers.Context, resp handlers.JsonResponse)
    }
    User interface {
        Create (req handlers.Context, resp handlers.JsonResponse)
    }
    Session interface {
        Login (req handlers.Context, resp handlers.JsonResponse)
        Logout(req handlers.Context, resp handlers.JsonResponse)
    }
    Room  interface {
        List (req handlers.Context, resp handlers.JsonResponse)
        Create (req handlers.Context, resp handlers.JsonResponse)
        Delete (req handlers.Context, resp handlers.JsonResponse)
        Edit (req handlers.Context, resp handlers.JsonResponse)
    }
    Message interface {
        List (req handlers.Context, resp handlers.JsonResponse)
        Create (req handlers.Context, resp handlers.JsonResponse)
        Edit (req handlers.Context, resp handlers.JsonResponse)
        Delete (req handlers.Context, resp handlers.JsonResponse)
    }
    Channel interface {
        Listen(context handlers.ChannelContext)
    }
)

type Handlers struct {
    Authentication Authentication
    User           User
    Session        Session
    Room           Room
    Message        Message
    Channel        Channel
}

func New(handlers *Handlers) *gin.Engine {
    r := gin.Default()
    r.POST("/session", WrapContext(handlers.Session.Login))
    r.POST("/user",  WrapContext(handlers.User.Create))
    r.DELETE("/session",  WrapContext(handlers.Session.Logout))

    authenticated := r.Group("")
    authenticated.Use(WrapContext(handlers.Authentication.Middleware)) //TODO:CheckWrap for middlewares
    authenticated.GET("/room", WrapContext(handlers.Room.List))
    authenticated.GET("/user/:user/room/:room/messages", WrapContext(handlers.Message.List))
    authenticated.GET("/user/:user/room/:room/listen", WrapChannelContext(handlers.Channel.Listen))
    authenticated.POST("/user/:user/room", WrapContext(handlers.Room.Create))
    authenticated.DELETE("/user/:user/room/:room", WrapContext(handlers.Room.Delete))
    authenticated.POST("/user/:user/room/:room/message", WrapContext(handlers.Message.Create))
    authenticated.PATCH("/user/:user/room/:room/message/:message", WrapContext(handlers.Message.Edit))
    authenticated.DELETE("/user/:user/room/:room/message/:message", WrapContext(handlers.Message.Delete))
    return r
}


