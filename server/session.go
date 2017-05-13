package server

import (

    "net/http"
    "gopkg.in/gin-gonic/gin.v1"
)

//LoginBody is used to parse Login's handler body
type LoginBody struct {
    Username [255]byte `json:"username"`
    Password [255]byte `json:"password"`
}

func Login(c *gin.Context) {
    body := new(LoginBody)
    if err := c.BindJSON(&body); err != nil {
        c.AbortWithStatus(http.StatusBadRequest)
    } else {
        notImplemented(c)
    }

}
func Logout(c *gin.Context) {
    notImplemented(c)
}
func CreateUser(c *gin.Context) {
    notImplemented(c)
}