package handlers

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

var connections = make(map[string]map[net.Conn]bool)

func WebsocketHandler(c *gin.Context) {
	conn, _, _, err := ws.UpgradeHTTP(c.Request, c.Writer)
	if err != nil {
		log.Println(err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	estName := c.GetString("est_name")
	user, _ := c.Params.Get("user")

	wsConnString := fmt.Sprint(user, estName)

	if _, ok := connections[wsConnString]; !ok {
		connections[wsConnString] = make(map[net.Conn]bool)
	}

	connections[wsConnString][conn] = true

	go func() {
		defer conn.Close()
		defer delete(connections[wsConnString], conn)

		for {
			msg, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				log.Println(err)
				break
			}

			err = wsutil.WriteServerMessage(conn, op, msg)
			if err != nil {
				log.Println(err)
				break
			}

		}
	}()
}

func broadcastData(estName string, msg string, broadcastType string) {
	var wsConnString string

	message := []byte(msg)

	// broadcastType means where is the data sent from
	if broadcastType == "qm" {
		// for example this data is sent from the qm to all users
		wsConnString = fmt.Sprintf("user%s", estName)
	} else {
		// vice versa this is when a user joins and qm needs to be updated
		wsConnString = fmt.Sprintf("qm%s", estName)
	}

	log.Println("Triggered: ", wsConnString)
	if _, ok := connections[wsConnString]; !ok {
		return
	}

	for conn := range connections[wsConnString] {
		_ = wsutil.WriteServerText(conn, message)
	}
}
