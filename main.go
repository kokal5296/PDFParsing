package main

import (
	"PDFStoring/web/server"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {

	// Load the environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Get the connection string and database name from the environment variables
	connStr := os.Getenv("POSTGRESQL_URI")
	dbName := os.Getenv("DB_NAME")

	createServer := server.CreateServer(connStr, dbName)
	log.Println("Server started")

	err = createServer.Start()
	if err != nil {
		log.Fatalf("Error starting createServer: %v", err)
	}

}
