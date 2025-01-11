package schema

import (
	"github.com/uptrace/bun"
)

// schematic representation of the security setting

// IMPORTANT: when you add a new field to this struct, you must also add it to the SQL query in the FByKeyVal and FByMap methods
type Template struct {
	bun.BaseModel `bun:"table:templates"`

	Datum

	ID     string `bun:"id" json:"id"`
	UserID string `bun:"user_id" json:"user_id"`

	Subject string `bun:"subject" json:"subject"`
	Body    string `bun:"body" json:"body"`
	Key     string `bun:"key" json:"key"`
}
