package middlewares

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/sshaparenko/restApiOnGo/internal/utils"
)

func CreateMiddleware() func(*fiber.Ctx) error {
	config := jwtware.Config{
		SigningKey: []byte(utils.GetValue("JWT_SECRET_KEY")),
		ContextKey: "jwt",
		ErrorHandler: jwtError,
	}

	return jwtware.New(config)
}
//jwtError returns error handler for JWT middleware
func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": err.Error(),
	})
}