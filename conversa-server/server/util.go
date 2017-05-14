package server

import (
    "net/http"
    "github.com/pkg/errors"

    "gopkg.in/gin-gonic/gin.v1"
)


func notImplemented(c *gin.Context) {
    c.AbortWithError(
        http.StatusNotImplemented,
        errors.New("Not yet implemented!"),
    )
}