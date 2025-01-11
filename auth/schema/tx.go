package schema

import (
	"time"

	"cendit.io/garage/primer/enum"
	"cendit.io/garage/schema"

	"github.com/uptrace/bun"
)

type Transaction struct {
	bun.BaseModel `bun:"table:transactions" rsf:"false"`

	ID         string `bun:"id,pk" json:"id"`
	SavingsID  string `bun:"savings_id" json:"savings_id"`
	ProductID  string `bun:"product_id" json:"product_id"`
	CustomerID int    `bun:"customer_id" json:"customer_id"`
	// initiator
	UserID int `bun:"user_id" json:"user_id"`

	// CONFIGURATION

	Type      enum.TransactionType `bun:"type" json:"type"`
	Reference string               `bun:"reference" json:"reference"`
	Currency  enum.Currency        `bun:"currency" json:"currency"`
	// the payment method used
	Gateway  enum.PaymentGateway `bun:"gateway" json:"gateway"`
	Method   enum.PaymentMethod  `bun:"method" json:"method"`
	Paid     bool                `bun:"paid" json:"paid"`
	PaidAt   time.Time           `bun:"paid_at" json:"paid_at"`
	Failed   bool                `bun:"failed" json:"failed"`
	FailedAt time.Time           `bun:"failed_at" json:"failed_at"`
	// this will most not be used (actually, don't use it)
	Cancelled   bool      `bun:"cancelled" json:"cancelled"`
	CancelledAt time.Time `bun:"cancelled_at" json:"cancelled_at"`
	// SHA256 of the transaction.Invoice (prevents tampering & duplication)
	Checksum string `bun:"checksum" json:"checksum"`

	// available only if made via a gateway
	PaymentLink string `bun:"payment_link" json:"payment_link"`

	// configuration before the transaction
	InitialConfig map[string]interface{} `bun:"initial_config,type:jsonb" json:"initial_config" rsfr:"false"`
	// configuration after the transaction (mostly from payment gateway)
	FinalConfig map[string]interface{} `bun:"final_config,type:jsonb" json:"final_config" rsfr:"false"`

	History []schema.History `bun:"history,type:jsonb" json:"history" rsfr:"false"`

	// items that make up the total amount
	Invoice []Item `bun:"invoice,type:jsonb" json:"invoice" rsfr:"false"`

	// transactions amount only
	SavingsAmount float64 `bun:"savings_amount" json:"savings_amount"`

	// the total amount to be paid (including all possible fees)
	Amount float64 `bun:"amount" json:"amount"`

	// transactions account balance before and after tx
	PreBalance float64 `bun:"pre_balance" json:"pre_balance"`
	// transactions account balance after tx
	PostBalance float64 `bun:"post_balance" json:"post_balance"`

	Remark string `bun:"remark" json:"remark"`

	schema.Datum
}

type Transactions []Transaction

type Item struct {
	// the key of the item
	Key string `json:"key"`

	// the name of the item
	Name string `json:"name"`

	// the amount of the item
	Amount float64 `json:"amount"`

	// the quantity of the item
	Quantity int `json:"quantity"`

	// the metadata for the item (can also be a json string)
	Metadata string `json:"metadata"`
}
