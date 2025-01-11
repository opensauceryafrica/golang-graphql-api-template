package email

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"cendit.io/garage/primer/typing"
	"github.com/aymerick/raymond"
	"github.com/opensaucerer/goaxios"
)

type SendPulse struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
	TokenExpiry  time.Time
	BaseURL      string
}

// authenticate gets an Access token if there isn't a valid one
func (sp *SendPulse) Authenticate() error {
	if sp.AccessToken != "" && time.Now().Before(sp.TokenExpiry) {
		return nil
	}

	req := goaxios.GoAxios{
		Url:    fmt.Sprintf("%s/oauth/access_token", sp.BaseURL),
		Method: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     sp.ClientID,
			"client_secret": sp.ClientSecret,
		},
		ResponseStruct: &typing.SendPulseResponse{},
	}

	res := req.RunRest()
	if res.Error != nil {
		return res.Error
	}
	authResp, _ := res.Body.(*typing.SendPulseResponse)

	if authResp.AccessToken == "" {
		return fmt.Errorf("unable to authenticate email client")
	}

	sp.AccessToken = authResp.AccessToken
	sp.TokenExpiry = time.Now().Add(time.Duration(authResp.ExpiresIn) * time.Second)
	return nil
}

// SendEmail sends an email with the email payload
func (sp *SendPulse) SendEmail(payload typing.Email) error {
	if err := sp.Authenticate(); err != nil {
		return err
	}

	result, err := raymond.Render(payload.HTML, payload.Values)
	if err != nil {
		return err
	}

	req := goaxios.GoAxios{
		Url:    fmt.Sprintf("%s/smtp/emails", sp.BaseURL),
		Method: http.MethodPost,
		Headers: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %v", sp.AccessToken),
		},
		Body: map[string]map[string]any{
			"email": {
				"to":              payload.To,
				"subject":         payload.Subject,
				"html":            base64.StdEncoding.EncodeToString([]byte(result)),
				"from":            payload.From,
				"auto_plain_text": payload.AutoPlainText,
			},
		},
		ResponseStruct: &typing.SendPulseResponse{},
	}

	res := req.RunRest()
	if res.Error != nil {
		return res.Error
	}

	if res.Body.(*typing.SendPulseResponse).ErrorCode != 0 {
		return fmt.Errorf(res.Body.(*typing.SendPulseResponse).Message)
	}

	return nil
}
