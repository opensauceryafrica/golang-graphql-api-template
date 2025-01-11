package schema

import (
	"context"
	"fmt"

	"cendit.io/garage/database"
	"cendit.io/garage/primer/enum"

	"github.com/uptrace/bun"
)

type Factory struct {
	bun.BaseModel `bun:"table:factories"`

	Key   string `json:"key" bun:"key,pk"`
	Value int64  `json:"value" bun:"value"`
}

// GenerateKey generates a new account number for a new account
// by shifting the cursor by the given step. Account numbers are linearly
// random.
func GenerateKey(key string) (string, error) {
	database.FactoryTableMutex.Lock()
	defer database.FactoryTableMutex.Unlock()
	if enum.FactoryCursor == 0 {
		cursor, err := ShiftCursorForKey(enum.FactoryStep, key)
		if err != nil {
			return "", err
		}
		enum.FactoryCursor = cursor
		enum.FactoryPointer = cursor - enum.FactoryStep
	}
	if enum.FactoryPointer == enum.FactoryCursor {
		cursor, err := ShiftCursorForKey(enum.FactoryStep, key)
		if err != nil {
			return "", err
		}
		enum.FactoryCursor = cursor
		enum.FactoryPointer = cursor - enum.FactoryStep
	}
	enum.FactoryPointer++
	return fmt.Sprintf("%010d", enum.FactoryPointer), nil
}

// ShiftCursorForKey shifts the cursor by the given step if the field exists else it creates it.
func ShiftCursorForKey(step int64, key string) (int64, error) {
	var cursor int64
	err := database.DB.NewRaw(`INSERT INTO factories (key, value) VALUES (?, ?) ON CONFLICT (key) DO UPDATE SET value = factories.value + ? RETURNING value`, key, step, step).Scan(context.Background(), &cursor)
	if err != nil {
		return 0, err
	}
	return cursor, nil
}
