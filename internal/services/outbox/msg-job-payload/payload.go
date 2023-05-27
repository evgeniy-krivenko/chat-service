package msgjobpayload

import (
	"encoding/json"
	"fmt"

	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

type msgPayload struct {
	MessageID types.MessageID `json:"messageId"`
}

func MarshalPayload(messageID types.MessageID) (string, error) {
	if err := messageID.Validate(); err != nil {
		return "", fmt.Errorf("invalid message id: %v", err)
	}
	payload := msgPayload{MessageID: messageID}

	data, err := json.Marshal(&payload)
	if err != nil {
		return "", fmt.Errorf("send client message job: %v", err)
	}

	return string(data), nil
}

func UnmarshalPayload(payload string) (types.MessageID, error) {
	var pl msgPayload

	err := json.Unmarshal([]byte(payload), &pl)
	if err != nil {
		return types.MessageIDNil, fmt.Errorf("send client message job: %v", err)
	}

	return pl.MessageID, nil
}
