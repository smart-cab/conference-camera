package handlers

import (
	"conferencecam/ptz"
	"conferencecam/types"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true }, // for cors
	}
	hub     *websocket.Conn
	token   string
	clients map[*websocket.Conn]bool
)

func WebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		types.Logger.Errorln("fake websocket")
		c.Error(types.INCORRECT_WEBSOCKET)
		return
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if conn == hub {
				hub = nil

				for client := range clients {
					client.WriteMessage(websocket.TextMessage, []byte("disconnected"))
					client.Close()
				}
			} else if clients[conn] {
				if hub != nil {
					hub.WriteMessage(websocket.TextMessage, []byte("disconnected"))
				}

				delete(clients, conn)
			}
			continue
		}

		data := strings.Split(string(msg), ":")
		if data[0] == "hub" {
			// HUB ACTIONS

			// Check conn is from hub or not???
			if hub != nil && hub != conn {
				types.Logger.Errorf("user %s try start already hub", conn.LocalAddr())

				conn.WriteMessage(websocket.TextMessage, []byte("error:already"))
				return
			}
			hub = conn

			switch data[1] {
			case "init":
				// Отправялем клиент, включена ли автоматическая генерация QR кода или нужно через апи
				conn.WriteMessage(websocket.TextMessage, []byte("autoqr:"+os.Getenv("AUTO_QR_CODE")))

				types.Logger.Infof("user %s start hub successful", conn.LocalAddr())
			case "token":
				// Генерируем новый токен и отсылаем клиенту
				if os.Getenv("AUTO_QR_CODE") == "1" {
					token = randToken(32)
					conn.WriteMessage(websocket.TextMessage, []byte("token:"+token))
					types.Logger.Infof("sent token %s to hub to user: %s", token, conn.LocalAddr())
				}
			default:
				types.Logger.Errorf("from user %s recieved unknown action!", conn.LocalAddr())
			}
		} else if data[0] == "user" {
			// USER ACTIONS
			switch data[1] {
			case "init": // user:init
				types.Logger.Infof("user %s init connection tab", conn.LocalAddr())
			case "connect": // user:connect:TOKEN
				// Проверяем токен для авторизации
				types.Logger.Infof("user %s try auth to hub: %s", conn.LocalAddr(), data[2])
				if data[2] != token {
					// Неверный токен, отсылаем wrong
					types.Logger.Errorf("user %s wrong token: %s", conn.LocalAddr(), data[2])
					conn.WriteMessage(websocket.TextMessage, []byte("wrong"))
				} else {
					// Верный токен, записываем подключение в Client и отсылаем хабу о клиенте
					types.Logger.Errorf("user %s successful connected!", conn.LocalAddr())

					clients[conn] = true

					conn.WriteMessage(websocket.TextMessage, []byte("connected"))
					hub.WriteMessage(websocket.TextMessage, []byte("connected:"+conn.RemoteAddr().String()))

					// Prepare list of devices
					devices := []string{}
					for _, d := range ptz.GetActiveDevices() {
						devices = append(devices, d.Name()+":"+d.Capability().Card)
					}

					devices[0], devices[1] = devices[1], devices[0]
					selectedCamera := ""
					conn.WriteMessage(websocket.TextMessage, []byte("devices:"+strings.Join(devices, "|")))

					if ptz.PTZ != nil {
						selectedCamera = ptz.PTZ.Device.Name()
					} else {
						selectedCamera = ""
					}
					conn.WriteMessage(websocket.TextMessage, []byte("selected-device:"+selectedCamera))
					types.Logger.Debugf("devices list: %s --- selected device: %s", strings.Join(devices, " | "), selectedCamera)
				}
			case "device":
				// Ставим новую камеру
				types.Logger.Infof("user %s set new camera %s", conn.LocalAddr(), data[2])

				ptz.PTZ.Close()

				camera, err := ptz.Init(data[2])
				if err != nil {
					types.Logger.Errorf("\n\n\n------------------------\nuser %s error set new camera: %s\n------------------------\n\n\n", conn.LocalAddr(), err.Error())
				}

			case "move":
				// Движение PTZ камеры
				types.Logger.Infof("user %s move camera to %s", conn.LocalAddr(), data[2])
				var cmd uint32
				var value int32

				switch data[2] {
				case "left":
					cmd = ptz.CTRL_HORIZONTAL
					value = -300
				case "right":
					cmd = ptz.CTRL_HORIZONTAL
					value = 300
				case "top":
					cmd = ptz.CTRL_VERTICAL
					value = 200
				case "bottom":
					cmd = ptz.CTRL_VERTICAL
					value = -200
				}

				ptz.SendCmd(cmd, value)
			}
		}
	}
}

func randToken(n int) string {
	rand.Seed(time.Now().UnixNano()) // fix repeat tokens

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
