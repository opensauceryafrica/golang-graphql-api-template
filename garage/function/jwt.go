package function

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"cendit.io/garage/primer/constant"
	"cendit.io/garage/primer/enum"
	"cendit.io/garage/primer/typing"

	"cendit.io/garage/redis"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(session *typing.Session, secret string) (string, error) {
	type Claims struct {
		ID    string    `json:"id"`
		Email string    `json:"email"`
		Role  enum.Role `json:"role"`
		jwt.RegisteredClaims
	}

	claims := &Claims{
		ID:    session.ID,
		Email: session.Email,
		Role:  session.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	var sessionData bytes.Buffer
	err = gob.NewEncoder(&sessionData).Encode(session)
	if err != nil {
		return "", err
	}

	if err = redis.Ral.Set(fmt.Sprintf("%s-%s", constant.UserRedisKey, tokenString), sessionData.Bytes(), time.Minute*5); err != nil {
		return "", err
	}

	return tokenString, nil
}
