package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"paselista/database"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
	// Schedule & deviation
	EntrySchedule string `json:"entry_schedule,omitempty"` // "HH:MM"
	ExitSchedule  string `json:"exit_schedule,omitempty"`  // "HH:MM"
	DeviationMin  int    `json:"deviation_min"`            // positive = late/overtime, negative = early
	// Daily comment (same for all records of that user on that day)
	DayComment string `json:"day_comment,omitempty"`
	// GPS reverse-geocoded address
	GpsAddress string `json:"gps_address,omitempty"`
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
	AvatarURL   string    `json:"avatar_url"`
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
		AND COALESCE(closed_by_admin, false) = false
		AND user_id NOT IN (
			SELECT user_id FROM check_records WHERE type='exit' AND timestamp >= $1 AND timestamp < $2
		)`, dateStart, dateEnd).Scan(&stats.ActiveNow)

	c.JSON(http.StatusOK, stats)
}

func AdminGetRecords(c *gin.Context) {
	dateStr := c.Query("date")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	projectFilter := strings.TrimSpace(c.Query("project"))
	typeFilter := c.Query("type")
	search := strings.TrimSpace(c.Query("search"))

	args := []interface{}{}
	where := []string{"u.role = 'user'"}
	argIdx := 1

	// Single date (legacy) takes priority over range
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
	} else {
		if dateFrom != "" {
			df, err := time.Parse("2006-01-02", dateFrom)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date_from"})
				return
			}
			where = append(where, fmt.Sprintf("cr.timestamp >= $%d", argIdx))
			args = append(args, time.Date(df.Year(), df.Month(), df.Day(), 0, 0, 0, 0, time.Local))
			argIdx++
		}
		if dateTo != "" {
			dt, err := time.Parse("2006-01-02", dateTo)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date_to"})
				return
			}
			where = append(where, fmt.Sprintf("cr.timestamp < $%d", argIdx))
			args = append(args, time.Date(dt.Year(), dt.Month(), dt.Day(), 0, 0, 0, 0, time.Local).Add(24*time.Hour))
			argIdx++
		}
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

	_ = argIdx // suppress unused warning

	query := `SELECT cr.id, u.id, u.first_name, u.last_name, u.project_name, u.email,
	           cr.type, cr.timestamp,
	           (cr.photo_site_path IS NOT NULL AND cr.photo_site_path != '') AS has_site_photo,
	           (cr.photo_selfie_path IS NOT NULL AND cr.photo_selfie_path != '') AS has_selfie_photo,
	           (SELECT COUNT(*) FROM location_points lp WHERE lp.check_record_id = cr.id),
	           cr.is_suspicious, COALESCE(cr.suspicious_reason,''), COALESCE(cr.ip_country,''), COALESCE(cr.ip_city,''),
	           CASE WHEN EXTRACT(DOW FROM cr.timestamp) = 6
	                THEN COALESCE(ws.sat_entry_time,'')
	                ELSE COALESCE(ws.entry_time,'') END,
	           CASE WHEN EXTRACT(DOW FROM cr.timestamp) = 6
	                THEN COALESCE(ws.sat_exit_time,'')
	                ELSE COALESCE(ws.exit_time,'') END,
	           COALESCE((SELECT rc.comment FROM record_comments rc
	                     WHERE rc.user_id = u.id AND rc.record_date = cr.timestamp::date),''),
	           COALESCE(cr.gps_address,'')
	           FROM check_records cr
	           JOIN users u ON cr.user_id = u.id
	           LEFT JOIN work_schedules ws ON ws.user_id = u.id`
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	query += " ORDER BY cr.timestamp DESC LIMIT 1000"

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
			&r.IsSuspicious, &r.SuspiciousReason, &r.IPCountry, &r.IPCity,
			&r.EntrySchedule, &r.ExitSchedule, &r.DayComment, &r.GpsAddress)

		// Calculate deviation in minutes
		r.DeviationMin = calcDeviation(r.Timestamp, r.Type, r.EntrySchedule, r.ExitSchedule)

		records = append(records, r)
	}

	if records == nil {
		records = []AdminRecordRow{}
	}
	c.JSON(http.StatusOK, records)
}

// calcDeviation returns the deviation in minutes for a check record vs the work schedule.
// For entries: positive = arrived late, negative = arrived early.
// For exits:   positive = left late (overtime), negative = left early.
func calcDeviation(ts time.Time, checkType, entrySchedule, exitSchedule string) int {
	var schedule string
	switch checkType {
	case "entry":
		schedule = entrySchedule
	case "exit":
		schedule = exitSchedule
	default:
		return 0
	}
	if len(schedule) != 5 {
		return 0
	}

	var sh, sm int
	fmt.Sscanf(schedule, "%d:%d", &sh, &sm)
	scheduled := time.Date(ts.Year(), ts.Month(), ts.Day(), sh, sm, 0, 0, ts.Location())
	diff := int(ts.Sub(scheduled).Minutes())
	return diff
}

func AdminGetUsers(c *gin.Context) {
	projectFilter := strings.TrimSpace(c.Query("project"))

	args := []interface{}{}
	query := `SELECT u.id, u.first_name, u.last_name, u.project_name, u.email, u.created_at, COALESCE(u.avatar_url, ''),
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
		rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.ProjectName, &u.Email, &u.CreatedAt, &u.AvatarURL, &u.TotalChecks)
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

