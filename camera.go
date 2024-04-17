package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"time"

	pigo "github.com/esimov/pigo/core"
	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

const CTRL_HORIZONTAL uint32 = 0x009a0904
const CTRL_VERTICAL uint32 = 0x009a0905
const CTRL_ZOOM uint32 = 0x009a090d

type Camera struct {
	fps    uint32
	width  uint32
	height uint32
	path   string

	device  *device.Device
	context context.Context
	cancel  context.CancelFunc
	isPtz   bool

	faceEnabled bool
	faceFinder  *pigo.Pigo

	frames <-chan []byte
}

func NewCamera(fps, width, height uint32, path string) *Camera {
	return &Camera{
		fps:    fps,
		width:  width,
		height: height,
		path:   path,
	}
}

func (c *Camera) init() error {
	var err error
	c.device, err = device.Open(
		c.path,
		device.WithPixFormat(v4l2.PixFormat{PixelFormat: v4l2.PixelFmtMJPEG, Width: c.width, Height: c.height}),
		device.WithFPS(c.fps),
	)
	if err != nil {
		return err
	}

	c.context, c.cancel = context.WithCancel(context.Background())

	if err := c.device.Start(c.context); err != nil {
		log.Fatalf("stream capture: %s", err)
	}

	c.frames = c.device.GetOutput()
	if err := c.sendCommand(CTRL_ZOOM, 100); err == nil {
		c.isPtz = true
	}

	if c.isPtz {
		if err := c.initFace(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Camera) sendCommand(cmd uint32, value int32) error {
	if err := c.device.SetControlValue(cmd, value); err != nil {
		return err
	}

	if value != 0 {
		c.sendCommand(CTRL_HORIZONTAL, 0)
		c.sendCommand(CTRL_VERTICAL, 0)
	}

	return nil
}

func (c *Camera) initFace() error {
	model, err := os.ReadFile("./facefinder.model")
	if err != nil {
		return fmt.Errorf("failed to load face finder model: %s", err)
	}
	p := pigo.NewPigo()
	c.faceFinder, err = p.Unpack(model)
	if err != nil {
		c.faceFinder = nil
		return fmt.Errorf("failed to initialize face classifier: %s", err)
	}
	return nil
}

func (c *Camera) runFaceDetect(frame []byte) error {
	if !c.faceEnabled {
		return nil
	}

	img, _, err := image.Decode(bytes.NewReader(frame))
	if err != nil {
		log.Println(err)
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

	dets := c.faceFinder.RunCascade(params, 0.0)
	dets = c.faceFinder.ClusterDetections(dets, 0)

	// w.Write(frame)

	// drawer := gg.NewContext(bounds.Max.X, bounds.Max.Y)
	// drawer.DrawImage(img, 0, 0)

	if len(dets) == 0 {
		return nil
	}

	for _, person := range dets {
		if person.Q < 5.0 {
			continue
		}

		centerX, centerY := int(c.width/2), int(c.height/2)
		x, y := person.Col-person.Scale/2, person.Row-person.Scale/2
		x, y = x/100*100, y/100*100
		log.Println(x, y)

		var step int32 = 100
		timeWait := 50

		moveX, moveY := float64(centerX-x)/100, float64(centerY-y)/100
		if moveX < 0 {
			for i := 0; i < int(math.Abs(moveX)); i++ {
				c.sendCommand(CTRL_HORIZONTAL, step)
				time.Sleep(time.Millisecond * time.Duration(timeWait))
			}
		}
		if moveX > 0 {
			for i := 0; i < int(math.Abs(moveX)); i++ {
				c.sendCommand(CTRL_HORIZONTAL, -step)
				time.Sleep(time.Millisecond * time.Duration(timeWait))
			}
		}
		if moveY < 0 {
			for i := 0; i < int(math.Abs(moveY)); i++ {
				c.sendCommand(CTRL_VERTICAL, step)
				time.Sleep(time.Millisecond * time.Duration(timeWait))
			}
		}
		if moveY > 0 {
			for i := 0; i < int(math.Abs(moveY)); i++ {
				c.sendCommand(CTRL_VERTICAL, -step)
				time.Sleep(time.Millisecond * time.Duration(timeWait))
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
