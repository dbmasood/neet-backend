package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/evrone/go-clean-template/internal/entity"
	adminusecase "github.com/evrone/go-clean-template/internal/usecase/admin"
	"github.com/gofiber/fiber/v2"
)

func registerAdminUsersRoutes(api fiber.Router, r *Routes) {
	api.Get("", r.adminListUsers)
	api.Post("", r.adminCreateUser)
	api.Patch("/:id", r.adminUpdateUser)
	api.Delete("/:id", r.adminDeleteUser)
	api.Post("/bulk-status", r.adminBulkStatus)
	api.Post("/bulk-delete", r.adminBulkDelete)
	api.Post("/invite", r.adminInviteUser)
}

// @Summary List admin users
// @Tags Admin: Users
// @Security AdminAuth
// @Produce json
// @Param page query int false "Page"
// @Param pageSize query int false "Page size"
// @Param status query string false "Comma separated statuses"
// @Param role query string false "Role"
// @Param username query string false "Username search"
// @Success 200 {object} entity.AdminUserList
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /admin/users [get]
func (r *Routes) adminListUsers(ctx *fiber.Ctx) error {
	filter := entity.AdminUserFilter{
		Page:     parseQueryInt(ctx, "page", 1),
		PageSize: parseQueryInt(ctx, "pageSize", 20),
		Username: ctx.Query("username"),
	}

	if statusRaw := ctx.Query("status"); statusRaw != "" {
		statuses, err := parseUserStatuses(statusRaw)
		if err != nil {
			return errorResponse(ctx, http.StatusBadRequest, err.Error())
		}
		filter.Statuses = statuses
	}

	if roleRaw := ctx.Query("role"); roleRaw != "" {
		role := entity.AdminUserRole(strings.ToLower(roleRaw))
		if !isValidAdminRole(role) {
			return errorResponse(ctx, http.StatusBadRequest, "invalid role")
		}
		filter.Role = &role
	}

	list, err := r.uc.Admin.ListUsers(ctx.UserContext(), filter)
	if err != nil {
		r.l.Error(err, "http - v1 - adminListUsers - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to load users")
	}

	return ctx.Status(http.StatusOK).JSON(list)
}

// @Summary Create admin user
// @Tags Admin: Users
// @Security AdminAuth
// @Accept json
// @Produce json
// @Param request body entity.AdminUserCreateRequest true "User payload"
// @Success 201 {object} entity.AdminUser
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /admin/users [post]
func (r *Routes) adminCreateUser(ctx *fiber.Ctx) error {
	var payload entity.AdminUserCreateRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	if err := r.v.Struct(payload); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	user, err := r.uc.Admin.CreateUser(ctx.UserContext(), payload)
	if err != nil {
		switch {
		case errors.Is(err, adminusecase.ErrDuplicateUsername), errors.Is(err, adminusecase.ErrDuplicateEmail):
			return errorResponse(ctx, http.StatusConflict, err.Error())
		default:
			r.l.Error(err, "http - v1 - adminCreateUser - usecase")
			return errorResponse(ctx, http.StatusInternalServerError, "unable to create user")
		}
	}

	return ctx.Status(http.StatusCreated).JSON(user)
}

// @Summary Update admin user
// @Tags Admin: Users
// @Security AdminAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body entity.AdminUserUpdateRequest true "Update payload"
// @Success 200 {object} entity.AdminUser
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /admin/users/{id} [patch]
func (r *Routes) adminUpdateUser(ctx *fiber.Ctx) error {
	id, err := parseUUID(ctx, "id")
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid id")
	}

	var payload entity.AdminUserUpdateRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	user, err := r.uc.Admin.UpdateUser(ctx.UserContext(), id, payload)
	if err != nil {
		switch {
		case errors.Is(err, adminusecase.ErrUserNotFound):
			return errorResponse(ctx, http.StatusNotFound, "user not found")
		case errors.Is(err, adminusecase.ErrDuplicateUsername), errors.Is(err, adminusecase.ErrDuplicateEmail):
			return errorResponse(ctx, http.StatusConflict, err.Error())
		default:
			r.l.Error(err, "http - v1 - adminUpdateUser - usecase")
			return errorResponse(ctx, http.StatusInternalServerError, "unable to update user")
		}
	}

	return ctx.Status(http.StatusOK).JSON(user)
}

