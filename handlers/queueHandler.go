package handlers

import (
	"fmt"
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
	client := database.DbAccess(c)

	var QueueDB models.QueueModel

	estType := c.GetString("est_type")
	estName := c.GetString("est_name")

	doc, err := client.Collection(estType).Doc(estName).Get(c)

	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"Message": "Something went wrong in the backend.",
		})
	}

	scanError := doc.DataTo(&QueueDB)

	t := c.GetString("token")

	if scanError != nil {
		log.Println(scanError)

		c.HTML(500, "serverError", gin.H{
			"Message": "Something went wrong in the backend.",
		})

		return
	}

	QueueList := make([]int, 0, 20)

	for _, v := range QueueDB.QueueList {
		QueueList = append(QueueList, v)
	}

	sort.SliceStable(QueueList, func(i, j int) bool { return QueueList[i] < QueueList[j] })

	if t == "NoToken" {
		NextQueue := QueueList[len(QueueList)-1]
		tokenString := jwthandler.GenerateToken(int64(NextQueue))
		QueueDB.QueueList[tokenString] = NextQueue + 1

		_, setErr := client.Collection(estType).Doc(estName).Set(c, map[string]interface{}{
			"QueueList": map[string]interface{}{
				tokenString: NextQueue + 1,
			},
		}, firestore.MergeAll)

		if setErr != nil {
			log.Println(setErr)
			c.JSON(500, gin.H{
				"err": setErr,
			})
			return
		}

		c.SetCookie("SID", tokenString, 604800, "/", "localhost", true, true)

		c.HTML(http.StatusOK, "customer", gin.H{
			"data": fmt.Sprintf("%v nd/rd/th idk", NextQueue+1),
		})

		return
	}

	CustomerQueueNumber, _ := jwthandler.ParseCookies(c)

	LineLen := int64(CustomerQueueNumber) - QueueDB.CurrentQueueNumber

	c.HTML(http.StatusOK, "customer", gin.H{
		"data": fmt.Sprintf("%v nd/rd/th idk", LineLen),
	})
}
