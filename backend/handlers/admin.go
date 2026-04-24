package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"paselista/database"

	"github.com/gin-gonic/gin"
)

type AdminRecordRow struct {
	RecordID         string    `json:"record_id"`
	UserID           string    `json:"user_id"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	ProjectName      string    `json:"project_name"`
	Email            string    `json:"email"`
	Type             string    `json:"type"`
	Timestamp        time.Time `json:"timestamp"`
	HasSitePhoto     bool      `json:"has_site_photo"`
	HasSelfiePhoto   bool      `json:"has_selfie_photo"`
	LocationCount    int       `json:"location_count"`
	IsSuspicious     bool      `json:"is_suspicious"`
	SuspiciousReason string    `json:"suspicious_reason,omitempty"`
	IPCountry        string    `json:"ip_country,omitempty"`
	IPCity           string    `json:"ip_city,omitempty"`
}

type AdminStats struct {
	TotalUsers   int `json:"total_users"`
	RecordsToday int `json:"records_today"`
	EntriesTotal int `json:"entries_total"`
	ExitsTotal   int `json:"exits_total"`
	ActiveNow    int `json:"active_now"`
}

type AdminUserRow struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	ProjectName string    `json:"project_name"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	TotalChecks int       `json:"total_checks"`
}

func AdminGetStats(c *gin.Context) {
	dateStr := c.Query("date")
	var date time.Time
	if dateStr != "" {
		var err error
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
			return
		}
	} else {
		date = time.Now()
	}

	loc := time.Local
	dateStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
	dateEnd := dateStart.Add(24 * time.Hour)

	var stats AdminStats
	database.DB.QueryRow(`SELECT COUNT(*) FROM users WHERE role = 'user'`).Scan(&stats.TotalUsers)
	database.DB.QueryRow(`SELECT COUNT(*) FROM check_records WHERE timestamp >= $1 AND timestamp < $2`, dateStart, dateEnd).Scan(&stats.RecordsToday)
	database.DB.QueryRow(`SELECT COUNT(*) FROM check_records WHERE type='entry' AND timestamp >= $1 AND timestamp < $2`, dateStart, dateEnd).Scan(&stats.EntriesTotal)
	database.DB.QueryRow(`SELECT COUNT(*) FROM check_records WHERE type='exit' AND timestamp >= $1 AND timestamp < $2`, dateStart, dateEnd).Scan(&stats.ExitsTotal)
	database.DB.QueryRow(`
		SELECT COUNT(DISTINCT user_id) FROM check_records
		WHERE type='entry' AND timestamp >= $1 AND timestamp < $2
		AND user_id NOT IN (
			SELECT user_id FROM check_records WHERE type='exit' AND timestamp >= $1 AND timestamp < $2
		)`, dateStart, dateEnd).Scan(&stats.ActiveNow)

	c.JSON(http.StatusOK, stats)
}

