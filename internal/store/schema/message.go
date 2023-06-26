package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

// Message holds the schema definition for the Message entity.
type Message struct {
	ent.Schema
}

// Fields of the Message.
func (Message) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.MessageID{}).
			Default(types.NewMessageID).
			Unique(),
		field.UUID("author_id", types.UserID{}).
			Optional(),
		field.UUID("chat_id", types.ChatID{}),
		field.UUID("problem_id", types.ProblemID{}).
			Optional(),
		field.Bool("is_visible_for_client").
			Optional(),
		field.Bool("is_visible_for_manager").
			Optional(),
		field.String("body").
			MaxLen(2000).
			NotEmpty(),
		field.Time("checked_at").
			Optional(),
		field.Bool("is_blocked").
			Optional(),
		field.Bool("is_service").
			Optional(),
		field.UUID("initial_request_id", types.RequestID{}).
			Optional().
			Unique(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the Message.
func (Message) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("chat", Chat.Type).
			Ref("messages").
			Unique().
			Field("chat_id").
			Required(),
		edge.From("problem", Problem.Type).
			Ref("messages").
			Unique().
			Field("problem_id"),
	}
}

// Indexes of the Message.
func (Message) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("created_at"),
	}
}
