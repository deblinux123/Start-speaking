package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWTKey = []byte("my_secrete_key")

func GenerateToken(userId uint) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userId,
			"exp":     time.Now().Add(time.Minute * 1).Unix(),
		},
	)

	tokenString, err := token.SignedString(JWTKey)

	return tokenString, err
}

func GenerateRefreshToken(userId uint) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userId,
			"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
		},
	)

	return token.SignedString(JWTKey)
}
