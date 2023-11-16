package ptz

import (
	"context"
	"log"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

type Camera struct {
	Device *device.Device
	Frames <-chan []byte
	Status bool
	stop   context.CancelFunc
}

var PTZ *Camera = &Camera{
	nil, nil, false, nil,
}

const CTRL_HORIZONTAL uint32 = 0x009a0904
const CTRL_VERTICAL uint32 = 0x009a0905

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

func Init(path string) (*Camera, error) {
	device, err := device.Open(
		path,
		device.WithPixFormat(v4l2.PixFormat{PixelFormat: v4l2.PixelFmtMJPEG, Width: 640, Height: 330}),
		device.WithFPS(60),
	)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.TODO())

	if err := device.Start(ctx); err != nil {
		log.Fatalf("stream capture: %s", err)
	}

	return &Camera{
		Device: device,
		Frames: device.GetOutput(),
		Status: true,
		stop:   cancel,
	}, nil
}

func (c *Camera) Close() error {
	return c.Device.Close()
}

func (c *Camera) Control(cmd uint32, value int32) error {
	if err := c.Device.SetControlValue(cmd, value); err != nil {
		return err
	}

	return nil
}
