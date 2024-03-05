package middlewares

import (
	jwthandler "gantre/JWTHandlers"

	"github.com/gin-gonic/gin"
)

func CheckCookies(c *gin.Context) {
	token, err := c.Cookie("SID")

	if err != nil {
		c.Set("token", "NoToken")
		c.Next()
		return
	}

	c.Set("token", token)
	jwthandler.ParseCookies(c)
}
