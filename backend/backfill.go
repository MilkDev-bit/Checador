package main

import (
	"fmt"
	"log"

	"paselista/database"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}
	database.Connect()

	// 1. Verónica Hernández Sánchez
	// Moved to CALDERAS on June 3.
	res1, err := database.DB.Exec(`
		UPDATE check_records 
		SET project_name = 'CEDIS BAJIO SECOS' 
		WHERE user_id = (SELECT id FROM users WHERE first_name ILIKE '%Verónica%' AND last_name ILIKE '%Hernández Sánchez%' LIMIT 1) 
		AND timestamp < '2026-06-03 00:00:00'
	`)
	if err != nil {
		log.Fatal(err)
	}
	rows1, _ := res1.RowsAffected()

	res2, err := database.DB.Exec(`
		UPDATE check_records 
		SET project_name = 'CALDERAS' 
		WHERE user_id = (SELECT id FROM users WHERE first_name ILIKE '%Verónica%' AND last_name ILIKE '%Hernández Sánchez%' LIMIT 1) 
		AND timestamp >= '2026-06-03 00:00:00' AND project_name = ''
	`)
	if err != nil {
		log.Fatal(err)
	}
	rows2, _ := res2.RowsAffected()
	fmt.Printf("Verónica: %d records updated to CEDIS BAJIO SECOS, %d to CALDERAS\n", rows1, rows2)

	// 2. Miguel Bolainas
	// Moved to BAE JUAN PABLO SEGUNDO on June 23 (today).
	res3, err := database.DB.Exec(`
		UPDATE check_records 
		SET project_name = 'BALCONES DE HUENTITLAN' 
		WHERE user_id = (SELECT id FROM users WHERE first_name ILIKE '%Miguel%' AND last_name ILIKE '%Bolainas%' LIMIT 1) 
		AND timestamp < '2026-06-23 00:00:00'
	`)
	if err != nil {
		log.Fatal(err)
	}
	rows3, _ := res3.RowsAffected()

	res4, err := database.DB.Exec(`
		UPDATE check_records 
		SET project_name = 'BAE JUAN PABLO SEGUNDO' 
		WHERE user_id = (SELECT id FROM users WHERE first_name ILIKE '%Miguel%' AND last_name ILIKE '%Bolainas%' LIMIT 1) 
		AND timestamp >= '2026-06-23 00:00:00' AND project_name = ''
	`)
	if err != nil {
		log.Fatal(err)
	}
	rows4, _ := res4.RowsAffected()
	fmt.Printf("Miguel: %d records updated to BALCONES DE HUENTITLAN, %d to BAE JUAN PABLO SEGUNDO\n", rows3, rows4)
}
