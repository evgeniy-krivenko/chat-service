//go:build integration

package chatsrepo_test

import (
	"testing"

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
