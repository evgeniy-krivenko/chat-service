package messagesrepo

import (
	"context"
	"fmt"
	"time"

	storemessage "github.com/evgeniy-krivenko/chat-service/internal/store/message"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

func (r *Repo) MarkAsVisibleForManager(ctx context.Context, msgID types.MessageID) error {
	err := r.db.Message(ctx).Update().
		SetIsVisibleForManager(true).
		SetCheckedAt(time.Now()).
		Where(storemessage.ID(msgID)).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("mark msg as visible for manager: %v", err)
	}
	return nil
}

func (r *Repo) BlockMessage(ctx context.Context, msgID types.MessageID) error {
	err := r.db.Message(ctx).Update().
		SetIsBlocked(true).
		SetCheckedAt(time.Now()).
		Where(storemessage.ID(msgID)).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("block msg for manager: %v", err)
	}
	return nil
}
