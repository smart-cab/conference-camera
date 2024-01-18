package ptz

import (
	"context"
	"log"
	"time"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

type ICamera struct {
	FPS    uint32
	Width  uint32
	Height uint32

	Device  *device.Device
	Context context.Context
	Cancel  context.CancelFunc

	Frames <-chan []byte

	CurrentX int32
	CurrentY int32
}

const CTRL_HORIZONTAL uint32 = 0x009a0904
const CTRL_VERTICAL uint32 = 0x009a0905
const CTRL_ZOOM uint32 = 0x009a090d

var Camera = ICamera{}

func (c *ICamera) Init(path string) error {
	var err error

	if c.Device != nil {
		if c.Device.Name() != path {
			c.Cancel()
			c.Device.Close()
		}
		return nil
	}

	c.Device, err = device.Open(
		path,
		device.WithPixFormat(v4l2.PixFormat{PixelFormat: v4l2.PixelFmtMJPEG, Width: c.Width, Height: c.Height}),
		device.WithFPS(c.FPS),
	)
	if err != nil {
		return err
	}

	c.Context, c.Cancel = context.WithCancel(context.Background())

	if err := c.Device.Start(c.Context); err != nil {
		log.Fatalf("stream capture: %s", err)
	}

	c.Frames = c.Device.GetOutput()

	// Fix move camera
	if err := c.SendCmd(CTRL_HORIZONTAL, 0); err == nil {
		c.SendCmd(CTRL_VERTICAL, 0)
	}

	return nil
}

func (c *ICamera) Close() error {
	if c.Device == nil {
		return nil
	}

	c.Cancel()

	return c.Device.Close()
}

func (c *ICamera) SendCmd(cmd uint32, value int32) error {
	// TODO
	if err := c.Device.SetControlValue(cmd, value); err != nil {
		log.Printf("ERROR PTZ CAMERA COMMAND: %s", err.Error())
		return err
	}

	if value != 0 {
		log.Printf("Reset horizontal pos to zero")
		c.SendCmd(CTRL_HORIZONTAL, 0)
		c.SendCmd(CTRL_VERTICAL, 0)
	}

	if cmd == CTRL_HORIZONTAL {
		c.CurrentX += value
	} else if cmd == CTRL_VERTICAL {
		c.CurrentY += value
	}

	return nil
}

func (c *ICamera) CenterCamera() {
	c.SendCmd(CTRL_HORIZONTAL, -c.CurrentX)
	time.Sleep(time.Second) // small fix
	c.SendCmd(CTRL_VERTICAL, -c.CurrentY)
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
