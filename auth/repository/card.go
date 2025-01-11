package repository

import (
	"cendit.io/auth/schema"
	"cendit.io/garage/database"
	"cendit.io/garage/xiao"
)

var Card = func() *xiao.Xiao[schema.Card] {
	x := xiao.NewXiao[schema.Card](database.DB.DB)
	x.TableName = "cards"
	x.Preloaders = []string{"id", "user_id", "token"}
	return x
}
