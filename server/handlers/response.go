package handlers

import (
    "net/http"
    "gopkg.in/gin-gonic/gin.v1"

    "github.com/dnp1/conversa/server/errors"
)
type JsonResponse interface {
    WriteJSON(c *gin.Context)
    SetError(err errors.Error)
    SetMessage(msg string)
    SetStatus(status int)
    SetData(data interface{})
}

func NewResponse() *response {
    return &response{}
}

type response struct {
    status  int
    success bool
    message string
    data    interface{}
    err     errors.Error
}


func (resp *response) WriteJSON(c *gin.Context) {
    if resp.err != nil {
        if resp.message == "" {
            resp.message = resp.err.Error()
        }
        c.Error(resp.err)
    } else {
        resp.success = true
    }
    if resp.status == 0 {
        resp.status = http.StatusNotImplemented
    }
    type js struct {
        Success bool `json:"success"`
        Message string `json:"message"`
        Data interface{} `json:"data"`
    }
    c.JSON(resp.status, js{
        Success: resp.success,
        Message: resp.message,
        Data:    resp.data,
    })
}

func (resp *response) SetError(err errors.Error) {
    resp.err = err
    switch {
    case err.Validation():
        resp.status = http.StatusBadRequest
    case err.Conflict():
        resp.status = http.StatusConflict
    case err.Authentication():
        resp.status = http.StatusUnauthorized
    case err.Authorization():
        resp.status = http.StatusForbidden
    case err.Empty():
        resp.status = http.StatusNotFound
    case err.Internal():
        resp.status = http.StatusInternalServerError
    default:
        resp.status = http.StatusNotImplemented
    }
}

func (resp *response) SetMessage(msg string) {
    resp.message = msg
}

func (resp *response) SetStatus(status int) {
    resp.status = status
}

func (resp *response) SetData(data interface{}) {
    resp.data = data
}

