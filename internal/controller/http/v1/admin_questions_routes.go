package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/gofiber/fiber/v2"
)

func registerAdminQuestionRoutes(api fiber.Router, r *Routes) {
	api.Get("", r.adminListQuestions)
	api.Post("", r.adminCreateQuestion)
	api.Get("/:id", r.adminGetQuestion)
	api.Patch("/:id", r.adminUpdateQuestion)
	api.Delete("/:id", r.adminDeleteQuestion)
}

// @Summary List questions
// @Tags Admin: Questions
// @Security AdminAuth
// @Produce json
// @Param exam query string false "Exam"
// @Param subjectId query string false "Subject ID"
// @Param topicId query string false "Topic ID"
// @Success 200 {array} entity.Question
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/questions [get]
func (r *Routes) adminListQuestions(ctx *fiber.Ctx) error {
	filter := repo.QuestionFilter{}

	if query := ctx.Query("exam"); query != "" {
		exam := entity.ExamCategory(query)
		filter.Exam = &exam
	}

	if value, err := parseQueryUUID(ctx, "subjectId"); err != nil {
		r.l.Error(err, "http - v1 - adminListQuestions - subject")
		return errorResponse(ctx, http.StatusBadRequest, "invalid subjectId")
	} else {
		filter.SubjectID = value
	}

	if value, err := parseQueryUUID(ctx, "topicId"); err != nil {
		r.l.Error(err, "http - v1 - adminListQuestions - topic")
		return errorResponse(ctx, http.StatusBadRequest, "invalid topicId")
	} else {
		filter.TopicID = value
	}

	list, err := r.uc.Question.AdminList(ctx.UserContext(), filter)
	if err != nil {
		r.l.Error(err, "http - v1 - adminListQuestions - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to list questions")
	}

	return ctx.Status(http.StatusOK).JSON(list)
}

// @Summary Create question
// @Tags Admin: Questions
// @Security AdminAuth
// @Accept json
// @Produce json
// @Param request body entity.QuestionCreateRequest true "Question payload"
// @Success 201 {object} entity.Question
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/questions [post]
func (r *Routes) adminCreateQuestion(ctx *fiber.Ctx) error {
	var payload entity.QuestionCreateRequest
	if err := ctx.BodyParser(&payload); err != nil {
		r.l.Error(err, "http - v1 - adminCreateQuestion - parse")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	if err := r.v.Struct(payload); err != nil {
		r.l.Error(err, "http - v1 - adminCreateQuestion - validation")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	created, err := r.uc.Question.AdminCreate(ctx.UserContext(), payload)
	if err != nil {
		r.l.Error(err, "http - v1 - adminCreateQuestion - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to create question")
	}

	return ctx.Status(http.StatusCreated).JSON(created)
}

// @Summary Get question
// @Tags Admin: Questions
// @Security AdminAuth
// @Produce json
// @Param id path string true "Question ID"
// @Success 200 {object} entity.Question
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/questions/{id} [get]
func (r *Routes) adminGetQuestion(ctx *fiber.Ctx) error {
	id, err := parseUUID(ctx, "id")
	if err != nil {
		r.l.Error(err, "http - v1 - adminGetQuestion")
		return errorResponse(ctx, http.StatusBadRequest, "invalid id")
	}

	question, err := r.uc.Question.AdminGet(ctx.UserContext(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - adminGetQuestion - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to find question")
	}

	return ctx.Status(http.StatusOK).JSON(question)
}

// @Summary Update question
// @Tags Admin: Questions
// @Security AdminAuth
// @Accept json
// @Produce json
// @Param id path string true "Question ID"
// @Param request body entity.QuestionUpdateRequest true "Update payload"
// @Success 200 {object} entity.Question
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/questions/{id} [patch]
func (r *Routes) adminUpdateQuestion(ctx *fiber.Ctx) error {
	id, err := parseUUID(ctx, "id")
	if err != nil {
		r.l.Error(err, "http - v1 - adminUpdateQuestion")
		return errorResponse(ctx, http.StatusBadRequest, "invalid id")
	}

	var payload entity.QuestionUpdateRequest
	if err := ctx.BodyParser(&payload); err != nil {
		r.l.Error(err, "http - v1 - adminUpdateQuestion - parse")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	updated, err := r.uc.Question.AdminUpdate(ctx.UserContext(), id, payload)
	if err != nil {
		r.l.Error(err, "http - v1 - adminUpdateQuestion - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to update question")
	}

	return ctx.Status(http.StatusOK).JSON(updated)
}

// @Summary Delete question
// @Tags Admin: Questions
// @Security AdminAuth
// @Param id path string true "Question ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /admin/questions/{id} [delete]
func (r *Routes) adminDeleteQuestion(ctx *fiber.Ctx) error {
	id, err := parseUUID(ctx, "id")
	if err != nil {
		r.l.Error(err, "http - v1 - adminDeleteQuestion")
		return errorResponse(ctx, http.StatusBadRequest, "invalid id")
	}

	if err := r.uc.Question.AdminDelete(ctx.UserContext(), id); err != nil {
		r.l.Error(err, "http - v1 - adminDeleteQuestion - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to delete question")
	}

	return ctx.SendStatus(http.StatusNoContent)
}
