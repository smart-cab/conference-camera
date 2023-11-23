package middlewares

import (
	"conferencecam/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, e := range c.Errors {
			if apiErr, ok := e.Err.(*types.ApiError); ok {
				c.JSON(apiErr.Code, types.RESPONSE{
					Success: false,
					Error:   apiErr.Msg,
				})
			} else {
				c.JSON(http.StatusInternalServerError, types.RESPONSE{
					Success: false,
					Error:   "system error",
				})
			}
		}

	}
}
