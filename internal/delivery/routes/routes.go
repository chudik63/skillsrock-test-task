package routes

import (
	"skillsrock-test-task/internal/delivery/http/v1/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

func RegistrateRoutes(app *fiber.App, c *handler.Handler) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Content-Type,Authorization",
	}))

	api := app.Group("/api")
	v1 := api.Group("/v1")

	_ = v1

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/docs/swagger.json",
	}))

	app.Get("/docs/*", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})
}
