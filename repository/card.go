package repository

import (
	"blacheapi/dal"
	"blacheapi/primer/enum"
	"blacheapi/primer/function"
	"blacheapi/primer/typing"
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

// schematic representation of a blacher's credit/debit card
type Card struct {
	bun.BaseModel `bun:"table:cards" rsf:"false"`

	ID         string `bun:"id,pk" json:"id"`
	CustomerID int    `bun:"customer_id" json:"customer_id"`
	// owner
	UserID int `bun:"user_id" json:"user_id"`

	// CONFIGURATION

	First6digits string `bun:"first_6digits" json:"first_6digits"`
	Last4digits  string `bun:"last_4digits" json:"last_4digits"`
	Issuer       string `bun:"issuer" json:"issuer"`
	Country      string `bun:"country" json:"country"`
	Type         string `bun:"type" json:"type"`
	Expiry       string `bun:"expiry" json:"expiry"`
	Checksum     string `bun:"checksum" json:"checksum"`

	// AUTHORIZATION
	Token string `bun:"token" json:"token"`

	CreatedAt bun.NullTime `bun:"created_at" json:"created_at" rsfr:"false"`
	UpdatedAt bun.NullTime `bun:"updated_at" json:"updated_at" rsfr:"false"`
}

type Cards []Card

/* Fields returns the struct fields as a slice of interface{} values */
func (c *Card) Fields() []interface{} {
	return function.ReturnStructFields(c)
}

/*
Date loads the created_at and updated_at fields of the card if not already present, otherwise, it loads the updated_at field only.

If the "pessimistic" parameter is set to true, it loads both fields regardless
*/
func (c *Card) Date(pessimistic ...bool) {
	if len(pessimistic) > 0 && !pessimistic[0] {
		if c.CreatedAt.IsZero() {
			c.CreatedAt = schema.NullTime{Time: time.Now()}
			c.UpdatedAt = schema.NullTime{Time: time.Now()}
			return
		}
		c.UpdatedAt = schema.NullTime{Time: time.Now()}
		return
	}
	c.CreatedAt = schema.NullTime{Time: time.Now()}
	c.UpdatedAt = schema.NullTime{Time: time.Now()}
}

/*
Exist checks if a card with the given key/value pair exists in the database

# If returns true if it exists, false otherwise

It returns an error if any
*/
func (c *Card) Exist(key string, value interface{}) (bool, error) {
	var exists interface{}
	err := dal.Dal.BlacheDB.NewRaw(`SELECT EXISTS(SELECT 1 FROM cards WHERE `+key+` = ?)`, value).Scan(context.Background(), &exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists.(bool), nil
}

/*
Create inserts a new card into the database

It returns an error if any
*/
func (c *Card) Create() error {
	if _, err := dal.Dal.BlacheDB.NewRaw(`INSERT INTO cards (id, user_id, customer_id, first_6digits, last_4digits, issuer, country, type, expiry, checksum, token, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, c.ID, c.UserID, c.CustomerID, c.First6digits, c.Last4digits, c.Issuer, c.Country, c.Type, c.Expiry, c.Checksum, c.Token, c.CreatedAt, c.UpdatedAt).Exec(context.Background()); err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

/*
CreateTx inserts a new card into the database with the given transaction

It returns an error if any
*/
func (c *Card) CreateTx(tx bun.Tx) error {
	if _, err := tx.NewRaw(`INSERT INTO cards (id, user_id, customer_id, first_6digits, last_4digits, issuer, country, type, expiry, checksum, token, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, c.ID, c.UserID, c.CustomerID, c.First6digits, c.Last4digits, c.Issuer, c.Country, c.Type, c.Expiry, c.Checksum, c.Token, c.CreatedAt, c.UpdatedAt).Exec(context.Background()); err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

/*
FByKeyVal finds and returns a card matching the key/value pair

By default, only the id and product_id fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded

It returns an error if any
*/
func (c *Card) FByKeyVal(key string, val interface{}, preloadandjoin ...bool) error {
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT * FROM cards WHERE `+key+` = ?`, val).Scan(context.Background(), c)
	}
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT * FROM cards WHERE `+key+` = ?`, val).Scan(context.Background(), c)
	}
	return dal.Dal.BlacheDB.NewRaw(`SELECT id, product_id FROM cards WHERE cards.`+key+` = ?`, val).Scan(context.Background(), c)
}

/*
FByKeyVal finds and returns all cards matching the key/value pair

By default, only the id and product_id fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded

It returns an error if any
*/
func (c *Cards) FByKeyVal(key string, val interface{}, limit, offset int, sort string, preloadandjoin ...bool) error {
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT * FROM cards LEFT JOIN users ON users.id = cards.user_id WHERE cards.`+key+` = ? ORDER BY cards.updated_at `+sort+` LIMIT ? OFFSET ?`, val, limit, offset).Scan(context.Background(), c)
	}
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT * FROM cards WHERE `+key+` = ? ORDER BY cards.updated_at `+sort+` LIMIT ? OFFSET ?`, val, limit, offset).Scan(context.Background(), c)
	}
	return dal.Dal.BlacheDB.NewRaw(`SELECT id, product_id FROM cards WHERE `+key+` = ? ORDER BY cards.updated_at `+sort+` LIMIT ? OFFSET ?`, val, limit, offset).Scan(context.Background(), c)
}

