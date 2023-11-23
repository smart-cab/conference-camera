package handlers

import (
	"conferencecam/ptz"
	"fmt"
	"mime/multipart"
	"net/textproto"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ServeVideoStream(c *gin.Context) {
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
		partWriter, err := mimeWriter.CreatePart(partHeader)
		if err != nil {
			log.Printf("failed to create multi-part writer: %s", err)
			return
		}

		if _, err := partWriter.Write(frame); err != nil {
			log.Printf("failed to write image: %s", err)
			return
		}
	}
}
