package server_test

import (
    "testing"

    "net/http/httptest"
    "net/http"
    "github.com/stretchr/testify/assert"
    "github.com/dnp1/conversa/server"
    "strings"
    "io"
    "gopkg.in/gin-gonic/gin.v1"
)

func init() {
    gin.SetMode(gin.TestMode)
}

func TestLogin(t *testing.T) {
    type Case struct {
        router *gin.Engine
        body io.Reader
        status int
    }

    cases := [...]Case {
        {
            server.NewRouter(),
            strings.NewReader(""),
            http.StatusBadRequest,
        },
        {
            server.NewRouter(),
            strings.NewReader(`{"user_name": "json","password"}`),
            http.StatusBadRequest,
        },
        {
            server.NewRouter(),
            strings.NewReader(`{"username": "user", "password": "passphrase"}`),
            http.StatusUnauthorized,
        },
        {
            server.NewRouter(),
            strings.NewReader(`{"username":"user", "password":"passphrase"}`),
            http.StatusOK,
        },
    }
    for i, c := range cases {
        req, err := http.NewRequest("POST", "/session", c.body)
        if assert.NoError(t, err) {
            resp :=  httptest.NewRecorder()
            c.router.ServeHTTP(resp, req)
            if !assert.Exactly(t, c.status, resp.Code){
                t.Logf("Case %d", i)
            }
        }
    }
}

