package getchathistory_test

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/evgeniy-krivenko/chat-service/internal/cursor"
	messagesrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/messages"
	problemsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/problems"
	"github.com/evgeniy-krivenko/chat-service/internal/testingh"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
	getchathistory "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/get-chat-history"
	getchathistorymocks "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/get-chat-history/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl         *gomock.Controller
	msgRepo      *getchathistorymocks.MockmessagesRepository
	problemsRepo *getchathistorymocks.MockproblemsRepository
	uCase        getchathistory.UseCase
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.msgRepo = getchathistorymocks.NewMockmessagesRepository(s.ctrl)
	s.problemsRepo = getchathistorymocks.NewMockproblemsRepository(s.ctrl)

	var err error
	s.uCase, err = getchathistory.New(getchathistory.NewOptions(s.msgRepo, s.problemsRepo))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *UseCaseSuite) TestRequestValidationError() {
	// Arrange.
	req := getchathistory.Request{}

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, getchathistory.ErrInvalidRequest)
	s.Empty(resp.Messages)
	s.Empty(resp.NextCursor)
}

func (s *UseCaseSuite) TestCursorDecodingError() {
	// Arrange.
	req := getchathistory.Request{
		ID:        types.NewRequestID(),
		ManagerID: types.NewUserID(),
		ChatID:    types.NewChatID(),
		Cursor:    "eyJwYWdlX3NpemUiOjEwMA==", // {"page_size":100
	}

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, getchathistory.ErrInvalidCursor)
	s.Empty(resp.Messages)
	s.Empty(resp.NextCursor)
}

func (s *UseCaseSuite) TestOpenProblemForChat_AnyError() {
	// Arrange.
	errExpected := errors.New("any error")
	req := getchathistory.Request{
		ID:        types.NewRequestID(),
		ManagerID: types.NewUserID(),
		ChatID:    types.NewChatID(),
		PageSize:  10,
	}

	s.problemsRepo.EXPECT().GetOpenProblemForChat(gomock.Any(), req.ChatID, req.ManagerID).
		Return(nil, errExpected)

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.Empty(resp.Messages)
	s.Empty(resp.NextCursor)
}

func (s *UseCaseSuite) TestGetProblemMessages_InvalidCursor() {
	// Arrange.
	clientID := types.NewUserID()
	problemID := types.NewProblemID()

	c := messagesrepo.Cursor{PageSize: -1, LastCreatedAt: time.Now()}
	cursorWithNegativePageSize, err := cursor.Encode(c)
	s.Require().NoError(err)

	req := getchathistory.Request{
		ID:        types.NewRequestID(),
		ManagerID: clientID,
		ChatID:    types.NewChatID(),
		PageSize:  0,
		Cursor:    cursorWithNegativePageSize,
	}

	p := problemsrepo.Problem{
		ID: problemID,
	}

	s.problemsRepo.EXPECT().GetOpenProblemForChat(gomock.Any(), req.ChatID, req.ManagerID).
		Return(&p, nil)
	s.msgRepo.EXPECT().GetProblemMessages(s.Ctx, p.ID, 0, messagesrepo.NewCursorMatcher(c)).
		Return(nil, nil, messagesrepo.ErrInvalidCursor)

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.ErrorIs(err, getchathistory.ErrInvalidCursor)
	s.Empty(resp.Messages)
	s.Empty(resp.NextCursor)
}

func (s *UseCaseSuite) TestGetProblemMessages_Error() {
	// Arrange.
	managerID := types.NewUserID()
	problemID := types.NewProblemID()
	errExpected := errors.New("any error")

	req := getchathistory.Request{
		ID:        types.NewRequestID(),
		ManagerID: managerID,
		ChatID:    types.NewChatID(),
		PageSize:  10,
	}

	p := problemsrepo.Problem{
		ID: problemID,
	}

	s.problemsRepo.EXPECT().GetOpenProblemForChat(gomock.Any(), req.ChatID, req.ManagerID).
		Return(&p, nil)
	s.msgRepo.EXPECT().GetProblemMessages(s.Ctx, p.ID, 10, (*messagesrepo.Cursor)(nil)).
		Return(nil, nil, errExpected)

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.Empty(resp.Messages)
	s.Empty(resp.NextCursor)
}

