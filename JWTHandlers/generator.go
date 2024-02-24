package jwthandler

import (
	"gantre/models"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(currQueue int64) string {
	currQueue++

	key := []byte(os.Getenv("SECRET_KEY"))
	claims := models.QueueClaimsStruct{
		QueueNumber: currQueue,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "Gantr-e",
			Subject:   "Customer",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	SignedToken, err := token.SignedString(key)

	if err != nil {
		log.Println(err)
		return ""
	}

	return SignedToken
}
