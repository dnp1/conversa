package controller

import (
    "gopkg.in/gin-gonic/gin.v1"
    "github.com/dnp1/conversa/server/model/user"
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
    var resp ResponseBody
    defer resp.WriteJSON(c)

    if err := c.BindJSON(&body); err != nil {
        const msg = "body sent is not a valid json"
        resp.Fill(http.StatusBadRequest, msg)
    } else if err:= uc.User.Create(body.Username, body.Password, body.PasswordConfirmation); err != nil {
        switch err {
        case user.ErrPasswordConfirmationDoesNotMatch: fallthrough
        case user.ErrUsernameHasInvalidCharacters: fallthrough
        case user.ErrUsernameWrongLength: fallthrough
        case user.ErrPasswordTooShort:
            resp.Fill(http.StatusBadRequest, err.Error())
        case user.ErrUsernameAlreadyTaken:
            resp.Fill(http.StatusConflict, err.Error())
        default:
            resp.FillWithUnexpected(err)
        }
    } else {
        const msg = "user created with success"
        resp.Fill(http.StatusOK, msg)
    }
}