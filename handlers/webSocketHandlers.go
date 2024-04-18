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

func broadcastToUsers(estName string, msg string) {
	data := fmt.Sprintf(`<h1 id="current-queue-number" class="text-9xl text-center font-bold"> %s </h1>`, msg)
	message := []byte(data)
	wsConnString := fmt.Sprintf("user%s", estName)

	log.Println("Broadcasting to: ", wsConnString)
	if _, ok := connections[wsConnString]; !ok {
		return
	}

	for conn := range connections[wsConnString] {
		_ = wsutil.WriteServerText(conn, message)
	}
}

func broadcastToFrontDesk(estName string, msg string) {
	data := fmt.Sprintf(`<p id="queue-number-%s">%s</p>`, msg, msg)
	message := []byte(data)
	wsConnString := fmt.Sprintf("qm%s", estName)

	log.Println("Broadcasting to: ", wsConnString)
	if _, ok := connections[wsConnString]; !ok {
		return
	}

	for conn := range connections[wsConnString] {
		_ = wsutil.WriteServerText(conn, message)
	}
}
