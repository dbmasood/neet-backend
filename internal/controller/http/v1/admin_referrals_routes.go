package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

var _ = entity.AdminReferralSummary{}

func registerAdminReferralRoutes(api fiber.Router, r *Routes) {
	api.Get("/summary", r.adminReferralSummary)
}

// @Summary Referral summary
// @Tags Admin: Referrals
// @Security AdminAuth
// @Produce json
// @Param range query string false "today|7d|30d"
// @Success 200 {object} entity.AdminReferralSummary
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/referrals/summary [get]
func (r *Routes) adminReferralSummary(ctx *fiber.Ctx) error {
	resp, err := r.uc.Admin.ReferralSummary(ctx.UserContext(), ctx.Query("range"))
	if err != nil {
		r.l.Error(err, "http - v1 - adminReferralSummary - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to load referral summary")
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}
