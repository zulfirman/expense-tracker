package main

import (
	"expenses-tracker/src/registry"
	"expenses-tracker/src/route"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize registry (database, repositories, handlers)
	reg, err := registry.NewRegistry()
	if err != nil {
		log.Fatal("Failed to initialize registry:", err)
	}
	//defer reg.Close()

	// Setup Echo
	e := echo.New()

	// Setup routes
	route.SetupRoutes(e, reg)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
