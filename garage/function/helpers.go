package function

import (
	"math/rand"
	"strings"
	"time"
)

func GenerateRandomNumber(length int) string {
	if length <= 0 {
		return ""
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(byte(r.Intn(10) + '0'))
	}
	return sb.String()
}
