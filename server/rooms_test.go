package server_test

import (
    "gopkg.in/gin-gonic/gin.v1"
    "testing"
    "io"
    "github.com/golang/mock/gomock"
    "strings"
    "github.com/stretchr/testify/assert"
    "net/http/httptest"
    "net/http"
    "github.com/dnp1/conversa/server"
    "github.com/twinj/uuid"
    "github.com/dnp1/conversa/server/mock_session"
    "github.com/dnp1/conversa/server/session"
)

func init() {
    gin.SetMode(gin.TestMode)
}

func TestCreateRoom(t *testing.T) {
    type Case struct {
        router *gin.Engine
        user   string
        body   io.Reader
        status int
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

    const bodyExample = `{"name":"golang"}`
    cases := []Case{
        {//no token
            server.NewRouter(),
            "dnp1",
            strings.NewReader(bodyExample),
            http.StatusBadRequest,
        },
        {//invalid token
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                s.EXPECT().Valid(tokens[1].String()).Return(session.ErrTokenNotFound)
                r := routerForAuthenticationTest(s)
                return r
            }(),
            "dnp1",
            strings.NewReader(bodyExample),
            http.StatusUnauthorized,
        },
        {//trying create room to other user
            server.NewRouter(),
            "dnp1",
            strings.NewReader(bodyExample),
            http.StatusUnauthorized,
        },
        {//wrong json
            server.NewRouter(),
            "dnp1",
            strings.NewReader(`{"name":"golang",`),
            http.StatusBadRequest,
        },
        {//couldn't insert
            server.NewRouter(),
            "dnp1",
            strings.NewReader(bodyExample),
            http.StatusBadRequest,
        },
        {//Everything ok
            server.NewRouter(),
            "dnp1",
            strings.NewReader(bodyExample),
            http.StatusOK,
        },
    }

    for i, c := range cases {
        req, err := http.NewRequest("DELETE", "/session", strings.NewReader(""))
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