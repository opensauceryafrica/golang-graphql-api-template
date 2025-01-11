package xiao

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"cendit.io/garage/function"
)

/*
Xiao is a base struct for postgres data models

xiao supports any driver that implements the database/sql interface

xiao uses the lowercase value of your struct name as the table name by default but you can change this by overriding the TableName property of the Xiao instance

xiao provides a Preloaders property to define a limited set of columns to return when fetching data with preload set to true. it defaults to an empty array and returns all columns

xiao assumes a postgres dialect by default but you can change this by overriding the Dialect property of the Xiao instance
*/
type Xiao[T any] struct {
	TableName  string
	Preloaders []string
	Pool       *sql.DB
	Dialect    SQLDialect
}

// NewXiao creates a new instance of Xiao
func NewXiao[T any](db *sql.DB) *Xiao[T] {
	return &Xiao[T]{
		TableName:  strings.ToLower(fmt.Sprintf("%T", new(T))),
		Preloaders: []string{},
		Pool:       db,
		Dialect:    Postgres,
	}
}

// Tx creates a new transaction
func (x *Xiao[T]) Tx(ctx context.Context) (*sql.Tx, error) {
	return x.Pool.BeginTx(ctx, nil)
}

// Exists checks if a record exists in the table by a key-value pair
func (x *Xiao[T]) Exists(ctx context.Context, key string, value interface{}) (bool, error) {
	query := fmt.Sprintf("%s %s(%s 1 %s %s %s %s = ?)",
		SQLSelect, SQLExists, SQLSelect, SQLFrom, x.TableName, SQLWhere, key)
	row := x.Pool.QueryRowContext(ctx, Parametrize(query, string(x.Dialect)), value)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

// Create inserts a new record into the table
func (x *Xiao[T]) Create(ctx context.Context, m SQLMaps) error {
	query, args := MapsToIQuery(m)
	_, err := x.Pool.ExecContext(ctx, fmt.Sprintf("%s %s %s %s", SQLInsert, SQLInto, x.TableName, Parametrize(query, string(x.Dialect))), args...)
	return err
}

// CreateTx inserts a new record into the table within a transaction
func (x *Xiao[T]) CreateTx(ctx context.Context, tx *sql.Tx, m SQLMaps) error {
	query, args := MapsToIQuery(m)
	_, err := tx.ExecContext(ctx, fmt.Sprintf("%s %s %s %s", SQLInsert, SQLInto, x.TableName, Parametrize(query, string(x.Dialect))), args...)
	return err
}

// FindByKeyVal finds a record by a key-value pair
func (x *Xiao[T]) FindByKeyVal(ctx context.Context, key string, val interface{}, preload bool) (*T, error) {
	selectColumns := "*"
	if !preload {
		selectColumns = strings.Join(x.Preloaders, ", ")
	}
	query := fmt.Sprintf("%s %s %s %s %s %s = ?", SQLSelect, selectColumns, SQLFrom, x.TableName, SQLWhere, key)
	row := x.Pool.QueryRowContext(ctx, Parametrize(query, string(x.Dialect)), val)
	var result T
	if err := row.Scan(function.ReturnStructFields(&result)...); err != nil {
		return &result, err
	}
	return &result, nil
}

// FindAllByKeyVal finds all records by a key-value pair
func (x *Xiao[T]) FindAllByKeyVal(ctx context.Context, key string, val interface{}, preload bool) ([]T, error) {
	selectColumns := "*"
	if !preload {
		selectColumns = strings.Join(x.Preloaders, ", ")
	}
	query := fmt.Sprintf("%s %s %s %s %s %s = ?", SQLSelect, selectColumns, SQLFrom, x.TableName, SQLWhere, key)
	rows, err := x.Pool.QueryContext(ctx, Parametrize(query, string(x.Dialect)), val)
	if err != nil {
		return []T{}, err
	}
	defer rows.Close()

	var results []T
	for rows.Next() {
		var result T
		if err := rows.Scan(function.ReturnStructFields(&result)...); err != nil {
			return []T{}, err
		}
		results = append(results, result)
	}
	return results, nil
}

// FindAndLockByKeyVal finds a record in the table by a key-value pair and locks it within a transaction
func (x *Xiao[T]) FindAndLockByKeyVal(ctx context.Context, tx *sql.Tx, key string, val interface{}, preload bool) (*T, error) {
	selectColumns := "*"
	if !preload {
		selectColumns = strings.Join(x.Preloaders, ", ")
	}
	query := fmt.Sprintf("%s %s %s %s %s = ? %s %s", SQLSelect, selectColumns, SQLFrom, x.TableName, key, SQLFor, SQLUpdate)
	row := tx.QueryRowContext(ctx, Parametrize(query, string(x.Dialect)), val)
	var result T
	if err := row.Scan(function.ReturnStructFields(&result)...); err != nil {
		return &result, err
	}
	return &result, nil
}

// FindAllAndLockByKeyVal finds all records in the table by a key-value pair and locks them within a transaction
func (x *Xiao[T]) FindAllAndLockByKeyVal(ctx context.Context, tx *sql.Tx, key string, val interface{}, preload bool) ([]T, error) {
	selectColumns := "*"
	if !preload {
		selectColumns = strings.Join(x.Preloaders, ", ")
	}
	query := fmt.Sprintf("%s %s %s %s %s = ? %s %s", SQLSelect, selectColumns, SQLFrom, x.TableName, key, SQLFor, SQLUpdate)
	rows, err := tx.QueryContext(ctx, Parametrize(query, string(x.Dialect)), val)
	if err != nil {
		return []T{}, err
	}
	defer rows.Close()

	var results []T
	for rows.Next() {
		var result T
		if err := rows.Scan(function.ReturnStructFields(&result)...); err != nil {
			return []T{}, err
		}
		results = append(results, result)
	}
	return results, nil
}

// FindByMap finds a record in the table by a map of conditions
func (x *Xiao[T]) FindByMap(ctx context.Context, m SQLMaps, preload bool) (*T, error) {
	query, args := MapsToWQuery(m)
	selectColumns := "*"
	if !preload {
		selectColumns = strings.Join(x.Preloaders, ", ")
	}
	fullQuery := fmt.Sprintf("%s %s %s %s %s %s", SQLSelect, selectColumns, SQLFrom, x.TableName, SQLWhere, query)
	row := x.Pool.QueryRowContext(ctx, Parametrize(fullQuery, string(x.Dialect)), args...)
	var result T
	if err := row.Scan(function.ReturnStructFields(&result)...); err != nil {
		return &result, err
	}
	return &result, nil
}

// FindAllByMap finds all records in the table by a map of conditions
func (x *Xiao[T]) FindAllByMap(ctx context.Context, m SQLMaps, preload bool) ([]T, error) {
	query, args := MapsToWQuery(m)
	oquery := MapsToOQuery(m)
	lquery := MapsToLQuery(m)
	selectColumns := "*"
	if !preload {
		selectColumns = strings.Join(x.Preloaders, ", ")
	}
	fullQuery := fmt.Sprintf("%s %s %s %s %s %s %s %s", SQLSelect, selectColumns, SQLFrom, x.TableName, SQLWhere, query, oquery, lquery)
	rows, err := x.Pool.QueryContext(ctx, Parametrize(fullQuery, string(x.Dialect)), args...)
	if err != nil {
		return []T{}, err
	}
	defer rows.Close()

	var results []T
	for rows.Next() {
		var result T
		if err := rows.Scan(function.ReturnStructFields(&result)...); err != nil {
			return []T{}, err
		}
		results = append(results, result)
	}
	return results, nil
}

// FindAndLockByMap finds a record in the table by a map of conditions and locks it within a transaction
func (x *Xiao[T]) FindAndLockByMap(ctx context.Context, tx *sql.Tx, m SQLMaps, preload bool) (*T, error) {
	query, args := MapsToWQuery(m)
	selectColumns := "*"
	if !preload {
		selectColumns = strings.Join(x.Preloaders, ", ")
	}
	fullQuery := fmt.Sprintf("%s %s %s %s %s %s %s %s", SQLSelect, selectColumns, SQLFrom, x.TableName, SQLWhere, query, SQLFor, SQLUpdate)
	row := tx.QueryRowContext(ctx, Parametrize(fullQuery, string(x.Dialect)), args...)
	var result T
	if err := row.Scan(function.ReturnStructFields(&result)...); err != nil {
		return &result, err
	}
	return &result, nil
}

// FindAllAndLockByMap finds all records in the table by a map of conditions and locks them within a transaction
func (x *Xiao[T]) FindAllAndLockByMap(ctx context.Context, tx *sql.Tx, m SQLMaps, preload bool) ([]T, error) {
	query, args := MapsToWQuery(m)
	oquery := MapsToOQuery(m)
	lquery := MapsToLQuery(m)
	selectColumns := "*"
	if !preload {
		selectColumns = strings.Join(x.Preloaders, ", ")
	}
	fullQuery := fmt.Sprintf("%s %s %s %s %s %s %s %s %s %s", SQLSelect, selectColumns, SQLFrom, x.TableName, SQLWhere, query, oquery, lquery, SQLFor, SQLUpdate)
	rows, err := tx.QueryContext(ctx, Parametrize(fullQuery, string(x.Dialect)), args...)
	if err != nil {
		return []T{}, err
	}
	defer rows.Close()

	var results []T
	for rows.Next() {
		var result T
		if err := rows.Scan(function.ReturnStructFields(&result)...); err != nil {
			return []T{}, err
		}
		results = append(results, result)
	}
	return results, nil
}

// UpdateByMap updates records in the table by a map of conditions and returns the updated records if requested
func (x *Xiao[T]) UpdateByMap(ctx context.Context, m SQLMaps) ([]T, error) {
	query, args := MapsToSQuery(m)
	fullQuery := fmt.Sprintf("%s %s %s", SQLUpdate, x.TableName, query)

	// If the query contains RETURNING, fetch and return the updated records
	if strings.Contains(query, SQLReturning.String()) {
		rows, err := x.Pool.QueryContext(ctx, Parametrize(fullQuery, string(x.Dialect)), args...)
		if err != nil {
			return []T{}, err
		}
		defer rows.Close()

		var results []T
		for rows.Next() {
			var result T
			if err := rows.Scan(function.ReturnStructFields(&result)...); err != nil {
				return []T{}, err
			}
			results = append(results, result)
		}
		return results, nil
	}

	// Otherwise, just execute the update without returning anything
	_, err := x.Pool.ExecContext(ctx, Parametrize(fullQuery, string(x.Dialect)), args...)
	return []T{}, err
}

// UpdateByMapTx updates records in the table by a map of conditions within a transaction and returns the updated records if requested
func (x *Xiao[T]) UpdateByMapTx(ctx context.Context, tx *sql.Tx, m SQLMaps) ([]T, error) {
	query, args := MapsToSQuery(m)
	fullQuery := fmt.Sprintf("%s %s %s", SQLUpdate, x.TableName, query)

	// If the query contains RETURNING, fetch and return the updated records
	if strings.Contains(query, SQLReturning.String()) {
		rows, err := tx.QueryContext(ctx, Parametrize(fullQuery, string(x.Dialect)), args...)
		if err != nil {
			return []T{}, err
		}
		defer rows.Close()

		var results []T
		for rows.Next() {
			var result T
			if err := rows.Scan(function.ReturnStructFields(&result)...); err != nil {
				return []T{}, err
			}
			results = append(results, result)
		}
		return results, nil
	}

	// Otherwise, just execute the update without returning anything
	_, err := tx.ExecContext(ctx, Parametrize(fullQuery, string(x.Dialect)), args...)
	return []T{}, err
}

// DeleteByMap deletes a record by a map of conditions
func (x *Xiao[T]) DeleteByMap(ctx context.Context, m SQLMaps) error {
	query, args := MapsToWQuery(m)
	_, err := x.Pool.ExecContext(ctx, fmt.Sprintf("%s %s %s %s %s", SQLDelete, SQLFrom, x.TableName, SQLWhere, Parametrize(query, string(x.Dialect))), args...)
	return err
}

// DeleteByMapTx deletes a record by a map of conditions within a transaction
func (x *Xiao[T]) DeleteByMapTx(ctx context.Context, tx *sql.Tx, m SQLMaps) error {
	query, args := MapsToWQuery(m)
	_, err := tx.ExecContext(ctx, fmt.Sprintf("%s %s %s %s %s", SQLDelete, SQLFrom, x.TableName, SQLWhere, Parametrize(query, string(x.Dialect))), args...)
	return err
}

// CountByMap counts the number of records matching a map of conditions
func (x *Xiao[T]) CountByMap(ctx context.Context, m SQLMaps) (int, error) {
	query, args := MapsToWQuery(m)
	row := x.Pool.QueryRowContext(ctx, fmt.Sprintf("%s %s %s %s %s %s", SQLSelect, SQLCount("*"), SQLFrom, x.TableName, SQLWhere, Parametrize(query, string(x.Dialect))), args...)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// Execute runs a raw SQL query on the database
func (x *Xiao[T]) Execute(ctx context.Context, query string, args ...interface{}) ([]T, error) {
	rows, err := x.Pool.QueryContext(ctx, Parametrize(query, string(x.Dialect)), args...)
	if err != nil {
		return []T{}, err
	}
	defer rows.Close()

	var results []T
	for rows.Next() {
		var result T
		if err := rows.Scan(function.ReturnStructFields(&result)...); err != nil {
			return []T{}, err
		}
		results = append(results, result)
	}
	return results, nil
}

// ExecuteTx runs a raw SQL query within a transaction
func (x *Xiao[T]) ExecuteTx(ctx context.Context, tx *sql.Tx, query string, args ...interface{}) ([]T, error) {
	rows, err := tx.QueryContext(ctx, Parametrize(query, string(x.Dialect)), args...)
	if err != nil {
		return []T{}, err
	}
	defer rows.Close()

	var results []T
	for rows.Next() {
		var result T
		if err := rows.Scan(function.ReturnStructFields(&result)...); err != nil {
			return []T{}, err
		}
		results = append(results, result)
	}
	return results, nil
}
