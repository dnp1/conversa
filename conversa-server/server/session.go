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
    var resp ResponseBody
    defer resp.WriteJSON(c)

    if err := c.BindJSON(&body); err != nil {
        const msg = "body sent is not a valid json"
        resp.Fill(http.StatusBadRequest, msg)
    } else if key, err := sc.Session.Create(body.Username, body.Password); err == session.ErrBadCredentials {
        resp.Fill(http.StatusUnauthorized, err.Error())
    } else if err != nil {
        resp.FillWithUnexpected(err)
    } else {
        cookie := http.Cookie{Name: TokenCookieName, Value:key, Expires: time.Now().Add(1 * 24 * time.Hour)}
        http.SetCookie(c.Writer, &cookie)
        const msg = "sucessful sign-in"
        resp.Fill(http.StatusCreated, msg)
    }
}

func (sc *SessionController) Logout(c *gin.Context) {
    var resp ResponseBody
    defer resp.WriteJSON(c)

    if token, err := c.Cookie(TokenCookieName); err == http.ErrNoCookie  {
        resp.Fill(http.StatusNoContent, err.Error())
    } else if err != nil {
        resp.FillWithUnexpected(err)
    } else {
        if err := sc.Session.Delete(token); err != nil {
            const msg = "looks like your cookie is outdated!"
            resp.Fill(http.StatusResetContent, msg)
        } else {
            const msg = "sucessful sign-out"
            resp.Fill(http.StatusOK, msg)
        }
        deleteCookie(c, TokenCookieName)
    }
}

