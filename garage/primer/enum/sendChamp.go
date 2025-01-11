package enum

// SendChampRoute represents the possible values for the route
type SendChampRoute string

const (
	DND           SendChampRoute = "dnd"
	NonDND        SendChampRoute = "non_dnd"
	International SendChampRoute = "international"
)

// TokenType represents the possible types for the otp
type TokenType string

const (
	Numeric      TokenType = "numeric"
	Alphanumeric TokenType = "alphanumeric"
)
