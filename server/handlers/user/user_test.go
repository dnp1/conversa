package user_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "net/http"
    "github.com/golang/mock/gomock"
    "github.com/dnp1/conversa/server/mock_handlers"
    "github.com/dnp1/conversa/server/handlers/user"
    "github.com/dnp1/conversa/server/handlers/mock_user"
    "github.com/dnp1/conversa/server/errors"
    "encoding/json"
)

//TestHandler_CreateUser case wrong json body
func TestHandler_CreateUser(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    u := mock_user.NewMockModel(ctrl)
    handler := user.New(u)
    expectedErr := errors.Validation(errors.FromString("wrong body!"))
    ctxt.EXPECT().BindJSON(gomock.Any()).Return(expectedErr)
    resp.EXPECT().SetError(expectedErr).Do(
        func(i errors.Error) {
            assert.True(t, i.Validation())
        })
    handler.Create(ctxt, resp)
}

//TestHandler_CreateUser when create fails for a password doesn't match, for example
func TestHandler_CreateUser1(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    u := mock_user.NewMockModel(ctrl)
    handler := user.New(u)
    expectedErr := errors.Validation(errors.FromString("password doesn't matc!"))
    ctxt.EXPECT().BindJSON(gomock.Any()).Do(
        func(i interface{}) {
            const jsonStr = `{"username":"user", "password": "senha", "passwordConfirmation":"senha"}`
            assert.NoError(t, json.Unmarshal([]byte(jsonStr), i))
        })
    u.EXPECT().Create("user", "senha", "senha").Return(expectedErr)
    resp.EXPECT().SetError(expectedErr).Do(
        func(i errors.Error) {
            assert.True(t, i.Validation())
        })
    handler.Create(ctxt, resp)
}


//TestHandler_CreateUser when everything ok
func TestHandler_CreateUser2(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    u := mock_user.NewMockModel(ctrl)
    handler := user.New(u)
    ctxt.EXPECT().BindJSON(gomock.Any()).Do(
        func(i interface{}) {
            const jsonStr = `{"username":"user", "password": "senha", "passwordConfirmation":"senha"}`
            assert.NoError(t, json.Unmarshal([]byte(jsonStr), i))
        })
    u.EXPECT().Create("user", "senha", "senha").Return(nil)
    resp.EXPECT().SetMessage(gomock.Any())
    resp.EXPECT().SetStatus(http.StatusCreated)
    handler.Create(ctxt, resp)
}
