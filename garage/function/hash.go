package function

import (
	"crypto/sha256"

	"golang.org/x/crypto/bcrypt"
)

// StringSha256 computes the sha256 of the given string
func StringSha256(input string) string {
	n := sha256.New()
	n.Write([]byte(input))
	return string(n.Sum([]byte("")))
}

// HashPasscode hashes the given password using bcrypt and returns the hashed password as a string.
func HashPasscode(password string) (string, error) {
	hashedPasscode, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPasscode), nil
}

// ComparePasscode compares a hashed password with a plain text password.
func ComparePasscode(hashedPasscode, passcode string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPasscode), []byte(passcode))
	return err == nil
}
