package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB opens the database connection and initializes the schema
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	err = createTables(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	userTableSQL := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT UNIQUE,
		password TEXT
	);`

	habitTableSQL := `CREATE TABLE IF NOT EXISTS habits (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		title TEXT,
		is_complete BOOLEAN,
		level_of_complete INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	_, err := db.Exec(userTableSQL)
	if err != nil {
		return err
	}

	_, err = db.Exec(habitTableSQL)
	if err != nil {
		return err
	}

	return nil
}

// SeedData populates the database with initial data
func SeedData(db *sql.DB) error {
	// Check if data already exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("Data already seeded, skipping...")
		return nil
	}

	// Create a user
	result, err := db.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)",
		"Billy", "billy@example.com", "hashedpassword") // In a real application, make sure to hash the password
	if err != nil {
		return err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Create habits for the user
	habits := []struct {
		Title           string
		IsComplete      bool
		LevelOfComplete int
	}{
		{"Going to the gym", true, 100},
		{"Coding for more than 4 hours a day", false, 100},
		{"Take a walk", false, 50},
	}

	for _, habit := range habits {
		_, err := db.Exec("INSERT INTO habits (user_id, title, is_complete, level_of_complete) VALUES (?, ?, ?, ?)",
			userID, habit.Title, habit.IsComplete, habit.LevelOfComplete)
		if err != nil {
			return err
		}
	}

	log.Println("Seed data created successfully")
	return nil
}
