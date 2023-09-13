package enum

type Currency string

// Currency describes the currency of a transaction
const (
	// USD
	USD Currency = "USD"

	// EUR
	EUR Currency = "EUR"

	// GBP
	GBP Currency = "GBP"

	// NGN
	NGN Currency = "NGN"
)

type PaymentMethod string

// PaymentMethod describes the payment method of a transaction
const (
	// Card
	Card PaymentMethod = "CARD"

	// BankTransfer
	BankTransfer PaymentMethod = "BANK_TRANSFER"

	// Wallet
	Wallet PaymentMethod = "WALLET"

	// 	Manual
	Cash PaymentMethod = "CASH"
)

type TransactionType string

// TransactionType describes the type of a transaction
const (
	// FUNDING
	FUNDING TransactionType = "FUNDING"

	// TRANSFER
	TRANSFER TransactionType = "TRANSFER"

	// WITHDRAWAL
	WITHDRAWAL TransactionType = "WITHDRAWAL"

	// INTEREST
	INTEREST TransactionType = "INTEREST"
)

type PaymentGateway string

const (
	// Flutterwave
	Flutterwave PaymentGateway = "FLUTTERWAVE"

	// Paystack
	Paystack PaymentGateway = "PAYSTACK"

	// MANUAL
	MANUAL PaymentGateway = "MANUAL"
)

var GatewayMethodChannel = map[PaymentGateway]map[PaymentMethod]interface{}{
	Flutterwave: {
		BankTransfer: "banktransfer",
		Card:         "card",
	},
	Paystack: {
		BankTransfer: []string{"bank_transfer"},
		Card:         []string{"card"},
	},
	PaymentGateway(Wallet): {
		BankTransfer: "wallet",
		Card:         "wallet",
	},
	PaymentGateway(Cash): {
		BankTransfer: "manual",
		Card:         "manual",
	},
}

var GatewayCurrencyFee = map[PaymentGateway]map[Currency]struct {
	Percentage float64
	Flat       float64
	Waiver     float64
	Cap        float64
}{
	Paystack: {
		NGN: {
			Percentage: 0.015,
			Flat:       100,
			Waiver:     2_500,
			Cap:        2_000,
		},
		USD: {
			Percentage: 0.039,
			Flat:       100,
			Waiver:     0,
			Cap:        0,
		},
	},
}

type PaymentGatewayStatus string

const (
	FlutterwaveChargeCompleted   PaymentGatewayStatus = "charge.completed"
	FlutterwaveTransferCompleted PaymentGatewayStatus = "transfer.completed"
	FlutterwaveSuccessful        PaymentGatewayStatus = "successful"
	FlutterwaveSuccess           PaymentGatewayStatus = "success"
	FlutterwaveFailed            PaymentGatewayStatus = "failed"
	FlutterwaveError             PaymentGatewayStatus = "error"
	PaystackSuccessful           PaymentGatewayStatus = "success"
	PaystackFailed               PaymentGatewayStatus = "failed"
	WalletFailed                 PaymentGatewayStatus = "failed"
	WalletSuccess                PaymentGatewayStatus = "success"
)
