// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/dnp1/conversa/server/session (interfaces: Session)

package mock_session

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of Session interface
type MockSession struct {
	ctrl     *gomock.Controller
	recorder *_MockSessionRecorder
}

// Recorder for MockSession (not exported)
type _MockSessionRecorder struct {
	mock *MockSession
}

func NewMockSession(ctrl *gomock.Controller) *MockSession {
	mock := &MockSession{ctrl: ctrl}
	mock.recorder = &_MockSessionRecorder{mock}
	return mock
}

func (_m *MockSession) EXPECT() *_MockSessionRecorder {
	return _m.recorder
}

func (_m *MockSession) Create(_param0 string, _param1 string) (string, error) {
	ret := _m.ctrl.Call(_m, "Create", _param0, _param1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockSessionRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Create", arg0, arg1)
}

func (_m *MockSession) Delete(_param0 string) error {
	ret := _m.ctrl.Call(_m, "Delete", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockSessionRecorder) Delete(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Delete", arg0)
}

func (_m *MockSession) Valid(_param0 string) error {
	ret := _m.ctrl.Call(_m, "Valid", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockSessionRecorder) Valid(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Valid", arg0)
}
