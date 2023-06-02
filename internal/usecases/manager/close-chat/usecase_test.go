package closechat_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	problemsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/problems"
	closechatjob "github.com/evgeniy-krivenko/chat-service/internal/services/outbox/jobs/close-chat"
	"github.com/evgeniy-krivenko/chat-service/internal/testingh"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
	closechat "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/close-chat"
	closechatmocks "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/close-chat/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl         *gomock.Controller
	problemsRepo *closechatmocks.MockproblemsRepository
	outboxSvc    *closechatmocks.MockoutboxService
	txtor        *closechatmocks.Mocktransactor
	uCase        closechat.UseCase
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.problemsRepo = closechatmocks.NewMockproblemsRepository(s.ctrl)
	s.outboxSvc = closechatmocks.NewMockoutboxService(s.ctrl)
	s.txtor = closechatmocks.NewMocktransactor(s.ctrl)

	var err error
	s.uCase, err = closechat.New(closechat.NewOptions(s.problemsRepo, s.outboxSvc, s.txtor))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()
	s.ContextSuite.TearDownTest()
}

func (s *UseCaseSuite) TestInvalidRequestError() {
	// Arrange.
	req := closechat.Request{}

	// Action.
	err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, closechat.ErrInvalidRequest)
}

func (s *UseCaseSuite) TestUseCase_ProblemNotFound() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()
	managerID := types.NewUserID()

	req := closechat.Request{
		ID:        reqID,
		ManagerID: managerID,
		ChatID:    chatID,
	}

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})
	s.problemsRepo.EXPECT().Resolve(gomock.Any(), managerID, chatID).
		Return(problemsrepo.ErrProblemNotFound)

	// Action
	err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, closechat.ErrProblemNotFound)
}

func (s *UseCaseSuite) TestUseCase_OutboxError() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()
	managerID := types.NewUserID()

	req := closechat.Request{
		ID:        reqID,
		ManagerID: managerID,
		ChatID:    chatID,
	}

	payload, err := closechatjob.MarshalPayload(req.ID, req.ManagerID, req.ChatID)
	s.Require().NoError(err)

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})
	s.problemsRepo.EXPECT().Resolve(gomock.Any(), managerID, chatID).
		Return(nil)
	s.outboxSvc.EXPECT().Put(gomock.Any(), closechatjob.Name, payload, gomock.Any()).
		Return(types.JobIDNil, errors.New("unexpected"))

	// Action
	err = s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestUseCaseTransactionError() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()
	managerID := types.NewUserID()

	req := closechat.Request{
		ID:        reqID,
		ManagerID: managerID,
		ChatID:    chatID,
	}

	payload, err := closechatjob.MarshalPayload(req.ID, req.ManagerID, req.ChatID)
	s.Require().NoError(err)

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			_ = f(ctx)
			return sql.ErrTxDone
		})
	s.problemsRepo.EXPECT().Resolve(gomock.Any(), managerID, chatID).
		Return(nil)
	s.outboxSvc.EXPECT().Put(gomock.Any(), closechatjob.Name, payload, gomock.Any()).
		Return(types.JobIDNil, nil)

	// Action
	err = s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
}

func (s *UseCaseSuite) TestUseCaseSuccess() {
	// Arrange.
	reqID := types.NewRequestID()
	chatID := types.NewChatID()
	managerID := types.NewUserID()

	req := closechat.Request{
		ID:        reqID,
		ManagerID: managerID,
		ChatID:    chatID,
	}

	payload, err := closechatjob.MarshalPayload(req.ID, req.ManagerID, req.ChatID)
	s.Require().NoError(err)

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})
	s.problemsRepo.EXPECT().Resolve(gomock.Any(), managerID, chatID).
		Return(nil)
	s.outboxSvc.EXPECT().Put(gomock.Any(), closechatjob.Name, payload, gomock.Any()).
		Return(types.JobIDNil, nil)

	// Action
	err = s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().NoError(err)
}
