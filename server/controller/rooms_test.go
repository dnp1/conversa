package controller_test

import (
    "gopkg.in/gin-gonic/gin.v1"
    "testing"
    "io"
    "github.com/golang/mock/gomock"
    "strings"
    "github.com/stretchr/testify/assert"
    "net/http/httptest"
    "net/http"
    "github.com/dnp1/conversa/server/controller"
    "github.com/twinj/uuid"
    "github.com/dnp1/conversa/server/mock_model/mock_session"
    "github.com/dnp1/conversa/server/model/session"
    "fmt"
    "github.com/dnp1/conversa/server/mock_model/mock_room"
    "github.com/dnp1/conversa/server/model/room"
    "github.com/pkg/errors"
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
        uuid.NewV4(),
    }

    const bodyExample = `{"name":"golang"}`
    cases := []Case{
        {//no token
            NoDependencyRouter(),
            "dnp1",
            strings.NewReader(bodyExample),
            http.StatusUnauthorized,
        },
        {//invalid token
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                s.EXPECT().Retrieve(tokens[1].String()).Return(nil, session.ErrTokenNotFound)
                rb := controller.RouterBuilder{
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
                rb := controller.RouterBuilder{
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
                rb := controller.RouterBuilder{
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
                r.EXPECT().Create(sessionData.Username, "golang").Return(errors.New("Unexpected error"))
                rb := controller.RouterBuilder{
                    Session:s,
                    Room:r,
                }
                return rb.Build()
            }(),
            "dnp1",
            strings.NewReader(bodyExample),
            http.StatusInternalServerError,
        },
        {//couldn't insert
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                r := mock_room.NewMockRoom(mockCtrl)
                sessionData := &session.Data{Username:"dnp1"}
                s.EXPECT().Retrieve(tokens[5].String()).Return(sessionData, nil)
                r.EXPECT().Create(sessionData.Username, "golang").Return(room.ErrRoomNameAlreadyExists)
                rb := controller.RouterBuilder{
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
                s.EXPECT().Retrieve(tokens[6].String()).Return(sessionData, nil)
                r.EXPECT().Create(sessionData.Username, "golang").Return(nil)
                rb := controller.RouterBuilder{
                    Session:s,
                    Room:r,
                }
                return rb.Build()
            }(),
            "dnp1",
            strings.NewReader(bodyExample),
            http.StatusCreated,
        },
    }

    for i, c := range cases {
        url := fmt.Sprintf("/users/%s/rooms", c.user)
        req, err := http.NewRequest("POST", url, c.body)
        if tokens[i] != nil {
            req.AddCookie(&http.Cookie{
                Name:     controller.TokenCookieName,
                Value:    tokens[i].String(),
                MaxAge:   24 * 60 * 60,
                Secure:   true,
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


func TestRoomController_DeleteRoom(t *testing.T) {
    type Case struct {
        router *gin.Engine
        user   string
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

    const roomName = "golang"
    cases := []Case{
        {//no token
            NoDependencyRouter(),
            "dnp1",
            http.StatusUnauthorized,
        },
        {//invalid token
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                s.EXPECT().Retrieve(tokens[1].String()).Return(nil, session.ErrTokenNotFound)
                rb := controller.RouterBuilder{
                    Session:s,
                }
                return rb.Build()
            }(),
            "dnp1",
            http.StatusUnauthorized,
        },
        {//trying delete other users room
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                s.EXPECT().Retrieve(tokens[2].String()).Return(nil, session.ErrTokenNotFound)
                rb := controller.RouterBuilder{
                    Session:s,
                }
                return rb.Build()
            }(),
            "dnp1",
            http.StatusUnauthorized,
        },
        {//couldn't delete
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                r := mock_room.NewMockRoom(mockCtrl)
                sessionData := &session.Data{Username:"dnp1"}
                s.EXPECT().Retrieve(tokens[3].String()).Return(sessionData, nil)
                r.EXPECT().Delete(sessionData.Username, roomName).Return(room.ErrCouldNotDelete)
                rb := controller.RouterBuilder{
                    Session:s,
                    Room:r,
                }
                return rb.Build()
            }(),
            "dnp1",
            http.StatusNoContent,
        },
        {//Everything ok
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                r := mock_room.NewMockRoom(mockCtrl)
                sessionData := &session.Data{Username:"dnp1"}
                s.EXPECT().Retrieve(tokens[4].String()).Return(sessionData, nil)
                r.EXPECT().Delete(sessionData.Username, roomName).Return(nil)
                rb := controller.RouterBuilder{
                    Session:s,
                    Room:r,
                }
                return rb.Build()
            }(),
            "dnp1",
            http.StatusOK,
        },
    }

    for i, c := range cases {
        url := fmt.Sprintf("/users/%s/rooms/%s", c.user, roomName)
        req, err := http.NewRequest("DELETE", url, nil)
        if tokens[i] != nil {
            req.AddCookie(&http.Cookie{
                Name:     controller.TokenCookieName,
                Value:    tokens[i].String(),
                MaxAge:   24 * 60 * 60,
                Secure:   true,
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


func TestRoomController_EditRoom(t *testing.T) {
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

    const oldRoomName = "go"
    const newRoomName = "golang"
    const bodyExample = `{"name":"`+ newRoomName +`"}`
    cases := []Case{
        {//no token
            NoDependencyRouter(),
            "dnp1",
            strings.NewReader(bodyExample),
            http.StatusUnauthorized,
        },
        {//invalid token
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                s.EXPECT().Retrieve(tokens[1].String()).Return(nil, session.ErrTokenNotFound)
                rb := controller.RouterBuilder{
                    Session:s,
                }
                return rb.Build()
            }(),
            "dnp1",
            strings.NewReader(bodyExample),
            http.StatusUnauthorized,
        },
        {//trying update room owned by other user
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                s.EXPECT().Retrieve(tokens[2].String()).Return(nil, session.ErrTokenNotFound)
                rb := controller.RouterBuilder{
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
                rb := controller.RouterBuilder{
                    Session:s,
                }
                return rb.Build()
            }(),
            "dnp1",
            strings.NewReader(`{"name":"golang",`),
            http.StatusBadRequest,
        },
        {//couldn't update
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                r := mock_room.NewMockRoom(mockCtrl)
                sessionData := &session.Data{Username:"dnp1"}
                s.EXPECT().Retrieve(tokens[4].String()).Return(sessionData, nil)
                r.EXPECT().Rename(sessionData.Username, oldRoomName, newRoomName).Return(room.ErrCouldNotRename)
                rb := controller.RouterBuilder{
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
                r.EXPECT().Rename(sessionData.Username, oldRoomName, newRoomName).Return(nil)
                rb := controller.RouterBuilder{
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
        url := fmt.Sprintf("/users/%s/rooms/%s", c.user, oldRoomName)
        req, err := http.NewRequest("PATCH", url, c.body)
        if tokens[i] != nil {
            req.AddCookie(&http.Cookie{
                Name:     controller.TokenCookieName,
                Value:    tokens[i].String(),
                MaxAge:   24 * 60 * 60,
                Secure:   true,
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

func TestRoomController_ListRooms(t *testing.T) {
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
    }

    const RoomName = "golang"
    const bodyExample = `{"name":"`+ RoomName +`"}`
    cases := []Case{
        {//no token
            NoDependencyRouter(),
            "dnp1",
            strings.NewReader(bodyExample),
            http.StatusUnauthorized,
        },
        {//invalid token
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                s.EXPECT().Retrieve(tokens[1].String()).Return(nil, session.ErrTokenNotFound)
                rb := controller.RouterBuilder{
                    Session:s,
                }
                return rb.Build()
            }(),
            "dnp1",
            strings.NewReader(bodyExample),
            http.StatusUnauthorized,
        },
        {//Model can't  retrieve data
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                r := mock_room.NewMockRoom(mockCtrl)
                s.EXPECT().Retrieve(tokens[2].String()).Return(&session.Data{Username: "dnp1"}, nil)
                r.EXPECT().All().Return(nil, room.ErrCouldNotRetrieveRooms)
                rb := controller.RouterBuilder{
                    Session:s,
                    Room:r,
                }
                return rb.Build()
            }(),
            "dnp1",
            strings.NewReader(bodyExample),
            http.StatusInternalServerError,
        },
        {//Everything ok
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                r := mock_room.NewMockRoom(mockCtrl)
                s.EXPECT().Retrieve(tokens[3].String()).Return(&session.Data{Username: "dnp1"}, nil)
                r.EXPECT().All().Return([]room.Data{}, nil)
                rb := controller.RouterBuilder{
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

        req, err := http.NewRequest("GET", "/rooms", nil)
        if tokens[i] != nil {
            req.AddCookie(&http.Cookie{
                Name:     controller.TokenCookieName,
                Value:    tokens[i].String(),
                MaxAge:   24 * 60 * 60,
                Secure:   true,
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


func TestRoomController_ListUserRooms(t *testing.T) {
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
    }

    const userName = "korra"
    const roomName = "golang"
    const bodyExample = `{"name":"`+ roomName +`"}`
    cases := []Case{
        {//no token
            NoDependencyRouter(),
            "dnp1",
            strings.NewReader(bodyExample),
            http.StatusUnauthorized,
        },
        {//invalid token
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                s.EXPECT().Retrieve(tokens[1].String()).Return(nil, session.ErrTokenNotFound)
                rb := controller.RouterBuilder{
                    Session:s,
                }
                return rb.Build()
            }(),
            "dnp1",
            strings.NewReader(bodyExample),
            http.StatusUnauthorized,
        },
        {//Model lookup failure
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                r := mock_room.NewMockRoom(mockCtrl)
                s.EXPECT().Retrieve(tokens[2].String()).Return(&session.Data{Username: "dnp1"}, nil)
                r.EXPECT().AllByUser(userName).Return(nil, room.ErrCouldNotRetrieveRooms)
                rb := controller.RouterBuilder{
                    Session:s,
                    Room:r,
                }
                return rb.Build()
            }(),
            "dnp1",
            strings.NewReader(bodyExample),
            http.StatusInternalServerError,
        },
        {//Everything ok
            func() *gin.Engine {
                s := mock_session.NewMockSession(mockCtrl)
                r := mock_room.NewMockRoom(mockCtrl)
                s.EXPECT().Retrieve(tokens[3].String()).Return(&session.Data{Username: "dnp1"}, nil)
                r.EXPECT().AllByUser(userName).Return([]room.Data{}, nil)
                rb := controller.RouterBuilder{
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
        url := fmt.Sprintf("/users/%s/rooms", userName)
        req, err := http.NewRequest("GET", url, nil)
        if tokens[i] != nil {
            req.AddCookie(&http.Cookie{
                Name:     controller.TokenCookieName,
                Value:    tokens[i].String(),
                MaxAge:   24 * 60 * 60,
                Secure:   true,
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