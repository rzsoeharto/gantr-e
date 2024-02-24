package handlers

import (
	"fmt"
	jwthandler "gantre/JWTHandlers"
	"gantre/database"
	"gantre/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func QueueHandler(c *gin.Context) {
	client := database.DbAccess(c)

	var QueueDB models.QueueModel

	estType := c.GetString("est_type")
	estName := c.GetString("est_name")

	doc, err := client.Collection(estType).Doc(estName).Get(c)

	if err != nil {
		c.JSON(500, gin.H{
			"Message": err,
		})
	}

	scanError := doc.DataTo(&QueueDB)

	t := c.GetString("token")

	if scanError != nil {
		fmt.Println(scanError)

		c.HTML(500, "serverError", gin.H{})

		return
	}

	if t == "NoToken" {
		tokenString := jwthandler.GenerateToken(QueueDB.CurrentQueueNumber)
		c.SetCookie("SID", tokenString, 604800, "/", "localhost", true, true)
	}

	c.HTML(http.StatusOK, "queuer", gin.H{
		"data": "8th",
	})
}
