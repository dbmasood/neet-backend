package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	adminusecase "github.com/evrone/go-clean-template/internal/usecase/admin"
	"github.com/gofiber/fiber/v2"
)

func registerAdminAnalyticsRoutes(api fiber.Router, r *Routes) {
	api.Get("/overview", r.adminAnalyticsOverview)
	api.Get("/time-series", r.adminAnalyticsTimeSeries)
	api.Get("/subject-accuracy", r.adminAnalyticsSubjectAccuracy)
	api.Get("/weak-topics", r.adminAnalyticsWeakTopics)
}

// @Summary Analytics overview
// @Tags Admin: Analytics
// @Security AdminAuth
// @Produce json
// @Param exam query string false "Exam"
// @Param range query string false "today,7d,30d"
// @Success 200 {object} entity.AnalyticsOverview
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
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

// @Summary Analytics time series
// @Tags Admin: Analytics
// @Security AdminAuth
// @Produce json
// @Param metric query string true "Metric (active_users|questions_answered)"
// @Param exam query string false "Exam"
// @Param range query string false "today|7d|30d"
// @Success 200 {object} entity.AnalyticsTimeSeries
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/analytics/time-series [get]
func (r *Routes) adminAnalyticsTimeSeries(ctx *fiber.Ctx) error {
	metric := ctx.Query("metric")
	if metric == "" {
		return errorResponse(ctx, http.StatusBadRequest, "metric is required")
	}

	var exam *entity.ExamCategory
	if val := ctx.Query("exam"); val != "" {
		value := entity.ExamCategory(val)
		exam = &value
	}

	series, err := r.uc.Admin.TimeSeries(ctx.UserContext(), metric, exam, ctx.Query("range"))
	if err != nil {
		switch {
		case errors.Is(err, adminusecase.ErrInvalidMetric), errors.Is(err, adminusecase.ErrInvalidRange):
			return errorResponse(ctx, http.StatusBadRequest, err.Error())
		default:
			r.l.Error(err, "http - v1 - adminAnalyticsTimeSeries - usecase")
			return errorResponse(ctx, http.StatusInternalServerError, "unable to load time series")
		}
	}

	return ctx.Status(http.StatusOK).JSON(series)
}

// @Summary Subject accuracy
// @Tags Admin: Analytics
// @Security AdminAuth
// @Produce json
// @Param exam query string false "Exam"
// @Success 200 {object} entity.SubjectAccuracyResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/analytics/subject-accuracy [get]
func (r *Routes) adminAnalyticsSubjectAccuracy(ctx *fiber.Ctx) error {
	var exam *entity.ExamCategory
	if val := ctx.Query("exam"); val != "" {
		value := entity.ExamCategory(val)
		exam = &value
	}

	resp, err := r.uc.Admin.SubjectAccuracy(ctx.UserContext(), exam)
	if err != nil {
		r.l.Error(err, "http - v1 - adminAnalyticsSubjectAccuracy - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to load subject accuracy")
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}

// @Summary Weak topics
// @Tags Admin: Analytics
// @Security AdminAuth
// @Produce json
// @Param exam query string false "Exam"
// @Param limit query int false "Limit"
// @Success 200 {object} entity.WeakTopicsResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/analytics/weak-topics [get]
func (r *Routes) adminAnalyticsWeakTopics(ctx *fiber.Ctx) error {
	var exam *entity.ExamCategory
	if val := ctx.Query("exam"); val != "" {
		value := entity.ExamCategory(val)
		exam = &value
	}
	limit := 5
	if raw := ctx.Query("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	resp, err := r.uc.Admin.WeakTopics(ctx.UserContext(), exam, limit)
	if err != nil {
		r.l.Error(err, "http - v1 - adminAnalyticsWeakTopics - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to load weak topics")
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}
