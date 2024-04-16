package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// loading .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed load env config: %s", err.Error())
	}

	devices := getActiveDevices()
	if len(devices) > 0 {
		camera := NewCamera(30, 600, 400, devices[0].Name())
		if err := camera.init(); err != nil {
			log.Fatalf("failed start camera: %s", err.Error())
		}

		server = NewServer(camera)
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
