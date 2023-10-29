package ptz

import (
	"context"
	"fmt"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

var (
	Camera *device.Device
	Frames <-chan []byte
)

func Init(path string) error {
	var err error
	Camera, err = device.Open(
		path,
		device.WithPixFormat(v4l2.PixFormat{PixelFormat: v4l2.PixelFmtMJPEG, Width: 1920, Height: 1080}),
	)
	if err != nil {
		return err
	}

	if err := Camera.Start(context.TODO()); err != nil {
		return err
	}

	Frames = Camera.GetOutput()
	return nil
}

func Close() error {
	if Camera == nil {
		return nil
	}
	return Camera.Close()
}

func SendCmd(cmd string) {
	// TODO
	fmt.Printf("Send command to PTZ camera: %s", cmd)
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
