package server_test

import (
    "testing"
    "gopkg.in/gin-gonic/gin.v1"
    "io"
    "github.com/stretchr/testify/assert"
    "net/http/httptest"
    "net/http"
    "github.com/dnp1/conversa/server"
    "strings"
    "github.com/dnp1/conversa/server/user"
    "github.com/golang/mock/gomock"
    "github.com/dnp1/conversa/server/mock_user"
    "errors"
)

func TestSessionController_CreateUser(t *testing.T) {
    type Case struct {
        router *gin.Engine
        body   io.Reader
        status int
    }

    mockCtrl := gomock.NewController(t)
    defer mockCtrl.Finish()
    cases := [...]Case {
        {
            server.NewRouter(),
            strings.NewReader(`{"user":sdasdas "`),
            http.StatusBadRequest,
        },
        {
            func() *gin.Engine {
                u := mock_user.NewMockUser(mockCtrl)
                u.EXPECT().Create("user", "senha","passphrase").Return(user.ErrPasswordConfirmationDoesNotMatch)
                rb := server.RouterBuilder{
                    User:u,
                }
                return rb.Build()
            }(),
            strings.NewReader(`{"username":"user", "password": "senha", "passwordConfirmation":"passphrase"}`),
            http.StatusBadRequest,
        },
        {
            func() *gin.Engine {
                u := mock_user.NewMockUser(mockCtrl)
                u.EXPECT().Create("user", "passphrase","passphrase").Return(user.ErrUsernameAlreadyTaken)
                rb := server.RouterBuilder{
                    User:u,
                }
                return rb.Build()
            }(),
            strings.NewReader(`{"username":"user", "password": "passphrase", "passwordConfirmation":"passphrase"}`),
            http.StatusConflict,
        },
        {
            func() *gin.Engine {
                u := mock_user.NewMockUser(mockCtrl)
                u.EXPECT().Create("user", "passphrase","passphrase").Return(errors.New("Unexpected error!!!"))
                rb := server.RouterBuilder{
                    User:u,
                }
                return rb.Build()
            }(),
            strings.NewReader(`{"username":"user", "password": "passphrase", "passwordConfirmation":"passphrase"}`),
            http.StatusInternalServerError,
        },
        {
            func() *gin.Engine {
                u := mock_user.NewMockUser(mockCtrl)
                u.EXPECT().Create("user", "passphrase","passphrase").Return(nil)
                rb := server.RouterBuilder{
                    User:u,
                }
                return rb.Build()
            }(),
            strings.NewReader(`{"username":"user", "password": "passphrase", "passwordConfirmation":"passphrase"}`),
            http.StatusOK,
        },
    }

    for i, c := range cases {
        req, err := http.NewRequest("POST", "/users", c.body)
        if assert.NoError(t, err) {
            resp := httptest.NewRecorder()
            c.router.ServeHTTP(resp, req)
            if !assert.Exactly(t, c.status, resp.Code) {
                t.Logf("Case %d", i)
            }
        }
    }

}