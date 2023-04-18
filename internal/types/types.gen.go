// Code generated by cmd/gen-types; DO NOT EDIT.

package types

import (
	"fmt"

	"github.com/google/uuid"
)


type ChatID struct {
	uuid.UUID
}

var ChatIDNil = ChatID{uuid.Nil}

func NewChatID() ChatID {
	return ChatID{
		UUID: uuid.New(),
	}
}

func (r ChatID) Validate() error {
	if r.UUID == uuid.Nil {
		return fmt.Errorf("validate error")
	}
	return nil
}

func (r ChatID) Matches(x any) bool {
	_, ok := x.(ChatID)
	return ok
}

func (r ChatID) IsZero() bool {
	return r.UUID == uuid.Nil
}

type MessageID struct {
	uuid.UUID
}

var MessageIDNil = MessageID{uuid.Nil}

func NewMessageID() MessageID {
	return MessageID{
		UUID: uuid.New(),
	}
}

func (r MessageID) Validate() error {
	if r.UUID == uuid.Nil {
		return fmt.Errorf("validate error")
	}
	return nil
}

func (r MessageID) Matches(x any) bool {
	_, ok := x.(MessageID)
	return ok
}

func (r MessageID) IsZero() bool {
	return r.UUID == uuid.Nil
}

type ProblemID struct {
	uuid.UUID
}

var ProblemIDNil = ProblemID{uuid.Nil}

func NewProblemID() ProblemID {
	return ProblemID{
		UUID: uuid.New(),
	}
}

func (r ProblemID) Validate() error {
	if r.UUID == uuid.Nil {
		return fmt.Errorf("validate error")
	}
	return nil
}

func (r ProblemID) Matches(x any) bool {
	_, ok := x.(ProblemID)
	return ok
}

func (r ProblemID) IsZero() bool {
	return r.UUID == uuid.Nil
}

type UserID struct {
	uuid.UUID
}

var UserIDNil = UserID{uuid.Nil}

func NewUserID() UserID {
	return UserID{
		UUID: uuid.New(),
	}
}

func (r UserID) Validate() error {
	if r.UUID == uuid.Nil {
		return fmt.Errorf("validate error")
	}
	return nil
}

func (r UserID) Matches(x any) bool {
	_, ok := x.(UserID)
	return ok
}

func (r UserID) IsZero() bool {
	return r.UUID == uuid.Nil
}
type Types interface {
	ChatID | MessageID | ProblemID | UserID
}

func Parse[T Types](s string) (T, error) {
	var t T
	u, err := uuid.Parse(s)
	if err != nil {
		return t, err
	}
	switch any(t).(type) {
	case ChatID:
		return T(ChatID{u}), nil
	case MessageID:
		return T(MessageID{u}), nil
	case ProblemID:
		return T(ProblemID{u}), nil
	case UserID:
		return T(UserID{u}), nil
	default:
		return t, fmt.Errorf("wrong type")
	}
}

func MustParse[T Types](s string) T {
	u := uuid.MustParse(s)

	var t T
	switch any(t).(type) {
	case ChatID:
		return T(ChatID{u})
	case MessageID:
		return T(MessageID{u})
	case ProblemID:
		return T(ProblemID{u})
	case UserID:
		return T(UserID{u})
	default:
		panic("wrong type")
	}
}
