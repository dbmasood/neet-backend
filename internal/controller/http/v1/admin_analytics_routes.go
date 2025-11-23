package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/gofiber/fiber/v2"
)

func registerAdminAnalyticsRoutes(api fiber.Router, r *Routes) {
	api.Get("/overview", r.adminAnalyticsOverview)
}

// @Summary Analytics overview
// @Tags Admin: Analytics
// @Security AdminAuth
// @Produce json
// @Param exam query string false "Exam"
// @Param range query string false "today,7d,30d"
// @Success 200 {object} entity.AnalyticsOverview
// @Failure 401 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /admin/analytics/overview [get]
func (r *Routes) adminAnalyticsOverview(ctx *fiber.Ctx) error {
	filter := repo.AnalyticsFilter{
		Range: ctx.Query("range"),
	}

	if exam := ctx.Query("exam"); exam != "" {
		value := entity.ExamCategory(exam)
		filter.Exam = &value
	}

	overview, err := r.uc.Analytics.Overview(ctx.UserContext(), filter)
	if err != nil {
		r.l.Error(err, "http - v1 - adminAnalyticsOverview")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to load analytics")
	}

	return ctx.Status(http.StatusOK).JSON(overview)
}
