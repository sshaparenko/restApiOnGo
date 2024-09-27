package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := runServer(ctx, app); err != nil {
		log.Fatal(err)
	}
}

func runServer(ctx context.Context, app *fiber.App) error {
	port := os.Getenv("PORT")
	timeout := 5 * time.Second

	if port == "" {
		port = DEFAULT_PORT
	}

	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
			log.Fatalf("main app.Listen: %s", err.Error())
		}
	}()

	log.Printf("listetning on %s\n", port)

	<-ctx.Done()

	log.Println("graceful shutdown started")

	if err := app.ShutdownWithTimeout(timeout); err != nil {
		return fmt.Errorf("shut down with timeout in main.runServer: %w", err)
	}

	if err := clearResources(); err != nil {
		return fmt.Errorf("clearing resources in main.runServer: %w", err)
	}

	log.Println("servers shut down was successful")
	return nil
}

func clearResources() error {
	sqlDB, err := database.DB.DB()
	if err != nil {
		return fmt.Errorf("obtain sql.DB in main.clearResources: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("closing database connection in main.clearResources: %w", err)
	}

	log.Println("resources were cleared successfully")
	return nil
}
