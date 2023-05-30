// Code generated by ent, DO NOT EDIT.

package store

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/evgeniy-krivenko/chat-service/internal/store/chat"
	"github.com/evgeniy-krivenko/chat-service/internal/store/message"
	"github.com/evgeniy-krivenko/chat-service/internal/store/problem"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

// Message is the model entity for the Message schema.
type Message struct {
	config `json:"-"`
	// ID of the ent.
	ID types.MessageID `json:"id,omitempty"`
	// AuthorID holds the value of the "author_id" field.
	AuthorID types.UserID `json:"author_id,omitempty"`
	// ChatID holds the value of the "chat_id" field.
	ChatID types.ChatID `json:"chat_id,omitempty"`
	// ProblemID holds the value of the "problem_id" field.
	ProblemID types.ProblemID `json:"problem_id,omitempty"`
	// IsVisibleForClient holds the value of the "is_visible_for_client" field.
	IsVisibleForClient bool `json:"is_visible_for_client,omitempty"`
	// IsVisibleForManager holds the value of the "is_visible_for_manager" field.
	IsVisibleForManager bool `json:"is_visible_for_manager,omitempty"`
	// Body holds the value of the "body" field.
	Body string `json:"body,omitempty"`
	// CheckedAt holds the value of the "checked_at" field.
	CheckedAt time.Time `json:"checked_at,omitempty"`
	// IsBlocked holds the value of the "is_blocked" field.
	IsBlocked bool `json:"is_blocked,omitempty"`
	// IsService holds the value of the "is_service" field.
	IsService bool `json:"is_service,omitempty"`
	// InitialRequestID holds the value of the "initial_request_id" field.
	InitialRequestID types.RequestID `json:"initial_request_id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the MessageQuery when eager-loading is set.
	Edges MessageEdges `json:"edges"`
}

// MessageEdges holds the relations/edges for other nodes in the graph.
type MessageEdges struct {
	// Chat holds the value of the chat edge.
	Chat *Chat `json:"chat,omitempty"`
	// Problem holds the value of the problem edge.
	Problem *Problem `json:"problem,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// ChatOrErr returns the Chat value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e MessageEdges) ChatOrErr() (*Chat, error) {
	if e.loadedTypes[0] {
		if e.Chat == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: chat.Label}
		}
		return e.Chat, nil
	}
	return nil, &NotLoadedError{edge: "chat"}
}

// ProblemOrErr returns the Problem value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e MessageEdges) ProblemOrErr() (*Problem, error) {
	if e.loadedTypes[1] {
		if e.Problem == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: problem.Label}
		}
		return e.Problem, nil
	}
	return nil, &NotLoadedError{edge: "problem"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Message) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case message.FieldIsVisibleForClient, message.FieldIsVisibleForManager, message.FieldIsBlocked, message.FieldIsService:
			values[i] = new(sql.NullBool)
		case message.FieldBody:
			values[i] = new(sql.NullString)
		case message.FieldCheckedAt, message.FieldCreatedAt:
			values[i] = new(sql.NullTime)
		case message.FieldChatID:
			values[i] = new(types.ChatID)
		case message.FieldID:
			values[i] = new(types.MessageID)
		case message.FieldProblemID:
			values[i] = new(types.ProblemID)
		case message.FieldInitialRequestID:
			values[i] = new(types.RequestID)
		case message.FieldAuthorID:
			values[i] = new(types.UserID)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Message", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Message fields.
