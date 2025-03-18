package data

import (
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(secretKey string, userId int64) (string, error) {

	secret, err := base64.StdEncoding.DecodeString(secretKey)
	if err != nil {
		return "", err
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "gym-ops",
		"sub": userId,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := t.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
