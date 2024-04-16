package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"math/rand"
	"time"

	"github.com/disintegration/imaging"
	"github.com/vladimirvivien/go4vl/device"
)

func randToken(n int) string {
	rand.NewSource(time.Now().UnixNano()) // fix repeat tokens

	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func getActiveDevices() []*device.Device {
	var result []*device.Device
	devices, err := device.GetAllDevicePaths()

	if err != nil {
		return nil
	}

	for _, d := range devices {
		if temp_device, err := device.Open(d); err == nil {
			result = append(result, temp_device)
		}
	}

	return result
}

func getActiveDevicesForWs() ([]string, error) {
	var result []string

	for _, d := range getActiveDevices() {
		result = append(result, d.Name()+":"+d.Capability().Card)
	}

	return result, nil
}

func resizeImage(img image.Image, width, height int) image.Image {
	return imaging.Resize(img, width, height, imaging.Lanczos)
}

func overlayImage(background, overlay image.Image, posX, posY int) image.Image {
	return imaging.Overlay(background, overlay, image.Pt(posX, posY), 1)
}

func merge(cameraImage, screenImage <-chan []byte) <-chan []byte {
	output := make(chan []byte)

	go func() {
		defer close(output)

		for {
			cameraBytes := <-cameraImage
			screenBytes := <-screenImage

			cameraImg, err := jpeg.Decode(bytes.NewReader(cameraBytes))
			if err != nil {
				continue
			}
			screenImg, err := jpeg.Decode(bytes.NewReader(screenBytes))
			if err != nil {
				continue
			}

			if cameraImg == nil || screenImg == nil {
				continue
			}

			resizedCamera := resizeImage(cameraImg, 640, 360)
			resizedScreen := resizeImage(screenImg, 1280, 720)

			resultImg := overlayImage(resizedScreen, resizedCamera, 0, 0)

			var buf bytes.Buffer
			err = jpeg.Encode(&buf, resultImg, nil)

			if err != nil {
				continue
			}

			output <- buf.Bytes()
		}
	}()

	return output
}
