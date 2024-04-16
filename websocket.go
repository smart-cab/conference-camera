package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var server *Server

type Server struct {
	upgrader websocket.Upgrader
	hub      *websocket.Conn
	user     *websocket.Conn
	camera   *Camera
	token    string
	mutex    sync.Mutex
	scene    string
	frame    chan []byte
}

func NewServer(camera *Camera) *Server {
	return &Server{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		camera: camera,
		token:  "12345678",
		scene:  "merge",
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	log.Println("client connected")

	go server.stream()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if server.user == conn {
				server.user = nil
			}

			log.Println("read failed:", err)
			break
		}

		data := strings.Split(string(message), ":")

		switch data[0] {
		case "hub":
			server.connect(conn, data)

			if data[1] == "token" {
				server.getToken(conn, data)
			}
		case "user":
			if data[1] == "connect" {
				log.Println("user try connect to hub")
				server.auth(conn, data)
			}
			if data[1] == "switch" {
				server.switchDevice(conn, data)
			}
			if data[1] == "move" {
				server.move(conn, data)
			}
			if data[1] == "zoom" {
				server.zoom(conn, data)
			}
			if data[1] == "scene" {
				server.changeScene(conn, data)
			}
			if data[1] == "face" {
				server.faceDetection(conn, data)
			}
		}
	}
}

func (s *Server) sendUser(message string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.user != nil {
		err := s.user.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println(err)
		}
		return err
	}

	return nil
}

func (s *Server) sendHub(message string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.hub != nil {
		return s.hub.WriteMessage(websocket.TextMessage, []byte(message))
	}
	return nil
}

func (s *Server) connect(conn *websocket.Conn, data []string) {
	if data[1] != "init" {
		// обычное действие, не инициализация
		return
	}

	if data[0] == "hub" {
		if s.hub != nil {
			if err := s.sendHub("ping"); err == nil {
				// хаб активен уже
				return
			}
		}

		s.hub = conn
		s.sendHub("autoqr:" + os.Getenv("AUTO_QR_CODE"))
		log.Println("hub connected")
	}
}

func (s *Server) auth(conn *websocket.Conn, data []string) {
	if data[2] != s.token {
		conn.WriteMessage(websocket.TextMessage, []byte("error:wrong"))

		return
	}

	if s.hub == nil {
		conn.WriteMessage(websocket.TextMessage, []byte("error:hub_disconnected"))
		return
	}
	if s.user != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("error:user_already"))
		return
	}

	if s.hub != nil && s.user == nil && data[0] == "user" {
		// если нет активного юзера, и идет запрос на подключение
		s.user = conn
		s.sendUser("connected")

		// отправка данных о камерах
		devices, _ := getActiveDevicesForWs()
		selectedCamera := ""
		s.sendUser("devices:" + strings.Join(devices, "|"))
		if s.camera.device != nil {
			selectedCamera = s.camera.device.Name() + ":" + strconv.FormatBool(s.camera.isPtz)
		} else {
			selectedCamera = ""
		}
		s.sendUser("selected-device:" + selectedCamera)

		// отправляем сообщение хабу
		s.sendHub("connected:" + conn.RemoteAddr().String())

		return
	}
}

func (s *Server) getToken(conn *websocket.Conn, data []string) {
	length, _ := strconv.Atoi(os.Getenv("TOKEN_LENGTH"))
	s.token = randToken(length)

	s.sendHub("token:" + s.token)
}

func (s *Server) switchDevice(conn *websocket.Conn, data []string) {
	if s.camera.device != nil {
		s.camera.cancel()
		s.camera.context.Done()
	}

	camera := NewCamera(30, 600, 400, data[2])
	if err := camera.init(); err != nil {
		s.sendUser("error:device")
		return
	}

	s.camera = camera
}

func (s *Server) move(conn *websocket.Conn, data []string) {
	var cmd uint32
	value, err := strconv.Atoi(data[3])
	if err != nil {
		return
	}

	switch data[2] {
	case "left":
		cmd = CTRL_HORIZONTAL
		value *= -1
	case "right":
		cmd = CTRL_HORIZONTAL
	case "top":
		cmd = CTRL_VERTICAL
		value *= -1
	case "bottom":
		cmd = CTRL_VERTICAL
	}

	s.camera.sendCommand(cmd, int32(value))
}

func (s *Server) zoom(conn *websocket.Conn, data []string) {
	value, err := strconv.Atoi(data[2])

	if err != nil {
		return
	}

	s.camera.sendCommand(CTRL_ZOOM, int32(value)*100)
}

func (s *Server) changeScene(conn *websocket.Conn, data []string) {
	log.Printf("scene change to %s", data[2])
	s.scene = data[2]
}

func (s *Server) faceDetection(conn *websocket.Conn, data []string) {
	log.Printf("user set face detector to: %s", data[2])
	if data[2] == "true" {
		s.camera.faceEnabled = true
	} else {
		s.camera.faceEnabled = false
	}
}

func (s *Server) stream() {
	log.Println("started stream video")
	screen := captureScreen()

	idx := 0
	for ; ; idx++ {
		if s.user == nil {
			log.Println("user is null")
			continue
		}

		var frame []byte

		switch s.scene {
		case "merge":
			frame = <-merge(s.camera.frames, screen)
		case "camera":
			frame = <-s.camera.frames
		case "screen":
			frame = <-screen
		}

		if idx%20 == 0 && s.scene == "camera" && s.camera.isPtz {
			s.camera.runFaceDetect(frame)
		}

		// for frame := range s.frame {
		str := base64.StdEncoding.EncodeToString(frame)
		log.Println(string(str))
		urldata := "data:image/jpeg;base64," + string(str)
		server.sendUser(urldata)
		// }
	}
}
