package repository

import (
	"blacheapi/dal"
	"blacheapi/primer/enum"
	"blacheapi/primer/typing"
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/schema"
)

// schematic representation of a user's activity
type Activity struct {
	bun.BaseModel `bun:"table:activities"`

	ID            string    `bun:"id,pk" json:"id"`
	By            *int      `bun:"by" json:"by"`
	Role          enum.Role `bun:"role" json:"role"`
	SavingsID     string    `bun:"savings_id" json:"savings_id"`
	TransactionID string    `bun:"transaction_id" json:"transaction_id"`

	// can be preloaded using graphql reference resolvers
	User User `json:"user" bun:"-"`

	Resolver    string      `bun:"resolver" json:"resolver"`
	Payload     string      `bun:"payload" json:"payload"`
	Description string      `bun:"description" json:"description"`
	Status      enum.Status `bun:"status" json:"status"`
	Error       string      `bun:"error" json:"error"`

	CreatedAt bun.NullTime `bun:"created_at" json:"created_at"`
	UpdatedAt bun.NullTime `bun:"updated_at" json:"updated_at"`
}

type Activities []Activity

/*
Date loads the created_at and updated_at fields of the activity if not already present, otherwise, it loads the updated_at field only.

If the "pessimistic" parameter is set to true, it loads both fields regardless
*/
func (a *Activity) Date(pessimistic ...bool) {
	if len(pessimistic) > 0 && !pessimistic[0] {
		if a.CreatedAt.IsZero() {
			a.CreatedAt = schema.NullTime{Time: time.Now()}
			a.UpdatedAt = schema.NullTime{Time: time.Now()}
			return
		}
		a.UpdatedAt = schema.NullTime{Time: time.Now()}
		return
	}
	a.CreatedAt = schema.NullTime{Time: time.Now()}
	a.UpdatedAt = schema.NullTime{Time: time.Now()}
}

/*
Create inserts a new activity into the database

It returns an error if any
*/
func (a *Activity) Create() error {
	if _, err := dal.Dal.BlacheDB.NewRaw(`INSERT INTO activities (id, by, role, savings_id, transaction_id, resolver, payload, description, status, error, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, a.ID, a.By, a.Role, a.SavingsID, a.TransactionID, a.Resolver, a.Payload, a.Description, a.Status, a.Error, a.CreatedAt, a.UpdatedAt).Exec(context.Background()); err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

/*
FByKeyVal finds and returns an activity matching the key/value pair

# By default, only the ID and by fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded

It returns an error if any
*/
func (a *Activity) FByKeyVal(key string, val interface{}, preloadandjoin ...bool) error {
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT * FROM activities LEFT JOIN products ON activities.product_id = products.id WHERE products.`+key+` = ?`, val).Scan(context.Background(), a)
	}
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT * FROM activities WHERE `+key+` = ?`, val).Scan(context.Background(), a)
	}
	return dal.Dal.BlacheDB.NewRaw(`SELECT id, by, product_id FROM activities WHERE id = ?`, a.ID).Scan(context.Background(), a)
}

/*
FByMap finds and returns an activity matching the key/value pairs provided in the map

# By default, only the ID and name fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded

It returns an error if any
*/
func (a *Activity) FByMap(m typing.SQLMaps, preloadandjoin ...bool) error {
	query, args := dal.MapsToWQuery(m)
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT * FROM activities LEFT JOIN products ON activities.product_id = products.id WHERE `+query, args...).Scan(context.Background(), a)
	}
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT * FROM activities WHERE `+query, args...).Scan(context.Background(), a)
	}
	return dal.Dal.BlacheDB.NewRaw(`SELECT id, by, product_id FROM activities WHERE `+query, args...).Scan(context.Background(), a)
}

/*
FByMap finds and returns all activities matching the key/value pairs provided in the map

By default, only the id and product_id fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded

It returns an error if any
*/
func (c *Activities) FByMap(m typing.SQLMaps, limit, offset int, sort string, preloadandjoin ...bool) error {
	query, args := dal.MapsToWQuery(m)
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT * FROM activities LEFT JOIN products ON activities.product_id = product.id WHERE `+query+` ORDER BY activities.updated_at `+sort+` LIMIT ? OFFSET ?`, append(args, limit, offset)...).Scan(context.Background(), c)
	}
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT * FROM activities WHERE `+query+` ORDER BY activities.updated_at `+sort+` LIMIT ? OFFSET ?`, append(args, limit, offset)...).Scan(context.Background(), c)
	}
	return dal.Dal.BlacheDB.NewRaw(`SELECT id, product_id FROM activities WHERE `+query+` ORDER BY activities.updated_at `+sort+` LIMIT ? OFFSET ?`, append(args, limit, offset)...).Scan(context.Background(), c)
}
