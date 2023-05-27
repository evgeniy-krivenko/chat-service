package managerassignedtoproblemjob

import (
	"encoding/json"
	"fmt"

	"github.com/evgeniy-krivenko/chat-service/internal/types"
	"github.com/evgeniy-krivenko/chat-service/internal/validator"
)

type payload struct {
	MessageID types.MessageID `json:"messageId" validate:"required"`
	ManagerID types.UserID    `json:"managerId" validate:"required"`
	ProblemID types.ProblemID `json:"problemId" validate:"required"`
	RequestID types.RequestID `json:"requestId" validate:"required"`
}

func (p payload) Validate() error {
	return validator.Validator.Struct(p)
}

func MarshalPayload(
	msgID types.MessageID,
	managerID types.UserID,
	problemID types.ProblemID,
	reqID types.RequestID,
) (string, error) {
	p := payload{
		MessageID: msgID,
		ManagerID: managerID,
		ProblemID: problemID,
		RequestID: reqID,
	}
	if err := p.Validate(); err != nil {
		return "", fmt.Errorf("validate: %v", err)
	}

	data, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("marshal: %v", err)
	}
	return string(data), nil
}

func unmarshalPayload(data string) (p payload, err error) {
	err = json.Unmarshal([]byte(data), &p)
	return
}
