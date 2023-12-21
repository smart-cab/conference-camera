package handlers

import (
	"bytes"
	"conferencecam/core/studio"
	"conferencecam/ptz"
	"fmt"
	"image/jpeg"
	"mime/multipart"
	"net/textproto"

	"github.com/gin-gonic/gin"
	"github.com/kbinani/screenshot"
	"github.com/sirupsen/logrus"
)

func ServeStudioStream(c *gin.Context) {
	log := c.MustGet("log").(*logrus.Logger)

	for {
		if ptz.Camera != nil {
			break
		}
	}

	mimeWriter := multipart.NewWriter(c.Writer)
	c.Writer.Header().Set("Content-Type", fmt.Sprintf("multipart/x-mixed-replace; boundary=%s", mimeWriter.Boundary()))
	partHeader := make(textproto.MIMEHeader)
	partHeader.Add("Content-Type", "image/jpeg")

	var frame []byte
	for frame = range ptz.Frames {
		captureScreenImg := captureScreen()
		image := studio.MergeImages(frame, captureScreenImg)

		// convert to []byte for send http
		buf := new(bytes.Buffer)
		err := jpeg.Encode(buf, image, nil)
		if err != nil {
			log.Printf("error decode merged image: %s", err)
			return
		}
		send_s3 := buf.Bytes()

		partWriter, err := mimeWriter.CreatePart(partHeader)
		if err != nil {
			log.Printf("failed to create multi-part writer: %s", err)
			return
		}

		if _, err := partWriter.Write(send_s3); err != nil {
			log.Printf("failed to write image: %s", err)
			return
		}
	}
}

func captureScreen() []byte {
	bounds := screenshot.GetDisplayBounds(0)

	img, err := screenshot.CaptureRect(bounds.Bounds())
	if err != nil {
		fmt.Println("Failed to capture screen:", err)
		return nil
	}

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, nil); err != nil {
		fmt.Println("Failed to encode image:", err)
		return nil
	}

	return buf.Bytes()
}
