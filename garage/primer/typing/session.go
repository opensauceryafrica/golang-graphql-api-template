package typing

import "cendit.io/garage/primer/enum"

type Session struct {
	Email string `json:"Email"`
	ID    string `json:"ID"`

	Role enum.Role
}
