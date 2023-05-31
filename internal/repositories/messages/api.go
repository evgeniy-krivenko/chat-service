package messagesrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/evgeniy-krivenko/chat-service/internal/store"
	storemessage "github.com/evgeniy-krivenko/chat-service/internal/store/message"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
	"github.com/evgeniy-krivenko/chat-service/pkg/pointer"
)

var ErrMsgNotFound = errors.New("message not found")

func (r *Repo) GetMessageByRequestID(ctx context.Context, reqID types.RequestID) (*Message, error) {
	msg, err := r.db.Message(ctx).Query().Where(storemessage.InitialRequestIDEQ(reqID)).First(ctx)
	if store.IsNotFound(err) {
		return nil, ErrMsgNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("query message by requies id %v: %v", reqID, err)
	}

	return pointer.Ptr(adaptStoreMessage(msg)), nil
}

// CreateClientVisible creates a message that is visible only to the client.
func (r *Repo) CreateClientVisible(
	ctx context.Context,
	reqID types.RequestID,
	problemID types.ProblemID,
	chatID types.ChatID,
	authorID types.UserID,
	msgBody string,
) (*Message, error) {
	msg, err := r.db.Message(ctx).Create().
		SetInitialRequestID(reqID).
		SetProblemID(problemID).
		SetChatID(chatID).
		SetAuthorID(authorID).
		SetBody(msgBody).
		SetIsVisibleForManager(false).
		SetIsVisibleForClient(true).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create msg for client by request id %v: %v", reqID, err)
	}

	return pointer.Ptr(adaptStoreMessage(msg)), nil
}

func (r *Repo) GetMessageByID(ctx context.Context, msgID types.MessageID) (*Message, error) {
	msg, err := r.db.Message(ctx).
		Query().
		Where(storemessage.ID(msgID)).
		WithProblem().
		Only(ctx)
	if store.IsNotFound(err) {
		return nil, ErrMsgNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get msg by id %v: %v", msgID, err)
	}

	return pointer.Ptr(adaptStoreMessage(msg)), nil
}

func (r *Repo) CreateClientService(
	ctx context.Context,
	problemID types.ProblemID,
	chatID types.ChatID,
	msgBody string,
) (*Message, error) {
	msg, err := r.db.Message(ctx).Create().
		SetProblemID(problemID).
		SetChatID(chatID).
		SetBody(msgBody).
		SetIsVisibleForManager(false).
		SetIsVisibleForClient(true).
		SetIsService(true).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create service msg for client: %v", err)
	}

	am := adaptStoreMessage(msg)
	return &am, nil
}

func (r *Repo) CreateFullVisible(
	ctx context.Context,
	reqID types.RequestID,
	problemID types.ProblemID,
	chatID types.ChatID,
	authorID types.UserID,
	msgBody string,
) (*Message, error) {
	msg, err := r.db.Message(ctx).Create().
		SetInitialRequestID(reqID).
		SetProblemID(problemID).
		SetChatID(chatID).
		SetAuthorID(authorID).
		SetBody(msgBody).
		SetIsVisibleForManager(true).
		SetIsVisibleForClient(true).
		SetIsService(false).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create full visible msg: %v", err)
	}

	am := adaptStoreMessage(msg)
	return &am, nil
}
