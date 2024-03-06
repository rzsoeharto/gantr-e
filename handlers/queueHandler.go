package handlers

import (
	"fmt"
	jwthandler "gantre/JWTHandlers"
	"gantre/database"
	"gantre/models"
	"gantre/utils"
	"log"
	"net/http"
	"sort"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

func QueueHandler(c *gin.Context) {
	client := database.DbAccess(c)

	estType := c.GetString("est_type")
	estName := c.GetString("est_name")

	doc, err := client.Collection(estType).Doc(estName).Get(c)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"Message": "Something went wrong in the backend.",
		})
	}

	var QueueDB models.QueueModel
	if err := doc.DataTo(&QueueDB); err != nil {
		log.Println("Error scanning data:", err)
		c.HTML(http.StatusInternalServerError, "serverError", gin.H{"Message": "Something went wrong in the backend."})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update data in database"})
			return
		}

		ordinalIndicator := utils.QueueFormat(int64(NextQueue))

		c.SetCookie("SID", tokenString, 604800, "/", "localhost", true, true)

		c.HTML(http.StatusOK, "customer", gin.H{
			"data": fmt.Sprintf("%v%v", NextQueue, ordinalIndicator),
		})

		return
	}

	CustomerQueueNumber, _ := jwthandler.ParseCookies(c)

	LineLen := int64(CustomerQueueNumber) - QueueDB.CurrentQueueNumber
	ordinalIndicator := utils.QueueFormat(LineLen)

	c.HTML(http.StatusOK, "customer", gin.H{
		"QueueNumber":  CustomerQueueNumber,
		"QueueOrdinal": fmt.Sprintf("%v%v", LineLen, ordinalIndicator),
	})
}
