package http

import (
	"bash_api/internal/bashService"
	"bash_api/pkg/tokenManager"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
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

// @Summary      GetCommand
// @Description  Get bash command
// @Tags         command
// @Accept       json
// @Produce      json
// @Param 		 Authorization 	header 	string	true  "Authorization"
// @Param        cmd_id			path 	int		true  "command id"
// @Success      200  {object}	bashService.Command
// @Failure      400  {object}	error
// @Failure      401  {object}	error
// @Failure      500  {object}  error
// @Router       /get_cmd/{cmd_id} [get]
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

		idStr := c.Params("cmd_id")
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

// @Summary      GetList
// @Description  Get list bash commands
// @Tags         command
// @Accept       json
// @Produce      json
// @Param 		 Authorization 	header 	string	true  "Authorization"
// @Param 		 Limit		header 	int	false  "Limit"
// @Param 		 Offset		header 	int	false  "Offset"
// @Param 		 AuthorId	header 	int	false  "Author Id"
// @Success      200  {array}	bashService.Command
// @Failure      400  {object}	error
// @Failure      401  {object}	error
// @Failure      500  {object}  error
// @Router       /list [get]
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

// @Summary      CreateCommand
// @Description  Add bash command
// @Tags         command
// @Accept       json
// @Produce      json
// @Param 		 Authorization 	header 	string	true  "Authorization"
// @Param        input			body	bashService.CreateCommandParams  true   "command data"
// @Success      200  {int}		int		"command id"
// @Failure      400  {object}	error
// @Failure      401  {object}	error
// @Failure      500  {object}  error
// @Router       /create_cmd [post]
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

// @Summary      DeleteCommand
// @Description  Delete bash command
// @Tags         command
// @Accept       json
// @Produce      json
// @Param 		 Authorization 	header 	string	true  "Authorization"
// @Param        cmd_id			path 	int		true  "command id"
// @Success      200  {object}	nil
// @Failure      400  {object}	error
// @Failure      401  {object}	error
// @Failure      500  {object}  error
// @Router       /delete/{cmd_id} [delete]
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

		idStr := c.Params("cmd_id")
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

// @Summary      RunCommand
// @Description  Run bash command
// @Tags         run
// @Accept       json
// @Produce      json
// @Param 		 Authorization 	header 	string	true  "Authorization"
// @Param        cmd_id			path 	int		true  "command id"
// @Success      201  {object}	int		"run id"
// @Failure      400  {object}	error
// @Failure      401  {object}	error
// @Failure      500  {object}  error
// @Router       /run/{cmd_id} [post]
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

		idStr := c.Params("cmd_id")
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

// @Summary      KillRun
// @Description  Kill run proccess
// @Tags         run
// @Accept       json
// @Produce      json
// @Param 		 Authorization 	header 	string	true  "Authorization"
// @Param        run_id			path 	int		true  "run id"
// @Success      200  {object}	nil
// @Failure      400  {object}	error
// @Failure      401  {object}	error
// @Failure      500  {object}  error
// @Router       /run/{run_id} [get]
func (h *BashServiceHandler) KillRun() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenData := c.Locals("tokenData")
		if tokenData == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "re-login"})
		}
		data, ok := tokenData.(*tokenManager.Data)
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "re-login"})
		}

		idStr := c.Params("run_id")
		if idStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no id"})
		}
		runId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad id"})
		}

		err = h.bashUC.KillRun(data.Id, data.Role, runId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		c.Status(fiber.StatusOK)
		return nil
	}
}

// @Summary      GetRun
// @Description  Get run results
// @Tags         run
// @Accept       json
// @Produce      json
// @Param 		 Authorization 	header 	string	true  "Authorization"
// @Param        run_id			path 	int		true  "run id"
// @Success      200  {object}	bashService.Result
// @Failure      400  {object}	error
// @Failure      401  {object}	error
// @Failure      500  {object}  error
// @Router       /get_run/{run_id} [get]
func (h *BashServiceHandler) GetRun() fiber.Handler {
	return func(c *fiber.Ctx) error {

		tokenData := c.Locals("tokenData")
		if tokenData == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "re-login"})
		}
		_, ok := tokenData.(*tokenManager.Data)

		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "re-login"})
		}

		idStr := c.Params("run_id")
		if idStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no id"})
		}
		runId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad id"})
		}

		content, err := h.bashUC.GetRunResult(runId)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				return c.Status(fiber.StatusOK).JSON(fiber.Map{"error": err.Error()})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(content)
	}
}

// @Summary      GetPersonResults
// @Description  Get list run results of a certain user (if you are admin or creator)
// @Tags         command
// @Accept       json
// @Produce      json
// @Param 		 Authorization 	header 	string	true  "Authorization"
// @Param 		 Limit		header 	int	false  "Limit"
// @Param 		 Offset		header 	int	false  "Offset"
// @Param 		 AuthorId	header 	int	false  "Author Id"
// @Success      200  {array}	bashService.Result
// @Failure      400  {object}	error
// @Failure      401  {object}	error
// @Failure      500  {object}  error
// @Router       /run_list [get]
func (h *BashServiceHandler) GetPersonResults() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenData := c.Locals("tokenData")
		if tokenData == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "re-login"})
		}
		data, ok := tokenData.(*tokenManager.Data)
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "re-login"})
		}

		headers := c.GetReqHeaders()
		limitStr, ok := headers["Limit"]
		if !ok {
			limitStr = []string{"5"}
		}
		limit, err := strconv.ParseInt(limitStr[0], 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad limit"})
		}
		offsetStr, ok := headers["Offset"]
		if !ok {
			offsetStr = []string{"0"}
		}
		offset, err := strconv.ParseInt(offsetStr[0], 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad offset"})
		}
		authorStr, ok := headers["Author-Id"]
		log.Println(authorStr, ok)
		if !ok {
			authorStr = []string{fmt.Sprintf("%d", data.Id)}
		}
		authorId, err := strconv.ParseInt(authorStr[0], 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad author id"})
		}

		list, err := h.bashUC.GetPersonResult(&bashService.GetListParams{
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
