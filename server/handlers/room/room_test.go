package room_test

import (
    "testing"
    "github.com/golang/mock/gomock"
    "github.com/dnp1/conversa/server/mock_handlers"
    "github.com/dnp1/conversa/server/handlers/mock_room"
    "github.com/dnp1/conversa/server/handlers/room"
    data "github.com/dnp1/conversa/server/data/room"
    "github.com/dnp1/conversa/server/errors"
    "github.com/stretchr/testify/assert"
    "encoding/json"
    "net/http"
)

//assertion failed
func TestHandler_Create(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    model := mock_room.NewMockModel(ctrl)
    handler := room.New(model)

    expectedErr := errors.Internal(errors.FromString("unexpected Err!"))
    ctxt.EXPECT().Param(gomock.Any()).Return("user")
    ctxt.EXPECT().ShouldGetString(gomock.Any()).Return("", expectedErr)
    resp.EXPECT().SetError(expectedErr)
    handler.Create(ctxt, resp)
}

//authorization error
func TestHandler_Create1(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    model := mock_room.NewMockModel(ctrl)
    handler := room.New(model)


    ctxt.EXPECT().Param(gomock.Any()).Return("user")
    ctxt.EXPECT().ShouldGetString(gomock.Any()).Return("", nil)
    resp.EXPECT().SetError(gomock.Any()).Do(
        func(err errors.Error) {
            assert.True(t, err.Authorization())
        })
    handler.Create(ctxt, resp)
}

//wrong body error
func TestHandler_Create2(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    model := mock_room.NewMockModel(ctrl)
    handler := room.New(model)


    ctxt.EXPECT().Param(gomock.Any()).Return("user")
    ctxt.EXPECT().ShouldGetString(gomock.Any()).Return("user", nil)
    expectedErr := errors.Validation(errors.FromString("wrong body"))
    ctxt.EXPECT().BindJSON(gomock.Any()).Return(expectedErr)
    resp.EXPECT().SetError(expectedErr)
    handler.Create(ctxt, resp)
}

//error when inserting
func TestHandler_Create3(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    model := mock_room.NewMockModel(ctrl)
    handler := room.New(model)


    ctxt.EXPECT().Param(gomock.Any()).Return("user")
    ctxt.EXPECT().ShouldGetString(gomock.Any()).Return("user", nil)
    expectedErr := errors.Validation(errors.FromString("wrong body"))
    ctxt.EXPECT().BindJSON(gomock.Any()).Do(
        func(i interface{}) {
            const jsonStr = `{"name": "new_room"}`
            json.Unmarshal([]byte(jsonStr), i)
        })
    model.EXPECT().Create("user", "new_room").Return(expectedErr)
    resp.EXPECT().SetError(expectedErr)
    handler.Create(ctxt, resp)
}


//ok
func TestHandler_Create4(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    model := mock_room.NewMockModel(ctrl)
    handler := room.New(model)

    ctxt.EXPECT().Param(gomock.Any()).Return("user")
    ctxt.EXPECT().ShouldGetString(gomock.Any()).Return("user", nil)
    ctxt.EXPECT().BindJSON(gomock.Any()).Do(
        func(i interface{}) {
            const jsonStr = `{"name": "new_room"}`
            json.Unmarshal([]byte(jsonStr), i)
        })
    model.EXPECT().Create("user", "new_room").Return(nil)
    resp.EXPECT().SetMessage(gomock.Any())
    resp.EXPECT().SetStatus(http.StatusCreated)
    handler.Create(ctxt, resp)
}

//assertion failed
func TestHandler_Delete(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    model := mock_room.NewMockModel(ctrl)
    handler := room.New(model)

    expectedErr := errors.Internal(errors.FromString("unexpected Err!"))
    ctxt.EXPECT().Param(gomock.Any()).Return("any").AnyTimes()
    ctxt.EXPECT().ShouldGetString(gomock.Any()).Return("", expectedErr)
    resp.EXPECT().SetError(expectedErr)
    handler.Delete(ctxt, resp)
}

//authorization
func TestHandler_Delete2(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    model := mock_room.NewMockModel(ctrl)
    handler := room.New(model)

    ctxt.EXPECT().Param(gomock.Any()).Return("any").AnyTimes()
    ctxt.EXPECT().ShouldGetString(gomock.Any()).Return("", nil)
    resp.EXPECT().SetError(gomock.Any()).Do(
        func(err errors.Error) {
            assert.True(t, err.Authorization())
        })
    handler.Delete(ctxt, resp)
}

//no content
func TestHandler_Delete3(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    model := mock_room.NewMockModel(ctrl)
    handler := room.New(model)

    ctxt.EXPECT().Param(gomock.Any()).Return("any").AnyTimes()
    ctxt.EXPECT().ShouldGetString(gomock.Any()).Return("any", nil)
    expectedErr := errors.Empty(errors.FromString("no content"))
    model.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(expectedErr)
    resp.EXPECT().SetStatus(http.StatusNoContent)
    resp.EXPECT().SetMessage(gomock.Any())
    handler.Delete(ctxt, resp)
}

