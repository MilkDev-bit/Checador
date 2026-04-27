package handlers

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"paselista/database"
	"paselista/models"

	"github.com/gin-gonic/gin"
)

type CheckInRequest struct {
	Type      string  `form:"type"      binding:"required"`
	Timestamp string  `form:"timestamp" binding:"required"`
	Latitude  float64 `form:"latitude"`
	Longitude float64 `form:"longitude"`
}

type LocationPointRequest struct {
	CheckRecordID string  `json:"check_record_id" binding:"required"`
	Latitude      float64 `json:"latitude"        binding:"required"`
	Longitude     float64 `json:"longitude"       binding:"required"`
	Accuracy      float64 `json:"accuracy"`
	RecordedAt    string  `json:"recorded_at"     binding:"required"`
}

// ipAPIResponse maps the fields we use from ip-api.com (free, no key needed).
type ipAPIResponse struct {
	Status  string  `json:"status"`
	Country string  `json:"country"`
	City    string  `json:"city"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Proxy   bool    `json:"proxy"`
	Hosting bool    `json:"hosting"`
	Mobile  bool    `json:"mobile"`
}

// getClientIP extracts the real client IP honoring X-Forwarded-For (set by nginx).
func getClientIP(c *gin.Context) string {
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		// XFF can be "clientIP, proxy1, proxy2" — take the leftmost
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	return c.ClientIP()
}

// haversineKm returns the distance in km between two lat/lon points.
func haversineKm(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371.0
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	return R * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

// runFraudCheck queries ip-api.com and updates the DB record if suspicious.
// Runs in a goroutine — does NOT block the HTTP response.
func runFraudCheck(recordID string, clientIP string, gpsLat, gpsLon float64) {
	// ip-api.com free endpoint (45 req/min, no key required)
	url := fmt.Sprintf(
		"http://ip-api.com/json/%s?fields=status,country,city,lat,lon,proxy,hosting,mobile",
		clientIP,
	)

	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var ipData ipAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&ipData); err != nil || ipData.Status != "success" {
		return
	}

	var reasons []string

	// 1. VPN / Proxy / Datacenter IP
	if ipData.Proxy {
		reasons = append(reasons, "IP detectada como proxy/VPN")
	}
	if ipData.Hosting {
		reasons = append(reasons, "IP de datacenter/hosting (posible VPN)")
	}

	// 2. GPS vs IP geolocation mismatch (only if GPS coords were submitted)
	if gpsLat != 0 || gpsLon != 0 {
		dist := haversineKm(gpsLat, gpsLon, ipData.Lat, ipData.Lon)
		if dist > 500 {
			reasons = append(reasons, fmt.Sprintf(
				"Ubicación GPS y IP separadas %.0f km (IP: %s, %s)",
				dist, ipData.City, ipData.Country,
			))
		}
	}

	isSuspicious := len(reasons) > 0
	reasonText := strings.Join(reasons, " | ")

	database.DB.Exec(
		`UPDATE check_records
		 SET is_suspicious = $1, suspicious_reason = $2, ip_country = $3, ip_city = $4
		 WHERE id = $5`,
		isSuspicious, reasonText, ipData.Country, ipData.City, recordID,
	)
}

func RegisterCheck(c *gin.Context) {
	userID := c.GetString("user_id")

	var req CheckInRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Type != "entry" && req.Type != "exit" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type must be 'entry' or 'exit'"})
		return
	}

	ts, err := time.Parse(time.RFC3339, req.Timestamp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid timestamp format, use RFC3339"})
		return
	}

	photoSiteData := ""
	photoSelfieData := ""

	siteFile, siteErr := c.FormFile("photo_site")
	if siteErr == nil {
		photoSiteData, err = encodeToDataURL(siteFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing site photo"})
			return
		}
	}

	selfieFile, selfieErr := c.FormFile("photo_selfie")
	if selfieErr == nil {
		photoSelfieData, err = encodeToDataURL(selfieFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing selfie photo"})
			return
		}
	}

	var record models.CheckRecord
	err = database.DB.QueryRow(
		`INSERT INTO check_records (user_id, type, timestamp, photo_site_path, photo_selfie_path)
 VALUES ($1, $2, $3, $4, $5)
 RETURNING id, user_id, type, timestamp,
           COALESCE(photo_site_path, ''), COALESCE(photo_selfie_path, ''),
           is_suspicious,
           COALESCE(suspicious_reason, ''), COALESCE(ip_country, ''), COALESCE(ip_city, ''),
           created_at`,
		userID, req.Type, ts, photoSiteData, photoSelfieData,
	).Scan(&record.ID, &record.UserID, &record.Type, &record.Timestamp,
		&record.PhotoSitePath, &record.PhotoSelfiePath,
		&record.IsSuspicious, &record.SuspiciousReason, &record.IPCountry, &record.IPCity,
		&record.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving check record"})
		return
	}

	// Launch fraud check asynchronously — doesn't delay the user response
	clientIP := getClientIP(c)
	go runFraudCheck(record.ID, clientIP, req.Latitude, req.Longitude)

	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Check %s registered successfully", req.Type),
		"record":  record,
	})
}

func AddLocationPoint(c *gin.Context) {
	userID := c.GetString("user_id")

	var req LocationPointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ownerID string
	err := database.DB.QueryRow(
		`SELECT user_id FROM check_records WHERE id = $1`, req.CheckRecordID,
	).Scan(&ownerID)
	if err != nil || ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Record not found or unauthorized"})
		return
	}

	ts, err := time.Parse(time.RFC3339, req.RecordedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid recorded_at format"})
		return
	}

	var point models.LocationPoint
	err = database.DB.QueryRow(
		`INSERT INTO location_points (check_record_id, latitude, longitude, accuracy, recorded_at)
 VALUES ($1, $2, $3, $4, $5)
 RETURNING id, check_record_id, latitude, longitude, accuracy, recorded_at`,
		req.CheckRecordID, req.Latitude, req.Longitude, req.Accuracy, ts,
	).Scan(&point.ID, &point.CheckRecordID, &point.Latitude, &point.Longitude, &point.Accuracy, &point.RecordedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving location point"})
		return
	}

	c.JSON(http.StatusCreated, point)
}

type BatchLocationPointRequest struct {
	CheckRecordID string `json:"check_record_id" binding:"required"`
	Points        []struct {
		Latitude   float64 `json:"latitude"   binding:"required"`
		Longitude  float64 `json:"longitude"  binding:"required"`
		Accuracy   float64 `json:"accuracy"`
		RecordedAt string  `json:"recorded_at" binding:"required"`
	} `json:"points" binding:"required"`
}

func AddLocationPointsBatch(c *gin.Context) {
	userID := c.GetString("user_id")

	var req BatchLocationPointRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify ownership
	var ownerID string
	err := database.DB.QueryRow(
		`SELECT user_id FROM check_records WHERE id = $1`, req.CheckRecordID,
	).Scan(&ownerID)
	if err != nil || ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Record not found or unauthorized"})
		return
	}

	if len(req.Points) == 0 {
		c.JSON(http.StatusCreated, gin.H{"inserted": 0})
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction error"})
		return
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		`INSERT INTO location_points (check_record_id, latitude, longitude, accuracy, recorded_at)
		 VALUES ($1, $2, $3, $4, $5)`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Prepare error"})
		return
	}
	defer stmt.Close()

	inserted := 0
	for _, p := range req.Points {
		ts, err := time.Parse(time.RFC3339, p.RecordedAt)
		if err != nil {
			continue // skip malformed timestamps
		}
		if _, err = stmt.Exec(req.CheckRecordID, p.Latitude, p.Longitude, p.Accuracy, ts); err == nil {
			inserted++
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Commit error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"inserted": inserted})
}

func GetMyRecords(c *gin.Context) {
	userID := c.GetString("user_id")

	rows, err := database.DB.Query(
		`SELECT id, user_id, type, timestamp, photo_site_path, photo_selfie_path, created_at
 FROM check_records WHERE user_id = $1 ORDER BY timestamp DESC`,
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	records := []models.CheckRecord{}
	for rows.Next() {
		var r models.CheckRecord
		rows.Scan(&r.ID, &r.UserID, &r.Type, &r.Timestamp,
			&r.PhotoSitePath, &r.PhotoSelfiePath, &r.CreatedAt)
		records = append(records, r)
	}

	c.JSON(http.StatusOK, records)
}

func GetRecordRoute(c *gin.Context) {
	userID := c.GetString("user_id")
	recordID := c.Param("id")

	var ownerID string
	err := database.DB.QueryRow(
		`SELECT user_id FROM check_records WHERE id = $1`, recordID,
	).Scan(&ownerID)
	if err != nil || ownerID != userID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	rows, err := database.DB.Query(
		`SELECT id, check_record_id, latitude, longitude, accuracy, recorded_at
 FROM location_points WHERE check_record_id = $1 ORDER BY recorded_at ASC`,
		recordID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	points := []models.LocationPoint{}
	for rows.Next() {
		var p models.LocationPoint
		rows.Scan(&p.ID, &p.CheckRecordID, &p.Latitude, &p.Longitude, &p.Accuracy, &p.RecordedAt)
		points = append(points, p)
	}

	c.JSON(http.StatusOK, points)
}

// encodeToDataURL reads a multipart file and returns a base64 data URL.
// Photos are stored directly in the database so they survive container restarts.
func encodeToDataURL(fh *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	mimeTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".webp": "image/webp",
	}
	mime, ok := mimeTypes[ext]
	if !ok {
		return "", fmt.Errorf("file type not allowed")
	}

	src, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return "", err
	}

	return "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(data), nil
}

// GetCheckStatus returns if the user currently has an active 'entry' without an 'exit'
func GetCheckStatus(c *gin.Context) {
	userID := c.GetString("user_id")

	var r models.CheckRecord
	// Fetch the most recent check record for this user
	err := database.DB.QueryRow(
		`SELECT id, type, timestamp FROM check_records 
		 WHERE user_id = $1 
		 ORDER BY timestamp DESC LIMIT 1`,
		userID,
	).Scan(&r.ID, &r.Type, &r.Timestamp)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{"active": false})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if r.Type == "entry" {
		// Only consider it "active" if it was done within the same day/session.
		// (Optional) We are doing simple type tracking, so if it's entry, they are active.
		c.JSON(http.StatusOK, gin.H{
			"active":     true,
			"record_id":  r.ID,
			"entry_time": r.Timestamp,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{"active": false})
	}
}
