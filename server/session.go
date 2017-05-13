package server

import (
    "github.com/gin-gonic/gin"
)

//LoginBody is used to parse Login's handler body
type LoginBody struct {
    Username [255]byte `json:"username"`
    Password [255]byte `json:"password"`
}

func Login(c *gin.Context) {
    notImplemented(c)
}
func Logout(c *gin.Context) {
    notImplemented(c)
}
func CreateUser(c *gin.Context) {
    notImplemented(c)
}