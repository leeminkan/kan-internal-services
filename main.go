package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"kan-internal-services/features/checkinout"
)

func main() {
	_ = godotenv.Load()
	app := fiber.New()

	checkinout.RegisterRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("[INFO] Starting server on :%s...", port)
	log.Fatal(app.Listen(":" + port))
}
