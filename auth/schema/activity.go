package schema

import (
	"context"
	"database/sql"

	"cendit.io/garage/database"
	"cendit.io/garage/primer/enum"
	"cendit.io/garage/schema"

	"github.com/uptrace/bun"
)

type Activity struct {
	bun.BaseModel `bun:"table:activities"`

	schema.Datum

	ID            string    `bun:"id,pk" json:"id"`
	By            string    `bun:"by" json:"by"`
	Role          enum.Role `bun:"role" json:"role"`
	SavingsID     string    `bun:"savings_id" json:"savings_id"`
	TransactionID string    `bun:"transaction_id" json:"transaction_id"`

	// can be preloaded using graphql reference resolvers
	User User `json:"user" bun:"-"`

	Resolver    string      `bun:"resolver" json:"resolver"`
	Payload     string      `bun:"payload" json:"payload"`
	Description string      `bun:"description" json:"description"`
	Status      enum.Status `bun:"status" json:"status"`
	Error       string      `bun:"error" json:"error"`
}

type Activities []Activity

/*
Create inserts a new activity into the database

It returns an error if any
*/
func (a *Activity) Insert() error {
	if _, err := database.DB.NewRaw(`INSERT INTO activities (id, by, role, savings_id, transaction_id, resolver, payload, description, status, error, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, a.ID, a.By, a.Role, a.SavingsID, a.TransactionID, a.Resolver, a.Payload, a.Description, a.Status, a.Error, a.CreatedAt, a.UpdatedAt).Exec(context.Background()); err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}
