package main

import (
    "github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
    //empty everything passing
    c.Next()
}