func (s *UseCaseSuite) TestGetProblemMessages_Success_SinglePage() {
	// Arrange.
	const messagesCount = 10
	const pageSize = messagesCount + 1

	chatID := types.NewChatID()
	clientID := types.NewUserID()
	expectedMsgs := s.createMessages(messagesCount, clientID, chatID)
	managerID := types.NewUserID()
	problemID := types.NewProblemID()

	req := getchathistory.Request{
		ID:        types.NewRequestID(),
		ManagerID: managerID,
		ChatID:    chatID,
		PageSize:  pageSize,
	}

	p := problemsrepo.Problem{
		ID: problemID,
	}

	s.problemsRepo.EXPECT().GetOpenProblemForChat(gomock.Any(), req.ChatID, req.ManagerID).
		Return(&p, nil)
	s.msgRepo.EXPECT().GetProblemMessages(s.Ctx, p.ID, pageSize, (*messagesrepo.Cursor)(nil)).
		Return(expectedMsgs, nil, nil)

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)
	s.Require().NoError(err)

	// Assert.
	s.Empty(resp.NextCursor)

	s.Require().Len(resp.Messages, messagesCount)
	for i := 0; i < messagesCount; i++ {
		s.Equal(expectedMsgs[i].ID, resp.Messages[i].ID)
		s.Equal(expectedMsgs[i].AuthorID, resp.Messages[i].AuthorID)
		s.Equal(expectedMsgs[i].Body, resp.Messages[i].Body)
		s.Equal(expectedMsgs[i].CreatedAt.Unix(), resp.Messages[i].CreatedAt.Unix())
	}
}

func (s *UseCaseSuite) TestGetProblemMessages_Success_FirstPage() {
	// Arrange.
	const messagesCount = 10
	const pageSize = messagesCount + 1

	chatID := types.NewChatID()
	clientID := types.NewUserID()
	expectedMsgs := s.createMessages(messagesCount, clientID, chatID)
	managerID := types.NewUserID()
	problemID := types.NewProblemID()
	lastMsg := expectedMsgs[len(expectedMsgs)-1]

	req := getchathistory.Request{
		ID:        types.NewRequestID(),
		ManagerID: managerID,
		ChatID:    chatID,
		PageSize:  pageSize,
	}

	p := problemsrepo.Problem{
		ID: problemID,
	}

	s.problemsRepo.EXPECT().GetOpenProblemForChat(gomock.Any(), req.ChatID, req.ManagerID).
		Return(&p, nil)

	nextCursor := &messagesrepo.Cursor{PageSize: pageSize, LastCreatedAt: lastMsg.CreatedAt}
	s.msgRepo.EXPECT().GetProblemMessages(s.Ctx, p.ID, pageSize, (*messagesrepo.Cursor)(nil)).
		Return(expectedMsgs, nextCursor, nil)

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)
	s.Require().NoError(err)

	// Assert.
	s.NotEmpty(resp.NextCursor)
	s.Require().Len(resp.Messages, messagesCount)
}

func (s *UseCaseSuite) TestGetProblemMessages_Success_LastPage() {
	// Arrange.
	const messagesCount = 10
	const pageSize = messagesCount + 1

	chatID := types.NewChatID()
	clientID := types.NewUserID()
	expectedMsgs := s.createMessages(messagesCount, clientID, chatID)
	managerID := types.NewUserID()
	problemID := types.NewProblemID()

	c := messagesrepo.Cursor{PageSize: pageSize, LastCreatedAt: time.Now()}
	cursorStr, err := cursor.Encode(c)
	s.Require().NoError(err)

	req := getchathistory.Request{
		ID:        types.NewRequestID(),
		ManagerID: managerID,
		ChatID:    chatID,
		Cursor:    cursorStr,
	}

	p := problemsrepo.Problem{
		ID: problemID,
	}

	s.problemsRepo.EXPECT().GetOpenProblemForChat(gomock.Any(), req.ChatID, req.ManagerID).
		Return(&p, nil)
	s.msgRepo.EXPECT().GetProblemMessages(s.Ctx, p.ID, 0, messagesrepo.NewCursorMatcher(c)).
		Return(expectedMsgs, nil, nil)

	// Action.
	resp, err := s.uCase.Handle(s.Ctx, req)
	s.Require().NoError(err)

	// Assert.
	s.Empty(resp.NextCursor)
	s.Require().Len(resp.Messages, messagesCount)
}

func (s *UseCaseSuite) createMessages(count int, authorID types.UserID, chatID types.ChatID) []messagesrepo.Message {
	s.T().Helper()

	result := make([]messagesrepo.Message, 0, count)
	for i := 0; i < count; i++ {
		result = append(result, messagesrepo.Message{
			ID:                  types.NewMessageID(),
			ChatID:              chatID,
			AuthorID:            authorID,
			Body:                uuid.New().String(),
			CreatedAt:           time.Now(),
			IsVisibleForClient:  true,
			IsVisibleForManager: true,
			IsBlocked:           false,
			IsService:           false,
		})
	}
	return result
}
