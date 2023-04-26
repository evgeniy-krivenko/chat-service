package chatsrepo

import (
	"context"
	"fmt"

	"github.com/evgeniy-krivenko/chat-service/internal/store"
	storechat "github.com/evgeniy-krivenko/chat-service/internal/store/chat"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

func (r *Repo) CreateIfNotExists(ctx context.Context, userID types.UserID) (types.ChatID, error) {
	existedChatID, err := r.db.Chat(ctx).
		Query().
		Where(storechat.ClientIDEQ(userID)).
		FirstID(ctx)

	if store.IsNotFound(err) {
		newChat, err := r.db.Chat(ctx).
			Create().
			SetClientID(userID).
			Save(ctx)
		if err != nil {
			return types.ChatIDNil, fmt.Errorf("create new chat: %v", err)
		}

		return newChat.ID, nil
	}
	if err != nil {
		return types.ChatIDNil, fmt.Errorf("search existed chat: %v", err)
	}

	return existedChatID, nil
}
