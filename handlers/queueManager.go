package handlers

import (
	"gantre/database"
	"gantre/models"
	"log"
	"sort"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func FrontDeskHandler(c *gin.Context) {
	var QueueDB models.QueueModel
	client := database.DbAccess(c)

	estType := c.GetString("est_type")
	estName := c.GetString("est_name")

	doc, err := client.Collection(estType).Doc(estName).Get(c)
	if err != nil {
		log.Println(err)
		c.HTML(500, "serverError", gin.H{
			"Message":  "Uh-oh something went wrong.",
			"Message2": "Unable to retrieve establishment data.",
		})
		return
	}

	if err := doc.DataTo(&QueueDB); err != nil {
		log.Println("Error scanning data:", err)
		c.HTML(500, "serverError", gin.H{
			"Message": "Uh-oh something went wrong.",
		})
		return
	}

	c.HTML(200, "frontDesk", gin.H{
		"EstType":            estType,
		"EstName":            estName,
		"RestaurantName":     QueueDB.RestaurantName,
		"CurrentQueueNumber": QueueDB.CurrentQueueNumber,
	})
}

func UpdateFromFrontDeskHandler(c *gin.Context) {
	var QueueDB models.QueueModel
	client := database.DbAccess(c)

	estType := c.GetString("est_type")
	estName := c.GetString("est_name")

	doc, err := client.Collection(estType).Doc(estName).Get(c)
	if err != nil {
		log.Println(err)
		c.HTML(500, "serverError", gin.H{
			"Message":  "Uh-oh something went wrong.",
			"Message2": "Unable to retrieve establishment data.",
		})
		return
	}

	if err := doc.DataTo(&QueueDB); err != nil {
		log.Println("Error scanning data:", err)
		c.HTML(500, "serverError", gin.H{
			"Message": "Uh-oh something went wrong.",
		})
		return
	}

	QueueList := []int64{}
	for _, v := range QueueDB.QueueList {
		QueueList = append(QueueList, int64(v))
	}

	sort.SliceStable(QueueList, func(i, j int) bool { return QueueList[i] < QueueList[j] })

	if QueueDB.CurrentQueueNumber >= int64(QueueList[len(QueueList)-1]) {
		c.HTML(200, "clientQueueNumber", gin.H{
			"CurrentQueueNumber": QueueDB.CurrentQueueNumber,
		})
		return
	}

	QueueDB.CurrentQueueNumber++

	if _, err := client.Collection(estType).Doc(estName).Set(c, map[string]interface{}{
		"CurrentQueueNumber": QueueDB.CurrentQueueNumber,
	}, firestore.MergeAll); err != nil {
		log.Println("Error updating firestore:", err)
		c.JSON(500, gin.H{
			"Error": "Failed to update database",
		})
	}

	c.HTML(200, "clientQueueNumber", gin.H{
		"CurrentQueueNumber": QueueDB.CurrentQueueNumber,
	})
}
