package typing

type Email struct {
	HTML          string              `json:"html,omitempty"`
	Text          string              `json:"text,omitempty"`
	Subject       string              `json:"subject"`
	From          map[string]string   `json:"from"`
	To            []map[string]string `json:"to"`
	AutoPlainText bool                `json:"auto_plain_text"`
	Values        map[string]string
}

type EmailClient interface {
	Authenticate() error
	SendEmail(payload Email) error
}
