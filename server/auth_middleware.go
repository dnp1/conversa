package server


import "gopkg.in/gin-gonic/gin.v1"

func AuthMiddleware(c *gin.Context) {
    //empty everything passing
    c.Next()
}
