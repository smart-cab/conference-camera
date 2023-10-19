package core

import (
	"conferencecam/handlers"
	"conferencecam/middlewares"
)

func (api *Api) Routes() {
	api.App.GET("/ws", handlers.WebSocket)

	r := api.App.Group("/api/v1")

	r.Use(middlewares.Handler(api.Log))
	r.Use(middlewares.ErrorHandler())
	r.Use(middlewares.Cors())

	{
		r.GET("/ping", handlers.Ping)
		r.GET("/hub", handlers.Validate)
	}
}
