package handlers

import (
	"PDFStoring/models"
	"PDFStoring/service"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type QueueApiStruct struct {
	queueService service.QueueService
}

type QueueApi interface {
	GetQueue(c *fiber.Ctx) error
	UploadFile(c *fiber.Ctx) error
}

// NewQueueApiService creates a new instance of QueueApiStruct, which implements the QueueApi interface
func NewQueueApiService(queueService service.QueueService) QueueApi {
	return &QueueApiStruct{
		queueService: queueService,
	}
}

func (s *QueueApiStruct) GetQueue(c *fiber.Ctx) error {

	fileId, fileData, err := s.queueService.GetNextFile(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).SendString("File ID: " + string(fileId) + " File Data: " + string(fileData))
}

func (s *QueueApiStruct) UploadFile(c *fiber.Ctx) error {
	id := c.Params("id")
	fileId, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	var parsedFileData models.Parser

	err = json.Unmarshal(c.Body(), &parsedFileData)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err = s.queueService.UploadParsedFile(c.Context(), fileId, parsedFileData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusCreated).SendString("File was successfully added to queue")
}
