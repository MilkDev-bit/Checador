package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"paselista/database"
	"paselista/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	FirstName      string `json:"first_name" binding:"required"`
	LastName       string `json:"last_name" binding:"required"`
	ProjectName    string `json:"project_name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=8"`
	RecaptchaToken string `json:"recaptcha_token" binding:"required"`
}

type LoginRequest struct {
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required"`
	RecaptchaToken string `json:"recaptcha_token" binding:"required"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Todos los campos son requeridos, incluyendo el reCAPTCHA."})
		return
	}

	if !verifyRecaptcha(req.RecaptchaToken) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verificación reCAPTCHA fallida. Por favor, inténtalo de nuevo."})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing password"})
		return
	}

	var user models.User
	err = database.DB.QueryRow(
		`INSERT INTO users (first_name, last_name, project_name, email, password_hash)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, first_name, last_name, project_name, email, role, created_at`,
		req.FirstName, req.LastName, req.ProjectName, req.Email, string(hash),
	).Scan(&user.ID, &user.FirstName, &user.LastName, &user.ProjectName, &user.Email, &user.Role, &user.CreatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created", "user": user})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Correo, contraseña y reCAPTCHA son requeridos."})
		return
	}

	if !verifyRecaptcha(req.RecaptchaToken) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verificación reCAPTCHA fallida. Por favor, inténtalo de nuevo."})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	var user models.User
	var avatarURL, coverURL sql.NullString
	err := database.DB.QueryRow(
		`SELECT id, first_name, last_name, project_name, email, role, password_hash, avatar_url, cover_url FROM users WHERE email = $1`,
		req.Email,
	).Scan(&user.ID, &user.FirstName, &user.LastName, &user.ProjectName, &user.Email, &user.Role, &user.PasswordHash, &avatarURL, &coverURL)

	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}
	if coverURL.Valid {
		user.CoverURL = coverURL.String
	}

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := generateToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	// Set HttpOnly, Secure, SameSite=Strict cookie
	c.SetCookie("jwt_token", token, int(24*time.Hour/time.Second), "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func Logout(c *gin.Context) {
	// Clear the cookie by setting MaxAge to -1
	c.SetCookie("jwt_token", "", -1, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func Me(c *gin.Context) {
	userID := c.GetString("user_id")

	var user models.User
	var avatarURL, coverURL sql.NullString
	err := database.DB.QueryRow(
		`SELECT id, first_name, last_name, project_name, email, role, avatar_url, cover_url, created_at FROM users WHERE id = $1`,
		userID,
	).Scan(&user.ID, &user.FirstName, &user.LastName, &user.ProjectName, &user.Email, &user.Role, &avatarURL, &coverURL, &user.CreatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}
	if coverURL.Valid {
		user.CoverURL = coverURL.String
	}

	c.JSON(http.StatusOK, user)
}

func generateToken(userID, email, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "changeme_secret_key_32chars_long!"
	}

	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

type RecaptchaResponse struct {
	Success bool     `json:"success"`
	Errors  []string `json:"error-codes"`
}

func verifyRecaptcha(token string) bool {
	secret := os.Getenv("RECAPTCHA_SECRET_KEY")
	if secret == "" {
		// Google's official testing secret key (always returns success)
		secret = "6LeIxAcTAAAAAGG-vFI1TnRWxMZNFuojJ4WifJWe"
	}

	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify",
		url.Values{"secret": {secret}, "response": {token}})
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var result RecaptchaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false
	}

	return result.Success
}
