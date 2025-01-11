package repository

import (
	"cendit.io/auth/schema"
	"cendit.io/garage/database"
	"cendit.io/garage/xiao"
)

var Address = func() *xiao.Xiao[schema.Address] {
	x := xiao.NewXiao[schema.Address](database.DB.DB)
	x.TableName = "addresses"
	x.Preloaders = []string{"id", "state", "city"}
	return x
}
