package ptz

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"io"
	"log"
	"math"
	"os"
	"time"

	pigo "github.com/esimov/pigo/core"
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

	FaceEnabled bool
	FaceFinder  *pigo.Pigo

	IsPTZ bool
}

const CTRL_HORIZONTAL uint32 = 0x009a0904
const CTRL_VERTICAL uint32 = 0x009a0905
const CTRL_ZOOM uint32 = 0x009a090d

var Camera = ICamera{
	FPS:    60,
	Width:  1280,
	Height: 720,
}

func (c *ICamera) Init(path string) error {
	var err error

	if c.Device != nil {
		if c.Device.Name() != path {
			c.Cancel()
			c.Device.Close()
		} else {
			return nil
		}
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
	if err := c.InitFaceDetect(); err != nil {
		return err
	}

	if err := c.Device.Start(c.Context); err != nil {
		log.Fatalf("stream capture: %s", err)
	}

	c.Frames = c.Device.GetOutput()
	c.IsPTZ = false

	// Fix move camera
	err = c.SendCmd(CTRL_HORIZONTAL, 0)
	if err == nil {
		c.IsPTZ = true
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
	if err := c.Device.SetControlValue(cmd, value); err != nil {
		// log.Printf("ERROR PTZ CAMERA COMMAND: %s", err.Error())
		return err
	}

	if value != 0 {
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

func (c *ICamera) InitFaceDetect() error {
	model, err := os.ReadFile("./ptz/facefinder.model")
	if err != nil {
		return fmt.Errorf("failed to load face finder model: %s", err)
	}
	p := pigo.NewPigo()
	c.FaceFinder, err = p.Unpack(model)
	if err != nil {
		c.FaceFinder = nil
		return fmt.Errorf("failed to initialize face classifier: %s", err)
	}
	return nil
}

func (c *ICamera) RunFaceDetect(w io.Writer, frame []byte) error {
	if !c.FaceEnabled {
		return nil
	}

	img, _, err := image.Decode(bytes.NewReader(frame))
	if err != nil {
		return err
	}

	src := img.(*image.YCbCr)
	bounds := img.Bounds()
	params := pigo.CascadeParams{
		MinSize:     100,
		MaxSize:     600,
		ShiftFactor: 0.15,
		ScaleFactor: 1.1,
		ImageParams: pigo.ImageParams{
			Pixels: src.Y,
			Rows:   bounds.Dy(),
			Cols:   bounds.Dx(),
			Dim:    bounds.Dx(),
		},
	}

	dets := c.FaceFinder.RunCascade(params, 0.0)
	dets = c.FaceFinder.ClusterDetections(dets, 0)

	w.Write(frame)

	// drawer := gg.NewContext(bounds.Max.X, bounds.Max.Y)
	// drawer.DrawImage(img, 0, 0)

	if len(dets) == 0 {
		return nil
	}

	for _, person := range dets {
		if person.Q < 5.0 {
			continue
		}

		centerX, centerY := 500, 300
		x, y := person.Col-person.Scale/2, person.Row-person.Scale/2
		x, y = x/100*100, y/100*100

		moveX, moveY := float64(centerX-x)/100, float64(centerY-y)/100
		if moveX < 0 {
			for i := 0; i < int(math.Abs(moveX)); i++ {
				c.SendCmd(CTRL_HORIZONTAL, 200)
				time.Sleep(time.Millisecond * 100)
			}
		}
		if moveX > 0 {
			for i := 0; i < int(math.Abs(moveX)); i++ {
				time.Sleep(time.Millisecond * 100)
				c.SendCmd(CTRL_HORIZONTAL, -200)
			}
		}
		if moveY < 0 {
			for i := 0; i < int(math.Abs(moveY)); i++ {
				c.SendCmd(CTRL_VERTICAL, 100)
				time.Sleep(time.Millisecond * 100)
			}
		}
		if moveY > 0 {
			for i := 0; i < int(math.Abs(moveY)); i++ {
				time.Sleep(time.Millisecond * 100)
				c.SendCmd(CTRL_VERTICAL, -100)
			}
		}

		break
	}

	// 485
	// 485 / 100 round * 100

	return nil
	// return nil
	// return drawer.EncodePNG(w)
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
