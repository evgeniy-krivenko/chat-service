package messagesrepo

import (
	"time"

	"github.com/evgeniy-krivenko/chat-service/internal/store"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

type Message struct {
	ID                  types.MessageID
	ChatID              types.ChatID
	AuthorID            types.UserID
	AuthorFirstName     string
	AuthorLastName      string
	InitialRequestID    types.RequestID
	Body                string
	CreatedAt           time.Time
	IsVisibleForClient  bool
	IsVisibleForManager bool
	IsBlocked           bool
	IsService           bool
	ManagerID           types.UserID
}

func adaptStoreMessage(m *store.Message) Message {
	managerID := types.UserIDNil
	p := m.Edges.Problem
	if p != nil {
		managerID = p.ManagerID
	}

	var authorFirstName, authorLastName string
	profile := m.Edges.Profile
	if profile != nil {
		authorFirstName = profile.FirstName
		authorLastName = profile.LastName
	}

	return Message{
		ID:                  m.ID,
		ChatID:              m.ChatID,
		AuthorID:            m.AuthorID,
		AuthorFirstName:     authorFirstName,
		AuthorLastName:      authorLastName,
		InitialRequestID:    m.InitialRequestID,
		Body:                m.Body,
		CreatedAt:           m.CreatedAt,
		IsVisibleForClient:  m.IsVisibleForClient,
		IsVisibleForManager: m.IsVisibleForManager,
		IsBlocked:           m.IsBlocked,
		IsService:           m.IsService,
		ManagerID:           managerID,
	}
}
