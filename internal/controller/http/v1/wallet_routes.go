package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/controller/http/middleware"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/gofiber/fiber/v2"
)

func registerWalletRoutes(api fiber.Router, r *Routes, auth *jwt.Service) {
	walletGroup := api.Group("/wallet")
	walletGroup.Use(middleware.UserAuth(auth))
	walletGroup.Get("", r.walletSummary)
	walletGroup.Get("/transactions", r.walletTransactions)

	couponGroup := api.Group("/coupons")
	couponGroup.Use(middleware.UserAuth(auth))
	couponGroup.Post("/redeem", r.redeemCoupon)

	referralGroup := api.Group("/referral")
	referralGroup.Use(middleware.UserAuth(auth))
	referralGroup.Get("", r.referralSummary)
}

// @Summary Wallet summary for current user
// @Tags App: Wallet
// @Security UserAuth
// @Produce json
// @Success 200 {object} entity.WalletSummary
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /wallet [get]
func (r *Routes) walletSummary(ctx *fiber.Ctx) error {
	userID, err := r.getUserID(ctx)
	if err != nil {
		r.l.Error(err, "http - v1 - walletSummary")
		return errorResponse(ctx, http.StatusUnauthorized, "missing user")
	}

	summary, err := r.uc.Wallet.Summary(ctx.UserContext(), userID)
	if err != nil {
		r.l.Error(err, "http - v1 - walletSummary - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to load wallet")
	}

	return ctx.Status(http.StatusOK).JSON(summary)
}

// @Summary Wallet transactions
// @Tags App: Wallet
// @Security UserAuth
// @Produce json
// @Success 200 {array} entity.WalletTransaction
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /wallet/transactions [get]
func (r *Routes) walletTransactions(ctx *fiber.Ctx) error {
	userID, err := r.getUserID(ctx)
	if err != nil {
		r.l.Error(err, "http - v1 - walletTransactions")
		return errorResponse(ctx, http.StatusUnauthorized, "missing user")
	}

	txs, err := r.uc.Wallet.ListTransactions(ctx.UserContext(), userID)
	if err != nil {
		r.l.Error(err, "http - v1 - walletTransactions - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to load transactions")
	}

	return ctx.Status(http.StatusOK).JSON(txs)
}

// @Summary Redeem coupon code
// @Tags App: Wallet
// @Security UserAuth
// @Accept json
// @Produce json
// @Param request body entity.CouponRedeemRequest true "Coupon payload"
// @Success 200 {object} entity.WalletSummary
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /coupons/redeem [post]
func (r *Routes) redeemCoupon(ctx *fiber.Ctx) error {
	var payload entity.CouponRedeemRequest
	if err := ctx.BodyParser(&payload); err != nil {
		r.l.Error(err, "http - v1 - redeemCoupon")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	if err := r.v.Struct(payload); err != nil {
		r.l.Error(err, "http - v1 - redeemCoupon - validation")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	userID, err := r.getUserID(ctx)
	if err != nil {
		r.l.Error(err, "http - v1 - redeemCoupon - user")
		return errorResponse(ctx, http.StatusUnauthorized, "missing user")
	}

	summary, err := r.uc.Coupon.Redeem(ctx.UserContext(), userID, payload)
	if err != nil {
		r.l.Error(err, "http - v1 - redeemCoupon - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to redeem")
	}

	return ctx.Status(http.StatusOK).JSON(summary)
}

// @Summary Referral summary
// @Tags App: Referral
// @Security UserAuth
// @Produce json
// @Success 200 {object} entity.ReferralSummary
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /referral [get]
func (r *Routes) referralSummary(ctx *fiber.Ctx) error {
	userID, err := r.getUserID(ctx)
	if err != nil {
		r.l.Error(err, "http - v1 - referralSummary")
		return errorResponse(ctx, http.StatusUnauthorized, "missing user")
	}

	summary, err := r.uc.Referral.Summary(ctx.UserContext(), userID)
	if err != nil {
		r.l.Error(err, "http - v1 - referralSummary - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to load referral")
	}

	return ctx.Status(http.StatusOK).JSON(summary)
}
