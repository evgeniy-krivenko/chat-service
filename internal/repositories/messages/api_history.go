package messagesrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/evgeniy-krivenko/chat-service/internal/store"
	storechat "github.com/evgeniy-krivenko/chat-service/internal/store/chat"
	storemessage "github.com/evgeniy-krivenko/chat-service/internal/store/message"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
	"github.com/evgeniy-krivenko/chat-service/pkg/utils"
)

var (
	ErrInvalidPageSize = errors.New("invalid page size")
	ErrInvalidCursor   = errors.New("invalid cursor")
)

type Cursor struct {
	LastCreatedAt time.Time
	PageSize      int
}

// GetClientChatMessages returns Nth page of messages in the chat for client side.
func (r *Repo) GetClientChatMessages(
	ctx context.Context,
	clientID types.UserID,
	pageSize int,
	cursor *Cursor,
) ([]Message, *Cursor, error) {
	query := r.db.Message(ctx).Query().
		Order(store.Desc(storemessage.FieldCreatedAt)).
		Limit(pageSize + 1)

	if cursor != nil {
		if !r.isValidCursor(cursor) {
			return nil, nil, ErrInvalidCursor
		}

		query = query.
			Where(func(s *sql.Selector) {
				t := sql.Table(storechat.Table)
				s.Join(t).On(s.C(storemessage.ChatColumn), t.C(storechat.FieldID))
				s.Where(sql.And(
					sql.IsTrue(storemessage.FieldIsVisibleForClient),
					sql.EQ(t.C(storechat.FieldClientID), clientID),
					sql.LT(s.C(storemessage.FieldCreatedAt), cursor.LastCreatedAt),
				))
			})

		return r.clientMessagesWithCursor(ctx, query, cursor.PageSize)
	}

	if !r.isValidPageSize(pageSize) {
		return nil, nil, ErrInvalidPageSize
	}

	query = query.
		Where(func(s *sql.Selector) {
			t := sql.Table(storechat.Table)
			s.Join(t).On(s.C(storemessage.ChatColumn), t.C(storechat.FieldID))
			s.Where(sql.And(
				sql.IsTrue(storemessage.FieldIsVisibleForClient),
				sql.EQ(t.C(storechat.FieldClientID), clientID),
			))
		})

	return r.clientMessagesWithCursor(ctx, query, pageSize)
}

func (r *Repo) clientMessagesWithCursor(
	ctx context.Context,
	query *store.MessageQuery,
	pageSize int,
) ([]Message, *Cursor, error) {
	messages, err := query.All(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("query messages: %v", err)
	}

	var nextCursor *Cursor
	if pageSize < len(messages) {
		lastMsg := messages[len(messages)-2]
		nextCursor = &Cursor{
			PageSize:      pageSize,
			LastCreatedAt: lastMsg.CreatedAt,
		}
		messages = messages[0 : len(messages)-1]
	}

	adaptedMessages := utils.Apply[*store.Message, Message](messages, adaptStoreMessage)
	return adaptedMessages, nextCursor, nil
}

func (r *Repo) isValidPageSize(pageSize int) bool {
	return pageSize >= 10 && pageSize <= 100
}

func (r *Repo) isValidCursor(cursor *Cursor) bool {
	return r.isValidPageSize(cursor.PageSize) && !cursor.LastCreatedAt.IsZero()
}
