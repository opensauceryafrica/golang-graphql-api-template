package repository

import (
	"cendit.io/auth/schema"
	"cendit.io/garage/database"
	"cendit.io/garage/xiao"
)

var KYC = func() *xiao.Xiao[schema.KYC] {
	x := xiao.NewXiao[schema.KYC](database.DB.DB)
	x.TableName = "kycs"
	x.Preloaders = []string{"id", "user_id", "bvn"}
	return x
}
