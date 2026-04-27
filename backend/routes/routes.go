package routes

import (
	"paselista/handlers"
	"paselista/middleware"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	api := r.Group("/api")

	// Health check (used by Railway)
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public routes
	api.POST("/auth/register", handlers.Register)
	api.POST("/auth/login", handlers.Login)
	api.POST("/auth/logout", handlers.Logout)
	api.GET("/qr", handlers.GenerateQR)
	api.GET("/checkin-url", handlers.GetCheckinURL)

	// User protected routes
	auth := api.Group("/")
	auth.Use(middleware.AuthRequired())
	{
		auth.GET("/auth/me", handlers.Me)
		auth.GET("/checks/status", handlers.GetCheckStatus)
		auth.PUT("/profile/email", handlers.UpdateEmail)
		auth.PUT("/profile/password", handlers.UpdatePassword)
		auth.PUT("/profile/avatar", handlers.UpdateAvatar)
		auth.PUT("/profile/cover", handlers.UpdateCover)
		auth.POST("/checks", handlers.RegisterCheck)
		auth.GET("/checks", handlers.GetMyRecords)
		auth.GET("/checks/:id/route", handlers.GetRecordRoute)
		auth.POST("/location-points", handlers.AddLocationPoint)
		auth.POST("/location-points/batch", handlers.AddLocationPointsBatch)
	}

	// Admin protected routes
	admin := api.Group("/admin")
	admin.Use(middleware.AdminRequired())
	{
		admin.GET("/stats", handlers.AdminGetStats)
		admin.GET("/records", handlers.AdminGetRecords)
		admin.GET("/records/:id/route", handlers.AdminGetRecordRoute)
		admin.GET("/records/:id/photos", handlers.AdminGetRecordPhotos)
		admin.GET("/users", handlers.AdminGetUsers)
		admin.GET("/projects", handlers.AdminGetProjects)
	}
}
