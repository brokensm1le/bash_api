package http

import (
	"github.com/gofiber/fiber/v2"
)

func MapRoutes(router fiber.Router, h *AuthHandler) {
	router.Post("/auth/signUp", h.SignUp())
	router.Post("/auth/signIn", h.SignIn())
	router.Post("/auth/refreshTokens", h.RefreshTokens())
	router.Post("/auth/beAdmin/:id", h.BeAdmin())
}
