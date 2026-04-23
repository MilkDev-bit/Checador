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
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "paselista"),
	)

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
			photo_site_path VARCHAR(500),
			photo_selfie_path VARCHAR(500),
			created_at TIMESTAMP DEFAULT NOW()
		)`,
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

	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			log.Fatalf("Migration error: %v", err)
		}
	}
	log.Println("Database migrated")
}

func SeedAdmin() {
	email := getEnv("ADMIN_EMAIL", "admin@paselista.com")
	password := getEnv("ADMIN_PASSWORD", "Admin1234!")

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

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
