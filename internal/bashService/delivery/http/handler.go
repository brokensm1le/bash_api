package http

import (
	"bash_api/internal/bashService"
	"bash_api/pkg/tokenManager"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type BashServiceHandler struct {
	bashUC       bashService.Usecase
	tokenManager tokenManager.TokenManager
}

func NewBashServiceHandler(bashUC bashService.Usecase, tokenManager tokenManager.TokenManager) *BashServiceHandler {
	return &BashServiceHandler{
		bashUC:       bashUC,
		tokenManager: tokenManager,
	}
}

func (h *BashServiceHandler) GetCommand() fiber.Handler {
	return func(c *fiber.Ctx) error {

		tokenData := c.Locals("tokenData")
		if tokenData == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "re-login"})
		}
		_, ok := tokenData.(*tokenManager.Data)

		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "re-login"})
		}

		idStr := c.Params("id")
		if idStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no id"})
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad id"})
		}

		content, err := h.bashUC.GetCommand(id)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": err.Error()})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(content)
	}
}

func (h *BashServiceHandler) GetList() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenData := c.Locals("tokenData")
		if tokenData == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "re-login"})
		}
		_, ok := tokenData.(*tokenManager.Data)
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "re-login"})
		}

		limitStr := c.GetRespHeader("Limit", "5")
		limit, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad limit"})
		}
		offsetStr := c.GetRespHeader("Offset", "0")
		offset, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad offset"})
		}
		authorStr := c.GetRespHeader("AuthorId", "-1")
		authorId, err := strconv.ParseInt(authorStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad author id"})
		}

		list, err := h.bashUC.GetList(&bashService.GetListParams{
			Limit:    limit,
			Offset:   offset,
			AuthorId: authorId,
		})
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": err.Error()})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(list)
	}
}

func (h *BashServiceHandler) CreateCommand() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			params bashService.CreateCommandParams
		)

		tokenData := c.Locals("tokenData")
		if tokenData == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "re-login"})
		}
		data, ok := tokenData.(*tokenManager.Data)
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "re-login"})
		}

		if err := c.BodyParser(&params); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad body"})
		}
		params.AuthorId = data.Id

		cmdId, err := h.bashUC.CreateCommand(&params)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"command_id": cmdId})
	}
}

func (h *BashServiceHandler) DeleteCommand() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenData := c.Locals("tokenData")
		if tokenData == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "re-login"})
		}
		data, ok := tokenData.(*tokenManager.Data)
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "re-login"})
		}

		idStr := c.Params("id")
		if idStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no id"})
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad id"})
		}

		err = h.bashUC.DeleteCommand(id, data.Role, data.Id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		c.Status(fiber.StatusOK)
		return nil
	}
}

func (h *BashServiceHandler) RunCommand() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenData := c.Locals("tokenData")
		if tokenData == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "re-login"})
		}
		data, ok := tokenData.(*tokenManager.Data)
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "re-login"})
		}

		idStr := c.Params("id")
		if idStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no id"})
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad id"})
		}

		runId, err := h.bashUC.RunCommand(id, data.Id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"run_id": runId})
	}
}
