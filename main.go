package main

import (
	"fmt"
	jwthandler "gantre/JWTHandlers"
	"gantre/database"
	"gantre/handlers"
	"gantre/middlewares"
	"gantre/models"
	"log"
	"sort"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/lpernett/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Panicln(err)
		return
	}
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templ/**/*")

	userGroup := r.Group("/:est_type")
	{
		userGroup.GET("/:est_name", middlewares.RequestParams, middlewares.CheckCookies, handlers.QueueHandler)
	}

	r.POST("/test", func(c *gin.Context) {
		cl := database.DbAccess(c)

		var QueueDB models.QueueModel

		doc, err := cl.Collection("restaurants").Doc("sushi-tei-pim-2").Get(c)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"err": err,
			})
			return
		}

		doc.DataTo(&QueueDB)

		QueueList := make([]int, 0, 20)

		for _, v := range QueueDB.QueueList {
			QueueList = append(QueueList, v)
		}

		sort.SliceStable(QueueList, func(i, j int) bool { return QueueList[i] < QueueList[j] })

		fmt.Println(QueueList)
		NextQueue := QueueList[len(QueueList)-1]

		tokenString := jwthandler.GenerateToken(int64(NextQueue))

		QueueDB.QueueList[tokenString] = NextQueue + 1
		_, setErr := cl.Collection("restaurants").Doc("sushi-tei-pim-2").Set(c, map[string]interface{}{
			"QueueList": map[string]interface{}{
				tokenString: NextQueue + 1,
			},
		}, firestore.MergeAll)

		if setErr != nil {
			fmt.Println(setErr)
			c.JSON(500, gin.H{
				"err": setErr,
			})
			return
		}

		c.JSON(200, gin.H{
			"msg": "everything should be good",
		})

	})

	r.Run()
}
