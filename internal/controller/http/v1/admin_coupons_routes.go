package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

func registerAdminCouponsRoutes(api fiber.Router, r *Routes) {
	api.Get("", r.adminListCoupons)
	api.Post("", r.adminCreateCoupon)
	api.Get("/:id", r.adminGetCoupon)
	api.Patch("/:id", r.adminUpdateCoupon)
	api.Delete("/:id", r.adminDeleteCoupon)
}

// @Summary List coupons
// @Tags Admin: Coupons
// @Security AdminAuth
// @Produce json
// @Success 200 {array} entity.Coupon
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/coupons [get]
func (r *Routes) adminListCoupons(ctx *fiber.Ctx) error {
	coupons, err := r.uc.Coupon.AdminList(ctx.UserContext())
	if err != nil {
		r.l.Error(err, "http - v1 - adminListCoupons")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to list coupons")
	}

	return ctx.Status(http.StatusOK).JSON(coupons)
}

// @Summary Create coupon
// @Tags Admin: Coupons
// @Security AdminAuth
// @Accept json
// @Produce json
// @Param request body entity.CouponCreateRequest true "Coupon payload"
// @Success 201 {object} entity.Coupon
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/coupons [post]
func (r *Routes) adminCreateCoupon(ctx *fiber.Ctx) error {
	var payload entity.CouponCreateRequest
	if err := ctx.BodyParser(&payload); err != nil {
		r.l.Error(err, "http - v1 - adminCreateCoupon - parse")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	if err := r.v.Struct(payload); err != nil {
		r.l.Error(err, "http - v1 - adminCreateCoupon - validation")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	coupon, err := r.uc.Coupon.AdminCreate(ctx.UserContext(), payload)
	if err != nil {
		r.l.Error(err, "http - v1 - adminCreateCoupon - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to create coupon")
	}

	return ctx.Status(http.StatusCreated).JSON(coupon)
}

// @Summary Get coupon
// @Tags Admin: Coupons
// @Security AdminAuth
// @Produce json
// @Param id path string true "Coupon ID"
// @Success 200 {object} entity.Coupon
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/coupons/{id} [get]
func (r *Routes) adminGetCoupon(ctx *fiber.Ctx) error {
	id, err := parseUUID(ctx, "id")
	if err != nil {
		r.l.Error(err, "http - v1 - adminGetCoupon")
		return errorResponse(ctx, http.StatusBadRequest, "invalid id")
	}

	coupon, err := r.uc.Coupon.AdminGet(ctx.UserContext(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - adminGetCoupon - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to load coupon")
	}

	return ctx.Status(http.StatusOK).JSON(coupon)
}

// @Summary Update coupon
// @Tags Admin: Coupons
// @Security AdminAuth
// @Accept json
// @Produce json
// @Param id path string true "Coupon ID"
// @Param request body entity.CouponCreateRequest true "Coupon payload"
// @Success 200 {object} entity.Coupon
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/coupons/{id} [patch]
func (r *Routes) adminUpdateCoupon(ctx *fiber.Ctx) error {
	id, err := parseUUID(ctx, "id")
	if err != nil {
		r.l.Error(err, "http - v1 - adminUpdateCoupon")
		return errorResponse(ctx, http.StatusBadRequest, "invalid id")
	}

	var payload entity.CouponCreateRequest
	if err := ctx.BodyParser(&payload); err != nil {
		r.l.Error(err, "http - v1 - adminUpdateCoupon - parse")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	if err := r.v.Struct(payload); err != nil {
		r.l.Error(err, "http - v1 - adminUpdateCoupon - validation")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	updated, err := r.uc.Coupon.AdminUpdate(ctx.UserContext(), id, payload)
	if err != nil {
		r.l.Error(err, "http - v1 - adminUpdateCoupon - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to update coupon")
	}

	return ctx.Status(http.StatusOK).JSON(updated)
}

// @Summary Delete coupon
// @Tags Admin: Coupons
// @Security AdminAuth
// @Param id path string true "Coupon ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /admin/coupons/{id} [delete]
func (r *Routes) adminDeleteCoupon(ctx *fiber.Ctx) error {
	id, err := parseUUID(ctx, "id")
	if err != nil {
		r.l.Error(err, "http - v1 - adminDeleteCoupon")
		return errorResponse(ctx, http.StatusBadRequest, "invalid id")
	}

	if err := r.uc.Coupon.AdminDelete(ctx.UserContext(), id); err != nil {
		r.l.Error(err, "http - v1 - adminDeleteCoupon - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to delete coupon")
	}

	return ctx.SendStatus(http.StatusNoContent)
}
