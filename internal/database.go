package internal

import (
	"database/sql"
	"log"
)

func InitMySqlDB(db *sql.DB) {
	log.Println("Creating database if it doesn't exist")

	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS final_project")
	if err != nil {
		log.Fatalf("Error creating database: %v", err)
	}

	// Select the database for use
	_, err = db.Exec("USE final_project")
	if err != nil {
		log.Fatalf("Error selecting database: %v", err)
	}

	// Create the table if it doesn't exist
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS url_shortener (
        id INT AUTO_INCREMENT PRIMARY KEY,
        original_url VARCHAR(2048) NOT NULL,
        short_url VARCHAR(5) NOT NULL
    );`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	log.Println("Successfully initialized database")
}

var sqlOpen = sql.Open

func ConnectToMySqlDB() (*sql.DB, error) {
	db, err := sqlOpen("mysql", ConnectionString)
	if err != nil {
		log.Printf("Failed to connect to MySQL: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to MySQL database")
	return db, nil
}
