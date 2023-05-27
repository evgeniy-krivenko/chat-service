package afcverdictsprocessor

import (
	"github.com/evgeniy-krivenko/chat-service/internal/types"
	"github.com/evgeniy-krivenko/chat-service/internal/validator"
)

type Verdict struct {
	ChatID    types.ChatID    `json:"chatId" validate:"required"`
	MessageID types.MessageID `json:"messageId" validate:"required"`
	Status    string          `json:"status" validate:"oneof=ok suspicious"`
}

func (v *Verdict) Valid() error {
	return validator.Validator.Struct(v)
}

func (v *Verdict) IsSuccess() bool {
	return v.Status == "ok"
}
