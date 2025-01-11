package schema

import (
	"cendit.io/garage/schema"
	"github.com/uptrace/bun"
)

type SecuritySetting struct {
	bun.BaseModel `bun:"table:security_settings" rsf:"false"`

	ID               string `bun:"id" json:"id"`
	UserID           string `bun:"user_id" json:"user_id"`
	BiometricEnabled bool   `bun:"biometric_enabled" json:"biometric_enabled"`
	PrivateMode      bool   `bun:"private_mode" json:"private_mode"`
	AllowScreenshot  bool   `bun:"allow_screenshot" json:"allow_screenshot"`
	ContextMenu      bool   `bun:"context_menu" json:"context_menu"` //Transact from outside the app

	SmsAlert bool `bun:"sms_alert" json:"sms_alert"`

	schema.Datum
}
