package server

import (

    "net/http"
    "gopkg.in/gin-gonic/gin.v1"
    "github.com/dnp1/conversa/server/session"
    "fmt"
)

type sessionController struct {
    Session session.Session
}
//LoginBody is used to parse Login's handler body
type LoginBody struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

func (sc *sessionController) Login(c *gin.Context) {
    var body LoginBody
    if err := c.BindJSON(&body); err != nil {
        c.AbortWithError(http.StatusBadRequest, err)
        fmt.Println("oi!")
    } else if key, err := sc.Session.Create(body.Username, body.Password); err != nil {
        c.AbortWithError(http.StatusUnauthorized, err)
    } else {
        c.SetCookie(
            "AUTH_TOKEN",
            key,
            24*60*60,
            "",
            "",
            true,
            true,
        )
        c.Status(http.StatusOK)
    }
}

func (sc *sessionController) Logout(c *gin.Context) {
    notImplemented(c)
}
func (sc *sessionController) CreateUser(c *gin.Context) {
    notImplemented(c)
}