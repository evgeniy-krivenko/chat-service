// Code generated by ent, DO NOT EDIT.

package message

import (
	"time"

	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

const (
	// Label holds the string label denoting the message type in the database.
	Label = "message"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldAuthorID holds the string denoting the author_id field in the database.
	FieldAuthorID = "author_id"
	// FieldChatID holds the string denoting the chat_id field in the database.
	FieldChatID = "chat_id"
	// FieldProblemID holds the string denoting the problem_id field in the database.
	FieldProblemID = "problem_id"
	// FieldIsVisibleForClient holds the string denoting the is_visible_for_client field in the database.
	FieldIsVisibleForClient = "is_visible_for_client"
	// FieldIsVisibleForManager holds the string denoting the is_visible_for_manager field in the database.
	FieldIsVisibleForManager = "is_visible_for_manager"
	// FieldBody holds the string denoting the body field in the database.
	FieldBody = "body"
	// FieldCheckedAt holds the string denoting the checked_at field in the database.
	FieldCheckedAt = "checked_at"
	// FieldIsBlocked holds the string denoting the is_blocked field in the database.
	FieldIsBlocked = "is_blocked"
	// FieldIsService holds the string denoting the is_service field in the database.
	FieldIsService = "is_service"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// EdgeChat holds the string denoting the chat edge name in mutations.
	EdgeChat = "chat"
	// EdgeProblem holds the string denoting the problem edge name in mutations.
	EdgeProblem = "problem"
	// Table holds the table name of the message in the database.
	Table = "messages"
	// ChatTable is the table that holds the chat relation/edge.
	ChatTable = "messages"
	// ChatInverseTable is the table name for the Chat entity.
	// It exists in this package in order to avoid circular dependency with the "chat" package.
	ChatInverseTable = "chats"
	// ChatColumn is the table column denoting the chat relation/edge.
	ChatColumn = "chat_id"
	// ProblemTable is the table that holds the problem relation/edge.
	ProblemTable = "messages"
	// ProblemInverseTable is the table name for the Problem entity.
	// It exists in this package in order to avoid circular dependency with the "problem" package.
	ProblemInverseTable = "problems"
	// ProblemColumn is the table column denoting the problem relation/edge.
	ProblemColumn = "problem_id"
)

// Columns holds all SQL columns for message fields.
var Columns = []string{
	FieldID,
	FieldAuthorID,
	FieldChatID,
	FieldProblemID,
	FieldIsVisibleForClient,
	FieldIsVisibleForManager,
	FieldBody,
	FieldCheckedAt,
	FieldIsBlocked,
	FieldIsService,
	FieldCreatedAt,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// BodyValidator is a validator for the "body" field. It is called by the builders before save.
	BodyValidator func(string) error
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() types.ChatID
)
