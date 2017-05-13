// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/dnp1/conversa/server/user (interfaces: User)

package mock_user

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of User interface
type MockUser struct {
	ctrl     *gomock.Controller
	recorder *_MockUserRecorder
}

// Recorder for MockUser (not exported)
type _MockUserRecorder struct {
	mock *MockUser
}

func NewMockUser(ctrl *gomock.Controller) *MockUser {
	mock := &MockUser{ctrl: ctrl}
	mock.recorder = &_MockUserRecorder{mock}
	return mock
}

func (_m *MockUser) EXPECT() *_MockUserRecorder {
	return _m.recorder
}

func (_m *MockUser) Create(_param0 string, _param1 string, _param2 string) error {
	ret := _m.ctrl.Call(_m, "Create", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockUserRecorder) Create(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Create", arg0, arg1, arg2)
}
