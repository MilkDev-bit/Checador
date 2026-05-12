package handlers

import (
	crand "crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"strings"
	"time"

	"paselista/database"
	"paselista/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

func bindingErrorMessage(err error) string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			field := fe.Field()
			tag := fe.Tag()
			switch field {
			case "Email":
				if tag == "email" {
					return "El formato del correo es inválido."
				}
				return "El correo es requerido."
			case "Password":
				if tag == "min" {
					return "La contraseña debe tener al menos 8 caracteres."
				}
				return "La contraseña es requerida."
			case "RecaptchaToken":
				return "El token de reCAPTCHA es requerido."
			case "FirstName":
				return "El nombre es requerido."
			case "LastName":
				return "El apellido es requerido."
			case "ProjectName":
				return "El nombre del proyecto es requerido."
			}
		}
	}
	return "Los datos proporcionados son inválidos."
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindingErrorMessage(err)})
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
			c.JSON(http.StatusConflict, gin.H{"error": "El correo ya está registrado. Inicia sesión."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	// Auto-login: generate token and set cookie so the user doesn't need to log in separately.
	token, err := generateToken(user.ID, user.Email, user.Role)
	if err != nil {
		// Token generation failed — still return success but without cookie.
		log.Printf("WARNING: token generation after register failed: %v", err)
		c.JSON(http.StatusCreated, gin.H{"message": "User created", "user": user})
		return
	}

	c.SetCookie("jwt_token", token, int(24*time.Hour/time.Second), "/", "", true, true)
	c.JSON(http.StatusCreated, gin.H{"message": "User created", "user": user, "token": token})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindingErrorMessage(err)})
		return
	}

	if !verifyRecaptcha(req.RecaptchaToken) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verificación reCAPTCHA fallida. Marca la casilla \"No soy un robot\" y completa el desafío si aparece. Si el problema persiste, recarga la página."})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	var user models.User
	var avatarURL, coverURL sql.NullString
	err := database.DB.QueryRow(
		`SELECT id, first_name, last_name, project_name, email, role, password_hash, avatar_url, cover_url FROM users WHERE email = $1`,
		req.Email,
	).Scan(&user.ID, &user.FirstName, &user.LastName, &user.ProjectName, &user.Email, &user.Role, &user.PasswordHash, &avatarURL, &coverURL)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Correo no registrado. Verifica que sea el mismo con el que te registraste o crea una cuenta nueva."})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno. Intenta de nuevo en unos segundos."})
		return
	}

	if avatarURL.Valid {
		user.AvatarURL = avatarURL.String
	}
	if coverURL.Valid {
		user.CoverURL = coverURL.String
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Contraseña incorrecta. Las contraseñas distinguen mayúsculas y minúsculas. Si olvidaste tu contraseña, contacta al administrador."})
		return
	}

	token, err := generateToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	// Set HttpOnly, Secure cookie AND return token in body for clients that can't use cookies.
	c.SetCookie("jwt_token", token, int(24*time.Hour/time.Second), "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
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
		// No secret key configured — skip reCAPTCHA check and log a warning.
		// Configure RECAPTCHA_SECRET_KEY in production to enforce it.
		log.Println("WARNING: RECAPTCHA_SECRET_KEY not set, skipping reCAPTCHA verification")
		return true
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.PostForm("https://www.google.com/recaptcha/api/siteverify",
		url.Values{"secret": {secret}, "response": {token}})
	if err != nil {
		// Network error reaching Google — fail open so users aren't locked out.
		log.Printf("WARNING: reCAPTCHA verification network error (fail-open): %v", err)
		return true
	}
	defer resp.Body.Close()

	var result RecaptchaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("WARNING: reCAPTCHA response decode error (fail-open): %v", err)
		return true
	}

	if !result.Success {
		log.Printf("reCAPTCHA rejected token, error-codes: %v", result.Errors)
	}
	return result.Success
}

// ─── Forgot / Reset Password ─────────────────────────────────────────────────

