package chatsrepo

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"

	storechat "github.com/evgeniy-krivenko/chat-service/internal/store/chat"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

func (r *Repo) CreateIfNotExists(ctx context.Context, userID types.UserID) (types.ChatID, error) {
	chatID, err := r.db.Chat(ctx).Create().
		SetClientID(userID).
		OnConflict(
			sql.ConflictColumns(storechat.FieldClientID),
			sql.ResolveWith(func(set *sql.UpdateSet) {
				set.SetIgnore(storechat.FieldCreatedAt)
			}),
		).
		ID(ctx)
	if err != nil {
		return types.ChatIDNil, fmt.Errorf("create new chat: %v", err)
	}

	return chatID, nil
}
