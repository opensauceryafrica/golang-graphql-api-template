package typing

type SendPulseResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	ErrorCode   int    `json:"error_code"`
	Message     string `json:"message"`
}
