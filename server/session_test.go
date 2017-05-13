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
    "github.com/golang/mock/gomock"
    "github.com/dnp1/conversa/server/mock_session"
    "github.com/dnp1/conversa/server/session"
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

    mockCtrl := gomock.NewController(t)
    defer mockCtrl.Finish()

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
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                s.EXPECT().Create("user", "password").Return("", session.ErrBadCredentials)
                rb := server.RouterBuilder{
                    Session:s,
                }
                return rb.Build()
            }(),
            strings.NewReader(`{"username": "user", "password": "password"}`),
            http.StatusUnauthorized,
        },
        {
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                s.EXPECT().Create("user", "passphrase").Return("my token!!!", nil)
                rb := server.RouterBuilder{
                    Session:s,
                }
                return rb.Build()
            }(),
            strings.NewReader(`{"username":"user", "password":"passphrase"}`),
            http.StatusOK,
        },
    }
    for i, c := range cases {
        req, err := http.NewRequest("POST", "/session", c.body)
        if assert.NoError(t, err) {
            resp :=  httptest.NewRecorder()
            c.router.ServeHTTP(resp, req)
            if !assert.Exactly(t, c.status, resp.Code) {
                t.Logf("Case %d", i)
            }
        }
    }
}

