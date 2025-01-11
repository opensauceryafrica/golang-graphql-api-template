package typing

import "cendit.io/garage/primer/enum"

type SendChampResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    data   `json:"data"`
}

type data struct {
	Reference     string `json:"reference"`
	Status        string `json:"status"`
	Phone         string `json:"phone"`
	WalletBalance string `json:"wallet_balance"`
}

// SendChampRequest represents the structure of the API request body
type SendChampRequest struct {
	Message    string              `json:"message"`
	SenderName string              `json:"sender_name"`
	Route      enum.SendChampRoute `json:"route"`
	To         []string            `json:"to"`
}

type SendChampOtpPayload struct {
	Reference string `json:"reference"`
	Code      string `json:"code"`
}
