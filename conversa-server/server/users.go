package server

import (
    "gopkg.in/gin-gonic/gin.v1"
    "github.com/dnp1/conversa/conversa-server/user"
    "net/http"
)

type UsersController struct {
    User user.User
}


type CreateUser struct {
    Username string `json:"username"`
    Password string `json:"password"`
    PasswordConfirmation string `json:"passwordConfirmation"`
}

func (uc *UsersController) CreateUser(c *gin.Context) {
    var body CreateUser
    if err := c.BindJSON(&body); err != nil {
        c.AbortWithError(http.StatusBadRequest, err)
    } else if err:= uc.User.Create(body.Username, body.Password, body.PasswordConfirmation); err != nil {
        switch err {
        case user.ErrPasswordConfirmationDoesNotMatch:
            c.AbortWithError(http.StatusBadRequest, err)
        case user.ErrUsernameAlreadyTaken:
            c.AbortWithError(http.StatusConflict, err)
        default:
            c.AbortWithError(http.StatusInternalServerError, err)
        }
    } else {
        c.Status(http.StatusOK)
    }
}