func AdminGetRecords(c *gin.Context) {
	dateStr := c.Query("date")
	projectFilter := strings.TrimSpace(c.Query("project"))
	typeFilter := c.Query("type")
	search := strings.TrimSpace(c.Query("search"))

	args := []interface{}{}
	where := []string{"u.role = 'user'"}
	argIdx := 1

	if dateStr != "" {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date"})
			return
		}
		dateStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
		dateEnd := dateStart.Add(24 * time.Hour)
		where = append(where, fmt.Sprintf("cr.timestamp >= $%d AND cr.timestamp < $%d", argIdx, argIdx+1))
		args = append(args, dateStart, dateEnd)
		argIdx += 2
	}

	if projectFilter != "" {
		where = append(where, fmt.Sprintf("LOWER(u.project_name) LIKE LOWER($%d)", argIdx))
		args = append(args, "%"+projectFilter+"%")
		argIdx++
	}

	if typeFilter == "entry" || typeFilter == "exit" {
		where = append(where, fmt.Sprintf("cr.type = $%d", argIdx))
		args = append(args, typeFilter)
		argIdx++
	}

	if search != "" {
		where = append(where, fmt.Sprintf(
			"(LOWER(u.first_name) LIKE LOWER($%d) OR LOWER(u.last_name) LIKE LOWER($%d) OR LOWER(u.email) LIKE LOWER($%d))",
			argIdx, argIdx, argIdx))
		args = append(args, "%"+search+"%")
		argIdx++
	}

	query := `SELECT cr.id, u.id, u.first_name, u.last_name, u.project_name, u.email,
	           cr.type, cr.timestamp,
	           (cr.photo_site_path IS NOT NULL AND cr.photo_site_path != '') AS has_site_photo,
	           (cr.photo_selfie_path IS NOT NULL AND cr.photo_selfie_path != '') AS has_selfie_photo,
	           (SELECT COUNT(*) FROM location_points lp WHERE lp.check_record_id = cr.id),
	           cr.is_suspicious, COALESCE(cr.suspicious_reason,''), COALESCE(cr.ip_country,''), COALESCE(cr.ip_city,'')
	           FROM check_records cr JOIN users u ON cr.user_id = u.id`
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	query += " ORDER BY cr.timestamp DESC LIMIT 500"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	records := []AdminRecordRow{}
	for rows.Next() {
		var r AdminRecordRow
		rows.Scan(&r.RecordID, &r.UserID, &r.FirstName, &r.LastName, &r.ProjectName, &r.Email,
			&r.Type, &r.Timestamp, &r.HasSitePhoto, &r.HasSelfiePhoto, &r.LocationCount,
			&r.IsSuspicious, &r.SuspiciousReason, &r.IPCountry, &r.IPCity)
		records = append(records, r)
	}

	if records == nil {
		records = []AdminRecordRow{}
	}
	c.JSON(http.StatusOK, records)
}

func AdminGetUsers(c *gin.Context) {
	projectFilter := strings.TrimSpace(c.Query("project"))

	args := []interface{}{}
	query := `SELECT u.id, u.first_name, u.last_name, u.project_name, u.email, u.created_at,
	           (SELECT COUNT(*) FROM check_records cr WHERE cr.user_id = u.id) as total_checks
	           FROM users u WHERE u.role = 'user'`
	if projectFilter != "" {
		query += " AND LOWER(u.project_name) LIKE LOWER($1)"
		args = append(args, "%"+projectFilter+"%")
	}
	query += " ORDER BY u.created_at DESC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	users := []AdminUserRow{}
	for rows.Next() {
		var u AdminUserRow
		rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.ProjectName, &u.Email, &u.CreatedAt, &u.TotalChecks)
		users = append(users, u)
	}

	if users == nil {
		users = []AdminUserRow{}
	}
	c.JSON(http.StatusOK, users)
}

func AdminGetProjects(c *gin.Context) {
	rows, err := database.DB.Query(`SELECT DISTINCT project_name FROM users WHERE role='user' ORDER BY project_name`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	projects := []string{}
	for rows.Next() {
		var p string
		rows.Scan(&p)
		projects = append(projects, p)
	}
	c.JSON(http.StatusOK, projects)
}

func AdminGetRecordRoute(c *gin.Context) {
	recordID := c.Param("id")
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

	type Point struct {
		ID            string    `json:"id"`
		CheckRecordID string    `json:"check_record_id"`
		Latitude      float64   `json:"latitude"`
		Longitude     float64   `json:"longitude"`
		Accuracy      float64   `json:"accuracy"`
		RecordedAt    time.Time `json:"recorded_at"`
	}

	points := []Point{}
	for rows.Next() {
		var p Point
		rows.Scan(&p.ID, &p.CheckRecordID, &p.Latitude, &p.Longitude, &p.Accuracy, &p.RecordedAt)
		points = append(points, p)
	}
	c.JSON(http.StatusOK, points)
}

func AdminGetRecordPhotos(c *gin.Context) {
	recordID := c.Param("id")
	var sitePath, selfiePath string
	err := database.DB.QueryRow(
		`SELECT COALESCE(photo_site_path,''), COALESCE(photo_selfie_path,'')
		 FROM check_records WHERE id = $1`,
		recordID,
	).Scan(&sitePath, &selfiePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"photo_site_path":   sitePath,
		"photo_selfie_path": selfiePath,
	})
}
