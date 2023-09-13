package enum

type Frequency string

// Frequency
const (
	// daily
	Daily Frequency = "DAILY"

	// weekly
	Weekly Frequency = "WEEKLY"

	// monthly
	Monthly Frequency = "MONTHLY"

	// quarterly
	Quarterly Frequency = "QUARTERLY"

	// annually
	Annually Frequency = "ANNUALLY"

	// one time
	OneTime Frequency = "ONE_TIME"

	// pro-rata
	ProRata Frequency = "PRO_RATA"
)

var FrequencyDuration = map[Frequency]int{
	Daily:     1,
	Weekly:    7,
	Monthly:   30,
	Quarterly: 90, // 30 * (12 / 4)
	Annually:  365,
}

type Interest string

// Interest
const (
	// simple
	Simple Interest = "SIMPLE"

	// compound
	Compound Interest = "COMPOUND"

	// fixed
	Fixed Interest = "FIXED"

	// variable
	Variable Interest = "VARIABLE"

	// tiered
	Tiered Interest = "TIERED"
)
