package ptz

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"io"
	"os"

	pigo "github.com/esimov/pigo/core"
	"github.com/fogleman/gg"
	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
)

var (
	Camera     *device.Device
	Frames     <-chan []byte
	faceFinder *pigo.Pigo
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

	if err := InitFaceDetect(); err != nil {
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

func RunFaceDetect(w io.Writer, frame []byte) error {
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

	dets := faceFinder.RunCascade(params, 0.0)
	dets = faceFinder.ClusterDetections(dets, 0)

	drawer := gg.NewContext(bounds.Max.X, bounds.Max.Y)
	drawer.DrawImage(img, 0, 0)

	for _, det := range dets {
		if det.Q >= 5.0 {
			drawer.DrawRectangle(
				float64(det.Col-det.Scale/2),
				float64(det.Row-det.Scale/2),
				float64(det.Scale),
				float64(det.Scale),
			)

			drawer.SetLineWidth(3.0)
			drawer.SetStrokeStyle(gg.NewSolidPattern(color.RGBA{R: 255, G: 0, B: 0, A: 255}))
			drawer.Stroke()
		}
	}

	return drawer.EncodePNG(w)
}

func InitFaceDetect() error {
	model, err := os.ReadFile("./ptz/facefinder.model")
	if err != nil {
		return fmt.Errorf("failed to load face finder model: %s", err)
	}
	p := pigo.NewPigo()
	faceFinder, err = p.Unpack(model)
	if err != nil {
		faceFinder = nil
		return fmt.Errorf("failed to initialize face classifier: %s", err)
	}
	return nil
}
