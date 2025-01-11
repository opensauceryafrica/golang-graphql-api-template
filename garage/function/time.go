package function

import (
	"time"

	"cendit.io/garage/primer/typing"
)

// TimeStep returns a new time.Time stepped forward or backward by the given step parameters
func UTCTimeStep(t time.Time, step typing.TimeStep) time.Time {
	return time.Date(t.Year()+step.Year, time.Month(int(t.Month())+step.Month), t.Day()+step.Day, t.Hour()+step.Hour, t.Minute()+step.Minute, t.Second()+step.Second, 0, t.Location()).UTC()
}
