package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

func GenerateQR(c *gin.Context) {
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	checkinURL := frontendURL + "/checkin"

	png, err := qrcode.Encode(checkinURL, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating QR"})
		return
	}

	c.Data(http.StatusOK, "image/png", png)
}

func GetCheckinURL(c *gin.Context) {
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	c.JSON(http.StatusOK, gin.H{"url": frontendURL + "/checkin"})
}
