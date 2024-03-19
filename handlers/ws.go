package handlers

import (
	"conferencecam/ptz"
	"conferencecam/types"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true }, // for cors
	}
	hub    *websocket.Conn
	token  string
	Client *websocket.Conn
)

func WebSocket(c *gin.Context) {
	log := c.MustGet("log").(*logrus.Logger)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorln("fake websocket")
		c.Error(types.INCORRECT_WEBSOCKET)
		return
	}

	for {
		if conn == nil {
			log.Errorln("websocket disconnected")
			return
		}
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if conn == hub {
				hub = nil
				if Client != nil {
					Client.WriteMessage(websocket.TextMessage, []byte("disconnected"))
				}
			} else if conn == Client {
				Client = nil
				if hub != nil {
					hub.WriteMessage(websocket.TextMessage, []byte("disconnected"))
				}
			}

			break
		}

		data := strings.Split(string(msg), ":")
		if data[0] == "hub" {
			// HUB ACTIONS

			// Check conn is from hub or not???
			if hub != nil && hub != conn {
				log.Errorf("user %s try start already hub", conn.LocalAddr())
				conn.WriteMessage(websocket.TextMessage, []byte("error:already"))
				continue
			}
			hub = conn

			switch data[1] {
			case "init":
				// Отправялем клиент, включена ли автоматическая генерация QR кода или нужно через апи
				conn.WriteMessage(websocket.TextMessage, []byte("autoqr:"+os.Getenv("AUTO_QR_CODE")))
				log.Infof("user %s start hub successful", conn.LocalAddr())
			case "token":
				// Генерируем новый токен и отсылаем клиенту
				if os.Getenv("AUTO_QR_CODE") == "1" {
					token = randToken(32)
					conn.WriteMessage(websocket.TextMessage, []byte("token:"+token))
					log.Infof("sent token %s to hub to user: %s", token, conn.LocalAddr())
				}
			default:
				fmt.Println("unknown action")
			}
		} else if data[0] == "user" {
			if Client != nil && Client != conn {
				// Если у нас уже есть авторизованный клиент, то никому не даем доступ
				// Client = conn
				conn.WriteMessage(websocket.TextMessage, []byte("already"))
				return
			}

			// USER ACTIONS
			switch data[1] {
			case "init": // user:init
				log.Infof("user %s init connection tab", conn.LocalAddr())
			case "connect": // user:connect:TOKEN
				// Проверяем токен для авторизации
				log.Infof("user %s try auth to hub: %s", conn.LocalAddr(), data[2])
				if data[2] != token {
					// Неверный токен, отсылаем wrong
					log.Errorf("user %s wrong token: %s", conn.LocalAddr(), data[2])
					conn.WriteMessage(websocket.TextMessage, []byte("wrong"))
				} else {
					// Верный токен, записываем подключение в Client и отсылаем хабу о клиенте
					log.Errorf("user %s successful connected!", conn.LocalAddr())
					Client = conn
					conn.WriteMessage(websocket.TextMessage, []byte("connected"))
					if hub != nil {
						hub.WriteMessage(websocket.TextMessage, []byte("connected:"+conn.RemoteAddr().String()))
					}
					devices, _ := ptz.GetDevices()
					selectedCamera := ""
					conn.WriteMessage(websocket.TextMessage, []byte("devices:"+strings.Join(devices, "|")))
					if ptz.Camera.Device != nil {
						selectedCamera = ptz.Camera.Device.Name() + ":" + strconv.FormatBool(ptz.Camera.IsPTZ)
					} else {
						selectedCamera = ""
					}
					conn.WriteMessage(websocket.TextMessage, []byte("selected-device:"+selectedCamera))
					log.Debugf("devices list: %s --- selected device: %s", strings.Join(devices, " | "), selectedCamera)
				}
			case "device":
				// Ставим новую камеру
				log.Infof("user %s set new camera %s", conn.LocalAddr(), data[2])

				if err := ptz.Camera.Init(data[2]); err != nil {
					log.Errorf("\n\n\n------------------------\nuser %s error set new camera: %s\n------------------------\n\n\n", conn.LocalAddr(), err.Error())
				}
			case "move":
				// Движение PTZ камеры
				log.Infof("user %s move camera to %s", conn.LocalAddr(), data[2])
				var cmd uint32

				if data[2] == "center" {
					ptz.Camera.CenterCamera()
					continue
				}

				value, err := strconv.Atoi(data[3])
				if err != nil {
					continue
				}

				switch data[2] {
				case "left":
					cmd = ptz.CTRL_HORIZONTAL
					value *= -1
				case "right":
					cmd = ptz.CTRL_HORIZONTAL
				case "top":
					cmd = ptz.CTRL_VERTICAL
					value *= -1
				case "bottom":
					cmd = ptz.CTRL_VERTICAL
				}

				ptz.Camera.SendCmd(cmd, int32(value))
			case "zoom":
				log.Infof("user %s zoom camera to %s", conn.LocalAddr(), data[2])

				value, err := strconv.Atoi(data[2])
				if err != nil {
					log.Fatalf("user %s send wrong value to zoom: %s", conn.LocalAddr(), data[2])
				}
				ptz.Camera.SendCmd(ptz.CTRL_ZOOM, int32(value)*100)
			case "face":
				log.Infof("user %s set face detector to: %s", conn.LocalAddr(), data[2])
				if data[2] == "true" {
					ptz.Camera.FaceEnabled = true
				} else {
					ptz.Camera.FaceEnabled = false
				}
			}
		}
	}
}

func randToken(n int) string {
	rand.NewSource(time.Now().UnixNano()) // fix repeat tokens

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
