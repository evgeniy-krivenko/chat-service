package chatsrepo

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"

	storechat "github.com/evgeniy-krivenko/chat-service/internal/store/chat"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
	"github.com/evgeniy-krivenko/chat-service/pkg/pointer"
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
	query := `
	select 
		"chats"."id",
		"chats"."client_id",
		"chats"."created_at",
		"profiles"."first_name",
		"profiles"."last_name"
	from "chats"
		left join "profiles" on "profiles"."id" = "chats"."client_id"
		left join "problems" p on "chats"."id" = "p"."chat_id"
	where p.manager_id = $1 and p.resolved_at is null;`

	rows, err := r.db.Job(ctx).QueryContext(ctx, query, managerID)
	if err != nil {
		return nil, fmt.Errorf("query manager chats: %v", err)
	}
	defer rows.Close()

	var chats []Chat
	for rows.Next() {
		var (
			c         Chat
			firstName sql.NullString
			lastName  sql.NullString
		)
		if err := rows.Scan(&c.ID, &c.ClientID, &c.CreatedAt, &firstName, &lastName); err != nil {
			return nil, fmt.Errorf("scan row while query manager chat: %v", err)
		}

		c.FirstName = firstName.String
		c.LastName = lastName.String

		chats = append(chats, c)
	}

	return chats, nil
}
