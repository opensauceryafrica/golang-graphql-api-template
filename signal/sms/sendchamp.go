package sms

import (
	"fmt"
	"net/http"
	"strconv"

	"cendit.io/garage/primer/enum"
	"cendit.io/garage/primer/typing"
	"github.com/opensaucerer/goaxios"
)

type SendChamp struct {
	PublicKey string
	Mode      string
	BaseURL   string
	Sender    string
	APIKey    string
}

// SendSMS sends an SMS with the SMS payload
func (sc *SendChamp) Authenticate() error {
	req := goaxios.GoAxios{
		Url:    fmt.Sprintf("%s/wallet/wallet_balance", sc.BaseURL),
		Method: http.MethodPost,
		Headers: map[string]string{
			"accept":        "application/json",
			"content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", sc.APIKey),
		},
		Body:           map[string]any{},
		ResponseStruct: &typing.SendChampResponse{},
	}

	res := req.RunRest()
	if res.Error != nil {
		return res.Error
	}

	if res.Body.(*typing.SendChampResponse).Code != 200 {
		return fmt.Errorf(res.Body.(*typing.SendChampResponse).Message)
	}

	// Check if available balance is lower than the cost of sending and confirming otp for 100 people
	balance := res.Body.(*typing.SendChampResponse).Data.WalletBalance
	availableBalance, _ := strconv.ParseFloat(balance, 64)

	thresholdBalance := 100 * 2.8
	if availableBalance < thresholdBalance {
		return fmt.Errorf("SendChamp wallet balance is lower than threshold")
	}

	return nil
}

// SendSMS sends an SMS with the SMS payload
func (sc *SendChamp) SendSMS(payload typing.SMS, route enum.SendChampRoute) error {
	req := goaxios.GoAxios{
		Url:    fmt.Sprintf("%s/sms/send", sc.BaseURL),
		Method: http.MethodPost,
		Headers: map[string]string{
			"accept":        "application/json",
			"content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", sc.APIKey),
		},
		Body: map[string]any{
			"to":          payload.To,
			"sender_name": sc.Sender,
			"route":       route,
			"message":     payload.Message,
		},
		ResponseStruct: &typing.SendChampResponse{},
	}

	res := req.RunRest()
	if res.Error != nil {
		return res.Error
	}

	if res.Body.(*typing.SendChampResponse).Code != 200 {
		return fmt.Errorf(res.Body.(*typing.SendChampResponse).Message)
	}

	return nil
}

func (sc *SendChamp) SendOTP(phone string, tokenType enum.TokenType) (string, error) {
	req := goaxios.GoAxios{
		Url:    fmt.Sprintf("%s/verification/create", sc.BaseURL),
		Method: http.MethodPost,
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", sc.APIKey),
			"Accept":        "application/json,text/plain,*/*",
			"Content-Type":  "application/json",
		},
		Body: map[string]any{
			"channel":                "sms",
			"sender":                 sc.Sender,
			"token_type":             tokenType,
			"token_length":           6,
			"expiration_time":        5,
			"customer_mobile_number": phone,
		},
		ResponseStruct: &typing.SendChampResponse{},
	}

	res := req.RunRest()
	if res.Error != nil {
		return "", res.Error
	}

	if res.Body.(*typing.SendChampResponse).Code != 200 {
		return "", fmt.Errorf(res.Body.(*typing.SendChampResponse).Message)
	}

	reference := res.Body.(*typing.SendChampResponse).Data.Reference
	return reference, nil
}

func (sc *SendChamp) ConfirmOTP(payload typing.SendChampOtpPayload) (string, error) {
	req := goaxios.GoAxios{
		Url:    fmt.Sprintf("%s/verification/confirm", sc.BaseURL),
		Method: http.MethodPost,
		Headers: map[string]string{
			"accept":        "application/json",
			"content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", sc.APIKey),
		},
		Body: map[string]any{
			"verification_reference": payload.Reference,
			"verification_code":      payload.Code,
		},
		ResponseStruct: &typing.SendChampResponse{},
	}

	res := req.RunRest()
	if res.Error != nil {
		return "", res.Error
	}

	if res.Body.(*typing.SendChampResponse).Code != 200 {
		return "", fmt.Errorf(res.Body.(*typing.SendChampResponse).Message)
	}
	phone := res.Body.(*typing.SendChampResponse).Data.Phone
	return phone, nil
}
