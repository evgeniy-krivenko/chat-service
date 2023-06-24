// Code generated by MockGen. DO NOT EDIT.
// Source: handlers.go

// Package managerv1mocks is a generated GoMock package.
package managerv1mocks

import (
	context "context"
	reflect "reflect"

	canreceiveproblems "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/can-receive-problems"
	closechat "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/close-chat"
	freehands "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/free-hands"
	getchathistory "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/get-chat-history"
	getchats "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/get-chats"
	managerlogin "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/manager-login"
	sendmessage "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/send-message"
	gomock "github.com/golang/mock/gomock"
)

// MockcanReceiveProblemsUseCase is a mock of canReceiveProblemsUseCase interface.
type MockcanReceiveProblemsUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockcanReceiveProblemsUseCaseMockRecorder
}

// MockcanReceiveProblemsUseCaseMockRecorder is the mock recorder for MockcanReceiveProblemsUseCase.
type MockcanReceiveProblemsUseCaseMockRecorder struct {
	mock *MockcanReceiveProblemsUseCase
}

// NewMockcanReceiveProblemsUseCase creates a new mock instance.
func NewMockcanReceiveProblemsUseCase(ctrl *gomock.Controller) *MockcanReceiveProblemsUseCase {
	mock := &MockcanReceiveProblemsUseCase{ctrl: ctrl}
	mock.recorder = &MockcanReceiveProblemsUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockcanReceiveProblemsUseCase) EXPECT() *MockcanReceiveProblemsUseCaseMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MockcanReceiveProblemsUseCase) Handle(ctx context.Context, req canreceiveproblems.Request) (canreceiveproblems.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", ctx, req)
	ret0, _ := ret[0].(canreceiveproblems.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Handle indicates an expected call of Handle.
func (mr *MockcanReceiveProblemsUseCaseMockRecorder) Handle(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockcanReceiveProblemsUseCase)(nil).Handle), ctx, req)
}

// MockfreeHandsUseCase is a mock of freeHandsUseCase interface.
type MockfreeHandsUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockfreeHandsUseCaseMockRecorder
}

// MockfreeHandsUseCaseMockRecorder is the mock recorder for MockfreeHandsUseCase.
type MockfreeHandsUseCaseMockRecorder struct {
	mock *MockfreeHandsUseCase
}

// NewMockfreeHandsUseCase creates a new mock instance.
func NewMockfreeHandsUseCase(ctrl *gomock.Controller) *MockfreeHandsUseCase {
	mock := &MockfreeHandsUseCase{ctrl: ctrl}
	mock.recorder = &MockfreeHandsUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockfreeHandsUseCase) EXPECT() *MockfreeHandsUseCaseMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MockfreeHandsUseCase) Handle(ctx context.Context, req freehands.Request) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// Handle indicates an expected call of Handle.
func (mr *MockfreeHandsUseCaseMockRecorder) Handle(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockfreeHandsUseCase)(nil).Handle), ctx, req)
}

// MockgetChatsUseCase is a mock of getChatsUseCase interface.
type MockgetChatsUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockgetChatsUseCaseMockRecorder
}

// MockgetChatsUseCaseMockRecorder is the mock recorder for MockgetChatsUseCase.
type MockgetChatsUseCaseMockRecorder struct {
	mock *MockgetChatsUseCase
}

// NewMockgetChatsUseCase creates a new mock instance.
func NewMockgetChatsUseCase(ctrl *gomock.Controller) *MockgetChatsUseCase {
	mock := &MockgetChatsUseCase{ctrl: ctrl}
	mock.recorder = &MockgetChatsUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockgetChatsUseCase) EXPECT() *MockgetChatsUseCaseMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MockgetChatsUseCase) Handle(ctx context.Context, req getchats.Request) (getchats.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", ctx, req)
	ret0, _ := ret[0].(getchats.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Handle indicates an expected call of Handle.
func (mr *MockgetChatsUseCaseMockRecorder) Handle(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockgetChatsUseCase)(nil).Handle), ctx, req)
}

// MockgetChatHistoryUseCase is a mock of getChatHistoryUseCase interface.
type MockgetChatHistoryUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockgetChatHistoryUseCaseMockRecorder
}

// MockgetChatHistoryUseCaseMockRecorder is the mock recorder for MockgetChatHistoryUseCase.
type MockgetChatHistoryUseCaseMockRecorder struct {
	mock *MockgetChatHistoryUseCase
}

// NewMockgetChatHistoryUseCase creates a new mock instance.
func NewMockgetChatHistoryUseCase(ctrl *gomock.Controller) *MockgetChatHistoryUseCase {
	mock := &MockgetChatHistoryUseCase{ctrl: ctrl}
	mock.recorder = &MockgetChatHistoryUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockgetChatHistoryUseCase) EXPECT() *MockgetChatHistoryUseCaseMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MockgetChatHistoryUseCase) Handle(ctx context.Context, req getchathistory.Request) (getchathistory.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", ctx, req)
	ret0, _ := ret[0].(getchathistory.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Handle indicates an expected call of Handle.
func (mr *MockgetChatHistoryUseCaseMockRecorder) Handle(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockgetChatHistoryUseCase)(nil).Handle), ctx, req)
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

// MockcloseChatUseCase is a mock of closeChatUseCase interface.
type MockcloseChatUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockcloseChatUseCaseMockRecorder
}

// MockcloseChatUseCaseMockRecorder is the mock recorder for MockcloseChatUseCase.
type MockcloseChatUseCaseMockRecorder struct {
	mock *MockcloseChatUseCase
}

// NewMockcloseChatUseCase creates a new mock instance.
func NewMockcloseChatUseCase(ctrl *gomock.Controller) *MockcloseChatUseCase {
	mock := &MockcloseChatUseCase{ctrl: ctrl}
	mock.recorder = &MockcloseChatUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockcloseChatUseCase) EXPECT() *MockcloseChatUseCaseMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MockcloseChatUseCase) Handle(ctx context.Context, req closechat.Request) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// Handle indicates an expected call of Handle.
func (mr *MockcloseChatUseCaseMockRecorder) Handle(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockcloseChatUseCase)(nil).Handle), ctx, req)
}

// MockloginUseCase is a mock of loginUseCase interface.
type MockloginUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockloginUseCaseMockRecorder
}

// MockloginUseCaseMockRecorder is the mock recorder for MockloginUseCase.
type MockloginUseCaseMockRecorder struct {
	mock *MockloginUseCase
}

// NewMockloginUseCase creates a new mock instance.
func NewMockloginUseCase(ctrl *gomock.Controller) *MockloginUseCase {
	mock := &MockloginUseCase{ctrl: ctrl}
	mock.recorder = &MockloginUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockloginUseCase) EXPECT() *MockloginUseCaseMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MockloginUseCase) Handle(ctx context.Context, req managerlogin.Request) (managerlogin.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", ctx, req)
	ret0, _ := ret[0].(managerlogin.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Handle indicates an expected call of Handle.
func (mr *MockloginUseCaseMockRecorder) Handle(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockloginUseCase)(nil).Handle), ctx, req)
}
