package sms

import (
	"cendit.io/garage/primer/typing"
)

var Client typing.SMSClient

// NewClient initializes the Client with the publicKey and mode
func NewClient(publicKey, baseURL, sender, apiKey string) typing.SMSClient {
	Client = &SendChamp{
		PublicKey: publicKey,
		Mode:      "live",
		BaseURL:   baseURL,
		Sender:    sender,
		APIKey:    apiKey,
	}
	return Client
}
