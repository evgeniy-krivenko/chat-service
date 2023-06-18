package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

type Profile struct {
	ent.Schema
}

func (Profile) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.UserID{}).Immutable(),
		field.String("first_name").Optional(),
		field.String("last_name").Optional(),
		field.Time("updated_at"),
		field.Time("created_at").Immutable().Default(time.Now),
	}
}

func (Profile) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("messages", Message.Type),
	}
}
