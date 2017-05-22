package rest_test

import (
    "testing"
    "net/http/httptest"
    "net/http"
    "github.com/stretchr/testify/assert"
    "github.com/dnp1/conversa/server/rest"
    "strings"
    "io"
    "gopkg.in/gin-gonic/gin.v1"
    "github.com/golang/mock/gomock"
    "github.com/dnp1/conversa/server/mock_session"
    "github.com/dnp1/conversa/server/session"

    "github.com/twinj/uuid"
    "github.com/pkg/errors"
)

func init() {
    gin.SetMode(gin.TestMode)
}

func TestSessionController_Login(t *testing.T) {
    type Case struct {
        router *gin.Engine
        body   io.Reader
        status int
    }

    mockCtrl := gomock.NewController(t)
    defer mockCtrl.Finish()

    cases := [...]Case{
        {// Bad request. empty body
            rest.NewRouter(nil),
            strings.NewReader(""),
            http.StatusBadRequest,
        },
        {//Bad wrong json
            rest.NewRouter(nil),
            strings.NewReader(`{"user_name": "json","password"}`),
            http.StatusBadRequest,
        },
        {
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                s.EXPECT().Create("user", "password").Return("", errors.New("Unexpected!"))
                rb := rest.RouterBuilder{
                    Session:s,
                }
                return rb.Build()
            }(),
            strings.NewReader(`{"username": "user", "password": "password"}`),
            http.StatusInternalServerError,
        },
        {
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                s.EXPECT().Create("user", "passworddsfsfds").Return("", session.ErrBadCredentials)
                rb := rest.RouterBuilder{
                    Session:s,
                }
                return rb.Build()
            }(),
            strings.NewReader(`{"username": "user", "password": "passworddsfsfds"}`),
            http.StatusUnauthorized,
        },
        {
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                s.EXPECT().Create("user", "passphrase").Return("my token!!!", nil)
                rb := rest.RouterBuilder{
                    Session:s,
                }
                return rb.Build()
            }(),
            strings.NewReader(`{"username":"user", "password":"passphrase"}`),
            http.StatusCreated,
        },
    }
    for i, c := range cases {
        req, err := http.NewRequest("POST", "/sessions", c.body)
        if assert.NoError(t, err) {
            resp := httptest.NewRecorder()
            c.router.ServeHTTP(resp, req)
            if !assert.Exactly(t, c.status, resp.Code) {
                t.Logf("Case %d", i)
            }
        }
    }
}

func TestSessionController_Logout(t *testing.T) {
    type Case struct {
        router *gin.Engine
        status int
    }
    tokens := [...]uuid.UUID {
        nil,
        uuid.NewV4(),
        uuid.NewV4(),
    }
    mockCtrl := gomock.NewController(t)
    defer mockCtrl.Finish()
    cases := [...]Case{
        {
            rest.NewRouter(nil),
            http.StatusNoContent,
        },
        {
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                s.EXPECT().Delete(tokens[1].String()).Return(session.ErrTokenNotFound)
                rb := rest.RouterBuilder{
                    Session:s,
                }
                return rb.Build()
            }(),
            http.StatusResetContent,
        },
        {
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                s.EXPECT().Delete(tokens[2].String()).Return(nil)
                rb := rest.RouterBuilder{
                    Session:s,
                }
                return rb.Build()
            }(),
            http.StatusOK,
        },
    }
    for i, c := range cases {
        req, err := http.NewRequest("DELETE", "/sessions", strings.NewReader(""))
        if tokens[i] != nil {
            req.AddCookie(&http.Cookie{
                Name: rest.TokenCookieName,
                Value: tokens[i].String(),
                MaxAge: 24 * 60 * 60,
                Secure: true,
                HttpOnly: true,
            })
        }
        if assert.NoError(t, err) {
            resp := httptest.NewRecorder()
            c.router.ServeHTTP(resp, req)
            if !assert.Exactly(t, c.status, resp.Code) {
                t.Logf("Case %d", i)
            }
        }
    }

}
