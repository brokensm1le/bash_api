package http

import (
	"github.com/gofiber/fiber/v2"
)

func MapRoutes(router fiber.Router, h *BashServiceHandler) {
	routerWithToken := router.Use("/", h.UserIdentity())
	routerWithToken.Post("/create_script", h.CreateCommand())
	routerWithToken.Get("/get_script/:id", h.GetCommand())
	routerWithToken.Delete("/delete/:id", h.DeleteCommand())
	routerWithToken.Get("/list", h.GetList())
	routerWithToken.Post("/run/:id", h.RunCommand())
}
