package v1

import (
	"errors"
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/usecase/auth"
	"github.com/gofiber/fiber/v2"
)

func registerAuthRoutes(api fiber.Router, r *Routes) {
	api.Post("/telegram", r.authTelegram)
	api.Post("/admin/login", r.adminLogin)
}

func registerAdminAuthRoutes(api fiber.Router, r *Routes) {
	api.Get("/me", r.adminProfile)
}

// @Summary Login or signup via Telegram
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body entity.TelegramAuthRequest true "Telegram payload"
// @Success 200 {object} entity.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
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

// @Summary Login as admin
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body entity.AdminLoginRequest true "Admin credentials"
// @Success 200 {object} entity.AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/admin/login [post]
func (r *Routes) adminLogin(ctx *fiber.Ctx) error {
	var payload entity.AdminLoginRequest

	if err := ctx.BodyParser(&payload); err != nil {
		r.l.Error(err, "http - v1 - adminLogin - parse")
		return errorResponse(ctx, http.StatusBadRequest, "invalid payload")
	}

	if err := r.v.Struct(payload); err != nil {
		r.l.Error(err, "http - v1 - adminLogin - validation")
		return errorResponse(ctx, http.StatusBadRequest, "invalid payload")
	}

	result, err := r.uc.Auth.AdminLogin(ctx.UserContext(), payload)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidAdminCredentials) {
			return errorResponse(ctx, http.StatusUnauthorized, "invalid credentials")
		}
		r.l.Error(err, "http - v1 - adminLogin - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "authentication failed")
	}

	return ctx.Status(http.StatusOK).JSON(result)
}

// @Summary Admin profile
// @Tags Auth
// @Security AdminAuth
// @Produce json
// @Success 200 {object} entity.AdminProfile
// @Failure 500 {object} ErrorResponse
// @Router /auth/admin/me [get]
func (r *Routes) adminProfile(ctx *fiber.Ctx) error {
	profile, err := r.uc.Admin.Profile(ctx.UserContext())
	if err != nil {
		r.l.Error(err, "http - v1 - adminProfile - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to load profile")
	}

	return ctx.Status(http.StatusOK).JSON(profile)
}
