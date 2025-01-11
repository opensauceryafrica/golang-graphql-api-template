package schema

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"cendit.io/garage/function"
)

type Scanny struct{}

// Scan implements the Scanner interface.
func (h *Scanny) Scan(src interface{}) error {
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
func (h Scanny) Value() (driver.Value, error) {
	b, err := json.Marshal(h)
	return string(b), err
}

type Datum struct {
	CreatedAt time.Time `bun:"created_at" json:"created_at" rsfr:"false"`
	UpdatedAt time.Time `bun:"updated_at" json:"updated_at" rsfr:"false"`
}

/*
Date loads the created_at and updated_at fields of the address if not already present, otherwise, it loads the updated_at field only.

If the "pessimistic" parameter is set to true, it loads both fields regardless
*/
func (d *Datum) Date(pessimistic ...bool) {
	if len(pessimistic) > 0 && !pessimistic[0] {
		if d.CreatedAt.IsZero() {
			d.CreatedAt = time.Now()
			d.UpdatedAt = time.Now()
			return
		}
		d.UpdatedAt = time.Now()
		return
	}
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()
}

/* Fields returns the struct fields as a slice of interface{} values */
func (d *Datum) Fields() []interface{} {
	return function.ReturnStructFields(d)
}
