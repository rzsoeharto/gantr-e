package middlewares

import (
	"fmt"
	"gantre/database"
	"gantre/models"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CheckCookies(c *gin.Context) {
	client := database.DbAccess(c)

	token, err := c.Cookie("SID")

	if err != nil {
		c.Set("token", "NoToken")
		c.Next()
		return
	}

	c.Set("token", token)
	ParseCookies(c, client)
}

func ParseCookies(c *gin.Context, client *firestore.Client) {
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
		return
	}

	claims := token.Claims.(*models.QueueClaimsStruct)

	qNum := claims.GetQueueNumber()

	c.Set("CustomerQueue", qNum)
	c.Next()
}
