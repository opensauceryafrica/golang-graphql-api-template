package migration

import (
	"context"

	"cendit.io/garage/logger"
	"github.com/uptrace/bun"
)

// CreateTables creates tables that do not already exist. Although we have connections to other DBs configure.Cendit should only handle migration for configure.Cendit DB.
func CreateTables(db *bun.DB, tables []interface{}) error {
	for _, m := range tables {
		_, err := db.NewCreateTable().
			IfNotExists().
			Model(m).Exec(context.TODO())
		if err != nil {
			logger.GetLogger().Sugar().Warnf("failed to create %v table", m)
			return err
		}
	}
	return nil
}

// migrate effects any database schema migration
func Migrate(db *bun.DB) error {
	return nil
}
