// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/edwynrrangel/go-libraries/pkg/custom-http (interfaces: CustomHttp)

// Package customhttp is a generated GoMock package.
package customhttp

import (
	reflect "reflect"

	customhttp "github.com/edwynrrangel/go-libraries/pkg/custom-http"
	gomock "github.com/golang/mock/gomock"
)

// MockCustomHttp is a mock of CustomHttp interface.
type MockCustomHttp struct {
	ctrl     *gomock.Controller
	recorder *MockCustomHttpMockRecorder
}

// MockCustomHttpMockRecorder is the mock recorder for MockCustomHttp.
type MockCustomHttpMockRecorder struct {
	mock *MockCustomHttp
}

// NewMockCustomHttp creates a new mock instance.
func NewMockCustomHttp(ctrl *gomock.Controller) *MockCustomHttp {
	mock := &MockCustomHttp{ctrl: ctrl}
	mock.recorder = &MockCustomHttpMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCustomHttp) EXPECT() *MockCustomHttpMockRecorder {
	return m.recorder
}

// DoRequest mocks base method.
func (m *MockCustomHttp) DoRequest(arg0 customhttp.ParamsRequest, arg1 int, arg2, arg3 interface{}) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DoRequest", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DoRequest indicates an expected call of DoRequest.
func (mr *MockCustomHttpMockRecorder) DoRequest(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DoRequest", reflect.TypeOf((*MockCustomHttp)(nil).DoRequest), arg0, arg1, arg2, arg3)
}
