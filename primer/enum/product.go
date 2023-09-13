// a symbolic reference for all configure.Blache enumeration
package enum

type Template string

// Templates
const (
	// a deposit account with limited access to funds
	// and restrictions on the number of withdrawals
	RegularDepositWithRestrictedWithdrawals Template = "RDWRW"

	// a deposit account that allows for varying deposit
	// amounts and withdrawal frequencies
	FlexibleRegularDeposit Template = "FRD"

	// specific amount of money is deposited for a fixed
	// period of time, higher interest and limited withdrawals
	FixedDeposit Template = "FD"
)

type Status string

// Statuses (product)
const (

	// open
	Open Status = "OPEN"

	// closed
	Closed Status = "CLOSED"

	// active
	Active Status = "ACTIVE"

	// inactive
	Inactive Status = "INACTIVE"

	// pending
	Pending Status = "PENDING"

	// deleted
	Deleted Status = "DELETED"

	// completed
	Completed Status = "COMPLETED"
)

// Statuses (activity)
const (
	// success
	Success Status = "SUCCESS"

	// failure
	Failure Status = "FAILURE"
)

type AccountCreationType string

// AccountCreationTypes
const (
	// account is created automatically
	Automatic AccountCreationType = "AUTOMATIC"

	// account is created manually
	Manual AccountCreationType = "MANUAL"
)

type Tier string

// Tiers
const (
	// TierOne is the first tier
	TierOne Tier = "one"

	// TierTwo is the second tier
	TierTwo Tier = "two"

	// TierThree is the third tier
	TierThree Tier = "three"

	// TierPrefix is the prefix for all tiers
	TierPrefix Tier = "tier_"
)

// Tier Limits
const (
	// TierOneMaxCumulativeBalance is the maximum cumulative balance for tier one
	TierOneMaxCumulativeBalance = 300_000

	// TierOneMaxSingleDepositPerDay is the maximum single deposit per day for tier one
	TierOneMaxSingleDepositPerDay = 50_000

	// TierTwoMaxCumulativeBalance is the maximum cumulative balance for tier two
	TierTwoMaxCumulativeBalance = 500_000

	// TierTwoMaxSingleDepositPerDay is the maximum single deposit per day for tier two
	TierTwoMaxSingleDepositPerDay = 100_000

	// TierThreeMaxSingleDepositPerTransaction is the maximum single deposit per transaction for tier three
	TierThreeMaxSingleDepositPerTransaction = 1_000_000

	// TierThreeMaxSingleDepositPerDay is the maximum single deposit per day for tier three
	TierThreeMaxSingleDepositPerDay = 0

	// TierThreeMaxCumulativeBalance is the maximum cumulative balance for tier three
	TierThreeMaxCumulativeBalance = 0
)
