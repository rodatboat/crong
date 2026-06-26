package resp

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func Send(ctx fiber.Ctx, resp APIResponse) error {
	log.Infof("Sending response: %v", resp)
	return ctx.Status(resp.Status).JSON(resp)
}
