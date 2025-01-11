package zksnark

import (
	"encoding/hex"
	"time"

	"cendit.io/garage/config"
	"cendit.io/garage/function"
	"cendit.io/garage/primer/typing"
)

// Witness takes in a signture and a message and returns a boolean value indicating if the signature is valid for the message
func Witness(signature, message string) bool {
	// Generate a 6 digit code
	code, _ := GeneratePasscode(function.UTCTimeStep(time.Now(), typing.TimeStep{
		Minute: 0,
	}), config.Environment().TOTPSecret)

	// pad the 6 digit code to generate 32 bytes
	pad := make([]byte, 32)
	copy(pad, code)
	key := hex.EncodeToString(pad)

	// Encrypt the message with the key
	witness, err := Encrypt(key, message)
	if err != nil {
		return false
	}

	return witness == signature
}
