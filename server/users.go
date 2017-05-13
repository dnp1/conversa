package server

import (
    "gopkg.in/gin-gonic/gin.v1"
    "github.com/dnp1/conversa/server/user"
)

type usersController struct {
    User user.User
}

func (uc *usersController) CreateUser(c *gin.Context) {

}