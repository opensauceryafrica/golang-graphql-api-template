package database

import (
	"database/sql"

	"cendit.io/garage/logger"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var DB *bun.DB

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
func New(url string, debug bool, connections int) (*bun.DB, error) {

	if DB != nil {
		return DB, nil
	}

	var err error

	DB, err = connectDB(url, debug, connections)
	if err != nil {
		logger.GetLogger().Sugar().Fatalf("[DB]: unable to intiate connection to DB: %v", err.Error())
		return nil, err
	}

	return DB, nil
}

// Close closes all database connections
func Close() {
	if DB != nil {
		DB.Close()
	}
}
