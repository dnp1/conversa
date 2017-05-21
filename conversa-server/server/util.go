package server

import (
    "net/http"
    "gopkg.in/gin-gonic/gin.v1"
)

type ResponseBody struct {
    Status  int `json:"-"`
    Message string `json:"message"`
    Data    interface{} `json:"data,ommitempty"`
    Err     error `json:"-"`
}

func (resp *ResponseBody) WriteJSON(c *gin.Context) {
    if resp.Err != nil {
        c.Error(resp.Err)
    }
    c.JSON(resp.Status, resp)
}

func (resp *ResponseBody) FillWithUnexpected(err error) {
    *resp = ResponseBody{
        Status: http.StatusInternalServerError,
        Message: "a server error has ocurred!",
        Err: err,
    }
}

func (resp *ResponseBody) Fill(status int, msg string) {
    *resp = ResponseBody{
        Status: status,
        Message: msg,
    }
}

func (resp *ResponseBody) FillWithData(status int, msg string, data interface{}) {
    *resp = ResponseBody{
        Status: status,
        Message: msg,
        Data: data,
    }
}

func GetString(c *gin.Context, name string) (string, bool) {
    if iValue, ok := c.Get(name); !ok {
        return "", false
    } else if value, ok := iValue.(string); !ok {
        return "", false
    } else {
        return value, true
    }
}