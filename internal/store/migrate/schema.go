// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// ChatsColumns holds the columns for the "chats" table.
	ChatsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "client_id", Type: field.TypeUUID, Unique: true},
		{Name: "created_at", Type: field.TypeTime},
	}
	// ChatsTable holds the schema information for the "chats" table.
	ChatsTable = &schema.Table{
		Name:       "chats",
		Columns:    ChatsColumns,
		PrimaryKey: []*schema.Column{ChatsColumns[0]},
	}
	// MessagesColumns holds the columns for the "messages" table.
	MessagesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "author_id", Type: field.TypeUUID, Nullable: true},
		{Name: "is_visible_for_client", Type: field.TypeBool, Nullable: true},
		{Name: "is_visible_for_manager", Type: field.TypeBool, Nullable: true},
		{Name: "body", Type: field.TypeString, Size: 2000},
		{Name: "checked_at", Type: field.TypeTime, Nullable: true},
		{Name: "is_blocked", Type: field.TypeBool, Nullable: true},
		{Name: "is_service", Type: field.TypeBool, Nullable: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "initial_request_id", Type: field.TypeUUID, Unique: true, Nullable: true},
		{Name: "chat_id", Type: field.TypeUUID},
		{Name: "problem_id", Type: field.TypeUUID, Nullable: true},
	}
	// MessagesTable holds the schema information for the "messages" table.
	MessagesTable = &schema.Table{
		Name:       "messages",
		Columns:    MessagesColumns,
		PrimaryKey: []*schema.Column{MessagesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "messages_chats_messages",
				Columns:    []*schema.Column{MessagesColumns[10]},
				RefColumns: []*schema.Column{ChatsColumns[0]},
				OnDelete:   schema.NoAction,
			},
			{
				Symbol:     "messages_problems_messages",
				Columns:    []*schema.Column{MessagesColumns[11]},
				RefColumns: []*schema.Column{ProblemsColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "message_created_at",
				Unique:  false,
				Columns: []*schema.Column{MessagesColumns[8]},
			},
		},
	}
	// ProblemsColumns holds the columns for the "problems" table.
	ProblemsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "manager_id", Type: field.TypeUUID, Nullable: true},
		{Name: "resolved_at", Type: field.TypeTime, Nullable: true},
		{Name: "created_at", Type: field.TypeTime},
		{Name: "chat_id", Type: field.TypeUUID},
	}
	// ProblemsTable holds the schema information for the "problems" table.
	ProblemsTable = &schema.Table{
		Name:       "problems",
		Columns:    ProblemsColumns,
		PrimaryKey: []*schema.Column{ProblemsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "problems_chats_problems",
				Columns:    []*schema.Column{ProblemsColumns[4]},
				RefColumns: []*schema.Column{ChatsColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "problem_chat_id_resolved_at",
				Unique:  true,
				Columns: []*schema.Column{ProblemsColumns[4], ProblemsColumns[2]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		ChatsTable,
		MessagesTable,
		ProblemsTable,
	}
)

func init() {
	MessagesTable.ForeignKeys[0].RefTable = ChatsTable
	MessagesTable.ForeignKeys[1].RefTable = ProblemsTable
	ProblemsTable.ForeignKeys[0].RefTable = ChatsTable
}
