package service

import (
	"PDFStoring/database"
	er "PDFStoring/error"
	"PDFStoring/models"
	"context"
	"log"
	"time"
)

type UserServiceStruct struct {
	dbService database.DatabaseService
}

// UserService interface defines methods for user-related operations
type UserService interface {
	CreateUser(ctx context.Context) (int, error)
	GetUserFiles(ctx context.Context, userId int) ([]models.UserFile, error)
}

// NewUserService creates a new instance of UserServiceStruct, implementing UserService
func NewUserService(dbService database.DatabaseService) UserService {
	return &UserServiceStruct{
		dbService: dbService,
	}
}

// CreateUser creates a new user in the database
func (s *UserServiceStruct) CreateUser(ctx context.Context) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `INSERT INTO user () VALUES () RETURNING id`

	var userId int
	err := s.dbService.GetPool().QueryRow(ctx, query).Scan(&userId)
	if err != nil {
		if er.HandleDeadlineExceededError(err) != nil {
			return 0, err
		}
		log.Printf("Error creating user: %v", err)
		return 0, err
	}

	return userId, nil
}

// GetUserFiles retrieves all files uploaded by a user
func (s *UserServiceStruct) GetUserFiles(ctx context.Context, userId int) ([]models.UserFile, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
	SELECT uf.filename, uf.upload_date, f.status, f.parsed_file 
	FROM user_files uf
	INNER JOIN files f ON uf.file_id = f.id
	WHERE uf.user_id = $1
	ORDER BY uf.upload_date DESC
	`

	rows, err := s.dbService.GetPool().Query(ctx, query, userId)
	if err != nil {
		if er.HandleDeadlineExceededError(err) != nil {
			return nil, err
		}
		log.Printf("Error fetching user files: %v", err)
		return nil, err
	}
	defer rows.Close()

	var userFiles []models.UserFile
	for rows.Next() {
		var userFile models.UserFile
		err := rows.Scan(&userFile.Filename, &userFile.UploadDate, &userFile.Status)
		if err != nil {
			log.Printf("Error scanning user files: %v", err)
			return nil, err
		}
		userFiles = append(userFiles, userFile)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return nil, err
	}

	return userFiles, nil
}
