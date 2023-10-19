package handlers

import (
	"conferencecam/types"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.JSON(200, types.RESPONSE{
		Success: true,
		Data: types.JSON{
			"message": "pong",
		},
	})
}

func Connect(c *gin.Context) {

	c.JSON(200, types.RESPONSE{
		Success: true,
		Data: types.JSON{
			"message": "pong",
		},
	})
}
