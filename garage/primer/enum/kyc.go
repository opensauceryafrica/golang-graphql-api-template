package enum

// Enum for ID Document
type IDDocument string

const (
	NINSlip        IDDocument = "NIN_SLIP"
	IDCard         IDDocument = "ID_CARD"
	Passport       IDDocument = "PASSPORT"
	VotersCard     IDDocument = "VOTERS_CARD"
	DriversLicense IDDocument = "DRIVERS_LICENSE"
)

// Enum for Status
type BVNStatus string

const (
	NotVerified BVNStatus = "NOT_VERIFIED"
	BVNVerified BVNStatus = "BVN_VERIFIED"
	IDVerified  BVNStatus = "ID_VERIFIED"
	BVNRejected BVNStatus = "BVN_REJECTED"
	IDRejected  BVNStatus = "ID_REJECTED"
)

type KYCStatus string

const (
	KYCPending  KYCStatus = "KYC_PENDING"
	KYCVerified KYCStatus = "KYC_VERIFIED"
	KYCRejected KYCStatus = "KYC_REJECTED"
)

type KYCOption string

const (
	NIN      KYCOption = "NIN"
	BVN      KYCOption = "BVN"
	IDUPLOAD KYCOption = "IDUPLOAD"
)
