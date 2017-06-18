package session_test

import (
    "testing"
    "net/http"
    "github.com/stretchr/testify/assert"
    "github.com/golang/mock/gomock"
    "github.com/dnp1/conversa/server/handlers/mock_session"
    "github.com/dnp1/conversa/server/mock_handlers"
    "github.com/dnp1/conversa/server/handlers/session"
    "github.com/dnp1/conversa/server/errors"
    "encoding/json"
)

//wrong json
func TestHandler_Login(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    s := mock_session.NewMockModel(ctrl)
    handler := session.New(s)
    expectedErr := errors.Validation(errors.FromString("wrong body!"))
    ctxt.EXPECT().BindJSON(gomock.Any()).Return(expectedErr)
    resp.EXPECT().SetError(expectedErr)
    handler.Login(ctxt, resp)
}

//unauthorized
func TestHandler_Login1(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    s := mock_session.NewMockModel(ctrl)
    handler := session.New(s)
    ctxt.EXPECT().BindJSON(gomock.Any()).Do(
        func(i interface{}) {
            const jsonStr = `{"username": "user", "password": "password"}`
            assert.NoError(t, json.Unmarshal([]byte(jsonStr), i))
        })
    expectedErr := errors.Validation(errors.FromString("wrong body!"))
    s.EXPECT().Create("user", "password").Return("", expectedErr)
    resp.EXPECT().SetError(gomock.Any()).Do(func(err errors.Error) {
        assert.True(t, err.Authentication())
    })
    handler.Login(ctxt, resp)
}

//ok
func TestHandler_Login2(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    s := mock_session.NewMockModel(ctrl)
    handler := session.New(s)
    ctxt.EXPECT().BindJSON(gomock.Any()).Do(
        func(i interface{}) {
            const jsonStr = `{"username": "user", "password": "password"}`
            assert.NoError(t, json.Unmarshal([]byte(jsonStr), i))
        })

    s.EXPECT().Create("user", "password").Return("", nil)
    resp.EXPECT().SetStatus(http.StatusCreated)
    resp.EXPECT().SetMessage(gomock.Any())
    ctxt.EXPECT().SetCookie(gomock.Any()).Do(
        func(cookie *http.Cookie) {
            assert.NotNil(t, cookie)
        })
    handler.Login(ctxt, resp)
}

//no cookie
func TestHandler_Logout(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    s := mock_session.NewMockModel(ctrl)
    handler := session.New(s)
    expectedErr := http.ErrNoCookie
    ctxt.EXPECT().Cookie(gomock.Any()).Return("", expectedErr)
    resp.EXPECT().SetStatus(http.StatusNoContent)
    resp.EXPECT().SetMessage(gomock.Any())
    handler.Logout(ctxt, resp)
}

//internal error
func TestHandler_Logout1(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    s := mock_session.NewMockModel(ctrl)
    handler := session.New(s)
    expectedErr := errors.FromString("unexpected error")
    ctxt.EXPECT().Cookie(gomock.Any()).Return("", expectedErr)
    resp.EXPECT().SetError(gomock.Any()).Do(
        func(err errors.Error) {
            assert.True(t, err.Internal())
        })
    handler.Logout(ctxt, resp)
}


//model internal error
func TestHandler_Logout2(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    s := mock_session.NewMockModel(ctrl)
    handler := session.New(s)
    ctxt.EXPECT().Cookie(gomock.Any()).Return("token", nil)
    expectedErr := errors.Internal(errors.FromString("unexpected error"))
    s.EXPECT().Delete(gomock.Any()).Return(expectedErr)
    resp.EXPECT().SetError(expectedErr)
    handler.Logout(ctxt, resp)
}

//model empty error
func TestHandler_Logout3(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    s := mock_session.NewMockModel(ctrl)
    handler := session.New(s)
    ctxt.EXPECT().Cookie(gomock.Any()).Return("token", nil)
    expectedErr := errors.Empty(errors.FromString("empty"))
    s.EXPECT().Delete(gomock.Any()).Return(expectedErr)
    resp.EXPECT().SetStatus(http.StatusResetContent)
    resp.EXPECT().SetMessage(gomock.Any())
    ctxt.EXPECT().DeleteCookie(gomock.Any())
    handler.Logout(ctxt, resp)
}

//ok
func TestHandler_Logout4(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    s := mock_session.NewMockModel(ctrl)
    handler := session.New(s)
    ctxt.EXPECT().Cookie(gomock.Any()).Return("", nil)
    s.EXPECT().Delete(gomock.Any()).Return(nil)
    resp.EXPECT().SetStatus(http.StatusOK)
    resp.EXPECT().SetMessage(gomock.Any())
    ctxt.EXPECT().DeleteCookie(gomock.Any())
    handler.Logout(ctxt, resp)
}

