package http

import (
	"pbmap_api/src/internal/domain"
	"pbmap_api/src/internal/delivery/http/middleware"
	"pbmap_api/src/internal/delivery/http/v1"
	"pbmap_api/src/internal/repository"
	"pbmap_api/src/pkg/auth"

	"github.com/gofiber/fiber/v2"
)

// Handlers holds all v1 HTTP handlers.
type Handlers struct {
	Alarm         *v1.AlarmHandler
	Auth          *v1.AuthHandler
	User          *v1.UserHandler
	Notification  *v1.NotificationHandler
}

// Router registers all routes and returns the Fiber app.
func Router(h *Handlers, jwtService *auth.JWTService, tokenRepo repository.TokenRepository) *fiber.App {
	app := fiber.New()

	api := app.Group("/api")
	api.Get("/health", healthCheck)

	v1 := api.Group("/v1")
	dispatch := v1.Group("/dispatch")
	dispatch.Post("/alarm", h.Alarm.Alarm)

	authGroup := api.Group("/auth")
	authGroup.Post("/login", h.Auth.LoginWithSocial)
	authGroup.Post("/logout", middleware.Protected(jwtService, tokenRepo), h.Auth.Logout)
	authGroup.Post("/refresh", h.Auth.RefreshToken)

	users := api.Group("/users")
	users.Post("/", h.User.Create)
	users.Get("/", h.User.List)
	users.Get("/me", middleware.Protected(jwtService, tokenRepo), h.User.Me)
	users.Get("/:id", h.User.Get)
	users.Put("/:id", h.User.Update)
	users.Delete("/:id", h.User.Delete)

	notifications := api.Group("/notifications")
	notifications.Post("/broadcast", h.Notification.Broadcast)
	notifications.Post("/subscribe", h.Notification.Subscribe)
	notifications.Post("/unsubscribe", h.Notification.Unsubscribe)

	return app
}

func healthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(domain.APIResponse{
		Status:  fiber.StatusOK,
		Message: "OK",
	})
}
