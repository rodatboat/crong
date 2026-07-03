package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/rodatboat/crong/internal/middleware"
	"github.com/rodatboat/crong/internal/models"
	"github.com/rodatboat/crong/internal/resp"
	"github.com/rodatboat/crong/internal/services"
	"github.com/rodatboat/crong/internal/utils"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) LoginUser(c fiber.Ctx) error {
	var req models.UserLogin
	if err := c.Bind().Body(&req); err != nil {
		return resp.Send(c, resp.BadRequest())
	}

	// Validate request
	if validationErrors, err := utils.ValidateStruct(&req); err != nil {
		return resp.HandleValidationError(c, err, validationErrors)
	}

	user, err := h.userService.LoginUser(req.Email, req.Password)
	if err != nil {
		return resp.ErrorResponse(c, err)
	}

	return resp.Send(c, resp.Success(user))
}

func (h *UserHandler) RegisterUser(c fiber.Ctx) error {
	var req models.UserRegister
	if err := c.Bind().Body(&req); err != nil {
		return resp.Send(c, resp.BadRequest())
	}

	// Validate request
	if validationErrors, err := utils.ValidateStruct(&req); err != nil {
		return resp.HandleValidationError(c, err, validationErrors)
	}

	user, err := h.userService.RegisterUser(req.Email, req.FirstName, req.LastName, req.Password)
	if err != nil {
		if errors.Is(err, resp.ErrUserAlreadyExists) {
			return resp.Send(c, resp.Response(fiber.StatusConflict, resp.ErrUserAlreadyExists.Error(), nil))
		}
		return resp.ErrorResponse(c, err)
	}

	return resp.Send(c, resp.Success(user))
}

func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	var req models.UserUpdate
	if err := c.Bind().Body(&req); err != nil {
		return resp.Send(c, resp.BadRequest())
	}

	// Validate request
	if validationErrors, err := utils.ValidateStruct(&req); err != nil {
		return resp.HandleValidationError(c, err, validationErrors)
	}

	auth := c.Locals(middleware.AuthContextKey).(*middleware.AuthContext)

	user, err := h.userService.UpdateUser(auth.UserID, req)
	if err != nil {
		return resp.ErrorResponse(c, err)
	}

	return resp.Send(c, resp.Success(user))
}

func (h *UserHandler) GetUserDetails(c fiber.Ctx) error {
	auth := c.Locals(middleware.AuthContextKey).(*middleware.AuthContext)

	return resp.Send(c, resp.Success(auth.User))
}
