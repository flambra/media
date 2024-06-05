package internal

import (
	"os"

	"github.com/flambra/account/internal/image"
	"github.com/flambra/account/internal/middleware"
	"github.com/flambra/account/internal/video"
	"github.com/gofiber/fiber/v2"
)

func InitializeRoutes(app *fiber.App) {
	app.Get("/", middleware.Auth, func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"project":     os.Getenv("PROJECT"),
			"environment": os.Getenv("ENV"),
			"version":     os.Getenv("BUILD_VERSION"),
		})
	})

	// Image
	app.Post("/image/upload", middleware.Auth, image.Upload)

	// Video
	app.Post("/video/upload", middleware.Auth, video.Upload)
}
