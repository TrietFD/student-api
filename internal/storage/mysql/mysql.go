package mysql

import (
	"database/sql"
	"fmt"
	"log"
	
	"github.com/TrietFD/student-api/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*MySQL, error) {
	// Construct MySQL connection string without specifying database
	// This allows connection to MySQL server to create database
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/", 
		cfg.DBUser, 
		cfg.DBPassword, 
		cfg.DBHost, 
		cfg.DBPort,
	)

	// Open connection to MySQL
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database server: %v", err)
	}

	// Check if database exists, if not create it
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", cfg.DBName))
	if err != nil {
		return nil, fmt.Errorf("error creating database: %v", err)
	}

	// Close the initial connection
	db.Close()

	// Now connect to the specific database
	dataSourceName = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", 
		cfg.DBUser, 
		cfg.DBPassword, 
		cfg.DBHost, 
		cfg.DBPort, 
		cfg.DBName,
	)

	// Reopen connection to the specific database
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening specific database: %v", err)
	}

	// Create students table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255),
		email VARCHAR(255) UNIQUE,
		age INT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)

	if err != nil {
		return nil, fmt.Errorf("error creating students table: %v", err)
	}

	log.Println("Database and table setup completed successfully")

	return &MySQL{
		Db: db,
	}, nil
}

// Close method to properly close the database connection
func (m *MySQL) Close() error {
	return m.Db.Close()
}