package handlers

import (
	jwthandler "gantre/JWTHandlers"
	"gantre/database"
	"gantre/models"
	"log"
	"net/http"
	"sort"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func QueueHandler(c *gin.Context) {
	var QueueDB models.QueueModel

	client := database.DbAccess(c)

	estType := c.GetString("est_type")
	estName := c.GetString("est_name")

	hostName := c.Request.Host

	doc, err := client.Collection(estType).Doc(estName).Get(c)
	if err != nil {
		log.Println(err)
		c.HTML(500, "serverError", gin.H{
			"Message": "Uh-oh something went wrong.",
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

	QueueList := make([]int, len(QueueDB.QueueList))
	for _, v := range QueueDB.QueueList {
		QueueList = append(QueueList, v)
	}

	sort.SliceStable(QueueList, func(i, j int) bool { return QueueList[i] < QueueList[j] })

	t := c.GetString("token")
	if t == "NoToken" {
		var NextQueue int

		if len(QueueDB.QueueList) == 0 && QueueDB.CurrentQueueNumber != 1 {
			NextQueue = int(QueueDB.CurrentQueueNumber)
		} else if len(QueueDB.QueueList) == 0 {
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
			"Hostname":           hostName,
			"CurrentQueueNumber": QueueDB.CurrentQueueNumber,
		})

		// Notify front desk
		broadcastToFrontDesk(estName, strconv.Itoa(NextQueue))
		return
	}

	CustomerQueueNumber, _ := jwthandler.ParseCookies(c)

	c.HTML(http.StatusOK, "customer", gin.H{
		"EstType":            estType,
		"EstName":            estName,
		"UserType":           "user",
		"QueueNumber":        CustomerQueueNumber,
		"Hostname":           hostName,
		"CurrentQueueNumber": QueueDB.CurrentQueueNumber,
	})
}

func ClearCookie(c *gin.Context) {
	c.SetCookie("SID", "", 0, "/", "localhost", true, true)
	c.HTML(200, "customerAdmitted", gin.H{})
}
