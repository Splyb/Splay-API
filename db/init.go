package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	url := os.Getenv("SUPABASE_DB_URL")
	DB, err = sql.Open("postgres", url)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS tracks (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		artist VARCHAR(255) NOT NULL,
		file_url TEXT NOT NULL,
		uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatalf("Error al crear la tabla: %v", err)
	}
}
