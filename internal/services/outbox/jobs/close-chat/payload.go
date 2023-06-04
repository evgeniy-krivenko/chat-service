package closechatjob

import (
	"encoding/json"
	"fmt"

	"github.com/evgeniy-krivenko/chat-service/internal/types"
	"github.com/evgeniy-krivenko/chat-service/internal/validator"
)

type payload struct {
	RequestID   types.RequestID `json:"requestId" validate:"required"`
	ManagerID   types.UserID    `json:"managerId" validate:"required"`
	ChatID      types.ChatID    `json:"chatId" validate:"required"`
	ClientMsgID types.MessageID `json:"clientMsgId" validate:"required"`
}

func (p payload) Validate() error {
	return validator.Validator.Struct(p)
}

func MarshalPayload(
	reqID types.RequestID,
	managerID types.UserID,
	chatID types.ChatID,
	clientMsgID types.MessageID,
) (string, error) {
	p := payload{
		RequestID:   reqID,
		ManagerID:   managerID,
		ChatID:      chatID,
		ClientMsgID: clientMsgID,
	}

	if err := p.Validate(); err != nil {
		return "", fmt.Errorf("validate close chat job payload: %v", err)
	}

	data, err := json.Marshal(&p)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %v", err)
	}
	return string(data), nil
}

func unmarshalPayload(data string) (p payload, err error) {
	err = json.Unmarshal([]byte(data), &p)
	return
}
