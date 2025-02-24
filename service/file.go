package service

import (
	"PDFStoring/database"
	er "PDFStoring/error"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"time"
)

type FileServiceStruct struct {
	dbService    database.DatabaseService
	QueueService QueueServiceStruct
}

// FileService interface defines methods for user-related operations
type FileService interface {
	UploadFile(ctx context.Context, userId int, file *multipart.FileHeader) (int, error)
	DeleteFile(ctx context.Context, userId int, fileId int) error
	ImportFile(ctx context.Context, userId int, fileId int) error
}

// NewFileService creates a new instance of FileServiceStruct, implementing FileService
func NewFileService(dbService database.DatabaseService) FileService {
	return &FileServiceStruct{
		dbService: dbService,
	}
}

func (s *FileServiceStruct) UploadFile(ctx context.Context, userId int, file *multipart.FileHeader) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	uploadedFile, err := file.Open()
	if err != nil {
		if err == er.HandleDeadlineExceededError(err) {
			log.Println("Deadline exceeded while opening file")
			return 0, err
		}
		log.Printf("Error while opening file: %v", err)
		return 0, err
	}
	defer uploadedFile.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, uploadedFile); err != nil {
		if err == er.HandleDeadlineExceededError(err) {
			log.Println("Deadline exceeded while hashing file")
			return 0, err
		}
		log.Printf("Error while hashing file: %v", err)
		return 0, err
	}

	fileHash := hex.EncodeToString(hasher.Sum(nil))

	var fileId int
	query := `SELECT id FROM files WHERE hash = $1`
	err = s.dbService.GetPool().QueryRow(ctx, query, fileHash).Scan(&fileId)
	if err != nil {
		if err == er.HandleDeadlineExceededError(err) {
			log.Println("Deadline exceeded while checking if file exists")
			return 0, err
		}
		log.Printf("Error while checking if file exists: %v", err)
		return 0, err
	}
	if err == nil {
		err = s.insertUserFile(ctx, userId, fileId, file.Filename)
		if err != nil {
			if err == er.HandleDeadlineExceededError(err) {
				log.Println("Deadline exceeded while chechking user file")
				return 0, err
			}
			log.Printf("Error while inserting user file: %v", err)
			return 0, err
		}
		return fileId, nil
	}

	query = `INSERT INTO files (hash, filename, status, upload_date) VALUES ($1, $2, $3 $4) RETURNING id`
	err = s.dbService.GetPool().QueryRow(ctx, query, fileHash, file.Filename, InQueue, time.Now().Format("2006-01-02 15:04:05")).Scan(&fileId)
	if err != nil {
		if err == er.HandleDeadlineExceededError(err) {
			log.Println("Deadline exceeded while inserting file")
			return 0, err
		}
		log.Printf("Error while inserting file: %v", err)
		return 0, err
	}

	err = s.insertUserFile(ctx, userId, fileId, file.Filename)
	if err != nil {
		if err == er.HandleDeadlineExceededError(err) {
			log.Println("Deadline exceeded while inserting user file")
			return 0, err
		}
		log.Printf("Error while inserting user file: %v", err)
		return 0, err
	}

	err = s.QueueService.AddFileToQueue(ctx, fileId, uploadedFile)
	if err != nil {
		if err == er.HandleDeadlineExceededError(err) {
			log.Println("Deadline exceeded while adding file to queue")
			return 0, err
		}
		log.Printf("Error while adding file to queue: %v", err)
		return 0, err
	}

	return fileId, nil
}

