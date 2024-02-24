package middlewares

import (
	"log"

	"github.com/gin-gonic/gin"
)

func RequestParams(c *gin.Context) {
	estType, ok := c.Params.Get("est_type")
	if !ok {
		log.Println("Unable to parse est_type params")
		c.JSON(500, gin.H{
			"Message": "The QR you scanned seems to be incorrect.",
		})
		c.Abort()
		return
	}

	estName, ok := c.Params.Get("est_name")
	if !ok {
		log.Println("Unable to parse est_type params")
		c.JSON(500, gin.H{
			"Message": "The QR you scanned seems to be incorrect.",
		})
	}

	c.Set("est_type", estType)
	c.Set("est_name", estName)

	c.Next()
}
