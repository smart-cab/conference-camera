package handlers

import (
	"conferencecam/types"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Validate(c *gin.Context) {
	token := c.Query("token")
	timestamp := time.Now().Unix()

	auth := false
	for i := 0; i <= 30; i++ {
		// check last 30 seconds tokens
		hasher := md5.New()
		io.WriteString(hasher, strconv.Itoa(int(timestamp)-i)+os.Getenv("AUTH_TOKEN"))
		hashBytes := hasher.Sum(nil)
		hashString := fmt.Sprintf("%x", hashBytes)
		if hashString == token {
			auth = true
			break
		}
	}

	c.JSON(200, types.RESPONSE{
		Success: true,
		Data: types.JSON{
			"message": auth,
		},
	})
}