// AdminResetUserPassword resets a user's password without requiring the current one.
// PUT /admin/users/:id/password
func AdminResetUserPassword(c *gin.Context) {
	userID := c.Param("id")

	var req struct {
		Password string `json:"password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "La contraseña debe tener al menos 8 caracteres."})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar la contraseña."})
		return
	}

	result, err := database.DB.Exec(
		`UPDATE users SET password_hash = $1 WHERE id = $2 AND role = 'user'`,
		string(hash), userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar la contraseña."})
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contraseña actualizada correctamente."})
}

// AdminGetOpenSessions returns all users with an open entry (no matching exit,
// not closed by admin) from the last 24 hours.
// GET /admin/open-sessions
func AdminGetOpenSessions(c *gin.Context) {
	const sessionTTL = 20 * time.Hour
	windowStart := time.Now().UTC().Add(-24 * time.Hour)

	type OpenSession struct {
		RecordID    string    `json:"record_id"`
		UserID      string    `json:"user_id"`
		FirstName   string    `json:"first_name"`
		LastName    string    `json:"last_name"`
		Email       string    `json:"email"`
		ProjectName string    `json:"project_name"`
		EntryTime   time.Time `json:"entry_time"`
		DurationMin int       `json:"duration_min"`
		AutoExpired bool      `json:"auto_expired"`
	}

	rows, err := database.DB.Query(
		`SELECT cr.id, u.id, u.first_name, u.last_name, u.email, u.project_name, cr.timestamp
		 FROM check_records cr
		 JOIN users u ON cr.user_id = u.id
		 WHERE cr.type = 'entry'
		   AND COALESCE(cr.closed_by_admin, false) = false
		   AND cr.timestamp >= $1
		   AND u.role = 'user'
		   AND NOT EXISTS (
		       SELECT 1 FROM check_records cr2
		       WHERE cr2.user_id = cr.user_id
		         AND cr2.type = 'exit'
		         AND cr2.timestamp > cr.timestamp
		   )
		 ORDER BY cr.timestamp ASC`,
		windowStart,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	now := time.Now()
	sessions := []OpenSession{}
	for rows.Next() {
		var s OpenSession
		rows.Scan(&s.RecordID, &s.UserID, &s.FirstName, &s.LastName, &s.Email, &s.ProjectName, &s.EntryTime)
		dur := now.Sub(s.EntryTime)
		s.DurationMin = int(dur.Minutes())
		s.AutoExpired = dur >= sessionTTL
		sessions = append(sessions, s)
	}

	c.JSON(http.StatusOK, sessions)
}

// AdminCloseSession manually marks an open entry as closed so the user can
// register a new entry without waiting for the auto-expiry (sessionTTL).
// POST /admin/open-sessions/:id/close
func AdminCloseSession(c *gin.Context) {
	recordID := c.Param("id")

	result, err := database.DB.Exec(
		`UPDATE check_records SET closed_by_admin = true WHERE id = $1 AND type = 'entry'`,
		recordID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Registro no encontrado"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sesión cerrada correctamente"})
}
