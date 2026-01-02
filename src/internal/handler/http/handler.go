package http

import (
	"pbmap_api/src/domain"
	"pbmap_api/src/internal/handler"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	UserHandler         *handler.UserHandler
	NotificationHandler *handler.NotificationHandler
	AuthHandler         *handler.AuthHandler
}

func NewHandler(userHandler *handler.UserHandler, notificationHandler *handler.NotificationHandler, authHandler *handler.AuthHandler) *Handler {
	return &Handler{
		UserHandler:         userHandler,
		NotificationHandler: notificationHandler,
		AuthHandler:         authHandler,
	}
}

func (h *Handler) HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(domain.APIResponse{
		Status:  fiber.StatusOK,
		Message: "OK",
	})
}
