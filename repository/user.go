package repository

import (
	"blacheapi/dal"
	"blacheapi/primer/enum"
	"blacheapi/primer/function"
	"blacheapi/primer/typing"
	"context"

	"github.com/uptrace/bun"
)

// schematic representation of a user

// IMPORTANT: when you add a new field to this struct, you must also add it to the SQL query in the FByKeyVal and FByMap methods
type User struct {
	// @TODO: id should be a string
	ID        int    `bun:"id" json:"id"`
	OrgID     int    `bun:"org_id" json:"org_id"`
	Email     string `bun:"email" json:"email"`
	Firstname string `bun:"firstname" json:"firstname"`
	Lastname  string `bun:"lastname" json:"lastname"`
	Phone     string `bun:"phone" json:"phone"`
	Gender    string `bun:"gender" json:"gender"`

	CreatedAt bun.NullTime `bun:"created_at" json:"created_at"`
	UpdatedAt bun.NullTime `bun:"updated_at" json:"updated_at"`

	// prepended during authentication
	Role enum.Role
}

type Users []User

/* Fields returns the struct fields as a slice of interface{} values */
func (u *User) Fields() []interface{} {
	return function.ReturnStructFields(u)
}

/*
FByKeyVal finds and returns a user matching the key/value pair

# By default, only the id and name fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded. The join makes a lookup on the savings table to only return users that have a savings account

It returns an error if any
*/
func (u *User) FByKeyVal(key string, val interface{}, preloadandjoin ...bool) error {
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT DISTINCT ON (core.users.id) core.users.id, core.users.org_id, core.users.email, core.users.firstname, core.users.lastname, core.users.phone, core.users.gender, core.users.created_at, core.users.updated_at FROM core.users JOIN savings ON core.users.id = savings.user_id WHERE core.users.`+key+` = ?`, val).Scan(context.Background(), u)
	}
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return dal.Dal.CoreDB.NewRaw(`SELECT core.users.id, core.users.org_id, core.users.email, core.users.firstname, core.users.lastname, core.users.phone, core.users.gender, core.users.created_at, core.users.updated_at FROM core.users WHERE `+key+` = ?`, val).Scan(context.Background(), u)
	}
	return dal.Dal.BlacheDB.NewRaw(`SELECT id, firstname, lastname FROM core.users WHERE core.users.`+key+` = ?`, val).Scan(context.Background(), u)
}

/*
FByKeyVal finds and returns all users matching the key/value pair

# By default, only the id and name fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded. The join makes a lookup on the savings table to only return users that have a savings account

It returns an error if any
*/
func (u *Users) FByKeyVal(key string, val interface{}, limit, offset int, sort string, preloadandjoin ...bool) error {
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT DISTINCT ON (core.users.id) core.users.id, core.users.org_id, core.users.email, core.users.firstname, core.users.lastname, core.users.phone, core.users.gender, core.users.created_at, core.users.updated_at FROM core.users JOIN savings ON core.users.id = savings.user_id WHERE core.users.`+key+` = ? ORDER BY core.users.id, core.users.updated_at `+sort+` LIMIT ? OFFSET ?`, val, limit, offset).Scan(context.Background(), u)
	}
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT core.users.id, core.users.org_id, core.users.email, core.users.firstname, core.users.lastname, core.users.phone, core.users.gender, core.users.created_at, core.users.updated_at FROM core.users WHERE `+key+` = ? ORDER BY core.users.updated_at `+sort+` LIMIT ? OFFSET ?`, val, limit, offset).Scan(context.Background(), u)
	}
	return dal.Dal.BlacheDB.NewRaw(`SELECT id, firstname, lastname FROM core.users WHERE `+key+` = ? ORDER BY core.users.updated_at `+sort+` LIMIT ? OFFSET ?`, val, limit, offset).Scan(context.Background(), u)
}

/*
FByMap finds and returns a user matching the key/value pairs provided in the map

# By default, only the id and name fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded. The join makes a lookup on the savings table to only return users that have a savings account

It returns an error if any
*/
func (u *User) FByMap(m typing.SQLMaps, preloadandjoin ...bool) error {
	query, args := dal.MapsToWQuery(m)
	join, jargs := dal.MapsToJQuery(m)
	args = append(args, jargs...)
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		if join == "" {
			join = "core.users.id = savings.user_id"
		}
		return dal.Dal.BlacheDB.NewRaw(`SELECT DISTINCT ON (core.users.id) core.users.id, core.users.org_id, core.users.email, core.users.firstname, core.users.lastname, core.users.phone, core.users.gender, core.users.created_at, core.users.updated_at FROM core.users JOIN savings ON `+join+` WHERE `+query, args...).Scan(context.Background(), u)
	}
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT core.users.id, core.users.org_id, core.users.email, core.users.firstname, core.users.lastname, core.users.phone, core.users.gender, core.users.created_at, core.users.updated_at FROM core.users WHERE `+query, args...).Scan(context.Background(), u)
	}
	return dal.Dal.BlacheDB.NewRaw(`SELECT id, firstname, lastname FROM core.users WHERE `+query, args...).Scan(context.Background(), u)
}

