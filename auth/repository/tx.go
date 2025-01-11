package repository

import (
	"cendit.io/auth/schema"
	"cendit.io/garage/database"
	"cendit.io/garage/xiao"
)

var Transaction = func() *xiao.Xiao[schema.Transaction] {
	x := xiao.NewXiao[schema.Transaction](database.DB.DB)
	x.TableName = "transactions"
	x.Preloaders = []string{"id", "user_id", "amount", "currency"}
	return x
}
