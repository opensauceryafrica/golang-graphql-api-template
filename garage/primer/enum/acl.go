package enum

type Role string

// Role describes the role of a user
const (
	// Customer
	User Role = "USER"

	// Organization
	Org Role = "ORG"
)

type UserStatus string

const (
	UserLocked  UserStatus = "USER_LOCKED"
	UserDeleted UserStatus = "USER_DELETED"
	UserActive  UserStatus = "USER_ACTIVE"
)