/*
FByMap finds and returns all users matching the key/value pairs provided in the map

# By default, only the id and name fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded. The join makes a lookup on the savings table to only return users that have a savings account

It returns an error if any
*/
func (u *Users) FByMap(m typing.SQLMaps, limit, offset int, sort string, preloadandjoin ...bool) error {
	query, args := dal.MapsToWQuery(m)
	join, jargs := dal.MapsToJQuery(m)
	args = append(jargs, args...)
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		if join == "" {
			join = "core.users.id = savings.user_id"
		}
		return dal.Dal.BlacheDB.NewRaw(`SELECT DISTINCT ON (core.users.id) core.users.id, core.users.org_id, core.users.email, core.users.firstname, core.users.lastname, core.users.phone, core.users.gender, core.users.created_at, core.users.updated_at FROM core.users JOIN savings ON `+join+` WHERE `+query+` ORDER BY core.users.id, core.users.updated_at `+sort+` LIMIT ? OFFSET ?`, append(args, limit, offset)...).Scan(context.Background(), u)
	}
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT core.users.id, core.users.org_id, core.users.email, core.users.firstname, core.users.lastname, core.users.phone, core.users.gender, core.users.created_at, core.users.updated_at FROM core.users WHERE `+query+` ORDER BY core.users.updated_at `+sort+` LIMIT ? OFFSET ?`, append(args, limit, offset)...).Scan(context.Background(), u)
	}
	return dal.Dal.BlacheDB.NewRaw(`SELECT id, firstname, lastname FROM core.users WHERE `+query+` ORDER BY core.users.updated_at `+sort+` LIMIT ? OFFSET ?`, append(args, limit, offset)...).Scan(context.Background(), u)
}

/*
CByMap finds and counts all users matching the key/value pairs provided in the map

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded. The join makes a lookup on the savings table to only count users that have a savings account

It returns an error if any
*/
func (s *Users) CByMap(m typing.SQLMaps, preloadandjoin ...bool) (int, error) {
	var count int
	query, args := dal.MapsToWQuery(m)
	join, jargs := dal.MapsToJQuery(m)
	args = append(jargs, args...)
	if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
		if join == "" {
			join = "core.users.id = savings.user_id"
		}
		err := dal.Dal.BlacheDB.NewRaw(`SELECT count(DISTINCT core.users.id) FROM core.users JOIN savings ON `+join+` WHERE `+query, args...).Scan(context.Background(), &count)
		return count, err
	}
	err := dal.Dal.BlacheDB.NewRaw(`SELECT count(*) FROM core.users WHERE `+query, args...).Scan(context.Background(), &count)
	return count, err
}

// schematic representation of a customer

// IMPORTANT: when you add a new field to this struct, you must also add it to the SQL query in the FByKeyVal and FByMap methods
type Customer struct {
	ID int `bun:"id" json:"id"`
	// @TODO: id should be a string
	OrgID           int    `bun:"org_id" json:"org_id"`
	SupportEmail    string `bun:"support_email" json:"support_email"`
	Country         string `bun:"country" json:"country"`
	DefaultCurrency string `bun:"default_currency" json:"default_currency"`
	CompanyName     string `bun:"company_name" json:"company_name" `
	CompanyNameSlug string `bun:"company_name_slug" json:"company_name_slug"`
	Logo            string `bun:"logo" json:"logo"`
	Website         string `bun:"website" json:"website"`

	CreatedAt bun.NullTime `bun:"created_at" json:"created_at"`
	UpdateAt  bun.NullTime `bun:"updated_at" json:"updated_at"`

	// prepended during authentication
	Role enum.Role
}

type Customers []Customer

/* Fields returns the struct fields as a slice of interface{} values */
func (c *Customer) Fields() []interface{} {
	return function.ReturnStructFields(c)
}

/*
FByKeyVal finds and returns a customer matching the key/value pair

# By default, only the id and name fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded

It returns an error if any
*/
func (c *Customer) FByKeyVal(key string, val interface{}, preloadandjoin ...bool) error {
	// if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
	// 	return dal.Dal.BlacheDB.NewRaw(`SELECT * FROM org_profiles LEFT JOIN configs ON org_profiles.id = configs.product_id WHERE org_profiles.`+key+` = ?`, val).Scan(context.Background(), p.Fields()...)
	// }
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return dal.Dal.CoreDB.NewRaw(`SELECT org_profiles.id, org_profiles.org_id, org_profiles.support_email, org_profiles.country, org_profiles.default_currency, org_profiles.company_name, org_profiles.company_name_slug, org_profiles.logo, org_profiles.website, org_profiles.created_at, org_profiles.updated_at	FROM org_profiles WHERE `+key+` = ?`, val).Scan(context.Background(), c)
	}
	return dal.Dal.BlacheDB.NewRaw(`SELECT org_id, name FROM org_profiles WHERE org_profiles.`+key+` = ?`, val).Scan(context.Background(), c)
}

/*
FByKeyVal finds and returns all customers matching the key/value pair

# By default, only the id and name fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded

It returns an error if any
*/
func (c *Customers) FByKeyVal(key string, val interface{}, limit, offset int, sort string, preloadandjoin ...bool) error {
	// if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
	// 	return dal.Dal.BlacheDB.NewRaw(`SELECT * FROM org_profiles LEFT JOIN configs ON org_profiles.id = configs.product_id WHERE org_profiles.`+key+` = ? ORDER BY org_profiles.updated_at `+sort+` LIMIT ? OFFSET ?`, val, limit, offset).Scan(context.Background(), c)
	// }
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT org_profiles.id, org_profiles.org_id, org_profiles.support_email, org_profiles.country, org_profiles.default_currency, org_profiles.company_name, org_profiles.company_name_slug, org_profiles.logo, org_profiles.website, org_profiles.created_at, org_profiles.updated_at FROM org_profiles WHERE `+key+` = ? ORDER BY org_profiles.updated_at `+sort+` LIMIT ? OFFSET ?`, val, limit, offset).Scan(context.Background(), c)
	}
	return dal.Dal.BlacheDB.NewRaw(`SELECT org_id, name FROM org_profiles WHERE `+key+` = ? ORDER BY org_profiles.updated_at `+sort+` LIMIT ? OFFSET ?`, val, limit, offset).Scan(context.Background(), c)
}

