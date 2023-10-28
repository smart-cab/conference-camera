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

func Init() error {
	camera, err := device.Open(
		"/dev/video0",
		device.WithPixFormat(v4l2.PixFormat{PixelFormat: v4l2.PixelFmtMJPEG, Width: 640, Height: 480}),
	)
	if err != nil {
		return err
	}

	if err := camera.Start(context.TODO()); err != nil {
		return err
	}

	Frames = camera.GetOutput()
	return nil
}

func SendCmd(cmd string) {
	// TODO
	fmt.Printf("Send command to PTZ camera: %s", cmd)
}
