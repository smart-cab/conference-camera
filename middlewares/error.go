package middlewares

import (
	"conferencecam/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		errors := []interface{}{}

		for _, e := range c.Errors {
			if apiErr, ok := e.Err.(*types.ApiError); ok {
				errors = append(errors, apiErr.Msg)
			} else {
				errors = append(errors, "internal error")
			}
		}

		c.JSON(http.StatusInternalServerError, types.RESPONSE{
			Success: false,
			Error:   errors,
		})
	}
}
