package main

import (
	"log"
	"os"
	"strings"

	"paselista/database"
	"paselista/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	database.Connect()
	database.Migrate()
	database.SeedAdmin()

	r := gin.Default()

	// CORS: allowed origins come exclusively from the FRONTEND_URL env var.
	// Set it to a comma-separated list of origins, e.g.:
	//   FRONTEND_URL=http://localhost:5173,https://myapp.up.railway.app
	allowedOrigins := []string{}
	if frontendURL := os.Getenv("FRONTEND_URL"); frontendURL != "" {
		for _, u := range strings.Split(frontendURL, ",") {
			u = strings.TrimSpace(u)
			if u != "" {
				allowedOrigins = append(allowedOrigins, u)
			}
		}
	}
	if len(allowedOrigins) == 0 {
		log.Fatal("Required environment variable FRONTEND_URL is not set")
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.Static("/uploads", "./uploads")

	routes.Register(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	r.Run(":" + port)
}
