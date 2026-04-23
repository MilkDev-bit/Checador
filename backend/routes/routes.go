package routes

import (
	"paselista/handlers"
	"paselista/middleware"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	api := r.Group("/api")

	// Public routes
	api.POST("/auth/register", handlers.Register)
	api.POST("/auth/login", handlers.Login)
	api.GET("/qr", handlers.GenerateQR)
	api.GET("/checkin-url", handlers.GetCheckinURL)

	// Protected routes
	auth := api.Group("/")
	auth.Use(middleware.AuthRequired())
	{
		auth.GET("/auth/me", handlers.Me)
		auth.POST("/checks", handlers.RegisterCheck)
		auth.GET("/checks", handlers.GetMyRecords)
		auth.GET("/checks/:id/route", handlers.GetRecordRoute)
		auth.POST("/location-points", handlers.AddLocationPoint)
	}
}
