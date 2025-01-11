// Zero-Knowledge Succinct Non-Interactive Argument of Knowledge and refers to a proof construction where one can prove possession of certain information, e.g., a secret key, without revealing that information, and without any interaction between the prover and verifier.
package zksnark

import (
	"bytes"
	"image/png"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"

	"time"
)

// GenerateTOTPKey generates a new TOTP key
func GenerateTOTPKey(issuer, accountName string) (*otp.Key, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
	})
	if err != nil {
		return nil, err
	}
	return key, nil
}

// GenerateQRCodeBytes generates a QR code from a TOTP key and returns the bytes of the PNG
func GenerateQRCodeBytes(key *otp.Key) ([]byte, error) {
	// Convert TOTP key into a PNG
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		return nil, err
	}
	png.Encode(&buf, img)

	return buf.Bytes(), nil
}

// GeneratePasscode generates a new passcode from a TOTP secret
func GeneratePasscode(t time.Time, secret string) (string, error) {
	code, err := totp.GenerateCodeCustom(secret, t, totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		return "", err
	}
	return code, nil
}

// ValidatePasscode validates a passcode against a TOTP secret
func ValidatePasscode(passcode, secret string) (bool, error) {
	return totp.ValidateCustom(passcode, secret, time.Now().UTC(), totp.ValidateOpts{
		Period:    30,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
}
