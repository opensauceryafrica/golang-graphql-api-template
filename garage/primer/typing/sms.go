package typing

import "cendit.io/garage/primer/enum"

type SMSClient interface {
	//the second parameter route can be set to enum.DND if unsure
	Authenticate() error
	SendSMS(payload SMS, route enum.SendChampRoute) error
	SendOTP(phone string, tokenType enum.TokenType) (string, error)
	ConfirmOTP(payload SendChampOtpPayload) (string, error)
}

type SMS struct {
	// From    string   `json:"from"`
	To      []string `json:"to"`
	Message string   `json:"message"`
}

type SendOtpResponse struct {
	Reference string `json:"reference"`
	Message   string `json:"message"`
}
