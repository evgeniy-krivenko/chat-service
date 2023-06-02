// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

// Package closechatmocks is a generated GoMock package.
package closechatmocks

import (
	context "context"
	reflect "reflect"
	time "time"

	messagesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/messages"
	problemsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/problems"
	types "github.com/evgeniy-krivenko/chat-service/internal/types"
	gomock "github.com/golang/mock/gomock"
)

// MockproblemsRepository is a mock of problemsRepository interface.
type MockproblemsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockproblemsRepositoryMockRecorder
}

// MockproblemsRepositoryMockRecorder is the mock recorder for MockproblemsRepository.
type MockproblemsRepositoryMockRecorder struct {
	mock *MockproblemsRepository
}

// NewMockproblemsRepository creates a new mock instance.
func NewMockproblemsRepository(ctrl *gomock.Controller) *MockproblemsRepository {
	mock := &MockproblemsRepository{ctrl: ctrl}
	mock.recorder = &MockproblemsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockproblemsRepository) EXPECT() *MockproblemsRepositoryMockRecorder {
	return m.recorder
}

// GetAssignedProblem mocks base method.
func (m *MockproblemsRepository) GetAssignedProblem(ctx context.Context, managerID types.UserID, chatID types.ChatID) (*problemsrepo.Problem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAssignedProblem", ctx, managerID, chatID)
	ret0, _ := ret[0].(*problemsrepo.Problem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAssignedProblem indicates an expected call of GetAssignedProblem.
func (mr *MockproblemsRepositoryMockRecorder) GetAssignedProblem(ctx, managerID, chatID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAssignedProblem", reflect.TypeOf((*MockproblemsRepository)(nil).GetAssignedProblem), ctx, managerID, chatID)
}

// Resolve mocks base method.
func (m *MockproblemsRepository) Resolve(ctx context.Context, managerID types.UserID, chatID types.ChatID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Resolve", ctx, managerID, chatID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Resolve indicates an expected call of Resolve.
func (mr *MockproblemsRepositoryMockRecorder) Resolve(ctx, managerID, chatID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Resolve", reflect.TypeOf((*MockproblemsRepository)(nil).Resolve), ctx, managerID, chatID)
}

// MockmessageRepository is a mock of messageRepository interface.
type MockmessageRepository struct {
	ctrl     *gomock.Controller
	recorder *MockmessageRepositoryMockRecorder
}

// MockmessageRepositoryMockRecorder is the mock recorder for MockmessageRepository.
type MockmessageRepositoryMockRecorder struct {
	mock *MockmessageRepository
}

// NewMockmessageRepository creates a new mock instance.
func NewMockmessageRepository(ctrl *gomock.Controller) *MockmessageRepository {
	mock := &MockmessageRepository{ctrl: ctrl}
	mock.recorder = &MockmessageRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockmessageRepository) EXPECT() *MockmessageRepositoryMockRecorder {
	return m.recorder
}

// CreateClientService mocks base method.
func (m *MockmessageRepository) CreateClientService(ctx context.Context, problemID types.ProblemID, chatID types.ChatID, msgBody string) (*messagesrepo.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateClientService", ctx, problemID, chatID, msgBody)
	ret0, _ := ret[0].(*messagesrepo.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateClientService indicates an expected call of CreateClientService.
func (mr *MockmessageRepositoryMockRecorder) CreateClientService(ctx, problemID, chatID, msgBody interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateClientService", reflect.TypeOf((*MockmessageRepository)(nil).CreateClientService), ctx, problemID, chatID, msgBody)
}

// MockoutboxService is a mock of outboxService interface.
type MockoutboxService struct {
	ctrl     *gomock.Controller
	recorder *MockoutboxServiceMockRecorder
}

// MockoutboxServiceMockRecorder is the mock recorder for MockoutboxService.
type MockoutboxServiceMockRecorder struct {
	mock *MockoutboxService
}

// NewMockoutboxService creates a new mock instance.
func NewMockoutboxService(ctrl *gomock.Controller) *MockoutboxService {
	mock := &MockoutboxService{ctrl: ctrl}
	mock.recorder = &MockoutboxServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockoutboxService) EXPECT() *MockoutboxServiceMockRecorder {
	return m.recorder
}

// Put mocks base method.
func (m *MockoutboxService) Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", ctx, name, payload, availableAt)
	ret0, _ := ret[0].(types.JobID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Put indicates an expected call of Put.
func (mr *MockoutboxServiceMockRecorder) Put(ctx, name, payload, availableAt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockoutboxService)(nil).Put), ctx, name, payload, availableAt)
}

// Mocktransactor is a mock of transactor interface.
type Mocktransactor struct {
	ctrl     *gomock.Controller
	recorder *MocktransactorMockRecorder
}

// MocktransactorMockRecorder is the mock recorder for Mocktransactor.
type MocktransactorMockRecorder struct {
	mock *Mocktransactor
}

// NewMocktransactor creates a new mock instance.
func NewMocktransactor(ctrl *gomock.Controller) *Mocktransactor {
	mock := &Mocktransactor{ctrl: ctrl}
	mock.recorder = &MocktransactorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mocktransactor) EXPECT() *MocktransactorMockRecorder {
	return m.recorder
}

// RunInTx mocks base method.
func (m *Mocktransactor) RunInTx(ctx context.Context, f func(context.Context) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunInTx", ctx, f)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunInTx indicates an expected call of RunInTx.
func (mr *MocktransactorMockRecorder) RunInTx(ctx, f interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunInTx", reflect.TypeOf((*Mocktransactor)(nil).RunInTx), ctx, f)
}
