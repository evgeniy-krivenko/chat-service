package messagesrepo

import (
	"context"
	"errors"
	"fmt"
	storeprofile "github.com/evgeniy-krivenko/chat-service/internal/store/profile"
	"time"

	"github.com/evgeniy-krivenko/chat-service/internal/store"
	storechat "github.com/evgeniy-krivenko/chat-service/internal/store/chat"
	storemessage "github.com/evgeniy-krivenko/chat-service/internal/store/message"
	storeproblem "github.com/evgeniy-krivenko/chat-service/internal/store/problem"
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
		Unique(false).
		WithProfile(func(query *store.ProfileQuery) {
			query.Where(storeprofile.IDNEQ(clientID))
		}).
		Where(storemessage.IsVisibleForClient(true)).
		Where(storemessage.HasChatWith(storechat.ClientID(clientID))).
		Order(store.Desc(storemessage.FieldCreatedAt))

	return r.getChatMessages(ctx, query, pageSize, cursor)
}

// GetProblemMessages returns Nth page of messages in the chat for manager side (specific problem).
func (r *Repo) GetProblemMessages(
	ctx context.Context,
	problemID types.ProblemID,
	pageSize int,
	cursor *Cursor,
) ([]Message, *Cursor, error) {
	query := r.db.Message(ctx).Query().
		Unique(false).
		Where(storemessage.IsVisibleForManager(true)).
		Where(storemessage.HasProblemWith(
			storeproblem.ID(problemID),
			storeproblem.ResolvedAtIsNil(),
		)).
		Order(store.Desc(storemessage.FieldCreatedAt))

	return r.getChatMessages(ctx, query, pageSize, cursor)
}

func (r *Repo) getChatMessages(
	ctx context.Context,
	query *store.MessageQuery,
	pageSize int,
	cursor *Cursor,
) ([]Message, *Cursor, error) {
	// for query all messages
	lastCreatedAt := time.Now().AddDate(100, 0, 0)

	if cursor != nil {
		if !r.isValidCursor(cursor) {
			return nil, nil, ErrInvalidCursor
		}
		pageSize, lastCreatedAt = cursor.PageSize, cursor.LastCreatedAt
	} else if !r.isValidPageSize(pageSize) {
		return nil, nil, ErrInvalidPageSize
	}

	msgs, err := query.
		Where(storemessage.CreatedAtLT(lastCreatedAt)).
		Limit(pageSize + 1).
		All(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("select messages: %v", err)
	}

	result := utils.Apply(msgs, adaptStoreMessage)

	if len(result) <= pageSize {
		return result, nil, nil
	}

	result = result[:len(result)-1]

	return result, &Cursor{
		LastCreatedAt: result[len(result)-1].CreatedAt,
		PageSize:      pageSize,
	}, nil
}

func (r *Repo) isValidPageSize(pageSize int) bool {
	return pageSize >= 10 && pageSize <= 100
}

func (r *Repo) isValidCursor(cursor *Cursor) bool {
	return r.isValidPageSize(cursor.PageSize) && !cursor.LastCreatedAt.IsZero()
}
