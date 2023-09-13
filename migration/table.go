// a database interface for moving parts while ensuring not to break things
package migration

import (
	"blacheapi/config"
	"blacheapi/dal"
	"blacheapi/logger"
	"blacheapi/repository"
	"context"
)

var Table = []interface{}{
	&repository.Activity{},
	&repository.Factory{},
	&repository.Transaction{},
}

// CreateTables creates tables that do not already exist. Although we have connections to other DBs configure.Blache should only handle migration for configure.Blache DB.
func CreateTables() error {
	for _, m := range Table {
		_, err := dal.Dal.BlacheDB.NewCreateTable().
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
func Migrate(cfg *config.Config) error {
	return nil
}
