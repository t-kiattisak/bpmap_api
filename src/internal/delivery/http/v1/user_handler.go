package v1

import (
	"pbmap_api/src/internal/domain/entities"
	"pbmap_api/src/internal/dto"
	"pbmap_api/src/internal/usecase"
	"pbmap_api/src/pkg/auth"
	"pbmap_api/src/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UserHandler handles user CRUD.
type UserHandler struct {
	usecase    usecase.UserUsecase
	validator  *validator.Wrapper
	jwtService *auth.JWTService
}

// NewUserHandler creates the user HTTP handler.
func NewUserHandler(usecase usecase.UserUsecase, v *validator.Wrapper, jwtService *auth.JWTService) *UserHandler {
	return &UserHandler{usecase: usecase, validator: v, jwtService: jwtService}
}

// Create handles POST /api/users.
func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if errors := h.validator.Validate(req); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Validation failed",
			Data:    errors,
		})
	}

	user := entities.User{
		Email:       &req.Email,
		DisplayName: req.DisplayName,
		Role:        req.Role,
	}

	if err := h.usecase.CreateUser(c.Context(), &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.APIResponse{
		Status:  fiber.StatusCreated,
		Message: "User created successfully",
		Data:    user,
	})
}

// Get handles GET /api/users/:id.
func (h *UserHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid ID format",
		})
	}

	user, err := h.usecase.GetUser(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Status:  fiber.StatusNotFound,
			Message: "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  fiber.StatusOK,
		Message: "User retrieved successfully",
		Data:    user,
	})
}

// Update handles PUT /api/users/:id.
func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid ID format",
		})
	}

	var user entities.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	user.ID = id
	if err := h.usecase.UpdateUser(c.Context(), &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  fiber.StatusOK,
		Message: "User updated successfully",
		Data:    user,
	})
}

// Delete handles DELETE /api/users/:id.
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid ID format",
		})
	}

	if err := h.usecase.DeleteUser(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  fiber.StatusOK,
		Message: "User deleted successfully",
	})
}

// List handles GET /api/users.
func (h *UserHandler) List(c *fiber.Ctx) error {
	users, err := h.usecase.ListUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Users retrieved successfully",
		Data:    users,
	})
}

// Me handles GET /api/users/me.
func (h *UserHandler) Me(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	user, err := h.usecase.GetUser(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.APIResponse{
			Status:  fiber.StatusNotFound,
			Message: "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Current user retrieved successfully",
		Data:    user,
	})
}
