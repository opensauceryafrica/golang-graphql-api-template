package repository

import (
	"blacheapi/dal"
	"blacheapi/primer/enum"
	"blacheapi/primer/function"
	"blacheapi/primer/typing"
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

// schematic representation of the configurations for a transactions account transaction
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
	PaidAt   bun.NullTime        `bun:"paid_at" json:"paid_at"`
	Failed   bool                `bun:"failed" json:"failed"`
	FailedAt bun.NullTime        `bun:"failed_at" json:"failed_at"`
	// this will most not be used (actually, don't use it)
	Cancelled   bool         `bun:"cancelled" json:"cancelled"`
	CancelledAt bun.NullTime `bun:"cancelled_at" json:"cancelled_at"`
	// SHA256 of the transaction.Invoice (prevents tampering & duplication)
	Checksum string `bun:"checksum" json:"checksum"`

	// available only if made via a gateway
	PaymentLink string `bun:"payment_link" json:"payment_link"`

	// configuration before the transaction
	InitialConfig map[string]interface{} `bun:"initial_config,type:jsonb" json:"initial_config" rsfr:"false"`
	// configuration after the transaction (mostly from payment gateway)
	FinalConfig map[string]interface{} `bun:"final_config,type:jsonb" json:"final_config" rsfr:"false"`

	History []History `bun:"history,type:jsonb" json:"history" rsfr:"false"`

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

	CreatedAt bun.NullTime `bun:"created_at" json:"created_at" rsfr:"false"`
	UpdatedAt bun.NullTime `bun:"updated_at" json:"updated_at" rsfr:"false"`
}

type Transactions []Transaction

// schematic representation of an item in a transaction't invoice
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

// Scan implements the Scanner interface.
func (i *Item) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, i)
	case string:
		return json.Unmarshal([]byte(v), i)
	case nil:
		return nil
	}
	return nil
}

// Value implements the driver Valuer interface.
func (i Item) Value() (driver.Value, error) {
	b, err := json.Marshal(i)
	return string(b), err
}

// FlutterwaveTx to be received via the webhook mutation
type FlutterwaveTx struct {
	Status string `json:"status"`
	Data   struct {
		Amount        float64 `json:"amount"`
		Status        string  `json:"status"`
		ID            int     `json:"id"`
		TxRef         string  `json:"tx_ref"`
		Currency      string  `json:"currency"`
		AmountSettled float64 `json:"amount_settled"`
		ChargedAmount float64 `json:"charged_amount"`
		Card          *struct {
			First6digits string `json:"first_6digits"`
			Last4digits  string `json:"last_4digits"`
			Issuer       string `json:"issuer"`
			Country      string `json:"country"`
			Type         string `json:"type"`
			Expiry       string `json:"expiry"`
			Token        string `json:"token"`
		} `json:"card"`
		Customer struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Phone     string `json:"phone_number"`
			Email     string `json:"email"`
			CreatedAt string `json:"created_at"`
		} `json:"customer"`
	} `json:"data"`
}

