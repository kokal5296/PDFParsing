package handlers

import (
	"PDFStoring/service"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"strconv"
)

type UserApiStruct struct {
	userService service.UserService
}

type UserApi interface {
	CreateUser(c *fiber.Ctx) error
	GetUserFiles(c *fiber.Ctx) error
}

// NewUserApiService creates a new instance of UserApiStruct, which implements the UserApi interface
func NewUserApiService(userService service.UserService) UserApi {
	return &UserApiStruct{
		userService: userService,
	}
}

// CreateUser handles the request to create a new user
func (s *UserApiStruct) CreateUser(c *fiber.Ctx) error {

	userId, err := s.userService.CreateUser(c.Context())
	if err != nil {
		log.Println("Error while creating user: %v", err)
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusCreated).SendString("User was successfully created with id: " + string(userId))
}

func (s *UserApiStruct) GetUserFiles(c *fiber.Ctx) error {

	id := c.Params("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("Error while converting id to int: %v", err)
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	userFiles, err := s.userService.GetUserFiles(c.Context(), userId)
	if err != nil {
		log.Printf("Error fetching user files: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch user files"})
	}

	return c.Status(http.StatusOK).JSON(userFiles)
}
