package server_test

import (
    "testing"
    "gopkg.in/gin-gonic/gin.v1"
    "github.com/stretchr/testify/assert"
    "net/http/httptest"
    "net/http"
    "github.com/dnp1/conversa/conversa-server/server"
    "github.com/golang/mock/gomock"
    "github.com/dnp1/conversa/conversa-server/server/session"
    "github.com/dnp1/conversa/conversa-server/server/mock_session"
    "github.com/twinj/uuid"
    "errors"
)

func routerForAuthenticationTest(s session.Session) *gin.Engine {
    auth := server.Authentication{Session: s}
    r := gin.New()
    r.Use(auth.Middleware)
    r.GET("/", func(c *gin.Context) {
        c.Status(http.StatusOK)
    })
    return r
}

func TestAuthentication_Middleware(t *testing.T) {
    type Case struct {
        router *gin.Engine
        status int
    }

    tokens := [...]uuid.UUID {
        nil,
        uuid.NewV4(),
        uuid.NewV4(),
        uuid.NewV4(),
    }
    mockCtrl := gomock.NewController(t)
    defer mockCtrl.Finish()
    cases := [...]Case{
        {
            routerForAuthenticationTest(session.Session(nil)),
            http.StatusBadRequest,
        },
        {
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                s.EXPECT().Valid(tokens[1].String()).Return(session.ErrTokenNotFound)
                r := routerForAuthenticationTest(s)
                return r
            }(),
            http.StatusUnauthorized,
        },
        {
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                s.EXPECT().Valid(tokens[2].String()).Return(errors.New("unexpected error"))
                r := routerForAuthenticationTest(s)
                return r
            }(),
            http.StatusInternalServerError,
        },
        {
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                s.EXPECT().Valid(tokens[3].String()).Return(nil)
                r := routerForAuthenticationTest(s)
                return r
            }(),
            http.StatusOK,
        },
    }

    for i, c := range cases {
        req, err := http.NewRequest("GET", "/", nil)
        if tokens[i] != nil {
            req.AddCookie(&http.Cookie{
                Name: server.TokenCookieName,
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