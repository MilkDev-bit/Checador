package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"paselista/database"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UpdateEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

type UpdateAvatarRequest struct {
	AvatarURL string `json:"avatar_url" binding:"required"`
}

// UpdateEmail updates the authenticated user's email.
func UpdateEmail(c *gin.Context) {
	userID := c.GetString("user_id")

	var req UpdateEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Correo electrónico válido requerido"})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	_, err := database.DB.Exec(
		`UPDATE users SET email = $1 WHERE id = $2`,
		req.Email, userID,
	)
	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			c.JSON(http.StatusConflict, gin.H{"error": "Ese correo ya está registrado por otra cuenta"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el correo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Correo actualizado correctamente", "email": req.Email})
}

// UpdatePassword verifies the current password and sets a new one.
func UpdatePassword(c *gin.Context) {
	userID := c.GetString("user_id")

	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Se requieren la contraseña actual y la nueva (mín. 8 caracteres)"})
		return
	}

	var currentHash string
	err := database.DB.QueryRow(`SELECT password_hash FROM users WHERE id = $1`, userID).Scan(&currentHash)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error de base de datos"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentHash), []byte(req.CurrentPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "La contraseña actual es incorrecta"})
		return
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar la contraseña"})
		return
	}

	_, err = database.DB.Exec(`UPDATE users SET password_hash = $1 WHERE id = $2`, string(newHash), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la contraseña"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contraseña actualizada correctamente"})
}

// UpdateAvatar stores a base64-encoded image data URL as the user's avatar.
func UpdateAvatar(c *gin.Context) {
	userID := c.GetString("user_id")

	var req UpdateAvatarRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Imagen requerida"})
		return
	}

	// Only accept base64 data URLs with image MIME types
	if !strings.HasPrefix(req.AvatarURL, "data:image/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de imagen no válido"})
		return
	}

	// ~700 KB base64 ≈ ~500 KB binary
	if len(req.AvatarURL) > 700_000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Imagen demasiado grande. Máximo 500 KB"})
		return
	}

	_, err := database.DB.Exec(`UPDATE users SET avatar_url = $1 WHERE id = $2`, req.AvatarURL, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar la imagen"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Foto de perfil actualizada", "avatar_url": req.AvatarURL})
}
