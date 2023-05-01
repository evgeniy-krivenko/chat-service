package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

// Problem holds the schema definition for the Problem entity.
type Problem struct {
	ent.Schema
}

// Fields of the Problem.
func (Problem) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.ProblemID{}).
			Default(types.NewProblemID).
			Unique(),
		field.UUID("chat_id", types.ChatID{}),
		field.UUID("manager_id", types.UserID{}).
			Optional(),
		field.Time("resolved_at").
			Optional(),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the Problem.
func (Problem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("messages", Message.Type),
		edge.From("chat", Chat.Type).
			Ref("problems").
			Unique().
			Field("chat_id").
			Required(),
	}
}

// Indexes of the Problem.
func (Problem) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("chat_id", "resolved_at").
			Unique(),
	}
}
