package chatsrepo

import (
	"time"

	"github.com/evgeniy-krivenko/chat-service/internal/store"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

type Chat struct {
	ID        types.ChatID
	ClientID  types.UserID
	CreatedAt time.Time
}

func adaptChat(c *store.Chat) *Chat {
	return &Chat{
		ID:        c.ID,
		ClientID:  c.ClientID,
		CreatedAt: c.CreatedAt,
	}
}
