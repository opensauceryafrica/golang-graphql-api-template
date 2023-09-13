package repository

import (
	"blacheapi/dal"
	"blacheapi/primer/enum"
	"context"
	"fmt"

	"github.com/uptrace/bun"
)

type Factory struct {
	bun.BaseModel `bun:"table:factories"`

	Key   string `json:"key" bun:"key,pk"`
	Value int64  `json:"value" bun:"value"`
}

// GenerateAccountNumber generates a new account number for a new account
// by shifting the cursor by the given step. Account numbers are linearly
// random.
func GenerateAccountNumber() (string, error) {
	dal.FactoryTableMutex.Lock()
	defer dal.FactoryTableMutex.Unlock()
	if enum.FactoryCursor == 0 {
		cursor, err := ShiftCursorForKey(enum.FactoryStep, "account_number")
		if err != nil {
			return "", err
		}
		enum.FactoryCursor = cursor
		enum.FactoryPointer = cursor - enum.FactoryStep
	}
	if enum.FactoryPointer == enum.FactoryCursor {
		cursor, err := ShiftCursorForKey(enum.FactoryStep, "account_number")
		if err != nil {
			return "", err
		}
		enum.FactoryCursor = cursor
		enum.FactoryPointer = cursor - enum.FactoryStep
	}
	enum.FactoryPointer++
	return fmt.Sprintf("%010d", enum.FactoryPointer), nil
}

// ShiftCursorForAccountNumber shifts the cursor by the given step if the field exists else it creates it.
func ShiftCursorForKey(step int64, key string) (int64, error) {
	var cursor int64
	err := dal.Dal.BlacheDB.NewRaw(`INSERT INTO factories (key, value) VALUES (?, ?) ON CONFLICT (key) DO UPDATE SET value = factories.value + ? RETURNING value`, key, step, step).Scan(context.Background(), &cursor)
	if err != nil {
		return 0, err
	}
	return cursor, nil
}
