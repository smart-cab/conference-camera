package handlers

import (
	"conferencecam/types"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"

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
	hub    *websocket.Conn
	token  string
)

func WebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Error(types.INCORRECT_WEBSOCKET)
		return
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if conn == hub {
				hub = nil
				if client != nil {
					client.WriteMessage(websocket.TextMessage, []byte("disconnected"))
				}
			} else if conn == client {
				client = nil
				if hub != nil {
					hub.WriteMessage(websocket.TextMessage, []byte("disconnected"))
				}
			}
			return
		}

		data := strings.Split(string(msg), ":")
		if data[0] == "hub" {
			// HUB ACTIONS

			// Check conn is from hub or not???
			if hub != nil && hub != conn {
				conn.WriteMessage(websocket.TextMessage, []byte("error:already"))
				return
			}
			hub = conn

			switch data[1] {
			case "init":
				// Отправялем клиент, включена ли автоматическая генерация QR кода или нужно через апи
				conn.WriteMessage(websocket.TextMessage, []byte("autoqr:"+os.Getenv("AUTO_QR_CODE")))
				fmt.Println("hub connected successful!")
			case "token":
				// Генерируем новый токен и отсылаем клиенту
				if os.Getenv("AUTO_QR_CODE") == "1" {
					token = randToken(32)
					conn.WriteMessage(websocket.TextMessage, []byte("token:"+token))
					fmt.Println("sent token to hub")
				}
			default:
				fmt.Println("unknown action")
			}
		} else if data[0] == "user" {
			if client != nil && client != conn {
				// Если у нас уже есть авторизованный клиент, то никому не даем доступ
				conn.WriteMessage(websocket.TextMessage, []byte("already"))
				return
			}

			// USER ACTIONS
			switch data[1] {
			case "init": // user:init
				fmt.Println("user open page for connect")
			case "connect": // user:connect:TOKEN
				// Проверяем токен для авторизации
				fmt.Println("user try connect " + data[2])
				if data[2] != token {
					// Неверный токен, отсылаем wrong
					conn.WriteMessage(websocket.TextMessage, []byte("wrong"))
				} else {
					// Верный токен, записываем подключение в client и отсылаем хабу о клиенте
					client = conn
					conn.WriteMessage(websocket.TextMessage, []byte("connected"))
					hub.WriteMessage(websocket.TextMessage, []byte("connected:"+conn.RemoteAddr().String()))
				}
			}
		}
	}
}

func randToken(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
