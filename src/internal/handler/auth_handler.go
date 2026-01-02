package handler

import (
	"pbmap_api/src/domain"
	"pbmap_api/src/internal/dto"
	"pbmap_api/src/internal/usecase"
	"pbmap_api/src/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService usecase.AuthService
	validator   *validator.Wrapper
}

func NewAuthHandler(authService usecase.AuthService, v *validator.Wrapper) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   v,
	}
}

func (h *AuthHandler) LoginWithSocial(c *fiber.Ctx) error {
	var req dto.SocialLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Cannot parse JSON",
		})
	}

	if errors := h.validator.Validate(req); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(domain.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Validation failed",
			Data:    errors,
		})
	}

	resp, err := h.authService.LoginWithSocial(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.APIResponse{
			Status:  fiber.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Login successfully",
		Data:    resp,
	})
}
