package main

import (
	"fmt"
	"log"
	"os"

	"github.com/flambra/account/internal"
	"github.com/flambra/account/internal/domain"
	"github.com/flambra/helpers/hDb"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}

	err = hDb.New()
	if err != nil {
		log.Fatal(err)
	}

	err = hDb.Migrate(
		&domain.Image{},
		&domain.Video{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate the database: %v", err)
		return
	}
}

func main() {
	app := fiber.New()

	fiber.SetParserDecoder(fiber.ParserConfig{
		IgnoreUnknownKeys: true,
		ZeroEmpty:         true,
	})

	internal.InitializeRoutes(app)

	port := os.Getenv("SERVER_PORT")
	if len(port) == 0 {
		port = "8080"
	}

	/* Start Server */
	err := app.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
}
