package server

import (
    "gopkg.in/gin-gonic/gin.v1"
    "github.com/dnp1/conversa/conversa-server/server/session"
    "net/http"
)

type Authorization struct {
    Session session.Session
}


func (auth *Authorization) Middleware(c *gin.Context) {
    if token, err := c.Cookie(TokenCookieName); err != nil {
        c.AbortWithError(http.StatusBadRequest, err)
        return
    } else if data, err := auth.Session.Retrieve(token); err == session.ErrTokenNotFound {
        c.AbortWithError(http.StatusUnauthorized, err)
    } else if err != nil {
        c.AbortWithError(http.StatusInternalServerError, err)
    } else if data.Username != c.Param("user") {
        c.AbortWithError(http.StatusUnauthorized, err)
    } else {
        c.Next()
    }
}