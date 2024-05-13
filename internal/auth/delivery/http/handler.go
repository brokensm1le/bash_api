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

// @Summary      SignUp
// @Description  Create account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        input	body	auth.User  true  "user data"
// @Success      201  {object}	nil
// @Failure      400  {object}	error
// @Failure      500  {object}  error
// @Router       /auth/signUp [post]
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

// @Summary      SignIn
// @Description  Login
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        input	body	auth.SignInParams  true  "login and password"
// @Success      200  {object}	auth.TokensResponse
// @Failure      400  {object}	error
// @Failure      500  {object}  error
// @Router       /auth/signIn [post]
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

// @Summary      RefreshTokens
// @Description  Refresh tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param 		 Refresh-Token 	header 	string true  "Refresh-Token"
// @Success      200  {object}	auth.TokensResponse
// @Failure      401  {object}	error
// @Failure      500  {object}  error
// @Router       /auth/refreshTokens [post]
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

// @Summary      BeAdmin
// @Description  You're gonna be an admin
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param 		 Secret-Key 	header 	string true  "Secret-Key"
// @Param        input	body	auth.SignInParams  true  "login and password"
// @Success      200  {object}	nil
// @Failure      400  {object}	error
// @Failure      500  {object}  error
// @Router       /auth/beAdmin/{id} [post]
func (h *AuthHandler) BeAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		headers := c.GetReqHeaders()
		secret, ok := headers["Secret-Key"]
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no secret key"})
		}
		if secret[0] != cconstant.SecterKey {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid secret key"})
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
