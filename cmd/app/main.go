package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/sshaparenko/restApiOnGo/pkg/database"
	"github.com/sshaparenko/restApiOnGo/pkg/routes"
)

const DEFAULT_PORT = "8080"

func NewFiberApp() *fiber.App {
	var app *fiber.App = fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	routes.SetupRoutes(app)
	return app
}

func main() {
	var app *fiber.App = NewFiberApp()

	err := database.InitDatasource(
		os.Getenv("POSTGRES_DATABASE"),
		os.Getenv("POSTGRES_PORT"),
	)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	var port string = os.Getenv("PORT")
	if port == "" {
		port = DEFAULT_PORT
	}

	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))

}
