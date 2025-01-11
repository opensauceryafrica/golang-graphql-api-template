package repository

import (
	"cendit.io/auth/schema"
	"cendit.io/garage/database"
	"cendit.io/garage/xiao"
)

var Wallet = func() *xiao.Xiao[schema.Wallet] {
	x := xiao.NewXiao[schema.Wallet](database.DB.DB)
	x.TableName = "wallets"
	x.Preloaders = []string{"id", "user_id", "balance", "currency"}
	return x
}
