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
)

type Handlers struct {
    Authentication Authentication
    User           User
    Session        Session
    Room           Room
    Message        Message
}

func New(handlers *Handlers) *gin.Engine {
    r := gin.Default()
    r.POST("/session", Wrap(handlers.Session.Login))
    r.POST("/user",  Wrap(handlers.User.Create))
    r.DELETE("/session",  Wrap(handlers.Session.Logout))

    authenticated := r.Group("")
    authenticated.Use(Wrap(handlers.Authentication.Middleware))
    authenticated.GET("/room", Wrap(handlers.Room.List))
    authenticated.GET("/user/:user/room/:room/messages", Wrap(handlers.Message.List))
    authenticated.POST("/user/:user/room", Wrap(handlers.Room.Create))
    authenticated.DELETE("/user/:user/room/:room", Wrap(handlers.Room.Delete))
    authenticated.PATCH("/user/:user/room/:room", Wrap(handlers.Room.Edit))
    authenticated.POST("/user/:user/room/:room/message", Wrap(handlers.Message.Create))
    authenticated.PATCH("/user/:user/room/:room/message/:message", Wrap(handlers.Message.Edit))
    authenticated.DELETE("/user/:user/room/:room/message/:message", Wrap(handlers.Message.Delete))

    return r
}


