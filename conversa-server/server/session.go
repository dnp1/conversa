package server

import (
    "net/http"
    "gopkg.in/gin-gonic/gin.v1"
    "github.com/dnp1/conversa/conversa-server/session"
    "time"
)

const TokenCookieName = "AUTH_TOKEN"

type SessionController struct {
    Session session.Session
}

//LoginBody is used to parse Login's handler body
type LoginBody struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func deleteCookie(c *gin.Context, name string) {
    cookie := http.Cookie{Name:name, Expires:time.Now().Add(-1 * 24 * time.Hour), Value: "deleted"}
    http.SetCookie(c.Writer, &cookie)
}

func (sc *SessionController) Login(c *gin.Context) {
    var body LoginBody
    if err := c.BindJSON(&body); err != nil {
        c.AbortWithError(http.StatusBadRequest, err)
    } else if key, err := sc.Session.Create(body.Username, body.Password); err == session.ErrBadCredentials {
        c.AbortWithError(http.StatusUnauthorized, err)
    } else if err!=nil{
        c.AbortWithError(http.StatusInternalServerError, err)
    } else {
        cookie := http.Cookie{Name: TokenCookieName, Value:key, Expires: time.Now().Add(1 * 24 * time.Hour)}
        http.SetCookie(c.Writer, &cookie)
        c.Status(http.StatusOK)
    }
}

func (sc *SessionController) Logout(c *gin.Context) {
    if token, err := c.Cookie(TokenCookieName); err != nil {
        c.AbortWithError(http.StatusNoContent, err)
        return
    } else {
        if err := sc.Session.Delete(token); err != nil {
            c.Status(http.StatusResetContent)
        } else {
            c.Status(http.StatusOK)
        }
        deleteCookie(c, TokenCookieName)
    }
}

