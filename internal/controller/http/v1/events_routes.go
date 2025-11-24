package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

var _ = entity.ExamSummary{}

func registerEventsRoutes(api fiber.Router, r *Routes) {
	api.Get("", r.listEvents)
}

// @Summary List exams/events
// @Tags App: Exams
// @Security UserAuth
// @Produce json
// @Success 200 {array} entity.ExamSummary
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /events [get]
func (r *Routes) listEvents(ctx *fiber.Ctx) error {
	events, err := r.uc.Exam.ListEvents(ctx.UserContext())
	if err != nil {
		r.l.Error(err, "http - v1 - listEvents")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to load events")
	}

	return ctx.Status(http.StatusOK).JSON(events)
}
