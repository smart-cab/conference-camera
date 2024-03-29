package core

import (
	"conferencecam/handlers"
	"conferencecam/middlewares"
)

func (api *Api) Routes() {
	r := api.App.Group("/api/v1")

	r.Use(middlewares.Handler(api.Log))
	r.Use(middlewares.ErrorHandler())
	r.Use(middlewares.Cors())

	{
		r.GET("/ws", handlers.WebSocket)
		r.GET("/ping", handlers.Ping)
		r.GET("/video", handlers.ServeVideoStream)
		r.GET("/studio", handlers.ServeStudioStream)
		r.GET("/hub", handlers.Validate)
	}
}
