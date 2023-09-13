package dal

import (
	"blacheapi/config"
	"blacheapi/logger"
	"database/sql"
	"reflect"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

// Database Access Layer
type DAL struct {
	BlacheDB *bun.DB
	CoreDB   *bun.DB
}

var Dal *DAL

// connectDB initiates a connection to the database
func connectDB(url string, debug bool, conn int) (*bun.DB, error) {
	pgDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(url)))

	pgDB.SetMaxOpenConns(conn)

	db := bun.NewDB(pgDB, pgdialect.New())

	if debug {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	return db, db.Ping()
}

// New initiates a connection to the database,
// process schema migration and setsup the database layer
func New(cfg *config.Config) (*DAL, error) {

	blachedb, connErr := connectDB(cfg.BlacheDatabaseURL, cfg.DebugDatabase, cfg.DatabaseConnectionLimit)
	if connErr != nil {
		logger.GetLogger().Sugar().Fatalf("[DB]: unable to intiate connection to BlacheDB: %v", connErr.Error())
		return nil, connErr
	}

	// you can connect to other databases here

	Dal = &DAL{
		BlacheDB: blachedb,
		// you can add other databases here
	}

	return Dal, nil
}

// Close closes all database connections
func (d *DAL) Close() {
	dal := reflect.ValueOf(d).Elem()
	for i := 0; i < dal.NumField(); i++ {
		if dal.Field(i).Kind() == reflect.Ptr {
			// assert that it's of type *bun.DB
			if _, ok := dal.Field(i).Interface().(*bun.DB); ok {
				err := dal.Field(i).Interface().(*bun.DB).Close()
				if err != nil {
					logger.GetLogger().Sugar().Fatal(err)
				}
			}
		}
	}
}