// ForgotPassword accepts an email and, if registered, sends a 6-digit reset code.
// Always returns 200 to prevent email enumeration.
func ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Correo inválido."})
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Process in background so response time is identical regardless of whether
	// the email is registered (prevents enumeration).
	go processForgotPassword(req.Email)

	c.JSON(http.StatusOK, gin.H{"message": "Si el correo está registrado recibirás un código en los próximos minutos."})
}

func processForgotPassword(email string) {
	var userID string
	if err := database.DB.QueryRow(`SELECT id FROM users WHERE email = $1`, email).Scan(&userID); err != nil {
		return // user not found — silently ignore
	}

	code, err := generateResetCode()
	if err != nil {
		log.Printf("Error generating reset code: %v", err)
		return
	}

	// Remove any previous tokens for this user before inserting the new one.
	database.DB.Exec(`DELETE FROM password_reset_tokens WHERE user_id = $1`, userID)

	expiresAt := time.Now().Add(15 * time.Minute)
	if _, err = database.DB.Exec(
		`INSERT INTO password_reset_tokens (user_id, code, expires_at) VALUES ($1, $2, $3)`,
		userID, code, expiresAt,
	); err != nil {
		log.Printf("Error storing reset token: %v", err)
		return
	}

	sendResetEmail(email, code)
}

// ResetPassword validates the 6-digit code and updates the password.
func ResetPassword(c *gin.Context) {
	var req struct {
		Email    string `json:"email"    binding:"required,email"`
		Code     string `json:"code"     binding:"required"`
		Password string `json:"password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos. La contraseña debe tener al menos 8 caracteres."})
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Code = strings.TrimSpace(req.Code)

	var userID string
	if err := database.DB.QueryRow(`SELECT id FROM users WHERE email = $1`, req.Email).Scan(&userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Código incorrecto o expirado."})
		return
	}

	var storedCode string
	var expiresAt time.Time
	err := database.DB.QueryRow(
		`SELECT code, expires_at FROM password_reset_tokens WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1`,
		userID,
	).Scan(&storedCode, &expiresAt)

	if err != nil || storedCode != req.Code || time.Now().After(expiresAt) {
		// Delete expired/wrong token to force re-request
		if err == nil {
			database.DB.Exec(`DELETE FROM password_reset_tokens WHERE user_id = $1`, userID)
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Código incorrecto o expirado. Solicita uno nuevo."})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error procesando contraseña."})
		return
	}
	if _, err = database.DB.Exec(`UPDATE users SET password_hash = $1 WHERE id = $2`, string(hash), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error actualizando contraseña."})
		return
	}

	database.DB.Exec(`DELETE FROM password_reset_tokens WHERE user_id = $1`, userID)
	c.JSON(http.StatusOK, gin.H{"message": "Contraseña actualizada correctamente."})
}

// generateResetCode returns a cryptographically random 6-digit string.
func generateResetCode() (string, error) {
	n, err := crand.Int(crand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// sendResetEmail sends the reset code via SMTP.
// Requires SMTP_HOST, SMTP_USER, SMTP_PASS env vars.
// Optional: SMTP_PORT (default 587), SMTP_FROM.
func sendResetEmail(toEmail, code string) {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")

	if host == "" || user == "" || pass == "" {
		log.Printf("WARNING: SMTP not configured — reset code for %s is: %s", toEmail, code)
		return
	}
	if port == "" {
		port = "587"
	}
	if from == "" {
		from = user
	}

	body := fmt.Sprintf(
		"Hola,\n\nTu código de recuperación de contraseña es:\n\n    %s\n\nEste código expira en 15 minutos.\n\nSi no solicitaste este cambio, ignora este correo.\n\n— PaseLista",
		code,
	)
	msg := []byte(fmt.Sprintf(
		"From: PaseLista <%s>\r\nTo: %s\r\nSubject: Código de recuperación - PaseLista\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		from, toEmail, body,
	))

	smtpAuth := smtp.PlainAuth("", user, pass, host)
	if err := smtp.SendMail(host+":"+port, smtpAuth, from, []string{toEmail}, msg); err != nil {
		log.Printf("Error sending reset email to %s: %v", toEmail, err)
	}
}
