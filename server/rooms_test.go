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
    "fmt"
    "github.com/dnp1/conversa/server/mock_room"
    "github.com/dnp1/conversa/server/room"
)

func init() {
    gin.SetMode(gin.TestMode)
}

func TestRoomController_CreateRoom(t *testing.T) {
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
                s.EXPECT().Retrieve(tokens[1].String()).Return(nil, session.ErrTokenNotFound)
                rb := server.RouterBuilder{
                    Session:s,
                }
                return rb.Build()
            }(),
            "dnp1",
            strings.NewReader(bodyExample),
            http.StatusUnauthorized,
        },
        {//trying create room to other user
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                s.EXPECT().Retrieve(tokens[2].String()).Return(nil, session.ErrTokenNotFound)
                rb := server.RouterBuilder{
                    Session:s,
                }
                return rb.Build()
            }(),
            "dnp1",
            strings.NewReader(bodyExample),
            http.StatusUnauthorized,
        },
        {//wrong json
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                s.EXPECT().Retrieve(tokens[3].String()).Return(&session.Data{Username:"dnp1"}, nil)
                rb := server.RouterBuilder{
                    Session:s,
                }
                return rb.Build()
            }(),
            "dnp1",
            strings.NewReader(`{"name":"golang",`),
            http.StatusBadRequest,
        },
        {//couldn't insert
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                r := mock_room.NewMockRoom(mockCtrl)
                sessionData := &session.Data{Username:"dnp1"}
                s.EXPECT().Retrieve(tokens[4].String()).Return(sessionData, nil)
                r.EXPECT().Create(sessionData.Username, "golang").Return(room.ErrCouldNotInsert)
                rb := server.RouterBuilder{
                    Session:s,
                    Room:r,
                }
                return rb.Build()
            }(),
            "dnp1",
            strings.NewReader(bodyExample),
            http.StatusConflict,
        },
        {//Everything ok
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                r := mock_room.NewMockRoom(mockCtrl)
                sessionData := &session.Data{Username:"dnp1"}
                s.EXPECT().Retrieve(tokens[5].String()).Return(sessionData, nil)
                r.EXPECT().Create(sessionData.Username, "golang").Return(nil)
                rb := server.RouterBuilder{
                    Session:s,
                    Room:r,
                }
                return rb.Build()
            }(),
            "dnp1",
            strings.NewReader(bodyExample),
            http.StatusOK,
        },
    }

    for i, c := range cases {
        url := fmt.Sprintf("/users/%s/rooms", c.user)
        req, err := http.NewRequest("POST", url, c.body)
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