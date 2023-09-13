package redis

import "blacheapi/primer/enum"

type Session struct {
	Email string `json:"Email"`
	OrgID int    `json:"OrgID"`
	ID    int    `json:"ID"`

	Role enum.Role
}
