package response

import "github.com/gofiber/fiber/v3"

func Success(ctx fiber.Ctx, data any) error {
	return ctx.Status(fiber.StatusOK).JSON(Response[any]{
		Status:  fiber.StatusOK,
		Message: "Success",
		Data:    &data,
	})
}

func Error(ctx fiber.Ctx, status int, err string) error {
	return ctx.Status(status).JSON(Response[any]{
		Status:  status,
		Message: "Error",
		Error:   err,
	})
}

func SendResponse[T any](ctx fiber.Ctx, status int, message string, data T, err string) error {
	return ctx.Status(status).JSON(Response[T]{
		Status:  status,
		Message: message,
		Error:   err,
		Data:    &data,
	})
}
