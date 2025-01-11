package email

import (
	"cendit.io/garage/primer/typing"
)

var Client typing.EmailClient

// NewEmailClient initializes the Client with the clientID and clientSecret
func NewClient(clientID, clientSecret string, baseURL string) typing.EmailClient {
	Client = &SendPulse{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		BaseURL:      baseURL,
	}
	return Client
}
