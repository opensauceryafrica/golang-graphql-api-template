package logic

import (
	"context"
	"database/sql"

	"cendit.io/gate/http/graphql/model"

	"cendit.io/auth/repository"
	"cendit.io/auth/schema"
	"cendit.io/garage/function"
	"cendit.io/garage/primitive"
	"cendit.io/garage/xiao"
)

// function for user to modify security setting
func ModifySecuritySetting(input model.SecuritySettingInput, user schema.User) (*schema.SecuritySetting, error) {

	ptx, err := repository.SecuritySetting().Tx(context.Background())
	if err != nil {
		return nil, err
	}
	defer ptx.Rollback()

	query := xiao.SQLMaps{
		WMaps: []xiao.SQLMap{
			{
				Map: map[string]interface{}{
					"user_id": user.ID,
				},
				JoinOperator:       xiao.And,
				ComparisonOperator: xiao.Equal,
			},
		},
		SMap: xiao.SQLMap{
			Map: map[string]interface{}{
				"updated_at": "now()",
			},
			JoinOperator:       xiao.Comma,
			ComparisonOperator: xiao.Equal,
		},
		RMap: xiao.SQLMap{
			Map: map[string]interface{}{
				"*": nil,
			},
		},
		WJoinOperator: xiao.And,
	}

	// find the setting by user ID
	s, err := repository.SecuritySetting().FindByMap(context.Background(), query, true)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	var ss []schema.SecuritySetting

	// if setting doesn't exist, create a new one
	if s.ID == "" {
		// create user
		function.Parse(input, &s)

		// generate unique user id
		s.ID = function.GenerateUUID()

		s.UserID = user.ID
		s.Date()

		err := repository.SecuritySetting().CreateTx(context.Background(), ptx, xiao.SQLMaps{
			IMaps: []xiao.SQLMap{
				{
					Map: map[string]interface{}{
						"id":                s.ID,
						"user_id":           s.UserID,
						"biometric_enabled": s.BiometricEnabled,
						"private_mode":      s.PrivateMode,
						"allow_screenshot":  s.AllowScreenshot,
						"context_menu":      s.ContextMenu,
						"sms_alert":         s.SmsAlert,
						"created_at":        s.CreatedAt,
						"updated_at":        s.UpdatedAt,
					},
				},
			},
		})
		if err != nil {
			return nil, err
		}

		goto commit

	}

	// update existing setting with new data
	function.LayerMap(function.StructToMapOfNonNils(&input, "json", primitive.Array{"biometric_enabled", "private_mode", "allow_screenshot", "sms_alert", "context_menu"}, map[string]string{}), query.SMap.Map)

	ss, err = repository.SecuritySetting().UpdateByMapTx(context.Background(), ptx, query)
	if err != nil {
		return nil, err
	}

commit:
	if err := ptx.Commit(); err != nil {
		return nil, err
	}

	if len(ss) > 0 {
		s = &ss[0]
	}

	return s, nil
}
