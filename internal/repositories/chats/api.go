package chatsrepo

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"

	storechat "github.com/evgeniy-krivenko/chat-service/internal/store/chat"
	storeproblem "github.com/evgeniy-krivenko/chat-service/internal/store/problem"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
	"github.com/evgeniy-krivenko/chat-service/pkg/pointer"
	"github.com/evgeniy-krivenko/chat-service/pkg/utils"
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

func (r *Repo) GetChatByID(ctx context.Context, chatID types.ChatID) (*Chat, error) {
	c, err := r.db.Chat(ctx).Get(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("find chat %v: %v", chatID, err)
	}

	return pointer.Ptr(adaptChat(c)), nil
}

func (r *Repo) GetManagerChatsWithProblems(ctx context.Context, managerID types.UserID) ([]Chat, error) {
	chats, err := r.db.Chat(ctx).Query().
		Where(storechat.HasProblemsWith(
			storeproblem.ManagerID(managerID),
			storeproblem.ResolvedAtIsNil(),
		)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("get chats for manager %v: %v", managerID, err)
	}

	return utils.Apply(chats, adaptChat), nil
}
