// Code generated by MockGen. DO NOT EDIT.
// Source: handlers.go

// Package clientv1mocks is a generated GoMock package.
package clientv1mocks

import (
	context "context"
	reflect "reflect"

	gethistory "github.com/evgeniy-krivenko/chat-service/internal/usecases/client/get-history"
	getuserprofile "github.com/evgeniy-krivenko/chat-service/internal/usecases/client/get-user-profile"
	sendmessage "github.com/evgeniy-krivenko/chat-service/internal/usecases/client/send-message"
	gomock "github.com/golang/mock/gomock"
)

// MockgetHistoryUseCase is a mock of getHistoryUseCase interface.
type MockgetHistoryUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockgetHistoryUseCaseMockRecorder
}

// MockgetHistoryUseCaseMockRecorder is the mock recorder for MockgetHistoryUseCase.
type MockgetHistoryUseCaseMockRecorder struct {
	mock *MockgetHistoryUseCase
}

// NewMockgetHistoryUseCase creates a new mock instance.
func NewMockgetHistoryUseCase(ctrl *gomock.Controller) *MockgetHistoryUseCase {
	mock := &MockgetHistoryUseCase{ctrl: ctrl}
	mock.recorder = &MockgetHistoryUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockgetHistoryUseCase) EXPECT() *MockgetHistoryUseCaseMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MockgetHistoryUseCase) Handle(ctx context.Context, req gethistory.Request) (gethistory.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", ctx, req)
	ret0, _ := ret[0].(gethistory.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Handle indicates an expected call of Handle.
func (mr *MockgetHistoryUseCaseMockRecorder) Handle(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockgetHistoryUseCase)(nil).Handle), ctx, req)
}

// MocksendMessageUseCase is a mock of sendMessageUseCase interface.
type MocksendMessageUseCase struct {
	ctrl     *gomock.Controller
	recorder *MocksendMessageUseCaseMockRecorder
}

// MocksendMessageUseCaseMockRecorder is the mock recorder for MocksendMessageUseCase.
type MocksendMessageUseCaseMockRecorder struct {
	mock *MocksendMessageUseCase
}

// NewMocksendMessageUseCase creates a new mock instance.
func NewMocksendMessageUseCase(ctrl *gomock.Controller) *MocksendMessageUseCase {
	mock := &MocksendMessageUseCase{ctrl: ctrl}
	mock.recorder = &MocksendMessageUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocksendMessageUseCase) EXPECT() *MocksendMessageUseCaseMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MocksendMessageUseCase) Handle(ctx context.Context, req sendmessage.Request) (sendmessage.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", ctx, req)
	ret0, _ := ret[0].(sendmessage.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Handle indicates an expected call of Handle.
func (mr *MocksendMessageUseCaseMockRecorder) Handle(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MocksendMessageUseCase)(nil).Handle), ctx, req)
}

// MockgetUserProfile is a mock of getUserProfile interface.
type MockgetUserProfile struct {
	ctrl     *gomock.Controller
	recorder *MockgetUserProfileMockRecorder
}

// MockgetUserProfileMockRecorder is the mock recorder for MockgetUserProfile.
type MockgetUserProfileMockRecorder struct {
	mock *MockgetUserProfile
}

// NewMockgetUserProfile creates a new mock instance.
func NewMockgetUserProfile(ctrl *gomock.Controller) *MockgetUserProfile {
	mock := &MockgetUserProfile{ctrl: ctrl}
	mock.recorder = &MockgetUserProfileMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockgetUserProfile) EXPECT() *MockgetUserProfileMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MockgetUserProfile) Handle(ctx context.Context, req getuserprofile.Request) (getuserprofile.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", ctx, req)
	ret0, _ := ret[0].(getuserprofile.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Handle indicates an expected call of Handle.
func (mr *MockgetUserProfileMockRecorder) Handle(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockgetUserProfile)(nil).Handle), ctx, req)
}