/*
FByMap finds and returns all cards matching the key/value pairs provided in the map

By default, only the id and product_id fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded

It returns an error if any
*/
func (c *Card) FByMap(m typing.SQLMaps, preloadandjoin ...bool) error {
	query, args := dal.MapsToWQuery(m)
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT * FROM cards WHERE `+query, args...).Scan(context.Background(), c)
	}
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT * FROM cards WHERE `+query, args...).Scan(context.Background(), c)
	}
	return dal.Dal.BlacheDB.NewRaw(`SELECT id, product_id FROM cards WHERE `+query, args...).Scan(context.Background(), c)
}

/*
FByMap finds and returns all cards matching the key/value pairs provided in the map

By default, only the id and product_id fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded

It returns an error if any
*/
func (c *Cards) FByMap(m typing.SQLMaps, limit, offset int, sort string, preloadandjoin ...bool) error {
	query, args := dal.MapsToWQuery(m)
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		rows, err := dal.Dal.BlacheDB.QueryContext(context.Background(), `SELECT * FROM cards LEFT JOIN users ON users.id = cards.user_id WHERE `+query+` ORDER BY cards.updated_at `+sort+` LIMIT ? OFFSET ?`, append(args, limit, offset)...)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var _c Card
			if err := rows.Scan(_c.Fields()...); err != nil {
				return err
			}
			*c = append(*c, _c)
		}
	}
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT * FROM cards WHERE `+query+` ORDER BY cards.updated_at `+sort+` LIMIT ? OFFSET ?`, append(args, limit, offset)...).Scan(context.Background(), c)
	}
	return dal.Dal.BlacheDB.NewRaw(`SELECT id, product_id FROM cards WHERE `+query+` ORDER BY cards.updated_at `+sort+` LIMIT ? OFFSET ?`, append(args, limit, offset)...).Scan(context.Background(), c)
}

/*
UByMap updates a card matching the key/value pairs provided in the map

It returns an error if any
*/
func (c *Card) UByMap(m typing.SQLMaps) error {
	query, args := dal.MapsToSQuery(m)
	if strings.Contains(query, string(enum.RETURNING)) {
		return dal.Dal.BlacheDB.NewRaw(`UPDATE cards `+query, args...).Scan(context.Background(), c)
	}
	_, err := dal.Dal.BlacheDB.NewRaw(`UPDATE cards `+query, args...).Exec(context.Background())
	return err
}

/*
DByMap deletes a card matching the key/value pairs provided in the map

It returns an error if any
*/
func (c *Card) DByMap(m typing.SQLMaps) error {
	query, args := dal.MapsToWQuery(m)
	_, err := dal.Dal.BlacheDB.NewRaw(`DELETE FROM cards WHERE `+query, args...).Exec(context.Background())
	return err
}