// @Summary Delete admin user
// @Tags Admin: Users
// @Security AdminAuth
// @Param id path string true "User ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /admin/users/{id} [delete]
func (r *Routes) adminDeleteUser(ctx *fiber.Ctx) error {
	id, err := parseUUID(ctx, "id")
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid id")
	}

	if err := r.uc.Admin.DeleteUser(ctx.UserContext(), id); err != nil {
		if errors.Is(err, adminusecase.ErrUserNotFound) {
			return errorResponse(ctx, http.StatusNotFound, "user not found")
		}
		r.l.Error(err, "http - v1 - adminDeleteUser - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to delete user")
	}

	return ctx.SendStatus(http.StatusNoContent)
}

// @Summary Bulk status update
// @Tags Admin: Users
// @Security AdminAuth
// @Accept json
// @Produce json
// @Param request body entity.AdminBulkStatusRequest true "Bulk status payload"
// @Success 200 {object} entity.AdminBulkStatusResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /admin/users/bulk-status [post]
func (r *Routes) adminBulkStatus(ctx *fiber.Ctx) error {
	var payload entity.AdminBulkStatusRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	if err := r.v.Struct(payload); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	updated, err := r.uc.Admin.BulkStatus(ctx.UserContext(), payload)
	if err != nil {
		r.l.Error(err, "http - v1 - adminBulkStatus - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to update users")
	}

	return ctx.Status(http.StatusOK).JSON(entity.AdminBulkStatusResponse{Updated: updated})
}

// @Summary Bulk delete users
// @Tags Admin: Users
// @Security AdminAuth
// @Accept json
// @Produce json
// @Param request body entity.AdminBulkDeleteRequest true "Bulk delete payload"
// @Success 200 {object} entity.AdminBulkDeleteResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /admin/users/bulk-delete [post]
func (r *Routes) adminBulkDelete(ctx *fiber.Ctx) error {
	var payload entity.AdminBulkDeleteRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}
	if err := r.v.Struct(payload); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	deleted, err := r.uc.Admin.BulkDelete(ctx.UserContext(), payload)
	if err != nil {
		r.l.Error(err, "http - v1 - adminBulkDelete - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to delete users")
	}

	return ctx.Status(http.StatusOK).JSON(entity.AdminBulkDeleteResponse{Deleted: deleted})
}

// @Summary Invite admin user
// @Tags Admin: Users
// @Security AdminAuth
// @Accept json
// @Produce json
// @Param request body entity.AdminInviteRequest true "Invite payload"
// @Success 202 {object} entity.AdminInviteResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /admin/users/invite [post]
func (r *Routes) adminInviteUser(ctx *fiber.Ctx) error {
	var payload entity.AdminInviteRequest
	if err := ctx.BodyParser(&payload); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}
	if err := r.v.Struct(payload); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	resp, err := r.uc.Admin.InviteUser(ctx.UserContext(), payload)
	if err != nil {
		r.l.Error(err, "http - v1 - adminInviteUser - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to invite user")
	}

	return ctx.Status(http.StatusAccepted).JSON(resp)
}

func parseQueryInt(ctx *fiber.Ctx, key string, def int) int {
	if raw := ctx.Query(key); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			return parsed
		}
	}
	return def
}

func parseUserStatuses(raw string) ([]entity.AdminUserStatus, error) {
	parts := strings.Split(raw, ",")
	statuses := make([]entity.AdminUserStatus, 0, len(parts))
	for _, part := range parts {
		status := entity.AdminUserStatus(strings.ToLower(strings.TrimSpace(part)))
		if status == "" {
			continue
		}
		if !isValidAdminStatus(status) {
			return nil, fmt.Errorf("invalid status %s", status)
		}
		statuses = append(statuses, status)
	}
	return statuses, nil
}

func isValidAdminStatus(status entity.AdminUserStatus) bool {
	switch status {
	case entity.AdminUserStatusActive, entity.AdminUserStatusInactive, entity.AdminUserStatusInvited, entity.AdminUserStatusSuspended:
		return true
	default:
		return false
	}
}

func isValidAdminRole(role entity.AdminUserRole) bool {
	switch role {
	case entity.AdminUserRoleSuperAdmin, entity.AdminUserRoleAdmin, entity.AdminUserRoleManager, entity.AdminUserRoleCashier:
		return true
	default:
		return false
	}
}
