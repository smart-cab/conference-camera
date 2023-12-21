package studio

import (
	"bytes"
	"image"
	"image/draw"
	"image/jpeg"

	"github.com/nfnt/resize"
)

func MergeImages(cameraImage []byte, monitorImage []byte) image.Image {
	var cameraImg, monitorImg image.Image

	cameraImg, _ = jpeg.Decode(bytes.NewReader(cameraImage))
	monitorImg, _ = jpeg.Decode(bytes.NewReader(monitorImage))

	cameraNewImg := resize.Resize(600, 0, cameraImg, resize.Lanczos3)
	monitorNewImg := resize.Resize(1280, 720, monitorImg, resize.Lanczos3)

	newImg := image.NewRGBA(image.Rect(0, 0, 1280, 720))

	draw.Draw(newImg, newImg.Bounds(), monitorNewImg, image.Point{0, 0}, draw.Src)
	draw.Draw(newImg, image.Rect(0, 0, 600, 400), cameraNewImg, image.Point{0, 0}, draw.Src)

	return newImg
}
