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

	// Create users table if not exists
	usersTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = DB.Exec(usersTableQuery)

	if err != nil {
		fmt.Println("Failed to create table for users!", err)
	}

	// // Create tasks table if not exists
	// tasksTableQuery := `
	// CREATE TABLE IF NOT EXISTS tasks (
	// 	id SERIAL PRIMARY KEY,
	// 	name TEXT NOT NULL,
	// 	url TEXT,
	// 	user_id INTEGER REFERENCES users(id),
	// 	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	// 	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	// );`

	// _, err = DB.Exec(tasksTableQuery)

	// if err != nil {
	// 	fmt.Println("Failed to create table for users!", err)
	// }

	println("Database initialized and tables created if they did not exist.")

	return nil
}
