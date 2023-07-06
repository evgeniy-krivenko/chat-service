package getchats_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	chatsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/chats"
	"github.com/evgeniy-krivenko/chat-service/internal/testingh"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
	getchats "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/get-chats"
	getchatsmocks "github.com/evgeniy-krivenko/chat-service/internal/usecases/manager/get-chats/mocks"
)

type UseCaseSuite struct {
	testingh.ContextSuite

	ctrl      *gomock.Controller
	chatsRepo *getchatsmocks.MockchatsRepository

	uCase getchats.UseCase
}

func TestUseCaseSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UseCaseSuite))
}

func (s *UseCaseSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.chatsRepo = getchatsmocks.NewMockchatsRepository(s.ctrl)

	var err error
	s.uCase, err = getchats.New(getchats.NewOptions(s.chatsRepo))
	s.Require().NoError(err)

	s.ContextSuite.SetupTest()
}

func (s *UseCaseSuite) TearDownTest() {
	s.ctrl.Finish()

	s.ContextSuite.TearDownTest()
}

func (s *UseCaseSuite) TestRequestValidationError() {
	// Arrange.
	req := getchats.Request{}

	// Action.
	_, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().Error(err)
	s.Require().ErrorIs(err, getchats.ErrInvalidRequest)
}

func (s *UseCaseSuite) TestGetChatsError() {
	// Arrange.
	reqID := types.NewRequestID()
	managerID := types.NewUserID()

	s.chatsRepo.EXPECT().
		GetManagerChatsWithProblems(gomock.Any(), managerID).
		Return(nil, errors.New("unexpected"))

	req := getchats.Request{
		ID:        reqID,
		ManagerID: managerID,
	}

	// Action.
	response, err := s.uCase.Handle(s.Ctx, req)

	s.Require().Error(err)
	s.Empty(response)
}

func (s *UseCaseSuite) TestGetChatsSuccess() {
	// Arrange.
	reqID := types.NewRequestID()
	managerID := types.NewUserID()

	chat := chatsrepo.Chat{
		ID:        types.NewChatID(),
		ClientID:  types.NewUserID(),
		FirstName: "Eric",
		LastName:  "Cartman",
	}

	s.chatsRepo.EXPECT().
		GetManagerChatsWithProblems(gomock.Any(), managerID).
		Return([]chatsrepo.Chat{chat}, nil)

	req := getchats.Request{
		ID:        reqID,
		ManagerID: managerID,
	}

	// Action.
	response, err := s.uCase.Handle(s.Ctx, req)

	// Assert.
	s.Require().NoError(err)
	s.Len(response.Chats, 1)
	s.Equal(chat.FirstName, response.Chats[0].FirstName)
	s.Equal(chat.LastName, response.Chats[0].LastName)
	s.IsType(getchats.Chat{}, response.Chats[0])
}
