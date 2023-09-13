package function

import "crypto/sha256"

// StringSha256 computes the sha256 of the given string
func StringSha256(input string) string {
	n := sha256.New()
	n.Write([]byte(input))
	return string(n.Sum([]byte("")))
}
