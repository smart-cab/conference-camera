package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	// loading .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed load env config: %s", err.Error())
	}

	devices := getActiveDevices()
	if len(devices) > 0 {
		var camera, screen *Camera
		for _, device := range devices {
			if strings.Contains(device.Capability().Card, "PTZ") {
				camera = NewCamera(30, 600, 400, device.Name())
				if err := camera.init(); err != nil {
					log.Fatalf("failed start camera: %s", err.Error())
				}
			}
			log.Println(device.Capability().Card)
			if strings.Contains(device.Capability().Card, "HDMI") {
				screen = NewCamera(30, 1280, 720, device.Name())
				if err := screen.init(); err != nil {
					log.Fatalf("failed start screen: %s", err.Error())
				}
			}
		}

		if camera == nil || screen == nil {
			log.Fatalf("cannot be started a project, because camera is %v and screen is %v", camera, screen)
		}

		server = NewServer(camera, screen)

		log.Println("start stream method")
		go server.stream()
	} else {
		log.Fatalf("server cannot be started without camera")
		return
	}

	// create server and handle a websocket
	http.HandleFunc("/", handler)

	// start a project
	log.Println("server successful started")
	if err := http.ListenAndServe(fmt.Sprintf("%s:8888", os.Getenv("IP")), nil); err != nil {
		log.Fatalf("server failed: %s", err.Error())
	}
}
