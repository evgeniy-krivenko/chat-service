//go:build integration

package problemsrepo_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	problemsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/problems"
	"github.com/evgeniy-krivenko/chat-service/internal/testingh"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

const msgBody = "go is awesome"

type ProblemsRepoScheduleAPISuite struct {
	testingh.DBSuite
	repo *problemsrepo.Repo
}

func TestProblemsRepoScheduleAPISuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ProblemsRepoScheduleAPISuite{DBSuite: testingh.NewDBSuite("TestProblemsRepoScheduleAPISuite")})
}

func (s *ProblemsRepoScheduleAPISuite) SetupSuite() {
	s.DBSuite.SetupSuite()

	var err error
	s.repo, err = problemsrepo.New(problemsrepo.NewOptions(s.Database))
	s.Require().NoError(err)
}

func (s *ProblemsRepoScheduleAPISuite) SetupSubTest() {
	s.DBSuite.SetupTest()

	s.Database.Problem(s.Ctx).Delete().ExecX(s.Ctx)
	s.Database.Message(s.Ctx).Delete().ExecX(s.Ctx)
	s.Database.Chat(s.Ctx).Delete().ExecX(s.Ctx)
}

func (s *ProblemsRepoScheduleAPISuite) Test_GetAvailableProblems() {
	s.Run("two problems with visible msgs", func() {
		// Arrange.
		_ = s.createMessage(true)
		_ = s.createMessage(true)

		// Action.
		problems, err := s.repo.GetAvailableProblems(s.Ctx)

		// Assert.
		s.Require().NoError(err)
		s.Require().Len(problems, 2)
	})

	s.Run("no messages visible for manager", func() {
		// Arrange.
		_ = s.createMessage(false)

		// Action.
		problems, err := s.repo.GetAvailableProblems(s.Ctx)

		// Assert.
		s.Require().NoError(err)
		s.Require().Empty(problems)
	})

	s.Run("get first problem if second with manager", func() {
		// Arrange.
		_ = s.createMessage(true)

		// Create one more problem with msg chat and manager
		clientID := types.NewUserID()

		s.Database.Profile(s.Ctx).Create().
			SetID(clientID).
			SetUpdatedAt(time.Now()).
			SaveX(s.Ctx)

		problemID, chatID := s.createProblemAndChatWithManager(clientID)
		_ = s.Database.Message(s.Ctx).Create().
			SetID(types.NewMessageID()).
			SetChatID(chatID).
			SetAuthorID(clientID).
			SetProblemID(problemID).
			SetBody(msgBody).
			SetIsBlocked(false).
			SetIsVisibleForClient(true).
			SetIsVisibleForManager(true).
			SetInitialRequestID(types.NewRequestID()).
			SaveX(s.Ctx)

		// Action
		problems, err := s.repo.GetAvailableProblems(s.Ctx)
		s.Require().NoError(err)
		s.Require().Len(problems, 1)
	})
}

func (s *ProblemsRepoScheduleAPISuite) Test_SetManagerForProblem() {
	s.Run("success set manager for problem", func() {
		// Arrange.
		authorID := types.NewUserID()
		managerID := types.NewUserID()

		s.Database.Profile(s.Ctx).Create().
			SetID(authorID).
			SetUpdatedAt(time.Now()).
			SaveX(s.Ctx)

		problemID, chatID := s.createProblemAndChat(authorID)
		_ = s.Database.Message(s.Ctx).Create().
			SetID(types.NewMessageID()).
			SetChatID(chatID).
			SetAuthorID(authorID).
			SetProblemID(problemID).
			SetBody(msgBody).
			SetIsBlocked(false).
			SetIsVisibleForClient(true).
			SetIsVisibleForManager(true).
			SetIsService(false).
			SetInitialRequestID(types.NewRequestID()).
			SaveX(s.Ctx)

		// Action.
		err := s.repo.SetManagerForProblem(s.Ctx, problemID, managerID)

		// Assert.
		s.Require().NoError(err)
		problem, err := s.Database.Problem(s.Ctx).Get(s.Ctx, problemID)
		s.Assert().NoError(err)
		s.Assert().Equal(managerID, problem.ManagerID)
	})

	s.Run("no found problem without manager id", func() {
		// Arrange.
		s.Database.Problem(s.Ctx).Delete().ExecX(s.Ctx)

		managerID := types.NewUserID()
		authorID := types.NewUserID()

		problemID, chatID := s.createProblemAndChatWithManager(managerID)

		s.Database.Profile(s.Ctx).Create().
			SetID(authorID).
			SetUpdatedAt(time.Now()).
			SaveX(s.Ctx)

		_ = s.Database.Message(s.Ctx).Create().
			SetID(types.NewMessageID()).
			SetChatID(chatID).
			SetAuthorID(authorID).
			SetProblemID(problemID).
			SetBody(msgBody).
			SetIsBlocked(false).
			SetIsVisibleForClient(true).
			SetIsVisibleForManager(true).
			SetIsService(false).
			SetInitialRequestID(types.NewRequestID()).
			SaveX(s.Ctx)

		// Action.
		err := s.repo.SetManagerForProblem(s.Ctx, problemID, managerID)

		// Assert.
		s.Require().ErrorIs(err, problemsrepo.ErrProblemNotFound)
	})
}

