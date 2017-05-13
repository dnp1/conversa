package server

import (
    "net/http"
    "gopkg.in/gin-gonic/gin.v1"
    "github.com/dnp1/conversa/server/session"
)

const TokenCookieName = "AUTH_TOKEN"

type sessionController struct {
    Session session.Session
}
//LoginBody is used to parse Login's handler body
type LoginBody struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func deleteCookie(c *gin.Context, name string) {
    c.SetCookie(
        name,
        "deleted",
        -1,
        "",
        "",
        true,
        true,
    )
}

func (sc *sessionController) Login(c *gin.Context) {
    var body LoginBody
    if err := c.BindJSON(&body); err != nil {
        c.AbortWithError(http.StatusBadRequest, err)
    } else if key, err := sc.Session.Create(body.Username, body.Password); err != nil {
        c.AbortWithError(http.StatusUnauthorized, err)
    } else {
        c.SetCookie(
            TokenCookieName,
            key,
            24 * 60 * 60,
            "",
            "",
            true,
            true,
        )
        c.Status(http.StatusOK)
    }
}

func (sc *sessionController) Logout(c *gin.Context) {
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

