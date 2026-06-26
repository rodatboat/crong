package resp

import "github.com/gofiber/fiber/v3"

type APIResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func ValidationError(errors map[string]string) APIResponse {
	return APIResponse{
		Status:  fiber.StatusBadRequest,
		Message: ErrValidation.Error(),
		Data:    errors,
	}
}

func BadRequest() APIResponse {
	return APIResponse{
		Status:  fiber.StatusBadRequest,
		Message: ErrBadRequest.Error(),
	}
}

func Success(data any) APIResponse {
	return APIResponse{
		Status:  fiber.StatusOK,
		Message: "Success",
		Data:    data,
	}
}

func Unauthorized() APIResponse {
	return APIResponse{
		Status:  fiber.StatusUnauthorized,
		Message: ErrUnauthorized.Error(),
	}
}

func InternalServerError() APIResponse {
	return APIResponse{
		Status:  fiber.StatusInternalServerError,
		Message: ErrInternalServerError.Error(),
	}
}

func NotFound() APIResponse {
	return APIResponse{
		Status:  fiber.StatusNotFound,
		Message: ErrNotFound.Error(),
	}
}

func Response(status int, message string, data any) APIResponse {
	return APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
