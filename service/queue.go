package service

import (
	"PDFStoring/database"
	er "PDFStoring/error"
	"PDFStoring/models"
	"context"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"time"
)

type QueueServiceStruct struct {
	dbService database.DatabaseService
}

// QueueService interface defines methods for user-related operations
type QueueService interface {
	AddFileToQueue(ctx context.Context, fileId int, file multipart.File) error
	GetNextFile(ctx context.Context) (int, []byte, error)
	UploadParsedFile(ctx context.Context, fileId int, parsedData models.Parser) error
}

// NewQueueService creates a new instance of QueueServiceStruct, implementing QueueService
func NewQueueService(dbService database.DatabaseService) QueueService {
	return &QueueServiceStruct{
		dbService: dbService,
	}
}

func (s *QueueServiceStruct) AddFileToQueue(ctx context.Context, fileId int, file multipart.File) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	fileData, err := io.ReadAll(file)
	if err != nil {
		if err == er.HandleDeadlineExceededError(err) {
			log.Println("Deadline exceeded while reading file data")
			return err
		}
		log.Printf("Error reading file data: %v", err)
		return err
	}

	query := `INSERT INTO queue (file_id, pdf_file) VALUES ($1, $2)`
	_, err = s.dbService.GetPool().Query(ctx, query, fileId, fileData)
	if err != nil {
		if err == er.HandleDeadlineExceededError(err) {
			log.Println("Deadline exceeded while adding file to queue")
			return err
		}
		log.Printf("Error adding file to queue: %v", err)
		return err
	}

	log.Println("File added to queue:", fileId)
	return nil
}

func (s *QueueServiceStruct) GetNextFile(ctx context.Context) (int, []byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `SELECT file_id, pdf_file FROM queue ORDER BY id LIMIT 1`

	var fileId int
	var hasedFile string
	err := s.dbService.GetPool().QueryRow(ctx, query).Scan(&fileId, &hasedFile)
	if err != nil {
		if err == er.HandleDeadlineExceededError(err) {
			log.Println("Deadline exceeded while getting next file from queue")
			return 0, nil, err
		}
		log.Printf("Error getting next file from queue: %v", err)
		return 0, nil, err
	}

	if fileId == 0 {
		return 0, nil, errors.New("File does not exist")
	}

	fileData, err := hex.DecodeString(hasedFile)
	if err != nil {
		log.Printf("Error decoding file data: %v", err)
		return 0, nil, err
	}

	query = `DELETE FROM queue WHERE file_id = $1`
	_, err = s.dbService.GetPool().Query(ctx, query, fileId)
	if err != nil {
		if err == er.HandleDeadlineExceededError(err) {
			log.Println("Deadline exceeded while deleting file from queue")
			return 0, nil, err
		}
		log.Printf("Error deleting file from queue: %v", err)
		return 0, nil, err
	}

	query = `UPDATE files SET status = $1 WHERE id = $2`
	_, err = s.dbService.GetPool().Query(ctx, query, Parsing, fileId)
	if err != nil {
		if err == er.HandleDeadlineExceededError(err) {
			log.Println("Deadline exceeded while updating file status")
			return 0, nil, err
		}
		log.Printf("Error updating file status: %v", err)
		return 0, nil, err
	}

	return fileId, fileData, nil
}

func (s *QueueServiceStruct) UploadParsedFile(ctx context.Context, fileId int, parsedData models.Parser) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if parsedData.ParsedError != "" {
		query := `UPDATE files SET status = $1, parsed_file = $2 WHERE id = $3`
		_, err := s.dbService.GetPool().Query(ctx, query, Error, parsedData, fileId)
		if err != nil {
			if err == er.HandleDeadlineExceededError(err) {
				log.Println("Deadline exceeded while updating file status")
				return err
			}
			log.Printf("Error updating file status: %v", err)
			return err
		}
	}

	query := `UPDATE files SET status = $1, parsed_file = $2 WHERE id = $3`
	_, err := s.dbService.GetPool().Query(ctx, query, parsedData.ParsedStatus, parsedData, fileId)
	if err != nil {
		if err == er.HandleDeadlineExceededError(err) {
			log.Println("Deadline exceeded while updating file status")
			return err
		}
		log.Printf("Error updating file status: %v", err)
		return err
	}

	return nil
}
