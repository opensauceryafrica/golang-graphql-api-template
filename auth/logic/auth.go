package logic

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"cendit.io/gate/http/graphql/exception"
	"cendit.io/gate/http/graphql/model"

	"cendit.io/auth/repository"
	"cendit.io/auth/schema"
	"cendit.io/garage/function"
	"cendit.io/garage/primer/enum"
	"cendit.io/garage/primitive"
	"cendit.io/garage/redis"
	"cendit.io/garage/xiao"
)

// function to create user profile
func SignupStepOne(input model.StepOneInput) (*schema.User, error) {

	ptx, err := repository.User().Tx(context.Background())
	if err != nil {
		return nil, err
	}
	defer ptx.Rollback()

	u := schema.User{}

	userExist, err := repository.User().Exists(context.Background(), "email", strings.ToLower(input.Email))
	if err != nil {
		return nil, err
	}
	if userExist {
		return nil, exception.MakeError("Email already in use!", 400)
	}

	// if referral code is provided, check if it exists
	if input.ReferralCode != nil {
		referral, err := repository.User().FindByKeyVal(context.Background(), "code", *input.ReferralCode, true)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, exception.MakeError(fmt.Sprintf(`Referral with code "%s" not found!`, *input.ReferralCode), 404)
			}
			return nil, err
		}
		if referral.ID == "" {
			return nil, exception.MakeError(fmt.Sprintf(`Referral with code "%s" not found!`, *input.ReferralCode), 404)
		}
	}

	// create user
	function.Parse(input, &u)

	// generate unique user id
	u.ID = function.GenerateUUID()

	u.Date()
	u.Email = strings.ToLower(u.Email)
	u.Role = enum.User
	u.Tier = enum.TierZero
	u.Status = enum.UserActive

	if err := repository.User().CreateTx(context.Background(), ptx, xiao.SQLMaps{
		IMaps: []xiao.SQLMap{
			{
				Map: map[string]interface{}{
					"id":         u.ID,
					"firstname":  u.Firstname,
					"middlename": u.Middlename,
					"lastname":   u.Lastname,
					"email":      u.Email,
					"phone":      u.Phone,
					"phone_code": u.PhoneCode,
					"gender":     u.Gender,
					"country":    u.Country,
					"username":   u.Username,
					"tier":       u.Tier,
					"dob":        u.DOB,
					"passcode":   u.Passcode,
					"pin":        u.Pin,
					"referrer":   u.Referrer,
					"referral":   u.Referral,
					"avatar":     u.Avatar,
					"internal":   u.Internal,
					"role":       u.Role,
					"status":     u.Status,
					"created_at": u.CreatedAt,
					"updated_at": u.UpdatedAt,
					"language":   enum.English,
				},
			},
		},
	}); err != nil {
		return nil, err
	}

	// Get user's local currency
	currency, code, err := function.GetCurrencyInfo(input.Country)
	if err != nil {
		return nil, exception.MakeError(fmt.Sprintf(`Invalid country code "%s"!`, input.Country), 400)
	}

	go func(currency, code, userID string) {

		// add FIAT wallet
		AddWallet(model.WalletInput{
			Currency: model.ECurrency(currency),
			Code:     &code,
		}, schema.User{ID: userID})

		// add CRYPTO wallet
		currencies := []enum.Currency{enum.BTC, enum.USDT}

		for _, c := range currencies {
			AddWallet(model.WalletInput{
				Currency: model.ECurrency(c),
			}, schema.User{ID: userID})
		}

	}(currency, code, u.ID)

	if err := ptx.Commit(); err != nil {
		return nil, err
	}

	return &u, nil
}

