// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/dnp1/conversa/server/handlers/user (interfaces: Model)

package mock_user

import (
	errors "github.com/dnp1/conversa/server/errors"
	gomock "github.com/golang/mock/gomock"
)

// Mock of Model interface
type MockModel struct {
	ctrl     *gomock.Controller
	recorder *_MockModelRecorder
}

// Recorder for MockModel (not exported)
type _MockModelRecorder struct {
	mock *MockModel
}

func NewMockModel(ctrl *gomock.Controller) *MockModel {
	mock := &MockModel{ctrl: ctrl}
	mock.recorder = &_MockModelRecorder{mock}
	return mock
}

func (_m *MockModel) EXPECT() *_MockModelRecorder {
	return _m.recorder
}

func (_m *MockModel) Create(_param0 string, _param1 string, _param2 string) errors.Error {
	ret := _m.ctrl.Call(_m, "Create", _param0, _param1, _param2)
	ret0, _ := ret[0].(errors.Error)
	return ret0
}

func (_mr *_MockModelRecorder) Create(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Create", arg0, arg1, arg2)
}
