package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sshaparenko/restApiOnGo/pkg/handlers"
	"github.com/sshaparenko/restApiOnGo/pkg/middlewares"
)

func SetupRoutes(app *fiber.App) {
	//public routes
	var publicRoutes fiber.Router = app.Group("/api/v1")

	publicRoutes.Post("/signup", handlers.Signup)
	publicRoutes.Post("/login", handlers.Login)
	publicRoutes.Get("/items", handlers.GetAllItems)
	publicRoutes.Get("/items/:id", handlers.GetItemByID)

	//private routes
	var privateRoutes fiber.Router = app.Group("/api/v1", middlewares.CreateMiddleware())

	privateRoutes.Post("/items", handlers.CreateItem)
	privateRoutes.Put("/items/:id", handlers.UpdateItem)
	privateRoutes.Delete("/items/:id", handlers.DeleteItem)
}
