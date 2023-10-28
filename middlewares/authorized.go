package middlewares

import (
	"conferencecam/handlers"
	"conferencecam/types"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authorized(c *gin.Context) {
	if handlers.Client != nil {
		data := strings.Split(handlers.Client.RemoteAddr().String(), ":")

		if data[0] == c.ClientIP() {
			c.Next()
			return
		}
	}

	c.Error(types.UNAUTHORIZED)
	c.Abort()
}
