package schema

import (
	"cendit.io/garage/schema"
	"github.com/uptrace/bun"
)

type Card struct {
	bun.BaseModel `bun:"table:cards" rsf:"false"`

	ID string `bun:"id,pk" json:"id"`

	// owner
	UserID int `bun:"user_id" json:"user_id"`

	// CONFIGURATION

	First6digits string `bun:"first_6digits" json:"first_6digits"`
	Last4digits  string `bun:"last_4digits" json:"last_4digits"`
	Issuer       string `bun:"issuer" json:"issuer"`
	Country      string `bun:"country" json:"country"`
	Type         string `bun:"type" json:"type"`
	Expiry       string `bun:"expiry" json:"expiry"`
	Checksum     string `bun:"checksum" json:"checksum"`

	// AUTHORIZATION
	Token string `bun:"token" json:"token"`

	schema.Datum
}

type Cards []Card