func (s *FileServiceStruct) DeleteFile(ctx context.Context, userId int, fileId int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var status string
	query := `SELECT status FROM files WHERE id = $1`
	err := s.dbService.GetPool().QueryRow(ctx, query, fileId).Scan(&status)
	if err != nil {
		if err == er.HandleDeadlineExceededError(err) {
			log.Println("Deadline exceeded while checking file status")
			return err
		}
		log.Printf("Error while checking file status: %v", err)
		return err
	}

	if status != "in_queue" {
		log.Println("File is not in queue")
		return nil
	}

	query = `DELETE FROM user_files WHERE user_id = $1 AND file_id = $2`
	_, err = s.dbService.GetPool().Query(ctx, query, userId, fileId)
	if err != nil {
		if err == er.HandleDeadlineExceededError(err) {
			log.Println("Deadline exceeded while deleting user file")
			return err
		}
		log.Printf("Error while deleting user file: %v", err)
		return err
	}

	var storedFiles int
	err = s.dbService.GetPool().QueryRow(ctx, "SELECT COUNT(*) FROM user_files WHERE file_id = $1", fileId).Scan(&storedFiles)
	if err != nil {
		if err == er.HandleDeadlineExceededError(err) {
			log.Println("Deadline exceeded while checking stored files")
			return err
		}
		log.Printf("Error while checking stored files: %v", err)
		return err
	}

	if storedFiles == 0 {
		query = `DELETE FROM files WHERE id = $1`
		_, err = s.dbService.GetPool().Query(ctx, query, fileId)
		if err != nil {
			if err == er.HandleDeadlineExceededError(err) {
				log.Println("Deadline exceeded while deleting file")
				return err
			}
			log.Printf("Error while deleting file: %v", err)
			return err
		}
	}

	return nil
}

func (s *FileServiceStruct) ImportFile(ctx context.Context, userId int, fileId int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var userCount int
	query := `SELECT COUNT (*) FROM users WHERE id = $1`
	err := s.dbService.GetPool().QueryRow(ctx, query, userId).Scan(&userCount)
	if err != nil {
		if err == er.HandleDeadlineExceededError(err) {
			log.Println("Deadline exceeded while checking user count")
			return err
		}
		log.Printf("Error while checking user count: %v", err)
		return err
	}
	if userCount == 0 {
		log.Println("User does not exist")
		return nil
	}

	exist, err := s.userFileAlreadyExists(ctx, userId, fileId)
	if err != nil {
		log.Printf("Error while checking if user file exists: %v", err)
		return err
	}
	if !exist {
		log.Printf("User with this file does not exist")
		return nil
	}

	var status FileStatus
	query = `SELECT status FROM files WHERE id = $1`
	err = s.dbService.GetPool().QueryRow(ctx, query, fileId).Scan(&status)
	if err != nil {
		if err == er.HandleDeadlineExceededError(err) {
			log.Println("Deadline exceeded while checking file status")
			return err
		}
		log.Printf("Error while checking file status: %v", err)
		return err
	}

	if status != Success {
		return errors.New("File is not parsed")
	}

	query = `UPDATE files SET status = $1 WHERE id = $2`
	_, err = s.dbService.GetPool().Exec(ctx, query, Imported, fileId)
	if err != nil {
		if err == er.HandleDeadlineExceededError(err) {
			log.Println("Deadline exceeded while updating file")
			return err
		}
		log.Printf("Error while updating file: %v", err)
		return err
	}

	return nil
}

func (s *FileServiceStruct) insertUserFile(ctx context.Context, userId int, fileId int, filename string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	exists, err := s.userFileAlreadyExists(ctx, userId, fileId)
	if err != nil {
		log.Printf("Error while checking if user file exists: %v", err)
		return err
	}
	if exists {
		log.Printf("User file already exists")
		return nil
	}

	_, err = s.dbService.GetPool().Exec(ctx, "INSERT INTO user_files (user_id, file_id, filename, upload_date ) VALUES ($1, $2, $3, $4)", userId, fileId, filename, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		if err == er.HandleDeadlineExceededError(err) {
			log.Println("Deadline exceeded while inserting user file")
			return err
		}
		log.Printf("Error while inserting user file: %v", err)
		return err
	}

	return nil
}

func (s *FileServiceStruct) userFileAlreadyExists(ctx context.Context, userId int, fileId int) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var exists int
	err := s.dbService.GetPool().QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM user_files WHERE user_id = $1 AND file_id = $2)", userId, fileId).Scan(&exists)
	if err != nil {
		if err == er.HandleDeadlineExceededError(err) {
			log.Println("Deadline exceeded while checking if user file exists")
			return false, err
		}
		log.Printf("Error while checking if user file exists: %v", err)
		return false, err
		if err == sql.ErrNoRows {
			return false, nil
		}
	}

	return exists > 0, nil
}
