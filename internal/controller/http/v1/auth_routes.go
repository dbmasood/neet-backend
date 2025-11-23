package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

func registerAuthRoutes(api fiber.Router, r *Routes) {
	api.Post("/telegram", r.authTelegram)
}

// @Summary Login or signup via Telegram
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body entity.TelegramAuthRequest true "Telegram payload"
// @Success 200 {object} entity.AuthResponse
// @Failure 400 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /auth/telegram [post]
func (r *Routes) authTelegram(ctx *fiber.Ctx) error {
	var payload entity.TelegramAuthRequest

	if err := ctx.BodyParser(&payload); err != nil {
		r.l.Error(err, "http - v1 - authTelegram")
		return errorResponse(ctx, http.StatusBadRequest, "invalid payload")
	}

	if err := r.v.Struct(payload); err != nil {
		r.l.Error(err, "http - v1 - authTelegram - validation")
		return errorResponse(ctx, http.StatusBadRequest, "invalid payload")
	}

	result, err := r.uc.Auth.TelegramAuth(ctx.UserContext(), payload)
	if err != nil {
		r.l.Error(err, "http - v1 - authTelegram - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "authentication failed")
	}

	return ctx.Status(http.StatusOK).JSON(result)
}