// PaystackTx to be received via the webhook mutation
type PaystackTx struct {
	Status bool `json:"status"`
	Data   struct {
		Reference       string  `json:"reference"`
		Amount          float64 `json:"amount"`
		RequestedAmount float64 `json:"requested_amount"`
		Status          string  `json:"status"`
		Fees            float64 `json:"fees"`
		Currency        string  `json:"currency"`
		Customer        struct {
			Email     string `json:"email"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Phone     string `json:"phone"`
		} `json:"customer"`
		Authorization *struct {
			AuthorizationCode string `json:"authorization_code"`
			Last4             string `json:"last4"`
			ExpMonth          string `json:"exp_month"`
			ExpYear           string `json:"exp_year"`
			Bin               string `json:"bin"`
			Bank              string `json:"bank"`
			Signature         string `json:"signature"`
			Brand             string `json:"brand"`
			Reuseable         bool   `json:"reuseable"`
			CountryCode       string `json:"country_code"`
			CardType          string `json:"card_type"`
		} `json:"authorization"`
	} `json:"data"`
}

// WalletTx to be received via the webhook mutation
type WalletTx struct {
	Status bool `json:"status"`
	Data   struct {
		Amount   float64 `json:"amount"`
		Status   string  `json:"status"`
		ID       int     `json:"id"`
		TxRef    string  `json:"tx_ref"`
		Currency string  `json:"currency"`
		Customer struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Phone     string `json:"phone_number"`
			Email     string `json:"email"`
			CreatedAt string `json:"created_at"`
		} `json:"customer"`
	} `json:"data"`
}

// GeneratePaymentObject generates a payment object compatible with paystack or flutterwave
func (t *Transaction) GeneratePaymentObject(initiator User, paymentOptions interface{}, publicKey, redirectURL string, customer Customer) map[string]interface{} {
	return map[string]interface{}{
		"public_key":      publicKey,
		"tx_ref":          t.Reference,
		"amount":          t.Amount,
		"currency":        t.Currency,
		"payment_options": paymentOptions,
		"redirect_url":    redirectURL,
		"customer": map[string]interface{}{
			"email":        initiator.Email,
			"name":         initiator.Firstname + " " + initiator.Lastname,
			"phone_number": initiator.Phone,
			"id":           initiator.ID,
		},
		"customizations": map[string]interface{}{
			"title":       customer.CompanyName,
			"description": "Savings account funding",
			"logo":        customer.Logo,
		},
	}
}

/* Fields returns the struct fields as a slice of interface{} values */
func (t *Transaction) Fields() []interface{} {
	return function.ReturnStructFields(t)
}

/*
Date loads the created_at and updated_at fields of the transaction if not already present, otherwise, it loads the updated_at field only.

If the "pessimistic" parameter is set to true, it loads both fields regardless
*/
func (t *Transaction) Date(pessimistic ...bool) {
	if len(pessimistic) > 0 && !pessimistic[0] {
		if t.CreatedAt.IsZero() {
			t.CreatedAt = schema.NullTime{Time: time.Now()}
			t.UpdatedAt = schema.NullTime{Time: time.Now()}
			return
		}
		t.UpdatedAt = schema.NullTime{Time: time.Now()}
		return
	}
	t.CreatedAt = schema.NullTime{Time: time.Now()}
	t.UpdatedAt = schema.NullTime{Time: time.Now()}
}

/*
Exist checks if a transaction with the given key/value pair exists in the database

# If returns true if it exists, false otherwise

It returns an error if any
*/
func (t *Transaction) Exist(key string, value interface{}) (bool, error) {
	var exists interface{}
	err := dal.Dal.BlacheDB.NewRaw(`SELECT EXISTS(SELECT 1 FROM transactions WHERE `+key+` = ?)`, value).Scan(context.Background(), &exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists.(bool), nil
}

/*
Create inserts a new transaction into the database

It returns an error if any
*/
func (t *Transaction) Create() error {
	if _, err := dal.Dal.BlacheDB.NewRaw(`INSERT INTO transactions (id, savings_id, product_id, customer_id, user_id, reference, currency, gateway, method, paid, paid_at, cancelled, cancelled_at, failed, failed_at, checksum, payment_link, initial_config, final_config, history, invoice, amount, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, t.ID, t.SavingsID, t.ProductID, t.CustomerID, t.UserID, t.Reference, t.Currency, t.Gateway, t.Method, t.Paid, t.PaidAt, t.Cancelled, t.CancelledAt, t.Failed, t.FailedAt, t.Checksum, t.PaymentLink, t.InitialConfig, t.FinalConfig, t.History, t.Invoice, t.Amount, t.CreatedAt, t.UpdatedAt).Exec(context.Background()); err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

/*
FByKeyVal finds and returns a transaction matching the key/value pair

By default, only the id and product_id fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded

It returns an error if any
*/
func (t *Transaction) FByKeyVal(key string, val interface{}, preloadandjoin ...bool) error {
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT * FROM transactions WHERE `+key+` = ?`, val).Scan(context.Background(), t)
	}
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT * FROM transactions WHERE `+key+` = ?`, val).Scan(context.Background(), t)
	}
	return dal.Dal.BlacheDB.NewRaw(`SELECT id, product_id FROM transactions WHERE transactions.`+key+` = ?`, val).Scan(context.Background(), t)
}

/*
FUByKeyVal finds and returns a transaction matching the key/value pair for the purpose of an update thereby causing the matching rows to be locked

By default, only the id and product_id fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded

It returns an error if any
*/
func (t *Transaction) FUByKeyVal(tx bun.Tx, key string, val interface{}, preloadandjoin ...bool) error {
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		return tx.NewRaw(`SELECT * FROM transactions WHERE `+key+` = ? FOR UPDATE`, val).Scan(context.Background(), t)
	}
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return tx.NewRaw(`SELECT * FROM transactions WHERE `+key+` = ? FOR UPDATE`, val).Scan(context.Background(), t)
	}
	return tx.NewRaw(`SELECT id, product_id FROM transactions WHERE id = ? FOR UPDATE`, t.ID).Scan(context.Background(), t)
}

/*
FByKeyVal finds and returns all transactions matching the key/value pair

By default, only the id and product_id fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded

It returns an error if any
*/
func (t *Transactions) FByKeyVal(key string, val interface{}, limit, offset int, sort string, preloadandjoin ...bool) error {
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT * FROM transactions LEFT JOIN transactions ON transactions.id = transactions.savings_id WHERE transactions.`+key+` = ? ORDER BY transactions.updated_at `+sort+` LIMIT ? OFFSET ?`, val, limit, offset).Scan(context.Background(), t)
	}
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT * FROM transactions WHERE `+key+` = ? ORDER BY transactions.updated_at `+sort+` LIMIT ? OFFSET ?`, val, limit, offset).Scan(context.Background(), t)
	}
	return dal.Dal.BlacheDB.NewRaw(`SELECT id, product_id FROM transactions WHERE `+key+` = ? ORDER BY transactions.updated_at `+sort+` LIMIT ? OFFSET ?`, val, limit, offset).Scan(context.Background(), t)
}

