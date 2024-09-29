package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sshaparenko/restApiOnGo/pkg/database"
	"github.com/sshaparenko/restApiOnGo/pkg/routes"
	"gorm.io/gorm"
)

const DEFAULT_PORT = "8080"

var isPrinted = false

func NewFiberApp() *fiber.App {
	var app *fiber.App = fiber.New()
	routes.SetupRoutes(app)
	return app
}

func main() {
	timeout := 10 * time.Second

	server := Server{
		App: NewFiberApp(),
	}

	dbCtx, dbCtxCancel := context.WithTimeout(context.Background(), timeout)
	defer dbCtxCancel()

	success := make(chan struct{}, 1)

	server.ConnectDB(dbCtx, success)
	handleConnection(dbCtx, success, &server)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := server.Run(ctx); err != nil {
		log.Fatal(err)
	}
}

// Server is a struct that wraps
// fiber app and database for better
// code extensibility
type Server struct {
	App *fiber.App
	DB  *gorm.DB
}

// Run is a function that starts the server
func (server *Server) Run(ctx context.Context) error {
	port := getPort()

	go func() {
		if err := server.App.Listen(fmt.Sprintf(":%s", port)); err != nil {
			log.Fatalf("main app.runServer: %s", err.Error())
		}
	}()

	log.Printf("listetning on %s\n", port)

	<-ctx.Done()

	if err := server.Shutdown(); err != nil {
		return err
	}

	if err := server.ClearResources(); err != nil {
		return fmt.Errorf("clearing resources in main.gracefulShutdown: %w", err)
	}

	return nil
}

// Sutdown is a function that finishes server gracefully
func (server *Server) Shutdown() error {
	timeout := 5 * time.Second

	log.Println("graceful shutdown started")

	if err := server.App.ShutdownWithTimeout(timeout); err != nil {
		return fmt.Errorf("shut down with timeout in main.gracefulShutdown: %w", err)
	}

	log.Println("server shutdown was successful")
	return nil
}

// ClearupResources closes the database connection
func (server *Server) ClearResources() error {
	sqlDB, err := server.DB.DB()
	if err != nil {
		return fmt.Errorf("obtain sql.DB in main.clearResources: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("closing database connection in main.clearResources: %w", err)
	}

	log.Println("resources were cleared successfully")
	return nil
}

// ConnectDB handles database connection on startup
// In case if database refuses connection to it, server will try
// to reconnect while timeout is not exceeded
func (server *Server) ConnectDB(ctx context.Context, success chan struct{}) {
	err := database.InitDatasource(
		os.Getenv("POSTGRES_DATABASE"),
		os.Getenv("POSTGRES_PORT"),
	)
	if err != nil {
		innerErr := errors.Unwrap(err)
		switch innerErr.(type) {
		case *pgconn.ConnectError:
			if !isPrinted {
				log.Println(innerErr)
				log.Println("trying to reconnect")
				isPrinted = true
			}
			server.ConnectDB(ctx, success)
		default:
			log.Fatalf("Failed to initialize database: %v", err)
		}
	} else {
		server.DB = database.DB
		success <- struct{}{}
	}
}

// handleConnection tracks if the timeout for database
// reconnect was exceeded
func handleConnection(ctx context.Context, successConn chan struct{}, server *Server) {
	select {
	case <-ctx.Done():
		log.Println("timeout was exceeded")
		if err := server.Shutdown(); err != nil {
			log.Fatalln(err.Error())
		}
		return
	case <-successConn:
		log.Println("connected to database")
	}
}

func getPort() string {
	port := os.Getenv("PORT")

	if port == "" {
		return DEFAULT_PORT
	}
	return port
}
