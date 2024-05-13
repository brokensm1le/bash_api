package http

import (
	"github.com/gofiber/fiber/v2"
)

func MapRoutes(router fiber.Router, h *BashServiceHandler) {
	routerWithToken := router.Use("/", h.UserIdentity())
	routerWithToken.Post("/create_cmd", h.CreateCommand())
	routerWithToken.Get("/get_cmd/:cmd_id", h.GetCommand())
	routerWithToken.Delete("/delete/:cmd_id", h.DeleteCommand())
	routerWithToken.Get("/list", h.GetList())

	routerWithToken.Post("/run/:cmd_id", h.RunCommand())
	routerWithToken.Get("/get_run/:run_id", h.GetRun())
	routerWithToken.Get("/run_list", h.GetPersonResults())
	routerWithToken.Get("/kill/:run_id", h.KillRun())
}