/*
FByMap finds and returns a transactions matching the key/value pairs provided in the map

By default, only the id and product_id fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded

It returns an error if any
*/
func (t *Transaction) FByMap(m typing.SQLMaps, preloadandjoin ...bool) error {
	query, args := dal.MapsToWQuery(m)
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT * FROM transactions WHERE `+query, args...).Scan(context.Background(), t)
	}
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT * FROM transactions WHERE `+query, args...).Scan(context.Background(), t)
	}
	return dal.Dal.BlacheDB.NewRaw(`SELECT id, product_id FROM transactions WHERE `+query, args...).Scan(context.Background(), t)
}

/*
FUByMap finds and returns a transactions matching the key/value pairs provided in the map for the purpose of an update thereby causing the matching rows to be locked

By default, only the id and product_id fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded

It returns an error if any
*/
func (t *Transaction) FUByMap(tx bun.Tx, m typing.SQLMaps, preloadandjoin ...bool) error {
	query, args := dal.MapsToWQuery(m)
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		return tx.NewRaw(`SELECT * FROM transactions WHERE `+query+` FOR UPDATE`, args...).Scan(context.Background(), t)
	}
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return tx.NewRaw(`SELECT * FROM transactions WHERE `+query+` FOR UPDATE`, args...).Scan(context.Background(), t)
	}
	return tx.NewRaw(`SELECT id, product_id FROM transactions WHERE `+query+` FOR UPDATE`, args...).Scan(context.Background(), t)
}

/*
FByMap finds and returns all transactions matching the key/value pairs provided in the map

By default, only the id and product_id fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded

It returns an error if any
*/
func (t *Transactions) FByMap(m typing.SQLMaps, limit, offset int, sort string, preloadandjoin ...bool) error {
	query, args := dal.MapsToWQuery(m)
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		rows, err := dal.Dal.BlacheDB.QueryContext(context.Background(), `SELECT * FROM transactions LEFT JOIN transactions ON transactions.id = transactions.savings_id WHERE `+query+` ORDER BY transactions.updated_at `+sort+` LIMIT ? OFFSET ?`, append(args, limit, offset)...)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var tx Transaction
			if err := rows.Scan(tx.Fields()...); err != nil {
				return err
			}
			*t = append(*t, tx)
		}
	}
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT * FROM transactions WHERE `+query+` ORDER BY transactions.updated_at `+sort+` LIMIT ? OFFSET ?`, append(args, limit, offset)...).Scan(context.Background(), t)
	}
	return dal.Dal.BlacheDB.NewRaw(`SELECT id, product_id FROM transactions WHERE `+query+` ORDER BY transactions.updated_at `+sort+` LIMIT ? OFFSET ?`, append(args, limit, offset)...).Scan(context.Background(), t)
}

/*
UByMap updates a transaction matching the key/value pairs provided in the map

It returns an error if any
*/
func (t *Transaction) UByMap(m typing.SQLMaps) error {
	query, args := dal.MapsToSQuery(m)
	if strings.Contains(query, string(enum.RETURNING)) {
		return dal.Dal.BlacheDB.NewRaw(`UPDATE transactions `+query, args...).Scan(context.Background(), t)
	}
	_, err := dal.Dal.BlacheDB.NewRaw(`UPDATE transactions `+query, args...).Exec(context.Background())
	return err
}

/*
UByMapTx updates a transaction matching the key/value pairs provided in the map using the provided transaction

It returns an error if any
*/
func (t *Transaction) UByMapTx(tx bun.Tx, m typing.SQLMaps) error {
	query, args := dal.MapsToSQuery(m)
	if strings.Contains(query, string(enum.RETURNING)) {
		return tx.NewRaw(`UPDATE transactions `+query, args...).Scan(context.Background(), t)
	}
	_, err := tx.NewRaw(`UPDATE transactions `+query, args...).Exec(context.Background())
	return err
}

/*
CByMap finds and counts all transactions matching the key/value pairs provided in the map

It returns an error if any
*/
func (t *Transaction) CByMap(m typing.SQLMaps) (int, error) {
	var count int
	query, args := dal.MapsToWQuery(m)
	err := dal.Dal.BlacheDB.NewRaw(`SELECT count(*) FROM transactions WHERE `+query, args...).Scan(context.Background(), &count)
	return count, err
}

/*
CByMap finds and counts all transactions matching the key/value pairs provided in the map

It returns an error if any
*/
func (t *Transactions) CByMap(m typing.SQLMaps) (int, error) {
	var count int
	query, args := dal.MapsToWQuery(m)
	err := dal.Dal.BlacheDB.NewRaw(`SELECT count(*) FROM transactions WHERE `+query, args...).Scan(context.Background(), &count)
	return count, err
}
