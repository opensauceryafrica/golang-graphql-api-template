package repository

import (
	"cendit.io/auth/schema"

	"cendit.io/garage/database"
	"cendit.io/garage/xiao"
)

var SecuritySetting = func() *xiao.Xiao[schema.SecuritySetting] {
	x := xiao.NewXiao[schema.SecuritySetting](database.DB.DB)
	x.TableName = "security_settings"
	x.Preloaders = []string{"id", "user_id", "sms_alert"}
	return x
}
