package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

func registerAdminAISettingsRoutes(api fiber.Router, r *Routes) {
	api.Get("", r.adminGetAISettings)
	api.Put("", r.adminUpdateAISettings)
}

// @Summary Get AI settings
// @Tags Admin: AI Settings
// @Security AdminAuth
// @Produce json
// @Success 200 {object} entity.AISettings
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/ai-settings [get]
func (r *Routes) adminGetAISettings(ctx *fiber.Ctx) error {
	settings, err := r.uc.AI.Get(ctx.UserContext())
	if err != nil {
		r.l.Error(err, "http - v1 - adminGetAISettings")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to load settings")
	}

	return ctx.Status(http.StatusOK).JSON(settings)
}

// @Summary Update AI settings
// @Tags Admin: AI Settings
// @Security AdminAuth
// @Accept json
// @Produce json
// @Param request body entity.AISettings true "Settings payload"
// @Success 200 {object} entity.AISettings
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/ai-settings [put]
func (r *Routes) adminUpdateAISettings(ctx *fiber.Ctx) error {
	var payload entity.AISettings
	if err := ctx.BodyParser(&payload); err != nil {
		r.l.Error(err, "http - v1 - adminUpdateAISettings - parse")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	if err := r.v.Struct(payload); err != nil {
		r.l.Error(err, "http - v1 - adminUpdateAISettings - validation")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	updated, err := r.uc.AI.Update(ctx.UserContext(), payload)
	if err != nil {
		r.l.Error(err, "http - v1 - adminUpdateAISettings - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to update settings")
	}

	return ctx.Status(http.StatusOK).JSON(updated)
}
