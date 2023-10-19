package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Handler(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("log", log)

		// Finish handler
		c.Next()
	}
}
