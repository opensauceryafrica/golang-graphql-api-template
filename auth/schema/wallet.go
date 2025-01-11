package schema

import (
	"cendit.io/garage/primer/enum"
	"cendit.io/garage/schema"

	"github.com/uptrace/bun"
)

type Wallet struct {
	bun.BaseModel `bun:"table:addresses" rsf:"false"`

	ID       string          `bun:"id" json:"id"`
	UserID   string          `bun:"user_id" json:"user_id"`
	Type     enum.WalletType `bun:"type" json:"type"`
	Name     string          `bun:"name" json:"name"`
	Currency enum.Currency   `bun:"currency" json:"currency"`
	Balance  int             `bun:"balance" json:"balance"`
	Address  []WalletAddress `bun:"address" json:"address"`

	schema.Datum
}

type WalletAddress struct {
	Address string             `bun:"address" json:"address"`
	Network enum.CryptoNetwork `bun:"network" json:"network"`
}

type Wallets []Wallet
