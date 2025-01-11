package schema

import (
	"time"

	"cendit.io/garage/primer/enum"
)

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Pages int `json:"pages"`
}

type History struct {
	Act string    `json:"act"`
	By  int       `json:"by"`
	At  time.Time `json:"at"`
}

// NextOfKin struct
type NextOfKin struct {
	FirstName    string `json:"first_name"`   // Not required
	LastName     string `json:"last_name"`    // Not required
	Email        string `json:"email"`        // Not required
	Address      string `json:"address"`      // Not required
	Phone        string `json:"phone"`        // Not required
	PhoneCode    string `json:"phonecode"`    // Not required
	Relationship string `json:"relationship"` // Not required
}

type Attachment struct {

	// the name of the file
	Name string `json:"name"`

	// the url of the file
	URL string `json:"url"`
}

type KYCTier struct {
	Zero KYCLimit `json:"zero" bun:"zero"`
	One  KYCLimit `json:"one" bun:"one"`
	Two  KYCLimit `json:"two" bun:"two"`
}

type KYCLimit struct {
	Transfer                  bool             `json:"transfer" bun:"transfer"`
	MaxExternalTransferPerTx  float64          `json:"max_external_transfer_per_tx" bun:"max_external_transfer_per_tx"`
	MaxExternalTransferPerDay float64          `json:"max_external_transfer_per_day" bun:"max_external_transfer_per_day"`
	Bill                      bool             `json:"bill" bun:"bill"`
	KYC                       []enum.KYCOption `json:"kyc" bun:"kyc"`
	USDSpend                  float64          `json:"usd_spend" bun:"usd_spend"`
	NGNSpend                  float64          `json:"ngn_spend" bun:"ngn_spend"`
}

type PhysicalNairaCard struct {
	Withdrawal          float64 `json:"withdrawal"`           // Integer field
	WithdrawalFrequency int     `json:"withdrawal_frequency"` // Integer field
	POSPayment          float64 `json:"pos_payment"`          // Integer field
}

type Limit struct {
	Transfer                  *bool
	MaxExternalTransferPerTx  *float64
	MaxExternalTransferPerDay *float64
	Bill                      *bool

	USDSpend           *float64 `json:"usd_spend"` // Integer field
	NGNSpend           *float64 `json:"ngn_spend"`
	DailyTransferLimit *float64 `json:"daily_transfer_limit"` // Integer field, Modifiable by the user but bounded by the sending external limit of the tier
	FreeTransfer       int      `json:"free_transfer"`        // Integer field
	LastFreeTransfer   time.Time
	PhysicalNairaCard  PhysicalNairaCard `json:"physical_naira_card"` // Nested struct
}
