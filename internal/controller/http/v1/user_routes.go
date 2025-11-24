package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

func registerUserRoutes(api fiber.Router, r *Routes) {
	api.Get("/me", r.me)
	api.Get("/subjects", r.subjects)
	api.Get("/topics", r.topics)
}

// @Summary Get current user and profile
// @Tags App: User
// @Security UserAuth
// @Produce json
// @Success 200 {object} entity.MeResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /me [get]
func (r *Routes) me(ctx *fiber.Ctx) error {
	userID, err := r.getUserID(ctx)
	if err != nil {
		r.l.Error(err, "http - v1 - me - missing user")
		return errorResponse(ctx, http.StatusUnauthorized, "missing user")
	}

	result, err := r.uc.User.Me(ctx.UserContext(), userID)
	if err != nil {
		r.l.Error(err, "http - v1 - me - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to load profile")
	}

	return ctx.Status(http.StatusOK).JSON(result)
}

// @Summary List subjects for current exam
// @Tags App: User
// @Security UserAuth
// @Produce json
// @Param exam query string false "Exam category"
// @Success 200 {array} entity.Subject
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /subjects [get]
func (r *Routes) subjects(ctx *fiber.Ctx) error {
	var exam *entity.ExamCategory
	if query := ctx.Query("exam"); query != "" {
		value := entity.ExamCategory(query)
		exam = &value
	}

	list, err := r.uc.User.ListSubjects(ctx.UserContext(), exam)
	if err != nil {
		r.l.Error(err, "http - v1 - subjects")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to list subjects")
	}

	return ctx.Status(http.StatusOK).JSON(list)
}

// @Summary List topics for subject
// @Tags App: User
// @Security UserAuth
// @Produce json
// @Param subjectId query string true "Subject ID"
// @Success 200 {array} entity.Topic
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /topics [get]
func (r *Routes) topics(ctx *fiber.Ctx) error {
	subjectID, err := parseQueryUUID(ctx, "subjectId")
	if err != nil {
		r.l.Error(err, "http - v1 - topics - parse")
		return errorResponse(ctx, http.StatusBadRequest, "invalid subjectId")
	}

	if subjectID == nil {
		return errorResponse(ctx, http.StatusBadRequest, "subjectId is required")
	}

	list, err := r.uc.User.ListTopics(ctx.UserContext(), *subjectID)
	if err != nil {
		r.l.Error(err, "http - v1 - topics - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to list topics")
	}

	return ctx.Status(http.StatusOK).JSON(list)
}
