package http

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

func (h *BashServiceHandler) UserIdentity() fiber.Handler {
	return func(c *fiber.Ctx) error {
		headers := c.GetReqHeaders()

		header, ok := headers["Authorization"]
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		headerParts := strings.Split(header[0], " ")
		if len(headerParts) != 2 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid auth header"})
		}

		tokenData, err := h.tokenManager.Parse(headerParts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid auth header"})
		}

		c.Locals("tokenData", tokenData)
		return c.Next()
	}
}
