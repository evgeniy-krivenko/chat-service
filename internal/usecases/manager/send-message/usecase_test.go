package sendmessage_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	messagesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/messages"
	problemsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/problems"
	sendmanagermessagejob "github.com/evgeniy-krivenko/chat-service/internal/services/outbox/jobs/send-manager-message"
	"github.com/evgeniy-krivenko/chat-service/internal/testingh"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
	sendmessage "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/send-message"
	sendmessagemocks "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/send-message/mocks"
)

const msgBody = "go is awesome"

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl        *gomock.Controller
	msgRepo     *sendmessagemocks.MockmessagesRepository
	problemRepo *sendmessagemocks.MockproblemsRepository
	txtor       *sendmessagemocks.Mocktransactor
	outBoxSvc   *sendmessagemocks.MockoutboxService
	uCase       sendmessage.UseCase
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.msgRepo = sendmessagemocks.NewMockmessagesRepository(s.ctrl)
	s.outBoxSvc = sendmessagemocks.NewMockoutboxService(s.ctrl)
	s.problemRepo = sendmessagemocks.NewMockproblemsRepository(s.ctrl)
	s.txtor = sendmessagemocks.NewMocktransactor(s.ctrl)

	var err error
	s.uCase, err = sendmessage.New(sendmessage.NewOptions(s.msgRepo, s.problemRepo, s.outBoxSvc, s.txtor))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *UseCaseSuite) TestInvalidRequestError() {
	// Arrange.
	req := sendmessage.Request{}

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, sendmessage.ErrInvalidRequest)
	s.Empty(resp)
}

func (s *UseCaseSuite) TestGetProblemNotFoundError() {
	// Arrange.
	req := sendmessage.Request{
		ID:          types.NewRequestID(),
		ManagerID:   types.NewUserID(),
		ChatID:      types.NewChatID(),
		MessageBody: msgBody,
	}

	s.problemRepo.EXPECT().GetAssignedProblem(gomock.Any(), req.ManagerID, req.ChatID).
		Return(nil, problemsrepo.ErrProblemNotFound)

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, sendmessage.ErrProblemNotFound)
	s.Empty(resp)
}

func (s *UseCaseSuite) TestGetProblemUnknownError() {
	// Arrange.
	req := sendmessage.Request{
		ID:          types.NewRequestID(),
		ManagerID:   types.NewUserID(),
		ChatID:      types.NewChatID(),
		MessageBody: msgBody,
	}

	s.problemRepo.EXPECT().GetAssignedProblem(gomock.Any(), req.ManagerID, req.ChatID).
		Return(nil, errors.New("unexpected"))

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.Empty(resp)
}

func (s *UseCaseSuite) TestCreateMessageError() {
	// Arrange.
	problemID := types.NewProblemID()

	req := sendmessage.Request{
		ID:          types.NewRequestID(),
		ManagerID:   types.NewUserID(),
		ChatID:      types.NewChatID(),
		MessageBody: msgBody,
	}

	s.problemRepo.EXPECT().GetAssignedProblem(gomock.Any(), req.ManagerID, req.ChatID).
		Return(&problemsrepo.Problem{ID: problemID}, nil)

	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})
	s.msgRepo.EXPECT().
		CreateFullVisible(gomock.Any(), req.ID, problemID, req.ChatID, req.ManagerID, req.MessageBody).
		Return(nil, errors.New("unexpected"))

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.Empty(resp)
}

func (s *UseCaseSuite) TestPubJobError() {
	// Arrange.
	problemID := types.NewProblemID()
	msgID := types.NewMessageID()
	createdAt := time.Now()

	req := sendmessage.Request{
		ID:          types.NewRequestID(),
		ManagerID:   types.NewUserID(),
		ChatID:      types.NewChatID(),
		MessageBody: msgBody,
	}

	s.problemRepo.EXPECT().GetAssignedProblem(gomock.Any(), req.ManagerID, req.ChatID).
		Return(&problemsrepo.Problem{ID: problemID}, nil)
	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})
	s.msgRepo.EXPECT().
		CreateFullVisible(gomock.Any(), req.ID, problemID, req.ChatID, req.ManagerID, req.MessageBody).
		Return(&messagesrepo.Message{
			ID:        msgID,
			CreatedAt: createdAt,
		}, nil)
	s.outBoxSvc.EXPECT().Put(gomock.Any(), sendmanagermessagejob.Name, gomock.Any(), gomock.Any()).
		Return(types.JobIDNil, errors.New("unexpected"))

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.Empty(resp)
}

func (s *UseCaseSuite) TestTransactionError() {
	// Arrange.
	problemID := types.NewProblemID()
	msgID := types.NewMessageID()
	createdAt := time.Now()

	req := sendmessage.Request{
		ID:          types.NewRequestID(),
		ManagerID:   types.NewUserID(),
		ChatID:      types.NewChatID(),
		MessageBody: msgBody,
	}

	s.problemRepo.EXPECT().GetAssignedProblem(gomock.Any(), req.ManagerID, req.ChatID).
		Return(&problemsrepo.Problem{ID: problemID}, nil)
	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			_ = f(ctx)
			return sql.ErrTxDone
		})
	s.msgRepo.EXPECT().
		CreateFullVisible(gomock.Any(), req.ID, problemID, req.ChatID, req.ManagerID, req.MessageBody).
		Return(&messagesrepo.Message{
			ID:        msgID,
			CreatedAt: createdAt,
		}, nil)
	s.outBoxSvc.EXPECT().Put(gomock.Any(), sendmanagermessagejob.Name, gomock.Any(), gomock.Any()).
		Return(types.NewJobID(), nil)

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.Empty(resp.CreatedAt)
	s.Empty(resp.MessageID)
}

func (s *UseCaseSuite) TestSuccess() {
	// Arrange.
	problemID := types.NewProblemID()
	msgID := types.NewMessageID()
	createdAt := time.Now()

	req := sendmessage.Request{
		ID:          types.NewRequestID(),
		ManagerID:   types.NewUserID(),
		ChatID:      types.NewChatID(),
		MessageBody: msgBody,
	}

	s.problemRepo.EXPECT().GetAssignedProblem(gomock.Any(), req.ManagerID, req.ChatID).
		Return(&problemsrepo.Problem{ID: problemID}, nil)
	s.txtor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, f func(ctx context.Context) error) error {
			return f(ctx)
		})
	s.msgRepo.EXPECT().
		CreateFullVisible(gomock.Any(), req.ID, problemID, req.ChatID, req.ManagerID, req.MessageBody).
		Return(&messagesrepo.Message{
			ID:        msgID,
			CreatedAt: createdAt,
		}, nil)
	s.outBoxSvc.EXPECT().Put(gomock.Any(), sendmanagermessagejob.Name, gomock.Any(), gomock.Any()).
		Return(types.NewJobID(), nil)

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().NoError(err)
	s.Equal(msgID, resp.MessageID)
	s.Equal(createdAt, resp.CreatedAt)
}