// function for user second signup stage (Personal info entry)
func SignupStepTwo(input model.StepTwoInput, user schema.User) (*schema.User, error) {

	ptx, err := repository.User().Tx(context.Background())
	if err != nil {
		return nil, err
	}
	defer ptx.Rollback()

	query := xiao.SQLMaps{
		WMaps: []xiao.SQLMap{
			{
				Map: map[string]interface{}{
					"email": user.Email,
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

	// find the user by email
	u, err := repository.User().FindByMap(context.Background(), query, true)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.MakeError(fmt.Sprintf(`User with email "%s" not found!`, user.Email), 404)
		}
		return nil, err
	}
	if u.ID == "" {
		return nil, exception.MakeError(fmt.Sprintf(`User with email "%s" not found!`, user.Email), 404)
	}

	// TODO: age verification using DOB
	diff := time.Since(u.DOB.Time)
	age := diff.Hours() / (24 * 365.25)
	if age < 16 {
		return nil, exception.MakeError("User must be at least 16 years old!", 400)
	}

	// update the Set Map with the provided personal information
	function.LayerMap(function.StructToMapOfNonNils(&input, "json", primitive.Array{"firstname", "lastname", "middlename", "gender", "dob"}, map[string]string{}), query.SMap.Map)

	// update user with the new personal information
	us, err := repository.User().UpdateByMapTx(context.Background(), ptx, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.MakeError(fmt.Sprintf(`User with email "%s" not found!`, user.Email), 404)
		}
		return nil, err
	}

	if err := ptx.Commit(); err != nil {
		return nil, err
	}

	if len(us) > 0 {
		u = &us[0]
	}

	return u, nil
}

// function for user signup third step (Address entry)
func SignupStepThree(input model.StepThreeInput, user schema.User) (*schema.Address, error) {

	ptx, err := repository.Address().Tx(context.Background())
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

	// find the address by user ID
	a, err := repository.Address().FindByMap(context.Background(), query, true)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	// if address exists and is verified, return an error
	if a.Verified {
		return nil, exception.MakeError("Address already verified and cannot be updated!", 400)
	}

	var as []schema.Address

	// if address doesn't exist, create a new one
	if a.ID == "" {
		// create user
		function.Parse(input, &a)

		// generate unique user id
		a.ID = function.GenerateUUID()

		a.Verified = false
		a.UserID = user.ID
		a.Date()

		err := repository.Address().CreateTx(context.Background(), ptx, xiao.SQLMaps{
			IMaps: []xiao.SQLMap{
				{
					Map: map[string]interface{}{
						"id":           a.ID,
						"user_id":      a.UserID,
						"state":        a.State,
						"city":         a.City,
						"house_number": a.HouseNumber,
						"street":       a.Street,
						"zip":          a.Zip,
						"verified":     a.Verified,
						"created_at":   a.CreatedAt,
						"updated_at":   a.UpdatedAt,
					},
				},
			},
		})
		if err != nil {
			return nil, err
		}

		goto commit

	}

	// update existing address with new data
	function.LayerMap(function.StructToMapOfNonNils(&input, "json", primitive.Array{"state", "city", "house_number", "street", "zip"}, map[string]string{}), query.SMap.Map)

	as, err = repository.Address().UpdateByMapTx(context.Background(), ptx, query)
	if err != nil {
		return nil, err
	}

commit:
	if err := ptx.Commit(); err != nil {
		return nil, err
	}

	if len(as) > 0 {
		a = &as[0]
	}

	return a, nil
}

// function for User Signin
func Signin(input model.SigninInput) (*schema.User, error) {

	val, err := redis.Ral.Get(fmt.Sprintf("%s-%s", strings.ToLower(input.Email), input.Code))
	if err != nil || val == nil {
		return nil, exception.MakeError("Code is invalid or expired!", 401)
	}

	err = redis.Ral.Del(fmt.Sprintf("%s-%s", strings.ToLower(input.Email), input.Code))
	if err != nil {
		return nil, exception.MakeError("Code is invalid or expired!", 401)
	}

	// find User
	u, err := repository.User().FindByKeyVal(context.Background(), "email", strings.ToLower(input.Email), true)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.MakeError(fmt.Sprintf(`Account with email "%s" not found!`, input.Email), 404)
		}
		return nil, err
	}
	if u.ID == "" {
		return nil, exception.MakeError(`Account not found!`, 404)
	}

	if u.Status != enum.UserActive {
		return nil, exception.MakeError(`Account not active! Please contact support.`, 403)
	}

	return u, nil
}

// function to Set Passcode
func SetPasscode(input model.SetPasscodeInput, user schema.User) (*schema.User, error) {

	ptx, err := repository.User().Tx(context.Background())
	if err != nil {
		return nil, err
	}
	defer ptx.Rollback()

	query := xiao.SQLMaps{
		WMaps: []xiao.SQLMap{
			{
				Map: map[string]interface{}{
					"email": user.Email,
				},
				JoinOperator:       xiao.And,
				ComparisonOperator: xiao.Equal,
			},
		},
		RMap: xiao.SQLMap{
			Map: map[string]interface{}{
				"*": nil,
			},
		},
		WJoinOperator: xiao.And,
	}

	// Find the user by email
	u, err := repository.User().FindAndLockByMap(context.Background(), ptx, query, true)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.MakeError(fmt.Sprintf(`Account with email "%s" not found!`, user.Email), 404)
		}
		return nil, err
	}
	if u.ID == "" {
		return nil, exception.MakeError(fmt.Sprintf(`Account with email "%s" not found!`, user.Email), 404)
	}

	if u.Passcode != "" {

		if input.Pastcode == nil {
			return nil, exception.MakeError("Please provide your previous passcode!", 400)
		}

		valid := function.ComparePasscode(u.Passcode, *input.Pastcode)
		if !valid {
			return nil, exception.MakeError("Passcode incorrect!", 400)
		}
	}

	passcode, _ := function.HashPasscode(input.Passcode)
	query.SMap = xiao.SQLMap{
		Map: map[string]interface{}{
			"updated_at": "now()",
			"passcode":   passcode,
		},
		JoinOperator:       xiao.Comma,
		ComparisonOperator: xiao.Equal,
	}

	// Update user with the passcode
	us, err := repository.User().UpdateByMapTx(context.Background(), ptx, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.MakeError(fmt.Sprintf(`Account with email "%s" not found!`, user.Email), 404)
		}
		return nil, err
	}

	if err := ptx.Commit(); err != nil {
		return nil, err
	}

	if len(us) > 0 {
		u = &us[0]
	}

	return u, nil
}

// function to StartSession
func StartSession(input model.StartSessionInput) (*schema.User, error) {
	// Find User
	u, err := repository.User().FindByKeyVal(context.Background(), "email", strings.ToLower(input.Email), true)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.MakeError(fmt.Sprintf(`Account with email "%s" not found!`, input.Email), 404)
		}
		return nil, err
	}

	valid := function.ComparePasscode(u.Passcode, input.Passcode)
	if !valid {
		return nil, exception.MakeError("Passcode incorrect!", 403)
	}

	if u.ID == "" {
		return nil, exception.MakeError(`Account not found!`, 404)
	}
	if u.Status != enum.UserActive {
		return nil, exception.MakeError(`Account not active! Please contact support.`, 403)
	}

	return u, nil
}
