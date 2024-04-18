package handlers

import (
	"fmt"
	"gantre/database"
	"gantre/models"
	"log"
	"sort"

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

	lineLeng := len(QueueDB.QueueList)
	line := []int64{}

	for a := range lineLeng {
		line = append(line, int64(a+1))
	}

	c.HTML(200, "frontDesk", gin.H{
		"EstType":            estType,
		"EstName":            estName,
		"UserType":           "qm",
		"QueueList":          line,
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

	// Brute forced solution. Fix later by improving db design
	// wost case is now (O)n2
	for k, v := range QueueDB.QueueList {
		if v == int(QueueDB.CurrentQueueNumber) {
			delete(QueueDB.QueueList, k)
			break
		}
	}

	QueueDB.CurrentQueueNumber++

	if _, err := client.Collection(estType).Doc(estName).Set(c, map[string]interface{}{
		"QueueList":          QueueDB.QueueList,
		"CurrentQueueNumber": QueueDB.CurrentQueueNumber,
	}); err != nil {
		log.Println("Error updating firestore:", err)
		c.JSON(500, gin.H{
			"Error": "Failed to update database",
		})
	}

	c.HTML(200, "frontDeskMain", gin.H{
		"EstType":            estType,
		"EstName":            estName,
		"UserType":           "qm",
		"QueueList":          QueueList,
		"RestaurantName":     QueueDB.RestaurantName,
		"CurrentQueueNumber": QueueDB.CurrentQueueNumber,
	})

	broadcastToUsers(estName, fmt.Sprint(QueueDB.CurrentQueueNumber))
}