func (s *ProblemsRepoScheduleAPISuite) Test_GetProblemReqID() {
	s.Run("req id for first message in problem", func() {
		// Arrange.
		authorID := types.NewUserID()
		requestID := types.NewRequestID()

		s.Database.Profile(s.Ctx).Create().
			SetID(authorID).
			SetUpdatedAt(time.Now()).
			SaveX(s.Ctx)

		problemID, chatID := s.createProblemAndChat(authorID)
		firstMsg := s.Database.Message(s.Ctx).Create().
			SetID(types.NewMessageID()).
			SetChatID(chatID).
			SetAuthorID(authorID).
			SetProblemID(problemID).
			SetBody(msgBody).
			SetIsBlocked(false).
			SetIsVisibleForClient(true).
			SetIsVisibleForManager(true).
			SetIsService(false).
			SetInitialRequestID(requestID).
			SaveX(s.Ctx)

		// second message
		_ = s.Database.Message(s.Ctx).Create().
			SetID(types.NewMessageID()).
			SetChatID(chatID).
			SetAuthorID(authorID).
			SetProblemID(problemID).
			SetBody(msgBody).
			SetIsBlocked(false).
			SetIsVisibleForClient(true).
			SetIsVisibleForManager(true).
			SetIsService(false).
			SetInitialRequestID(types.NewRequestID()).
			SaveX(s.Ctx)

		// Action.
		expReqID, err := s.repo.GetProblemRequestID(s.Ctx, problemID)

		// Assert.
		s.Require().NoError(err)
		s.Require().Equal(firstMsg.InitialRequestID, expReqID)
	})

	s.Run("not found request id if manager not nil", func() {
		// Arrange.
		authorID := types.NewUserID()
		requestID := types.NewRequestID()

		s.Database.Profile(s.Ctx).Create().
			SetID(authorID).
			SetUpdatedAt(time.Now()).
			SaveX(s.Ctx)

		problemID, chatID := s.createProblemAndChatWithManager(authorID)
		firstMsg := s.Database.Message(s.Ctx).Create().
			SetID(types.NewMessageID()).
			SetChatID(chatID).
			SetAuthorID(authorID).
			SetProblemID(problemID).
			SetBody(msgBody).
			SetIsBlocked(false).
			SetIsVisibleForClient(true).
			SetIsVisibleForManager(true).
			SetIsService(false).
			SetInitialRequestID(requestID).
			SaveX(s.Ctx)

		// Action.
		expReqID, err := s.repo.GetProblemRequestID(s.Ctx, problemID)

		// Assert.
		s.Require().ErrorIs(err, problemsrepo.ErrReqIDNotFount)
		s.Assert().NotEqual(firstMsg.InitialRequestID, expReqID)
	})

	s.Run("not found request id if not visible message", func() {
		// Arrange.
		authorID := types.NewUserID()
		requestID := types.NewRequestID()

		s.Database.Profile(s.Ctx).Create().
			SetID(authorID).
			SetUpdatedAt(time.Now()).
			SaveX(s.Ctx)

		problemID, chatID := s.createProblemAndChat(authorID)
		firstMsg := s.Database.Message(s.Ctx).Create().
			SetID(types.NewMessageID()).
			SetChatID(chatID).
			SetAuthorID(authorID).
			SetProblemID(problemID).
			SetBody(msgBody).
			SetIsBlocked(false).
			SetIsVisibleForClient(true).
			SetIsVisibleForManager(false).
			SetIsService(false).
			SetInitialRequestID(requestID).
			SaveX(s.Ctx)

		// Action.
		expReqID, err := s.repo.GetProblemRequestID(s.Ctx, problemID)

		// Assert.
		s.Require().ErrorIs(err, problemsrepo.ErrReqIDNotFount)
		s.Assert().NotEqual(firstMsg.InitialRequestID, expReqID)
	})

	s.Run("not found request id if not message", func() {
		// Arrange.
		authorID := types.NewUserID()
		problemID, _ := s.createProblemAndChat(authorID)

		// Action.
		expReqID, err := s.repo.GetProblemRequestID(s.Ctx, problemID)

		// Assert.
		s.Require().ErrorIs(err, problemsrepo.ErrReqIDNotFount)
		s.Assert().Equal(types.RequestIDNil, expReqID)
	})
}

func (s *ProblemsRepoScheduleAPISuite) createMessage(isVisibleForManager bool) types.MessageID {
	s.T().Helper()

	authorID := types.NewUserID()
	problemID, chatID := s.createProblemAndChat(authorID)
	msgID := types.NewMessageID()

	s.Database.Profile(s.Ctx).Create().
		SetID(authorID).
		SetUpdatedAt(time.Now()).
		SaveX(s.Ctx)

	_, err := s.Database.Message(s.Ctx).Create().
		SetID(msgID).
		SetChatID(chatID).
		SetAuthorID(authorID).
		SetProblemID(problemID).
		SetBody(msgBody).
		SetIsBlocked(false).
		SetIsVisibleForClient(true).
		SetIsVisibleForManager(isVisibleForManager).
		SetIsService(false).
		SetInitialRequestID(types.NewRequestID()).
		Save(s.Ctx)
	s.Require().NoError(err)

	return msgID
}

func (s *ProblemsRepoScheduleAPISuite) createProblemAndChat(clientID types.UserID) (types.ProblemID, types.ChatID) {
	s.T().Helper()

	chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
	s.Require().NoError(err)

	problem, err := s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).Save(s.Ctx)
	s.Require().NoError(err)

	return problem.ID, chat.ID
}

func (s *ProblemsRepoScheduleAPISuite) createProblemAndChatWithManager(clientID types.UserID) (types.ProblemID, types.ChatID) {
	s.T().Helper()

	managerID := types.NewUserID()

	chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
	s.Require().NoError(err)

	problem, err := s.Database.Problem(s.Ctx).Create().
		SetChatID(chat.ID).
		SetManagerID(managerID).
		Save(s.Ctx)
	s.Require().NoError(err)

	return problem.ID, chat.ID
}
