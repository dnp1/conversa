package main

import (
    "os"
    "net/http"
    "github.com/pkg/errors"
    "github.com/gin-gonic/gin"
)

func env(key string, defaultVal string) string {
    if val, exists := os.LookupEnv(key); exists {
        return val;
    } else {
        return defaultVal
    }
}

func notImplemented(c *gin.Context) {
    c.AbortWithError(
        http.StatusNotImplemented,
        errors.New("Not yet implemented!"),
    )
}