package handlers

import (
	"net/http"
	"time"

	"paselista/database"

	"github.com/gin-gonic/gin"
)

// ─── Work Schedule ────────────────────────────────────────────────────────────

type WorkSchedule struct {
	UserID    string    `json:"user_id"`
	EntryTime string    `json:"entry_time"` // "HH:MM" 24h
	ExitTime  string    `json:"exit_time"`  // "HH:MM" 24h
	UpdatedAt time.Time `json:"updated_at"`
}

type SetScheduleRequest struct {
	EntryTime string `json:"entry_time" binding:"required"`
	ExitTime  string `json:"exit_time"  binding:"required"`
}

// AdminGetUserSchedule returns the work schedule for a user.
// GET /admin/users/:id/schedule
func AdminGetUserSchedule(c *gin.Context) {
	userID := c.Param("id")

	var ws WorkSchedule
	err := database.DB.QueryRow(
		`SELECT user_id, entry_time, exit_time, updated_at
		 FROM work_schedules WHERE user_id = $1`, userID,
	).Scan(&ws.UserID, &ws.EntryTime, &ws.ExitTime, &ws.UpdatedAt)

	if err != nil {
		// No schedule set yet — return empty
		c.JSON(http.StatusOK, gin.H{"user_id": userID, "entry_time": "", "exit_time": ""})
		return
	}
	c.JSON(http.StatusOK, ws)
}

// AdminSetUserSchedule creates or updates the work schedule for a user.
// POST /admin/users/:id/schedule
func AdminSetUserSchedule(c *gin.Context) {
	userID := c.Param("id")

	var req SetScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "entry_time and exit_time are required (HH:MM)"})
		return
	}

	// Validate HH:MM format
	if len(req.EntryTime) != 5 || req.EntryTime[2] != ':' {
		c.JSON(http.StatusBadRequest, gin.H{"error": "entry_time must be HH:MM"})
		return
	}
	if len(req.ExitTime) != 5 || req.ExitTime[2] != ':' {
		c.JSON(http.StatusBadRequest, gin.H{"error": "exit_time must be HH:MM"})
		return
	}

	// Verify user exists
	var count int
	database.DB.QueryRow(`SELECT COUNT(*) FROM users WHERE id = $1 AND role = 'user'`, userID).Scan(&count)
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	_, err := database.DB.Exec(
		`INSERT INTO work_schedules (user_id, entry_time, exit_time, updated_at)
		 VALUES ($1, $2, $3, NOW())
		 ON CONFLICT (user_id) DO UPDATE
		   SET entry_time = EXCLUDED.entry_time,
		       exit_time  = EXCLUDED.exit_time,
		       updated_at = NOW()`,
		userID, req.EntryTime, req.ExitTime,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving schedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":    userID,
		"entry_time": req.EntryTime,
		"exit_time":  req.ExitTime,
	})
}

// ─── Daily Comments ───────────────────────────────────────────────────────────

type RecordComment struct {
	UserID     string    `json:"user_id"`
	RecordDate string    `json:"record_date"` // "YYYY-MM-DD"
	Comment    string    `json:"comment"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type SetCommentRequest struct {
	UserID     string `json:"user_id"     binding:"required"`
	RecordDate string `json:"record_date" binding:"required"` // "YYYY-MM-DD"
	Comment    string `json:"comment"     binding:"required"`
}

// AdminSetRecordComment creates or updates a daily comment for a user.
// POST /admin/records/comment
func AdminSetRecordComment(c *gin.Context) {
	var req SetCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id, record_date and comment are required"})
		return
	}

	_, err := database.DB.Exec(
		`INSERT INTO record_comments (user_id, record_date, comment, updated_at)
		 VALUES ($1, $2::date, $3, NOW())
		 ON CONFLICT (user_id, record_date) DO UPDATE
		   SET comment    = EXCLUDED.comment,
		       updated_at = NOW()`,
		req.UserID, req.RecordDate, req.Comment,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":     req.UserID,
		"record_date": req.RecordDate,
		"comment":     req.Comment,
	})
}

// AdminDeleteRecordComment deletes a daily comment.
// DELETE /admin/records/comment?user_id=...&date=...
func AdminDeleteRecordComment(c *gin.Context) {
	userID := c.Query("user_id")
	date := c.Query("date")
	if userID == "" || date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id and date are required"})
		return
	}

	database.DB.Exec(`DELETE FROM record_comments WHERE user_id = $1 AND record_date = $2::date`, userID, date)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