func (m *Message) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case message.FieldID:
			if value, ok := values[i].(*types.MessageID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				m.ID = *value
			}
		case message.FieldAuthorID:
			if value, ok := values[i].(*types.UserID); !ok {
				return fmt.Errorf("unexpected type %T for field author_id", values[i])
			} else if value != nil {
				m.AuthorID = *value
			}
		case message.FieldChatID:
			if value, ok := values[i].(*types.ChatID); !ok {
				return fmt.Errorf("unexpected type %T for field chat_id", values[i])
			} else if value != nil {
				m.ChatID = *value
			}
		case message.FieldProblemID:
			if value, ok := values[i].(*types.ProblemID); !ok {
				return fmt.Errorf("unexpected type %T for field problem_id", values[i])
			} else if value != nil {
				m.ProblemID = *value
			}
		case message.FieldIsVisibleForClient:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field is_visible_for_client", values[i])
			} else if value.Valid {
				m.IsVisibleForClient = value.Bool
			}
		case message.FieldIsVisibleForManager:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field is_visible_for_manager", values[i])
			} else if value.Valid {
				m.IsVisibleForManager = value.Bool
			}
		case message.FieldBody:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field body", values[i])
			} else if value.Valid {
				m.Body = value.String
			}
		case message.FieldCheckedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field checked_at", values[i])
			} else if value.Valid {
				m.CheckedAt = value.Time
			}
		case message.FieldIsBlocked:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field is_blocked", values[i])
			} else if value.Valid {
				m.IsBlocked = value.Bool
			}
		case message.FieldIsService:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field is_service", values[i])
			} else if value.Valid {
				m.IsService = value.Bool
			}
		case message.FieldInitialRequestID:
			if value, ok := values[i].(*types.RequestID); !ok {
				return fmt.Errorf("unexpected type %T for field initial_request_id", values[i])
			} else if value != nil {
				m.InitialRequestID = *value
			}
		case message.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				m.CreatedAt = value.Time
			}
		}
	}
	return nil
}

// QueryChat queries the "chat" edge of the Message entity.
func (m *Message) QueryChat() *ChatQuery {
	return NewMessageClient(m.config).QueryChat(m)
}

// QueryProblem queries the "problem" edge of the Message entity.
func (m *Message) QueryProblem() *ProblemQuery {
	return NewMessageClient(m.config).QueryProblem(m)
}

// Update returns a builder for updating this Message.
// Note that you need to call Message.Unwrap() before calling this method if this Message
// was returned from a transaction, and the transaction was committed or rolled back.
func (m *Message) Update() *MessageUpdateOne {
	return NewMessageClient(m.config).UpdateOne(m)
}

// Unwrap unwraps the Message entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (m *Message) Unwrap() *Message {
	_tx, ok := m.config.driver.(*txDriver)
	if !ok {
		panic("store: Message is not a transactional entity")
	}
	m.config.driver = _tx.drv
	return m
}

// String implements the fmt.Stringer.
func (m *Message) String() string {
	var builder strings.Builder
	builder.WriteString("Message(")
	builder.WriteString(fmt.Sprintf("id=%v, ", m.ID))
	builder.WriteString("author_id=")
	builder.WriteString(fmt.Sprintf("%v", m.AuthorID))
	builder.WriteString(", ")
	builder.WriteString("chat_id=")
	builder.WriteString(fmt.Sprintf("%v", m.ChatID))
	builder.WriteString(", ")
	builder.WriteString("problem_id=")
	builder.WriteString(fmt.Sprintf("%v", m.ProblemID))
	builder.WriteString(", ")
	builder.WriteString("is_visible_for_client=")
	builder.WriteString(fmt.Sprintf("%v", m.IsVisibleForClient))
	builder.WriteString(", ")
	builder.WriteString("is_visible_for_manager=")
	builder.WriteString(fmt.Sprintf("%v", m.IsVisibleForManager))
	builder.WriteString(", ")
	builder.WriteString("body=")
	builder.WriteString(m.Body)
	builder.WriteString(", ")
	builder.WriteString("checked_at=")
	builder.WriteString(m.CheckedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("is_blocked=")
	builder.WriteString(fmt.Sprintf("%v", m.IsBlocked))
	builder.WriteString(", ")
	builder.WriteString("is_service=")
	builder.WriteString(fmt.Sprintf("%v", m.IsService))
	builder.WriteString(", ")
	builder.WriteString("initial_request_id=")
	builder.WriteString(fmt.Sprintf("%v", m.InitialRequestID))
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(m.CreatedAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// Messages is a parsable slice of Message.
type Messages []*Message
