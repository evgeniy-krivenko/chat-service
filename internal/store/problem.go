// Code generated by ent, DO NOT EDIT.

package store

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/evgeniy-krivenko/chat-service/internal/store/chat"
	"github.com/evgeniy-krivenko/chat-service/internal/store/problem"
	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

// Problem is the model entity for the Problem schema.
type Problem struct {
	config `json:"-"`
	// ID of the ent.
	ID types.ProblemID `json:"id,omitempty"`
	// ChatID holds the value of the "chat_id" field.
	ChatID types.ChatID `json:"chat_id,omitempty"`
	// ManagerID holds the value of the "manager_id" field.
	ManagerID types.UserID `json:"manager_id,omitempty"`
	// ResolvedAt holds the value of the "resolved_at" field.
	ResolvedAt time.Time `json:"resolved_at,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ProblemQuery when eager-loading is set.
	Edges ProblemEdges `json:"edges"`
}

// ProblemEdges holds the relations/edges for other nodes in the graph.
type ProblemEdges struct {
	// Messages holds the value of the messages edge.
	Messages []*Message `json:"messages,omitempty"`
	// Chat holds the value of the chat edge.
	Chat *Chat `json:"chat,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// MessagesOrErr returns the Messages value or an error if the edge
// was not loaded in eager-loading.
func (e ProblemEdges) MessagesOrErr() ([]*Message, error) {
	if e.loadedTypes[0] {
		return e.Messages, nil
	}
	return nil, &NotLoadedError{edge: "messages"}
}

// ChatOrErr returns the Chat value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ProblemEdges) ChatOrErr() (*Chat, error) {
	if e.loadedTypes[1] {
		if e.Chat == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: chat.Label}
		}
		return e.Chat, nil
	}
	return nil, &NotLoadedError{edge: "chat"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Problem) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case problem.FieldResolvedAt, problem.FieldCreatedAt:
			values[i] = new(sql.NullTime)
		case problem.FieldChatID:
			values[i] = new(types.ChatID)
		case problem.FieldID:
			values[i] = new(types.ProblemID)
		case problem.FieldManagerID:
			values[i] = new(types.UserID)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Problem", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Problem fields.
func (pr *Problem) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case problem.FieldID:
			if value, ok := values[i].(*types.ProblemID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				pr.ID = *value
			}
		case problem.FieldChatID:
			if value, ok := values[i].(*types.ChatID); !ok {
				return fmt.Errorf("unexpected type %T for field chat_id", values[i])
			} else if value != nil {
				pr.ChatID = *value
			}
		case problem.FieldManagerID:
			if value, ok := values[i].(*types.UserID); !ok {
				return fmt.Errorf("unexpected type %T for field manager_id", values[i])
			} else if value != nil {
				pr.ManagerID = *value
			}
		case problem.FieldResolvedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field resolved_at", values[i])
			} else if value.Valid {
				pr.ResolvedAt = value.Time
			}
		case problem.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				pr.CreatedAt = value.Time
			}
		}
	}
	return nil
}

// QueryMessages queries the "messages" edge of the Problem entity.
func (pr *Problem) QueryMessages() *MessageQuery {
	return NewProblemClient(pr.config).QueryMessages(pr)
}

// QueryChat queries the "chat" edge of the Problem entity.
func (pr *Problem) QueryChat() *ChatQuery {
	return NewProblemClient(pr.config).QueryChat(pr)
}

// Update returns a builder for updating this Problem.
// Note that you need to call Problem.Unwrap() before calling this method if this Problem
// was returned from a transaction, and the transaction was committed or rolled back.
func (pr *Problem) Update() *ProblemUpdateOne {
	return NewProblemClient(pr.config).UpdateOne(pr)
}

// Unwrap unwraps the Problem entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (pr *Problem) Unwrap() *Problem {
	_tx, ok := pr.config.driver.(*txDriver)
	if !ok {
		panic("store: Problem is not a transactional entity")
	}
	pr.config.driver = _tx.drv
	return pr
}

// String implements the fmt.Stringer.
func (pr *Problem) String() string {
	var builder strings.Builder
	builder.WriteString("Problem(")
	builder.WriteString(fmt.Sprintf("id=%v, ", pr.ID))
	builder.WriteString("chat_id=")
	builder.WriteString(fmt.Sprintf("%v", pr.ChatID))
	builder.WriteString(", ")
	builder.WriteString("manager_id=")
	builder.WriteString(fmt.Sprintf("%v", pr.ManagerID))
	builder.WriteString(", ")
	builder.WriteString("resolved_at=")
	builder.WriteString(pr.ResolvedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(pr.CreatedAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// Problems is a parsable slice of Problem.
type Problems []*Problem
