package resp

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func Send(ctx fiber.Ctx, resp APIResponse) error {
	log.Infof("Sending response: %v", resp)
	return ctx.Status(resp.Status).JSON(resp)
}

func HandleError(ctx fiber.Ctx, err error) error {
	var response APIResponse
	if errors.Is(err, ErrNotFound) {
		response = NotFound()
	} else {
		response = InternalServerError()
	}
	return Send(ctx, response)
}

func HandleValidationError(ctx fiber.Ctx, err error, validationErrors map[string]string) error {
	if validationErrors == nil {
		return Send(ctx, InternalServerError())
	}
	return Send(ctx, ValidationError(validationErrors))
}
