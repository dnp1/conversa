package server

import (
    "gopkg.in/gin-gonic/gin.v1"
    "github.com/dnp1/conversa/server/session"
)

type Authorization struct {
    Session session.Session
}


func (auth *Authorization) Middleware(c *gin.Context) {

}