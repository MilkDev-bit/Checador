package models

import "time"

type User struct {
	ID           string    `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	ProjectName  string    `json:"project_name"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	AvatarURL    string    `json:"avatar_url,omitempty"`
	CoverURL     string    `json:"cover_url,omitempty"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type CheckRecord struct {
	ID               string          `json:"id"`
	UserID           string          `json:"user_id"`
	Type             string          `json:"type"`
	Timestamp        time.Time       `json:"timestamp"`
	PhotoSitePath    string          `json:"photo_site_path,omitempty"`
	PhotoSelfiePath  string          `json:"photo_selfie_path,omitempty"`
	IsSuspicious     bool            `json:"is_suspicious"`
	SuspiciousReason string          `json:"suspicious_reason,omitempty"`
	IPCountry        string          `json:"ip_country,omitempty"`
	IPCity           string          `json:"ip_city,omitempty"`
	LocationPoints   []LocationPoint `json:"location_points,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
}

type LocationPoint struct {
	ID            string    `json:"id"`
	CheckRecordID string    `json:"check_record_id"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	Accuracy      float64   `json:"accuracy"`
	RecordedAt    time.Time `json:"recorded_at"`
}
