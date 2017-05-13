package server

import (
    "gopkg.in/gin-gonic/gin.v1"
    "github.com/dnp1/conversa/server/session"
    "net/http"
)


type Authentication struct {
    Session session.Session
}

func (auth *Authentication) Middleware(c *gin.Context) {
    if token, err := c.Cookie(TokenCookieName); err != nil {
        c.AbortWithError(http.StatusBadRequest, err)
        return
    } else if err := auth.Session.Valid(token); err == session.ErrTokenNotFound {
        c.AbortWithError(http.StatusUnauthorized, err)
    } else if err != nil {
        c.AbortWithError(http.StatusInternalServerError, err)
    } else {
        c.Next()
    }
}
