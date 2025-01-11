package schema

import (
	"cendit.io/garage/schema"
	"github.com/uptrace/bun"
)

type Address struct {
	bun.BaseModel `bun:"table:addresses" rsf:"false"`

	ID          string `bun:"id,pk" json:"id"`
	State       string `bun:"state" json:"state"`
	City        string `bun:"city" json:"city"`
	HouseNumber string `bun:"house_number" json:"house_number"`
	Street      string `bun:"street" json:"street"`
	Zip         string `bun:"zip" json:"zip"`
	UserID      string `bun:"user_id" json:"user_id"`
	Verified    bool   `bun:"verified" json:"verified"`

	schema.Datum
}

type Addresses []Address
