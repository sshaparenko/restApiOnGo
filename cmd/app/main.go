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
	)

	var port string = os.Getenv("PORT")
	if port == "" {
		port = DEFAULT_PORT
	}

	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
