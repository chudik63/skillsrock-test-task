package middleware

import (
	"log"
	"skillsrock-test-task/pkg/logger"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LoggingMiddleware(logger logger.Logger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		startTime := time.Now()

		err := ctx.Next()

		endTime := time.Now()
		log.Printf("Request: %s %s, Duration: %v", ctx.Method(), ctx.Path(), endTime.Sub(startTime))

		return err
	}
}
