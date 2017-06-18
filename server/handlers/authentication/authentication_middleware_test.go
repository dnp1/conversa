package authentication_test

import (
    "testing"
    "github.com/golang/mock/gomock"
    "github.com/dnp1/conversa/server/mock_handlers"
    "github.com/dnp1/conversa/server/handlers/mock_authentication"
    "github.com/dnp1/conversa/server/handlers/authentication"
    "net/http"
    "github.com/dnp1/conversa/server/errors"
    "github.com/stretchr/testify/assert"
    "github.com/dnp1/conversa/server/data/session"
)

const cookieName = "Authorization"

//no cookie found in request jar
func TestMiddleware_Middleware(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    model := mock_authentication.NewMockModel(ctrl)
    middleware := authentication.New(cookieName, model)
    resp.EXPECT().SetError(gomock.Any()).Do(
        func(err errors.Error) {
            assert.True(t, err.Authentication())
        })
    ctxt.EXPECT().Cookie(gomock.Any()).Return("", http.ErrNoCookie)
    ctxt.EXPECT().Abort()
    middleware.Middleware(ctxt, resp)
}

//token not found in database
func TestMiddleware_Middleware1(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    model := mock_authentication.NewMockModel(ctrl)
    middleware := authentication.New(cookieName, model)
    ctxt.EXPECT().Cookie(gomock.Any()).Return("tkn", nil)
    ctxt.EXPECT().Abort()
    expectedError := errors.Empty(errors.FromString("not found"))
    resp.EXPECT().SetError(gomock.Any()).Do(
        func(err errors.Error) {
            assert.True(t, err.Authentication())
        })
    model.EXPECT().Retrieve(gomock.Any()).Return(nil, expectedError)
    middleware.Middleware(ctxt, resp)
}

//authorized
func TestMiddleware_Middleware2(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    model := mock_authentication.NewMockModel(ctrl)
    middleware := authentication.New(cookieName, model)
    ctxt.EXPECT().Cookie(gomock.Any()).Return("tkn", nil)
    data := &session.Data{
        UserID:1,
        Username: "any",
    }
    model.EXPECT().Retrieve(gomock.Any()).Return(data,nil)
    ctxt.EXPECT().Set("username", data.Username)
    ctxt.EXPECT().Next()
    middleware.Middleware(ctxt, resp)
}