//unexpected error
func TestHandler_Delete4(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    model := mock_room.NewMockModel(ctrl)
    handler := room.New(model)

    ctxt.EXPECT().Param(gomock.Any()).Return("any").AnyTimes()
    ctxt.EXPECT().ShouldGetString(gomock.Any()).Return("any", nil)
    expectedErr := errors.Internal(errors.FromString("no content"))
    model.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(expectedErr)
    resp.EXPECT().SetError(expectedErr)
    handler.Delete(ctxt, resp)
}

//ok
func TestHandler_Delete5(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    model := mock_room.NewMockModel(ctrl)
    handler := room.New(model)

    ctxt.EXPECT().Param(gomock.Any()).Return("any").AnyTimes()
    ctxt.EXPECT().ShouldGetString(gomock.Any()).Return("any", nil)
    model.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
    resp.EXPECT().SetStatus(http.StatusOK)
    resp.EXPECT().SetMessage(gomock.Any())
    handler.Delete(ctxt, resp)
}



//assertion failed
func TestHandler_Edit(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    model := mock_room.NewMockModel(ctrl)
    handler := room.New(model)

    expectedErr := errors.Internal(errors.FromString("unexpected Err!"))
    ctxt.EXPECT().Param(gomock.Any()).Return("any").AnyTimes()
    ctxt.EXPECT().ShouldGetString(gomock.Any()).Return("", expectedErr)
    resp.EXPECT().SetError(expectedErr)
    handler.Edit(ctxt, resp)
}

//authorization
func TestHandler_Edit2(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    model := mock_room.NewMockModel(ctrl)
    handler := room.New(model)

    ctxt.EXPECT().Param(gomock.Any()).Return("any").AnyTimes()
    ctxt.EXPECT().ShouldGetString(gomock.Any()).Return("", nil)
    resp.EXPECT().SetError(gomock.Any()).Do(
        func(err errors.Error) {
            assert.True(t, err.Authorization())
        })
    handler.Edit(ctxt, resp)
}

//wrong body error
func TestHandler_Edit3(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    model := mock_room.NewMockModel(ctrl)
    handler := room.New(model)


    ctxt.EXPECT().Param(gomock.Any()).Return("user")
    ctxt.EXPECT().ShouldGetString(gomock.Any()).Return("user", nil)
    expectedErr := errors.Validation(errors.FromString("wrong body"))
    ctxt.EXPECT().BindJSON(gomock.Any()).Return(expectedErr)
    resp.EXPECT().SetError(expectedErr)
    handler.Edit(ctxt, resp)
}

//error when inserting
func TestHandler_Edit4(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    model := mock_room.NewMockModel(ctrl)
    handler := room.New(model)


    ctxt.EXPECT().Param(gomock.Any()).Return("any").AnyTimes()
    ctxt.EXPECT().ShouldGetString(gomock.Any()).Return("any", nil)
    ctxt.EXPECT().BindJSON(gomock.Any()).Do(
        func(i interface{}) {
            const jsonStr = `{"name": "new_room"}`
            json.Unmarshal([]byte(jsonStr), i)
        })
    expectedErr := errors.Validation(errors.FromString("validation"))
    model.EXPECT().Rename(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedErr)
    resp.EXPECT().SetError(expectedErr)
    handler.Edit(ctxt, resp)
}

//ok
func TestHandler_Edit5(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    model := mock_room.NewMockModel(ctrl)
    handler := room.New(model)


    ctxt.EXPECT().Param(gomock.Any()).Return("any").AnyTimes()
    ctxt.EXPECT().ShouldGetString(gomock.Any()).Return("any", nil)
    ctxt.EXPECT().BindJSON(gomock.Any()).Do(
        func(i interface{}) {
            const jsonStr = `{"name": "new_room"}`
            json.Unmarshal([]byte(jsonStr), i)
        })
    model.EXPECT().Rename(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
    resp.EXPECT().SetStatus(http.StatusOK)
    resp.EXPECT().SetMessage(gomock.Any())
    handler.Edit(ctxt, resp)
}

//error
func TestHandler_List(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    model := mock_room.NewMockModel(ctrl)
    handler := room.New(model)

    expectedErr := errors.Internal(errors.FromString("dsad"))
    model.EXPECT().All().Return(nil, expectedErr)
    resp.EXPECT().SetError(expectedErr)
    handler.List(ctxt, resp)
}

//ok
func TestHandler_List1(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    ctxt := mock_handlers.NewMockContext(ctrl)
    resp := mock_handlers.NewMockJsonResponse(ctrl)
    model := mock_room.NewMockModel(ctrl)
    handler := room.New(model)

    dt := []data.Data{}
    model.EXPECT().All().Return(dt, nil)
    resp.EXPECT().SetMessage(gomock.Any())
    resp.EXPECT().SetData(dt)
    resp.EXPECT().SetStatus(http.StatusOK)
    handler.List(ctxt, resp)
}