package ptz

import (
	"context"
	"log"
	"time"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

var (
	Camera   *device.Device
	Frames   <-chan []byte
	Cancel   context.CancelFunc
	CurrentX int32
	CurrentY int32
)

const CTRL_HORIZONTAL uint32 = 0x009a0904
const CTRL_VERTICAL uint32 = 0x009a0905

func Init(path string) error {
	var err error

	if Camera != nil && Camera.Name() == path {
		return nil
	}

	if Camera != nil {
		Cancel()
		Camera.Close()
	}

	Camera, err = device.Open(
		path,
		device.WithPixFormat(v4l2.PixFormat{PixelFormat: v4l2.PixelFmtMJPEG, Width: 640, Height: 330}),
		device.WithFPS(60),
	)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	Cancel = cancel

	if err := Camera.Start(ctx); err != nil {
		log.Fatalf("stream capture: %s", err)
	}

	Frames = Camera.GetOutput()

	// Fix move camera
	if err := SendCmd(CTRL_HORIZONTAL, 0); err == nil {
		SendCmd(CTRL_VERTICAL, 0)
	}

	return nil
}

func Close() error {
	if Camera == nil {
		return nil
	}
	return Camera.Close()
}

func SendCmd(cmd uint32, value int32) error {
	// TODO
	if err := Camera.SetControlValue(cmd, value); err != nil {
		log.Printf("ERROR PTZ CAMERA COMMAND: %s", err.Error())
		return err
	}

	if value != 0 {
		log.Printf("Reset horizontal pos to zero")
		SendCmd(CTRL_HORIZONTAL, 0)
		SendCmd(CTRL_VERTICAL, 0)
	}

	if cmd == CTRL_HORIZONTAL {
		CurrentX += value
	} else if cmd == CTRL_VERTICAL {
		CurrentY += value
	}

	return nil
}

func CenterCamera() {
	SendCmd(CTRL_HORIZONTAL, -CurrentX)
	time.Sleep(time.Second * 1)
	SendCmd(CTRL_VERTICAL, -CurrentY)
}

func GetActiveDevices() []*device.Device {
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

func GetDevices() ([]string, error) {
	var result []string

	for _, d := range GetActiveDevices() {
		result = append(result, d.Name()+":"+d.Capability().Card)
	}

	return result, nil
}
