package routes

import (
	"PDFStoring/web/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, userHendler handlers.UserApi, fileHandler handlers.FileApi, queueHandler handlers.QueueApi) {
	setupUserRoutes(app, userHendler)
	setupFileRoutes(app, fileHandler)
	setupQueueRoutes(app, queueHandler)
}

func setupUserRoutes(app *fiber.App, handler handlers.UserApi) {
	app.Post("/user", handler.CreateUser)
	app.Get("/user/:id", handler.GetUserFiles)
}

func setupFileRoutes(app *fiber.App, handler handlers.FileApi) {
	app.Post("/file/:id", handler.UploadFile)
	app.Delete("/file/:user_id/file_id/delete", handler.DeleteFile)
	app.Post("/file/:user_id/:file_id/import", handler.ImportFile)
}

func setupQueueRoutes(app *fiber.App, handler handlers.QueueApi) {
	app.Get("/queue/", handler.GetQueue)
	app.Get("/queue/:id", handler.UploadFile)
}
