package http

import (
	"bash_api/internal/auth"
	"bash_api/internal/cconstant"
	"github.com/gofiber/fiber/v2"
	"log"
	"regexp"
	"strconv"
)

type AuthHandler struct {
	authUC auth.Usecase
}

func NewAuthHandler(authUC auth.Usecase) *AuthHandler {
	return &AuthHandler{
		authUC: authUC,
	}
}

func (h *AuthHandler) SignUp() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			params auth.SignUpParams
		)

		if err := c.BodyParser(&params); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		pattern, _ := regexp.Compile("[A-Za-z0-9@.]+")
		if !pattern.MatchString(params.Email) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email must contain the characters a-z, A-z, 0-9, @ and ."})
		}

		err := h.authUC.SignUp(&params)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		c.Status(fiber.StatusCreated)
		return nil
	}
}

func (h *AuthHandler) SignIn() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			params auth.SignInParams
		)

		if err := c.BodyParser(&params); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		pattern, _ := regexp.Compile("[A-Za-z0-9@.]+")
		if !pattern.MatchString(params.Email) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email must contain the characters a-z, A-z, 0-9, @ and ."})
		}

		tokens, err := h.authUC.SignIn(&params)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(tokens)
	}
}

func (h *AuthHandler) RefreshTokens() fiber.Handler {
	return func(c *fiber.Ctx) error {
		headers := c.GetReqHeaders()
		refreshToken, ok := headers["Refresh-Token"]
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "no refresh token"})
		}

		log.Println("RefreshToken:", refreshToken)

		tokens, err := h.authUC.RefreshTokens(refreshToken[0])
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(tokens)
	}
}

func (h *AuthHandler) BeAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		headers := c.GetReqHeaders()
		secret, ok := headers["Secret-Key"]
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "no secret key"})
		}
		if secret[0] != cconstant.SecterKey {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid secret key"})
		}

		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad id"})
		}

		err = h.authUC.BeAdmin(id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		c.Status(fiber.StatusOK)
		return nil
	}
}
