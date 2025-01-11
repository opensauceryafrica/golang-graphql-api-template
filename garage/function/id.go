package function

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// GenerateUUID returns a new universally unique identifier
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateRef returns a usable reference computed against the given label (or app name) from the current timestamp
func GenerateRef(label ...string) string {
	if len(label) == 0 {
		label = []string{"cendit"}
	}
	return fmt.Sprintf(`%s-%d`, strings.ToUpper(label[0]), time.Now().Nanosecond())
}
