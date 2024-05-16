package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/sshaparenko/restApiOnGo/internal/database"
	"github.com/sshaparenko/restApiOnGo/internal/routes"
)

const DEFAULT_PORT = "8080"

func NewFiberApp() *fiber.App {
	//create new fiber app
	var app *fiber.App = fiber.New()
	//set up routes
	routes.SetupRoutes(app)

	return app
}

func main() {
	var app *fiber.App = NewFiberApp()

	database.InitDatasource(
		os.Getenv("POSTGRES_DATABASE"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USERNAME"),
		os.Getenv("POSTGRES_PASSWORD"),
	)

	var PORT string = os.Getenv("PORT")
	if PORT == "" {
		PORT = DEFAULT_PORT
	}

	app.Listen(fmt.Sprintf(":%s", PORT))
}
