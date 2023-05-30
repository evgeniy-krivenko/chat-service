//go:build integration

package chatsrepo_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	chatsrepo "github.com/evgeniy-krivenko/chat-service/internal/repositories/chats"
	"github.com/evgeniy-krivenko/chat-service/internal/testingh"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

type ChatsRepoSuite struct {
	testingh.DBSuite
	repo *chatsrepo.Repo
}

func TestChatsRepoSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, &ChatsRepoSuite{DBSuite: testingh.NewDBSuite("TestChatsRepoSuite")})
}

func (s *ChatsRepoSuite) SetupSuite() {
	s.DBSuite.SetupSuite()

	var err error

	s.repo, err = chatsrepo.New(chatsrepo.NewOptions(s.Database))
	s.Require().NoError(err)
}

func (s *ChatsRepoSuite) Test_CreateIfNotExists() {
	s.Run("chat does not exist, should be created", func() {
		clientID := types.NewUserID()

		chatID, err := s.repo.CreateIfNotExists(s.Ctx, clientID)
		s.Require().NoError(err)
		s.NotEmpty(chatID)
	})

	s.Run("chat already exists", func() {
		clientID := types.NewUserID()

		// Create chat.
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		chatID, err := s.repo.CreateIfNotExists(s.Ctx, clientID)
		s.Require().NoError(err)
		s.Require().NotEmpty(chatID)
		s.Equal(chat.ID, chatID)
	})
}

func (s *ChatsRepoSuite) Test_GetChatByID() {
	s.Run("chat does not exists", func() {
		chatID := types.NewChatID()

		chat, err := s.repo.GetChatByID(s.Ctx, chatID)
		s.Require().Error(err)
		s.Require().Nil(chat)
	})

	s.Run("success chat exists", func() {
		clientID := types.NewUserID()

		// Create chat.
		chat, err := s.Database.Chat(s.Ctx).Create().SetClientID(clientID).Save(s.Ctx)
		s.Require().NoError(err)

		chatFromDB, err := s.repo.GetChatByID(s.Ctx, chat.ID)
		s.Require().NoError(err)
		s.Equal(chat.ID, chatFromDB.ID)
		s.Equal(chat.ClientID, chatFromDB.ClientID)
		s.True(chat.CreatedAt.Equal(chatFromDB.CreatedAt))
	})
}

func (s *ChatsRepoSuite) Test_GetManagerChatsWithProblems() {
	s.Run("one chat exists", func() {
		// Arrange.
		clientFirst := types.NewUserID()
		managerID := types.NewUserID()

		chatOne := s.Database.Chat(s.Ctx).Create().SetClientID(clientFirst).SaveX(s.Ctx)
		s.Database.Problem(s.Ctx).Create().SetChatID(chatOne.ID).SetManagerID(managerID).SaveX(s.Ctx)

		// Action.
		chats, err := s.repo.GetManagerChatsWithProblems(s.Ctx, managerID)

		// Assert.
		s.Require().NoError(err)

		s.Len(chats, 1)
		s.Equal(chatOne.ID, chats[0].ID)
		s.IsType(chatsrepo.Chat{}, chats[0])
	})

	s.Run("empty chats with resolved problems", func() {
		// Arrange.
		clientFirst := types.NewUserID()
		clientSecond := types.NewUserID()
		managerID := types.NewUserID()
		resolveTimeFirst := time.Now().Add(-2 * time.Second)
		resolveTimeSecond := time.Now().Add(-1 * time.Second)

		c := s.Database.Chat(s.Ctx).Create().SetClientID(clientFirst).SaveX(s.Ctx)
		s.Database.Problem(s.Ctx).Create().
			SetChatID(c.ID).
			SetManagerID(managerID).
			SetResolvedAt(resolveTimeFirst).
			SaveX(s.Ctx)
		c2 := s.Database.Chat(s.Ctx).Create().SetClientID(clientSecond).SaveX(s.Ctx)
		s.Database.Problem(s.Ctx).Create().
			SetChatID(c2.ID).
			SetManagerID(managerID).
			SetResolvedAt(resolveTimeSecond).
			SaveX(s.Ctx)

		// Action.
		chats, err := s.repo.GetManagerChatsWithProblems(s.Ctx, managerID)

		// Assert.
		s.Require().NoError(err)
		s.Empty(chats)
	})

	s.Run("empty chats with open problems for other manager", func() {
		// Arrange.
		clientFirst := types.NewUserID()
		clientSecond := types.NewUserID()
		managerID := types.NewUserID()
		otherManagerID := types.NewUserID()

		c := s.Database.Chat(s.Ctx).Create().SetClientID(clientFirst).SaveX(s.Ctx)
		s.Database.Problem(s.Ctx).Create().
			SetChatID(c.ID).
			SetManagerID(otherManagerID).
			SaveX(s.Ctx)
		c2 := s.Database.Chat(s.Ctx).Create().SetClientID(clientSecond).SaveX(s.Ctx)
		s.Database.Problem(s.Ctx).Create().
			SetChatID(c2.ID).
			SetManagerID(otherManagerID).
			SaveX(s.Ctx)

		// Action.
		chats, err := s.repo.GetManagerChatsWithProblems(s.Ctx, managerID)

		// Assert.
		s.Require().NoError(err)
		s.Empty(chats)
	})

	s.Run("empty chats with resolved problems for other manager", func() {
		// Arrange.
		clientFirst := types.NewUserID()
		clientSecond := types.NewUserID()
		managerID := types.NewUserID()
		otherManagerID := types.NewUserID()
		resolveTimeFirst := time.Now().Add(-2 * time.Second)
		resolveTimeSecond := time.Now().Add(-1 * time.Second)

		c := s.Database.Chat(s.Ctx).Create().SetClientID(clientFirst).SaveX(s.Ctx)
		s.Database.Problem(s.Ctx).Create().
			SetChatID(c.ID).
			SetManagerID(otherManagerID).
			SetResolvedAt(resolveTimeFirst).
			SaveX(s.Ctx)
		c2 := s.Database.Chat(s.Ctx).Create().SetClientID(clientSecond).SaveX(s.Ctx)
		s.Database.Problem(s.Ctx).Create().
			SetChatID(c2.ID).
			SetManagerID(otherManagerID).
			SetResolvedAt(resolveTimeSecond).
			SaveX(s.Ctx)

		// Action.
		chats, err := s.repo.GetManagerChatsWithProblems(s.Ctx, managerID)

		// Assert.
		s.Require().NoError(err)
		s.Empty(chats)
	})

	s.Run("empty chats without problems", func() {
		// Arrange.
		clientFirst := types.NewUserID()
		clientSecond := types.NewUserID()
		managerID := types.NewUserID()

		s.Database.Chat(s.Ctx).Create().SetClientID(clientFirst).SaveX(s.Ctx)

		s.Database.Chat(s.Ctx).Create().SetClientID(clientSecond).SaveX(s.Ctx)

		// Action.
		chats, err := s.repo.GetManagerChatsWithProblems(s.Ctx, managerID)

		// Assert.
		s.Require().NoError(err)
		s.Empty(chats)
	})
}
