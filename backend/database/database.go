package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func Connect() {
	// Railway provides DATABASE_URL as a full connection string (preferred)
	// Fall back to individual env vars for local development
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			mustEnv("DB_HOST"),
			getEnv("DB_PORT", "5432"),
			mustEnv("DB_USER"),
			mustEnv("DB_PASSWORD"),
			mustEnv("DB_NAME"),
		)
	}

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	log.Println("Database connected")
}

func Migrate() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			project_name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			role VARCHAR(10) NOT NULL DEFAULT 'user',
			created_at TIMESTAMP DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS check_records (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id),
			type VARCHAR(10) NOT NULL CHECK (type IN ('entry','exit')),
			timestamp TIMESTAMP NOT NULL,
			photo_site_path TEXT,
			photo_selfie_path TEXT,
			is_suspicious BOOLEAN NOT NULL DEFAULT FALSE,
			suspicious_reason TEXT,
			ip_country VARCHAR(100),
			ip_city VARCHAR(100),
			created_at TIMESTAMP DEFAULT NOW()
		)`,
		// Widen existing columns from VARCHAR(500) to TEXT in case the table already exists
		`ALTER TABLE check_records ALTER COLUMN photo_site_path TYPE TEXT`,
		`ALTER TABLE check_records ALTER COLUMN photo_selfie_path TYPE TEXT`,
		`CREATE TABLE IF NOT EXISTS location_points (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			check_record_id UUID NOT NULL REFERENCES check_records(id),
			latitude DOUBLE PRECISION NOT NULL,
			longitude DOUBLE PRECISION NOT NULL,
			accuracy DOUBLE PRECISION,
			recorded_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		)`,
	}

	// Fatal migrations — table creation must succeed
	createQueries := queries[:len(queries)-2]
	for _, q := range createQueries {
		if _, err := DB.Exec(q); err != nil {
			log.Fatalf("Migration error: %v", err)
		}
	}

	// Non-fatal migrations — ALTER TABLE to fix existing deployed columns
	alterQueries := []string{
		`ALTER TABLE check_records ALTER COLUMN photo_site_path TYPE TEXT`,
		`ALTER TABLE check_records ALTER COLUMN photo_selfie_path TYPE TEXT`,
		`ALTER TABLE check_records ADD COLUMN IF NOT EXISTS is_suspicious BOOLEAN NOT NULL DEFAULT FALSE`,
		`ALTER TABLE check_records ADD COLUMN IF NOT EXISTS suspicious_reason TEXT`,
		`ALTER TABLE check_records ADD COLUMN IF NOT EXISTS ip_country VARCHAR(100)`,
		`ALTER TABLE check_records ADD COLUMN IF NOT EXISTS ip_city VARCHAR(100)`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_url TEXT`,
	}
	for _, q := range alterQueries {
		if _, err := DB.Exec(q); err != nil {
			log.Printf("Migration warning (non-fatal): %v", err)
		}
	}

	log.Println("Database migrated")
}

func SeedAdmin() {
	email := mustEnv("ADMIN_EMAIL")
	password := mustEnv("ADMIN_PASSWORD")

	var count int
	DB.QueryRow(`SELECT COUNT(*) FROM users WHERE role = 'admin'`).Scan(&count)
	if count > 0 {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error creating admin: %v", err)
		return
	}

	_, err = DB.Exec(
		`INSERT INTO users (first_name, last_name, project_name, email, password_hash, role)
		 VALUES ('Admin', 'Sistema', 'PaseLista Admin', $1, $2, 'admin')`,
		email, string(hash),
	)
	if err != nil {
		log.Printf("Error seeding admin: %v", err)
		return
	}
	log.Printf("Admin user created: %s", email)
}

// getEnv returns the env var value or a non-sensitive default (e.g. port numbers).
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// mustEnv returns the env var value or terminates the process.
// Use this for any variable that contains credentials or required config.
func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("Required environment variable %q is not set", key)
	}
	return v
}
