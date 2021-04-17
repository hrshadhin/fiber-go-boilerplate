package middleware

import (
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func FiberMiddleware(a *fiber.App) {
	a.Use(
		// Add CORS to each route.
		cors.New(),
		// Add simple logger.
		logger.New(),
	)
}

func IsAdmin(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	isAdmin, ok := claims["admin"]
	if !ok || !isAdmin.(bool) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"msg": "Forbidden",
		})
	}

	return c.Next()
}
