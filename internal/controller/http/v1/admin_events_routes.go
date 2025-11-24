package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

func registerAdminEventsRoutes(api fiber.Router, r *Routes) {
	api.Get("/upcoming", r.adminUpcomingEvents)
}

// @Summary Upcoming events
// @Tags Admin: Events
// @Security AdminAuth
// @Produce json
// @Param exam query string false "Exam"
// @Success 200 {object} entity.AdminEventsResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/events/upcoming [get]
func (r *Routes) adminUpcomingEvents(ctx *fiber.Ctx) error {
	var exam *entity.ExamCategory
	if value := ctx.Query("exam"); value != "" {
		examValue := entity.ExamCategory(value)
		exam = &examValue
	}

	resp, err := r.uc.Admin.UpcomingEvents(ctx.UserContext(), exam)
	if err != nil {
		r.l.Error(err, "http - v1 - adminUpcomingEvents - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to load events")
	}

	return ctx.Status(http.StatusOK).JSON(resp)
}
