// Code generated by MockGen. DO NOT EDIT.
// Source: s3.go

// Package mock_s3 is a generated GoMock package.
package mock_s3

import (
	context "context"
	reflect "reflect"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	gomock "github.com/golang/mock/gomock"
)

// MockS3ObjectAPI is a mock of S3ObjectAPI interface.
type MockS3ObjectAPI struct {
	ctrl     *gomock.Controller
	recorder *MockS3ObjectAPIMockRecorder
}

// MockS3ObjectAPIMockRecorder is the mock recorder for MockS3ObjectAPI.
type MockS3ObjectAPIMockRecorder struct {
	mock *MockS3ObjectAPI
}

// NewMockS3ObjectAPI creates a new mock instance.
func NewMockS3ObjectAPI(ctrl *gomock.Controller) *MockS3ObjectAPI {
	mock := &MockS3ObjectAPI{ctrl: ctrl}
	mock.recorder = &MockS3ObjectAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockS3ObjectAPI) EXPECT() *MockS3ObjectAPIMockRecorder {
	return m.recorder
}

// PresignGetObject mocks base method.
func (m *MockS3ObjectAPI) PresignGetObject(ctx context.Context, input *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, input}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PresignGetObject", varargs...)
	ret0, _ := ret[0].(*v4.PresignedHTTPRequest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PresignGetObject indicates an expected call of PresignGetObject.
func (mr *MockS3ObjectAPIMockRecorder) PresignGetObject(ctx, input interface{}, optFns ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, input}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PresignGetObject", reflect.TypeOf((*MockS3ObjectAPI)(nil).PresignGetObject), varargs...)
}

// PutObject mocks base method.
func (m *MockS3ObjectAPI) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, params}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PutObject", varargs...)
	ret0, _ := ret[0].(*s3.PutObjectOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutObject indicates an expected call of PutObject.
func (mr *MockS3ObjectAPIMockRecorder) PutObject(ctx, params interface{}, optFns ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, params}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutObject", reflect.TypeOf((*MockS3ObjectAPI)(nil).PutObject), varargs...)
}
