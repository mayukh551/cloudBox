package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() error {
	// connection string for pg
	connStr := "host=localhost port=5432 user=pqgotest password=yourpassword dbname=pqgotest sslmode=disable"
	conn, err := sql.Open("postgres", connStr)

	if err != nil {
		return fmt.Errorf("failed to connect to postgres!")
	}

	// Actually test the connection
	if err = conn.Ping(); err != nil {
		return fmt.Errorf("failed to ping postgres: %w", err)
	}

	fmt.Println("Connected to PostgreSQL")

	DB = conn

	// Create users table if does not exist
	usersTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(255) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = DB.Exec(usersTableQuery)

	if err != nil {
		fmt.Println("Failed to create table for users!", err)
	}

	// Create tasks table if does not exist
	fileTableQuery := `
		CREATE TABLE IF NOT EXISTS files (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			type VARCHAR(50) NOT NULL,
			size BIGINT DEFAULT 0,
			userID VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err = DB.Exec(fileTableQuery)

	if err != nil {
		fmt.Println("Failed to create table for users!", err)
	}

	// create share table if does not exist
	shareTableQuery := `
		CREATE TABLE IF NOT EXISTS shares (
			id VARCHAR(255) PRIMARY KEY,
			sharedTo VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			sharedBy VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			fileID VARCHAR(255) NOT NULL REFERENCES files(id) ON DELETE CASCADE,
			createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);	
	`

	_, err = DB.Exec(shareTableQuery)

	if err != nil {
		fmt.Println("Failed to create table for share!", err)
	}

	println("Database initialized and tables created if they did not exist.")

	return nil
}
