package server

import (
    "gopkg.in/gin-gonic/gin.v1"
    "github.com/dnp1/conversa/conversa-server/session"
    "net/http"
)


type Authentication struct {
    Session session.Session
}

func (auth *Authentication) Middleware(c *gin.Context) {
    var resp ResponseBody
    if token, err := c.Cookie(TokenCookieName); err != nil {
        resp.Fill(http.StatusUnauthorized, err.Error())
        resp.WriteJSON(c)
    } else if data, err := auth.Session.Retrieve(token); err == session.ErrTokenNotFound {
        resp.Fill(http.StatusUnauthorized, err.Error())
        resp.WriteJSON(c)
    } else if err != nil {
        resp.FillWithUnexpected(err)
        resp.WriteJSON(c)
    } else{
        c.Set("username", data.Username)
        c.Next()
    }
}
