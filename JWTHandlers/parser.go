package jwthandler

import (
	"fmt"
	"gantre/models"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ParseCookies(c *gin.Context) (int64, int64) {
	t := c.GetString("token")
	key := []byte(os.Getenv("SECRET_KEY"))

	token, err := jwt.ParseWithClaims(t, &models.QueueClaimsStruct{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			signErr := fmt.Sprintf("Unexpected signing method %v", t.Header["alg"])
			log.Println(signErr)
			return nil, fmt.Errorf(signErr)
		}
		return key, nil
	})

	if err != nil || !token.Valid {
		log.Println(err)
		c.JSON(400, gin.H{
			"Message": "Given token is not valid.",
		})
		return 0, 0
	}

	claims := token.Claims.(*models.QueueClaimsStruct)

	queueNumber := claims.GetQueueNumber()

	return queueNumber, 0
}
