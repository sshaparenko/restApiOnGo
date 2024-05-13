package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/sshaparenko/restApiOnGo/internal/database"
	"github.com/sshaparenko/restApiOnGo/internal/routes"
	"github.com/sshaparenko/restApiOnGo/internal/utils"
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
		utils.GetValue("DB_NAME"),
		utils.GetValue("DB_PORT"),
		utils.GetValue("DB_USER"),
		utils.GetValue("DB_PASSWORD"),
	)

	var PORT string = os.Getenv("PORT")
	if PORT == "" {
		PORT = DEFAULT_PORT
	}

	app.Listen(fmt.Sprintf(":%s", PORT))
}