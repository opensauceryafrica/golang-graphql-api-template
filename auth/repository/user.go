package repository

import (
	"cendit.io/auth/schema"
	"cendit.io/garage/database"
	"cendit.io/garage/xiao"
)

var User = func() *xiao.Xiao[schema.User] {
	x := xiao.NewXiao[schema.User](database.DB.DB)
	x.TableName = "users"
	x.Preloaders = []string{"id", "firstname", "middlename", "lastname", "email", "phone", "phone_code"}
	return x
}
