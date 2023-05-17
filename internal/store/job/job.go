// Code generated by ent, DO NOT EDIT.

package job

import (
	"time"

	"github.com/evgeniy-krivenko/chat-service/internal/types"
)

const (
	// Label holds the string label denoting the job type in the database.
	Label = "job"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldPayload holds the string denoting the payload field in the database.
	FieldPayload = "payload"
	// FieldAttempts holds the string denoting the attempts field in the database.
	FieldAttempts = "attempts"
	// FieldAvailableAt holds the string denoting the available_at field in the database.
	FieldAvailableAt = "available_at"
	// FieldReservedUntil holds the string denoting the reserved_until field in the database.
	FieldReservedUntil = "reserved_until"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// Table holds the table name of the job in the database.
	Table = "jobs"
)

// Columns holds all SQL columns for job fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldPayload,
	FieldAttempts,
	FieldAvailableAt,
	FieldReservedUntil,
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
	// DefaultAttempts holds the default value on creation for the "attempts" field.
	DefaultAttempts int
	// AttemptsValidator is a validator for the "attempts" field. It is called by the builders before save.
	AttemptsValidator func(int) error
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() types.JobID
)