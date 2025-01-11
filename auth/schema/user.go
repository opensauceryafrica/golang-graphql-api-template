package schema

import (
	"database/sql"

	"cendit.io/garage/primer/enum"
	"cendit.io/garage/schema"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users" rsf:"false"`

	ID         string       `bun:"id" json:"id"`
	Firstname  string       `bun:"firstname" json:"firstname"`
	Middlename string       `bun:"middlename" json:"middlename"`
	Lastname   string       `bun:"lastname" json:"lastname"`
	Email      string       `bun:"email" json:"email"`
	Phone      string       `bun:"phone" json:"phone"`
	PhoneCode  string       `bun:"phone_code" json:"phone_code"`
	Gender     enum.Gender  `bun:"gender" json:"gender"`
	Country    string       `bun:"country" json:"country"`
	Username   string       `bun:"username" json:"username"`
	Tier       enum.Tier    `bun:"tier" json:"tier"`
	DOB        sql.NullTime `bun:"dob" json:"dob" rsfr:"false"`
	Passcode   string       `bun:"passcode" json:"passcode"`
	Pin        string       `bun:"pin" json:"pin"`
	Referrer   string       `bun:"referrer" json:"referrer"`
	Referral   string       `bun:"referral" json:"referral"`
	Avatar     string       `bun:"avatar" json:"avatar"`

	// defines whether or not the user is an internal team member
	Internal bool            `bun:"internal" json:"internal"`
	Role     enum.Role       `bun:"role" json:"role"`
	Status   enum.UserStatus `bun:"status" json:"status"`

	schema.Datum

	Language enum.Language `bun:"language" json:"language"`
}

type Users []User
