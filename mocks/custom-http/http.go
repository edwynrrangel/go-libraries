// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/edwynrrangel/go-libraries/pkg/custom-http (interfaces: Http)

// Package customhttp is a generated GoMock package.
package customhttp

import (
	reflect "reflect"

	customhttp "github.com/edwynrrangel/go-libraries/pkg/custom-http"
	gomock "github.com/golang/mock/gomock"
)

// MockHttp is a mock of Http interface.
type MockHttp struct {
	ctrl     *gomock.Controller
	recorder *MockHttpMockRecorder
}

// MockHttpMockRecorder is the mock recorder for MockHttp.
type MockHttpMockRecorder struct {
	mock *MockHttp
}

// NewMockHttp creates a new mock instance.
func NewMockHttp(ctrl *gomock.Controller) *MockHttp {
	mock := &MockHttp{ctrl: ctrl}
	mock.recorder = &MockHttpMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHttp) EXPECT() *MockHttpMockRecorder {
	return m.recorder
}

// DoRequest mocks base method.
func (m *MockHttp) DoRequest(arg0 customhttp.ParamsRequest, arg1 int, arg2, arg3 interface{}) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DoRequest", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DoRequest indicates an expected call of DoRequest.
func (mr *MockHttpMockRecorder) DoRequest(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DoRequest", reflect.TypeOf((*MockHttp)(nil).DoRequest), arg0, arg1, arg2, arg3)
}
