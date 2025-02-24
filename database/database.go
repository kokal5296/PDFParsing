package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type PostgreSQLConnection struct {
	Pool *pgxpool.Pool
}

type DatabaseService interface {
	NewDatabase(connStr string, dbName string) (*PostgreSQLConnection, error)
	Close()
	GetPool() *pgxpool.Pool
}

func NewDatabaseService() DatabaseService {
	return &PostgreSQLConnection{}
}

// NewDatabase initializes a connection to PostgreSQL, checks if the target database exists,
// creates it if needed, and sets up the connection pool to the specific database.
func (db *PostgreSQLConnection) NewDatabase(connStr string, dbName string) (*PostgreSQLConnection, error) {

	conn, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Println("Unable to connect to database")
		return nil, err
	}
	defer conn.Close()

	var exists bool
	err = conn.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM pg_database WHERE LOWER(datname) = LOWER($1))", dbName).Scan(&exists)
	if err != nil {
		log.Println("Unable to check if database exists")
		return nil, err
	}

	log.Printf("Database exists: %v\n", exists)

	if !exists {
		query := fmt.Sprintf(`CREATE DATABASE %s`, dbName)
		_, err = conn.Exec(context.Background(), query)
		if err != nil {
			log.Println("Unable to create database")
			return nil, err
		}
		log.Println("Database created")
	}

	finalConnStr := fmt.Sprintf("%s%s?sslmode=disable", connStr, dbName)
	log.Println("Connecting to database")

	pool, err := pgxpool.Connect(context.Background(), finalConnStr)
	if err != nil {
		log.Println("Unable to connect to database")
		return nil, err
	}
	db.Pool = pool

	log.Println("Database connection established")

	err = db.CreateTablesIfNotExist()
	if err != nil {
		log.Println("Unable to create tables")
		return nil, err
	}

	return &PostgreSQLConnection{Pool: pool}, nil
}

// CreateTablesIfNotExist creates tables if they do not exist
func (db *PostgreSQLConnection) CreateTablesIfNotExist() error {

	queries := []string{
		`CREATE TABLE users (
    	 id SERIAL  PRIMARY KEY
		 );`,

		`CREATE TABLE files (
    	 id SERIAL PRIMARY KEY,
    	 filename VARCHAR(255) NOT NULL,
    	 file_hash VARCHAR(64) UNIQUE NOT NULL,
     	 parsed_file BYTEA,
     	 status VARCHAR(20) CHECK (status IN ('queued', 'processing', 'parsed')) NOT NULL DEFAULT 'queued'
		 );`,

		`CREATE TABLE user_files (
    	user_id INT NOT NULL,
    	file_id INT NOT NULL,
    	filename VARCHAR(255) NOT NULL,
    	upload_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    	PRIMARY KEY (user_id, file_id),
    	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    	FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
		);`,

		`CREATE TABLE queue (
    	id SERIAL PRIMARY KEY,
    	file_id INT NOT NULL,
    	FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
		);`,
	}

	for _, query := range queries {
		_, err := db.Pool.Exec(context.Background(), query)
		if err != nil {
			log.Println("Unable to create tables")
			return err
		}
	}

	log.Println("Tables created or already exist")
	return nil
}

// Close closes the database connection
func (db *PostgreSQLConnection) Close() {
	db.Pool.Close()
	log.Println("Database connection closed")
}

// GetPool returns the database connection pool
func (db *PostgreSQLConnection) GetPool() *pgxpool.Pool {
	return db.Pool
}
