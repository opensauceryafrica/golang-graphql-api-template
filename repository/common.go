package repository

import (
	"blacheapi/dal"
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"

	"github.com/uptrace/bun"
)

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Pages int `json:"pages"`
}

type History struct {
	Act string       `json:"act"`
	By  int          `json:"by"`
	At  bun.NullTime `json:"at"`
}

// Scan implements the Scanner interface.
func (h *History) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, h)
	case string:
		return json.Unmarshal([]byte(v), h)
	case nil:
		return nil
	}
	return nil
}

// Value implements the driver Valuer interface.
func (h History) Value() (driver.Value, error) {
	b, err := json.Marshal(h)
	return string(b), err
}

// schematic representation of an attachment in a comment
type Attachment struct {
	// the name of the file
	Name string `json:"name"`

	// the url of the file
	URL string `json:"url"`
}

// Scan implements the Scanner interface.
func (a *Attachment) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, a)
	case string:
		return json.Unmarshal([]byte(v), a)
	case nil:
		return nil
	}
	return nil
}

// Value implements the driver Valuer interface.
func (a Attachment) Value() (driver.Value, error) {
	b, err := json.Marshal(a)
	return string(b), err
}

// BeginBlacheTx returns a new transaction for Blache database
func BeginBlacheTx() (bun.Tx, error) {
	return dal.Dal.BlacheDB.BeginTx(context.Background(), &sql.TxOptions{})
}

// BeginCoreTx returns a new transaction for Core database
func BeginCoreTx() (bun.Tx, error) {
	return dal.Dal.BlacheDB.BeginTx(context.Background(), &sql.TxOptions{})
}
