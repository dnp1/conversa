package controller_test

import (
    "gopkg.in/gin-gonic/gin.v1"
    "github.com/dnp1/conversa/server/controller"
)

func NoDependencyRouter() *gin.Engine {
    c := controller.RouterBuilder{}
    return c.Build()
}
