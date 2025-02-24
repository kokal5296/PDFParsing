package server

import (
	"PDFStoring/database"
	"PDFStoring/service"
	"PDFStoring/web/handlers"
	"PDFStoring/web/routes"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)

type Server struct {
	App        *fiber.App
	PostgreSQL *database.PostgreSQLConnection
}

// CreateServer initializes and confugures the server, database connection, services, handlers, and routes
func CreateServer(connStr, dbName string) *Server {
	app := fiber.New()

	// Initialize PostgreSQL connection
	databaseService := database.NewDatabaseService()
	db, err := databaseService.NewDatabase(connStr, dbName)
	if err != nil {
		log.Println(err.Error())
		panic("Cannot connect to PostgreSQL")
	}

	log.Println("Connected to PostgreSQL")

	// Service initialization
	userService := service.NewUserService(db)
	fileService := service.NewFileService(db)
	queueService := service.NewQueueService(db)

	// Handlers initialization
	userHandler := handlers.NewUserApiService(userService)
	fileHandler := handlers.NewFileApiService(fileService)
	queueHandler := handlers.NewQueueApiService(queueService)

	// Routes initialization
	routes.SetupRoutes(app, userHandler, fileHandler, queueHandler)

	// Server initialization
	server := &Server{
		App:        app,
		PostgreSQL: db,
	}

	return server
}

// Start begins the application server, listening on the configured port
func (s *Server) Start() error {
	if err := s.App.Listen(os.Getenv("PORT")); err != nil {
		log.Println("Could not initiates the server", err)
		return err
	}
	return nil
}

// Close gracefully shuts down the database connection when the server is stopped
func (s *Server) Close() {
	s.PostgreSQL.Close()
	log.Println("Server and database connection closed")
}
