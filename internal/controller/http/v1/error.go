package v1

import "github.com/gofiber/fiber/v2"

func errorResponse(ctx *fiber.Ctx, code int, msg string) error {
	return ctx.Status(code).JSON(ErrorResponse{Error: msg})
}
