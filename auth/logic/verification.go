package logic

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"cendit.io/auth/repository"
	"cendit.io/auth/schema"
	"cendit.io/garage/config"
	"cendit.io/garage/function"
	"cendit.io/garage/primer/constant"
	"cendit.io/garage/primer/enum"
	"cendit.io/garage/primer/typing"
	"cendit.io/garage/primitive"
	"cendit.io/garage/redis"
	"cendit.io/garage/xiao"
	"cendit.io/gate/http/graphql/exception"
	"cendit.io/gate/http/graphql/model"
	"cendit.io/signal/email"
	"cendit.io/signal/sms"
)

// Function to verify email
func EmailVerification(input model.EmailVerificationInput) (string, error) {

	code := function.GenerateRandomNumber(6)
	redis.Ral.Set(fmt.Sprintf("%s-%s", strings.ToLower(input.Email), code), code, time.Minute*5)

	values := map[string]string{
		"code": code,
	}

	env := config.Environment().Email

	body := typing.Email{
		To:            []map[string]string{{"name": "", "email": input.Email}},
		Subject:       constant.VerificationEmailSubject,
		HTML:          constant.VerificationEmailHtml,
		From:          map[string]string{"name": env.Name, "email": env.From},
		AutoPlainText: false,
		Values:        values,
	}

	err := email.Client.SendEmail(body)
	if err != nil {
		return "", err
	}

	return "Verification code sent", nil
}

// Function to confirm email
func EmailConfirmation(input model.EmailConfirmationInput) (string, error) {

	val, err := redis.Ral.Get(fmt.Sprintf("%s-%s", strings.ToLower(input.Email), input.Code))
	if err != nil || val == nil {
		return "", exception.MakeError("Code is invalid or expired!", 400)
	}

	err = redis.Ral.Del(fmt.Sprintf("%s-%s", strings.ToLower(input.Email), input.Code))
	if err != nil {
		return "", exception.MakeError("Code is invalid or expired!", 400)
	}

	return "Email verification successful", nil
}

// Function to verify phone number
func PhoneVerification(input model.PhoneVerificationInput) (string, error) {
	reference, err := sms.Client.SendOTP(input.Phone, enum.Numeric)
	if err != nil {
		return "", exception.MakeError(err.Error(), 400)
	}

	return reference, nil
}

// Function to confirm phone number
func PhoneConfirmation(input model.PhoneConfirmationInput, user schema.User) (*schema.User, error) {
	payload := typing.SendChampOtpPayload{Reference: input.Reference, Code: input.Code}
	phone, err := sms.Client.ConfirmOTP(payload)
	if err != nil {
		return nil, exception.MakeError(err.Error(), 400)
	}
	if phone != input.Phone {
		return nil, exception.MakeError("Code is invalid or expired", 400)
	}

	query := xiao.SQLMaps{
		WMaps: []xiao.SQLMap{
			{
				Map: map[string]interface{}{
					"id": user.ID,
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

	// find the user by id
	u := schema.User{}

	// update the Set Map with the provided phone
	function.LayerMap(function.StructToMapOfNonNils(&input, "json", primitive.Array{"phone"}, map[string]string{}), query.SMap.Map)

	// update user with the phone
	us, err := repository.User().UpdateByMap(context.Background(), query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.MakeError(fmt.Sprintf(`User with email "%s" not found!`, user.Email), 404)
		}
		return nil, err
	}

	if len(us) > 0 {
		u = us[0]
	}

	return &u, nil
}
