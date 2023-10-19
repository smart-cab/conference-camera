package handlers

import (
	"conferencecam/types"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true }, // for cors
	}
	client *websocket.Conn
)

func WebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Error(types.INCORRECT_WEBSOCKET)
		return
	}
	if client != nil {
		c.Error(types.ALREADY_CONNECTED)
		return
	}

	client = conn

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			client = nil
			return
		}
		fmt.Println(string(msg))
	}
}
