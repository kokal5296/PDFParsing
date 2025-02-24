package handlers

import (
	"PDFStoring/service"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"strconv"
)

type FileApiStruct struct {
	fileService service.FileService
}

type FileApi interface {
	UploadFile(c *fiber.Ctx) error
	DeleteFile(c *fiber.Ctx) error
	ImportFile(c *fiber.Ctx) error
}

// NewFileApiService creates a new instance of FileApiStruct, which implements the FileApi interface
func NewFileApiService(fileService service.FileService) FileApi {
	return &FileApiStruct{
		fileService: fileService,
	}
}

func (s *FileApiStruct) UploadFile(c *fiber.Ctx) error {

	id := c.Params("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Error while converting id to int: %v", err)
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("Error while getting file: %v", err)
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	const maxFileSize = 10 << 20
	if file.Size > maxFileSize {
		log.Println("File is too large")
		return c.Status(http.StatusBadRequest).SendString("File is too large, it must be less than 10MB")
	}

	if file.Header.Get("Content-Type") != "application/pdf" {
		log.Println("Invalid file type")
		return c.Status(http.StatusBadRequest).SendString("Invalid file type, only PDF files are allowed")
	}

	fileId, err := s.fileService.UploadFile(c.Context(), userId, file)
	if err != nil {
		log.Printf("Error uploading file: %v", err)
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusCreated).SendString("File was successfully uploaded with id: " + string(fileId))
}

func (s *FileApiStruct) DeleteFile(c *fiber.Ctx) error {

	fId := c.Params("file_id")
	uId := c.Params("user_id")
	fileId, err := strconv.Atoi(fId)
	if err != nil {
		log.Printf("Error while converting id to int: %v", err)
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}
	userId, err := strconv.Atoi(uId)
	if err != nil {
		log.Printf("Error while converting id to int: %v", err)
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	err = s.fileService.DeleteFile(c.Context(), userId, fileId)
	if err != nil {
		log.Printf("Error deleting file: %v", err)
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusOK).SendString("File was successfully deleted")
}

func (s *FileApiStruct) ImportFile(c *fiber.Ctx) error {

	uId := c.Params("user_id")
	fId := c.Params("file_id")
	userId, err := strconv.Atoi(uId)
	if err != nil {
		log.Printf("Error while converting id to int: %v", err)
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	fileId, err := strconv.Atoi(fId)
	if err != nil {
		log.Printf("Error while converting id to int: %v", err)
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	err = s.fileService.ImportFile(c.Context(), userId, fileId)
	if err != nil {
		log.Printf("Error importing file: %v", err)
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusCreated).SendString("File was successfully imported with id: " + string(fileId))
}
