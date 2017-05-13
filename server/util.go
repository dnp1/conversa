package server

import (
    "net/http"
    "github.com/pkg/errors"
    "github.com/gin-gonic/gin"
)


func notImplemented(c *gin.Context) {
    c.AbortWithError(
        http.StatusNotImplemented,
        errors.New("Not yet implemented!"),
    )
}