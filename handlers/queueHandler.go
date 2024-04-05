package handlers

import (
	jwthandler "gantre/JWTHandlers"
	"gantre/database"
	"gantre/models"
	"log"
	"net/http"
	"sort"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func QueueHandler(c *gin.Context) {
	var QueueDB models.QueueModel

	client := database.DbAccess(c)

	estType := c.GetString("est_type")
	estName := c.GetString("est_name")

	doc, err := client.Collection(estType).Doc(estName).Get(c)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"Message": "Uh-oh, Something went wrong.",
		})
	}

	if err := doc.DataTo(&QueueDB); err != nil {
		log.Println("Error scanning data:", err)
		c.HTML(500, "serverError", gin.H{})
		return
	}

	QueueList := make([]int, len(QueueDB.QueueList))
	for _, v := range QueueDB.QueueList {
		QueueList = append(QueueList, v)
	}

	sort.SliceStable(QueueList, func(i, j int) bool { return QueueList[i] < QueueList[j] })

	t := c.GetString("token")
	if t == "NoToken" {
		var NextQueue int

		if len(QueueDB.QueueList) == 0 {
			NextQueue = 1
		} else {
			NextQueue = QueueList[len(QueueList)-1] + 1
		}

		tokenString := jwthandler.GenerateToken(int64(NextQueue))
		QueueDB.QueueList[tokenString] = NextQueue

		if _, err := client.Collection(estType).Doc(estName).Set(c, map[string]interface{}{
			"QueueList": map[string]interface{}{
				tokenString: NextQueue,
			},
		}, firestore.MergeAll); err != nil {
			log.Println("Error updating Firestore:", err)
			c.JSON(500, gin.H{"Error": "Failed to update database"})
			return
		}

		c.SetCookie("SID", tokenString, 604800, "/", "localhost", true, true)

		c.HTML(http.StatusOK, "customer", gin.H{
			"EstType":            estType,
			"EstName":            estName,
			"UserType":           "user",
			"QueueNumber":        NextQueue,
			"CurrentQueueNumber": QueueDB.CurrentQueueNumber,
		})

		// Notify front desk
		broadcastToFrontDesk(estName, "hi, update from backend")
		return
	}

	CustomerQueueNumber, _ := jwthandler.ParseCookies(c)

	c.HTML(http.StatusOK, "customer", gin.H{
		"EstType":            estType,
		"EstName":            estName,
		"UserType":           "user",
		"QueueNumber":        CustomerQueueNumber,
		"CurrentQueueNumber": QueueDB.CurrentQueueNumber,
	})
}