/*
FByMap finds and returns a customer matching the key/value pairs provided in the map

# By default, only the id and name fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded

It returns an error if any
*/
func (c *Customer) FByMap(m typing.SQLMaps, preloadandjoin ...bool) error {
	query, args := dal.MapsToWQuery(m)
	// if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
	// 	return dal.Dal.BlacheDB.NewRaw(`SELECT * FROM org_profiles LEFT JOIN configs ON org_profiles.id = configs.product_id WHERE `+query, args...).Scan(context.Background(), p.Fields()...)
	// }
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT org_profiles.id, org_profiles.org_id, org_profiles.support_email, org_profiles.country, org_profiles.default_currency, org_profiles.company_name, org_profiles.company_name_slug, org_profiles.logo, org_profiles.website, org_profiles.created_at, org_profiles.updated_at FROM org_profiles WHERE `+query, args...).Scan(context.Background(), c)
	}
	return dal.Dal.BlacheDB.NewRaw(`SELECT org_id, name FROM org_profiles WHERE `+query, args...).Scan(context.Background(), c)
}

/*
FByMap finds and returns all customers matching the key/value pairs provided in the map

# By default, only the id and name fields are loaded

The	"preloadandjoin" parameter can be used to request that all the fields of the struct be loaded

It returns an error if any
*/
func (c *Customers) FByMap(m typing.SQLMaps, limit, offset int, sort string, preloadandjoin ...bool) error {
	query, args := dal.MapsToWQuery(m)
	// if len(preloadandjoin) > 1 && preloadandjoin[0] && preloadandjoin[1] {
	// 	rows, err := dal.Dal.BlacheDB.QueryContext(context.Background(), `SELECT * FROM org_profiles LEFT JOIN configs ON org_profiles.id = configs.product_id WHERE `+query+` ORDER BY org_profiles.updated_at `+sort+` LIMIT ? OFFSET ?`, append(args, limit, offset)...)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	defer rows.Close()
	// 	for rows.Next() {
	// 		var user User
	// 		if err := rows.Scan(user.Fields()...); err != nil {
	// 			return err
	// 		}
	// 		*u = append(*u, user)
	// 	}
	// }
	if len(preloadandjoin) > 0 && preloadandjoin[0] {
		return dal.Dal.BlacheDB.NewRaw(`SELECT org_profiles.id, org_profiles.org_id, org_profiles.support_email, org_profiles.country, org_profiles.default_currency, org_profiles.company_name, org_profiles.company_name_slug, org_profiles.logo, org_profiles.website, org_profiles.created_at, org_profiles.updated_at FROM org_profiles WHERE `+query+` ORDER BY org_profiles.updated_at `+sort+` LIMIT ? OFFSET ?`, append(args, limit, offset)...).Scan(context.Background(), c)
	}
	return dal.Dal.BlacheDB.NewRaw(`SELECT org_id, name FROM org_profiles WHERE `+query+` ORDER BY org_profiles.updated_at `+sort+` LIMIT ? OFFSET ?`, append(args, limit, offset)...).Scan(context.Background(), c)
}

/*
CByMap finds and counts all customers matching the key/value pairs provided in the map

It returns an error if any
*/
func (s *Customers) CByMap(m typing.SQLMaps, preloadandjoin ...bool) (int, error) {
	var count int
	query, args := dal.MapsToWQuery(m)
	err := dal.Dal.BlacheDB.NewRaw(`SELECT count(*) FROM org_profiles WHERE `+query, args...).Scan(context.Background(), &count)
	return count, err
}
