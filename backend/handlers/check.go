package handlers

import (
"fmt"
"io"
"mime/multipart"
"net/http"
"os"
"path/filepath"
"strings"
"time"

"paselista/database"
"paselista/models"

"github.com/gin-gonic/gin"
"github.com/google/uuid"
)

type CheckInRequest struct {
Type      string `form:"type" binding:"required"`
Timestamp string `form:"timestamp" binding:"required"`
}

type LocationPointRequest struct {
CheckRecordID string  `json:"check_record_id" binding:"required"`
Latitude      float64 `json:"latitude" binding:"required"`
Longitude     float64 `json:"longitude" binding:"required"`
Accuracy      float64 `json:"accuracy"`
RecordedAt    string  `json:"recorded_at" binding:"required"`
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

uploadsDir := "./uploads"
os.MkdirAll(uploadsDir, 0750)

photoSitePath := ""
photoSelfiePath := ""

siteFile, siteErr := c.FormFile("photo_site")
if siteErr == nil {
photoSitePath, err = saveUploadedFile(siteFile, uploadsDir, "site")
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving site photo"})
return
}
}

selfieFile, selfieErr := c.FormFile("photo_selfie")
if selfieErr == nil {
photoSelfiePath, err = saveUploadedFile(selfieFile, uploadsDir, "selfie")
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving selfie photo"})
return
}
}

var record models.CheckRecord
err = database.DB.QueryRow(
`INSERT INTO check_records (user_id, type, timestamp, photo_site_path, photo_selfie_path)
 VALUES ($1, $2, $3, $4, $5)
 RETURNING id, user_id, type, timestamp, photo_site_path, photo_selfie_path, created_at`,
userID, req.Type, ts, photoSitePath, photoSelfiePath,
).Scan(&record.ID, &record.UserID, &record.Type, &record.Timestamp,
&record.PhotoSitePath, &record.PhotoSelfiePath, &record.CreatedAt)

if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving check record"})
return
}

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

func saveUploadedFile(fh *multipart.FileHeader, dir, prefix string) (string, error) {
ext := filepath.Ext(fh.Filename)
allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
if !allowed[strings.ToLower(ext)] {
return "", fmt.Errorf("file type not allowed")
}

filename := fmt.Sprintf("%s_%s_%d%s", prefix, uuid.New().String(), time.Now().UnixMilli(), ext)
dst := filepath.Join(dir, filename)

src, err := fh.Open()
if err != nil {
return "", err
}
defer src.Close()

out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
if err != nil {
return "", err
}
defer out.Close()

if _, err = io.Copy(out, src); err != nil {
return "", err
}

return "/uploads/" + filename, nil
}
