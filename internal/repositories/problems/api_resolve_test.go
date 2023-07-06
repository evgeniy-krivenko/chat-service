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

type ProblemsResolveAPISuite struct {
	testingh.DBSuite
	repo *problemsrepo.Repo
}

func TestProblemsResolveAPISuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ProblemsResolveAPISuite{DBSuite: testingh.NewDBSuite("TestProblemsResolveAPISuite")})
}

func (s *ProblemsResolveAPISuite) SetupSuite() {
	s.DBSuite.SetupSuite()

	var err error
	s.repo, err = problemsrepo.New(problemsrepo.NewOptions(s.Database))
	s.Require().NoError(err)
}

func (s *ProblemsResolveAPISuite) TestNoProblemError() {
	s.Run("no problems", func() {
		// Arrange.
		managerID := types.NewUserID()
		chatID := types.NewChatID()
		reqID := types.NewRequestID()

		// Action.
		err := s.repo.Resolve(s.Ctx, reqID, managerID, chatID)

		// Assert.
		s.Require().Error(err)
		s.ErrorIs(err, problemsrepo.ErrProblemNotFound)
	})

	s.Run("problem assigned to other manager", func() {
		// Arrange.
		managerID := types.NewUserID()
		otherManagerID := types.NewUserID()
		reqID := types.NewRequestID()

		chatID, _ := s.createChatWithProblemAssignedTo(managerID)
		_, otherProblemID := s.createChatWithProblemAssignedTo(otherManagerID)

		// Action.
		err := s.repo.Resolve(s.Ctx, reqID, otherManagerID, chatID)

		// Assert.
		s.Require().Error(err)
		s.ErrorIs(err, problemsrepo.ErrProblemNotFound)

		// other problem not resolved
		dbProblem := s.Database.Problem(s.Ctx).GetX(s.Ctx, otherProblemID)
		s.Empty(dbProblem.ResolvedAt)
	})

	s.Run("problem assigned to other chat", func() {
		// Arrange.
		managerID := types.NewUserID()
		reqID := types.NewRequestID()

		_, _ = s.createChatWithProblemAssignedTo(managerID)

		// Action.
		err := s.repo.Resolve(s.Ctx, reqID, managerID, types.NewChatID())

		// Assert.
		s.Require().Error(err)
		s.ErrorIs(err, problemsrepo.ErrProblemNotFound)
	})

	s.Run("problem was resolved", func() {
		// Arrange.
		managerID := types.NewUserID()
		reqID := types.NewRequestID()

		chatID, problemID := s.createChatWithProblemAssignedTo(managerID)

		_, err := s.Database.Problem(s.Ctx).UpdateOneID(problemID).SetResolvedAt(time.Now()).Save(s.Ctx)
		s.Require().NoError(err)

		// Action.
		err = s.repo.Resolve(s.Ctx, reqID, managerID, chatID)

		// Assert.
		s.Require().Error(err)
		s.ErrorIs(err, problemsrepo.ErrProblemNotFound)
	})
}

func (s *ProblemsResolveAPISuite) TestSuccessResolveProblem() {
	// Arrange.
	managerID := types.NewUserID()
	reqID := types.NewRequestID()

	chatID, problemID := s.createChatWithProblemAssignedTo(managerID)

	// Action.
	err := s.repo.Resolve(s.Ctx, reqID, managerID, chatID)

	// Assert.
	s.Require().NoError(err)

	dbProblem := s.Database.Problem(s.Ctx).GetX(s.Ctx, problemID)

	s.NotEmpty(dbProblem.ResolvedAt)
	s.Equal(managerID, dbProblem.ManagerID)
	s.Equal(chatID, dbProblem.ChatID)
	s.Equal(reqID, dbProblem.ResolveRequestID)
}

func (s *ProblemsResolveAPISuite) createChatWithProblemAssignedTo(managerID types.UserID) (types.ChatID, types.ProblemID) {
	s.T().Helper()

	// 1 chat can have only 1 open problem.

	chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(types.NewUserID()).Save(s.Ctx)
	s.Require().NoError(err)

	p, err := s.Database.Problem(s.Ctx).Create().SetChatID(chat.ID).SetManagerID(managerID).Save(s.Ctx)
	s.Require().NoError(err)

	return chat.ID, p.ID
}
