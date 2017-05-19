package server_test

import (
    "testing"
    "gopkg.in/gin-gonic/gin.v1"
    "github.com/stretchr/testify/assert"
    "net/http/httptest"
    "net/http"
    "github.com/dnp1/conversa/conversa-server/server"
    "github.com/golang/mock/gomock"
    "github.com/dnp1/conversa/conversa-server/session"
    "github.com/dnp1/conversa/conversa-server/mock_session"
    "github.com/twinj/uuid"
    "errors"
    "fmt"
    "strings"
)

func routerForAuthenticationTest(s session.Session) *gin.Engine {
    auth := server.Authentication{Session: s}
    r := gin.New()
    r.Use(auth.Middleware)
    r.GET("/users/:user", func(c *gin.Context) {
        c.Status(http.StatusOK)
    })
    return r
}

func TestAuthentication_Middleware(t *testing.T) {
    type Case struct {
        router   *gin.Engine
        username string
        status   int
    }
    mockCtrl := gomock.NewController(t)
    defer mockCtrl.Finish()

    tokens := [...]uuid.UUID{
        nil,
        uuid.NewV4(),
        uuid.NewV4(),
        uuid.NewV4(),
        uuid.NewV4(),
    }

    cases := [...]Case{
        {
            routerForAuthenticationTest(nil),
            "user",
            http.StatusUnauthorized,
        },
        {
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                s.EXPECT().Retrieve(tokens[1].String()).Return(nil, session.ErrTokenNotFound)
                r := routerForAuthenticationTest(s)
                return r
            }(),
            "user",
            http.StatusUnauthorized,
        },
        {
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                s.EXPECT().Retrieve(tokens[2].String()).Return(nil, errors.New("Unexpected Error."))
                r := routerForAuthenticationTest(s)
                return r
            }(),
            "user",
            http.StatusInternalServerError,
        },
        {
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                s.EXPECT().Retrieve(tokens[3].String()).Return(&session.Data{
                    Username:"user0",
                }, nil)
                r := routerForAuthenticationTest(s)
                return r
            }(),
            "user",
            http.StatusUnauthorized,
        },
        {
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                s.EXPECT().Retrieve(tokens[4].String()).Return(&session.Data{
                    Username:"user",
                }, nil)
                r := routerForAuthenticationTest(s)
                return r
            }(),
            "user",
            http.StatusOK,
        },
    }

    for i, c := range cases {
        fmt.Println("case", i)
        url := fmt.Sprintf("/users/%s", c.username)
        req, err := http.NewRequest("GET", url, strings.NewReader(""))
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