// Code generated by MockGen. DO NOT EDIT.
// Source: libs/cloudfrontkeyvaluestore.go
//
// Generated by this command:
//
//	mockgen -source=libs/cloudfrontkeyvaluestore.go -destination=./libs/cloudfrontkeyvaluestore_mock.go -package=libs
//

// Package libs is a generated GoMock package.
package libs

import (
	context "context"
	reflect "reflect"

	cloudfrontkeyvaluestore "github.com/aws/aws-sdk-go-v2/service/cloudfrontkeyvaluestore"
	gomock "go.uber.org/mock/gomock"
)

// MockCloudFrontKeyValueStoreClient is a mock of CloudFrontKeyValueStoreClient interface.
type MockCloudFrontKeyValueStoreClient struct {
	ctrl     *gomock.Controller
	recorder *MockCloudFrontKeyValueStoreClientMockRecorder
}

// MockCloudFrontKeyValueStoreClientMockRecorder is the mock recorder for MockCloudFrontKeyValueStoreClient.
type MockCloudFrontKeyValueStoreClientMockRecorder struct {
	mock *MockCloudFrontKeyValueStoreClient
}

// NewMockCloudFrontKeyValueStoreClient creates a new mock instance.
func NewMockCloudFrontKeyValueStoreClient(ctrl *gomock.Controller) *MockCloudFrontKeyValueStoreClient {
	mock := &MockCloudFrontKeyValueStoreClient{ctrl: ctrl}
	mock.recorder = &MockCloudFrontKeyValueStoreClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCloudFrontKeyValueStoreClient) EXPECT() *MockCloudFrontKeyValueStoreClientMockRecorder {
	return m.recorder
}

// DeleteKey mocks base method.
func (m *MockCloudFrontKeyValueStoreClient) DeleteKey(ctx context.Context, params *cloudfrontkeyvaluestore.DeleteKeyInput, optFns ...func(*cloudfrontkeyvaluestore.Options)) (*cloudfrontkeyvaluestore.DeleteKeyOutput, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, params}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteKey", varargs...)
	ret0, _ := ret[0].(*cloudfrontkeyvaluestore.DeleteKeyOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteKey indicates an expected call of DeleteKey.
func (mr *MockCloudFrontKeyValueStoreClientMockRecorder) DeleteKey(ctx, params any, optFns ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, params}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteKey", reflect.TypeOf((*MockCloudFrontKeyValueStoreClient)(nil).DeleteKey), varargs...)
}

// DescribeKeyValueStore mocks base method.
func (m *MockCloudFrontKeyValueStoreClient) DescribeKeyValueStore(ctx context.Context, params *cloudfrontkeyvaluestore.DescribeKeyValueStoreInput, optFns ...func(*cloudfrontkeyvaluestore.Options)) (*cloudfrontkeyvaluestore.DescribeKeyValueStoreOutput, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, params}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DescribeKeyValueStore", varargs...)
	ret0, _ := ret[0].(*cloudfrontkeyvaluestore.DescribeKeyValueStoreOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeKeyValueStore indicates an expected call of DescribeKeyValueStore.
func (mr *MockCloudFrontKeyValueStoreClientMockRecorder) DescribeKeyValueStore(ctx, params any, optFns ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, params}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeKeyValueStore", reflect.TypeOf((*MockCloudFrontKeyValueStoreClient)(nil).DescribeKeyValueStore), varargs...)
}

// GetKey mocks base method.
func (m *MockCloudFrontKeyValueStoreClient) GetKey(ctx context.Context, params *cloudfrontkeyvaluestore.GetKeyInput, optFns ...func(*cloudfrontkeyvaluestore.Options)) (*cloudfrontkeyvaluestore.GetKeyOutput, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, params}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetKey", varargs...)
	ret0, _ := ret[0].(*cloudfrontkeyvaluestore.GetKeyOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKey indicates an expected call of GetKey.
func (mr *MockCloudFrontKeyValueStoreClientMockRecorder) GetKey(ctx, params any, optFns ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, params}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKey", reflect.TypeOf((*MockCloudFrontKeyValueStoreClient)(nil).GetKey), varargs...)
}

// ListKeys mocks base method.
func (m *MockCloudFrontKeyValueStoreClient) ListKeys(ctx context.Context, params *cloudfrontkeyvaluestore.ListKeysInput, optFns ...func(*cloudfrontkeyvaluestore.Options)) (*cloudfrontkeyvaluestore.ListKeysOutput, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, params}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListKeys", varargs...)
	ret0, _ := ret[0].(*cloudfrontkeyvaluestore.ListKeysOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListKeys indicates an expected call of ListKeys.
func (mr *MockCloudFrontKeyValueStoreClientMockRecorder) ListKeys(ctx, params any, optFns ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, params}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListKeys", reflect.TypeOf((*MockCloudFrontKeyValueStoreClient)(nil).ListKeys), varargs...)
}

// PutKey mocks base method.
func (m *MockCloudFrontKeyValueStoreClient) PutKey(ctx context.Context, params *cloudfrontkeyvaluestore.PutKeyInput, optFns ...func(*cloudfrontkeyvaluestore.Options)) (*cloudfrontkeyvaluestore.PutKeyOutput, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, params}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PutKey", varargs...)
	ret0, _ := ret[0].(*cloudfrontkeyvaluestore.PutKeyOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutKey indicates an expected call of PutKey.
func (mr *MockCloudFrontKeyValueStoreClientMockRecorder) PutKey(ctx, params any, optFns ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, params}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutKey", reflect.TypeOf((*MockCloudFrontKeyValueStoreClient)(nil).PutKey), varargs...)
}

// UpdateKeys mocks base method.
func (m *MockCloudFrontKeyValueStoreClient) UpdateKeys(ctx context.Context, params *cloudfrontkeyvaluestore.UpdateKeysInput, optFns ...func(*cloudfrontkeyvaluestore.Options)) (*cloudfrontkeyvaluestore.UpdateKeysOutput, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, params}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateKeys", varargs...)
	ret0, _ := ret[0].(*cloudfrontkeyvaluestore.UpdateKeysOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateKeys indicates an expected call of UpdateKeys.
func (mr *MockCloudFrontKeyValueStoreClientMockRecorder) UpdateKeys(ctx, params any, optFns ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, params}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateKeys", reflect.TypeOf((*MockCloudFrontKeyValueStoreClient)(nil).UpdateKeys), varargs...)
}
