package main

import (
	"golang-ai/internal/controller"
	"golang-ai/internal/pkg/serverutils"
	"golang-ai/internal/repository"
	"golang-ai/internal/service"
	"golang-ai/pkg/database"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024,
	})

	app.Use(serverutils.ErrorHandlerMiddleware())

	db := database.ConnectDB(os.Getenv("DB_CONNECTION_STRING"))

	exampleRepository := repository.NewExampleRepository(db)
	notebookRepository := repository.NewNotebookRepository(db)

	exampleService := service.NewExampleService(exampleRepository)
	notebookService := service.NewNotebookService(notebookRepository, db)

	exampleController := controller.NewExampleController(exampleService)
	notebookController := controller.NewNotebookController(notebookService)

	api := app.Group("/api")
	exampleController.RegisterRoutes(api)
	notebookController.RegisterRoutes(api)

	log.Fatal(app.Listen(":3000"))
}
