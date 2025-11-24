package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

func registerAdminExamsRoutes(api fiber.Router, r *Routes) {
	api.Get("", r.adminListExams)
	api.Post("", r.adminCreateExam)
	api.Get("/:id", r.adminGetExam)
	api.Patch("/:id", r.adminUpdateExam)
	api.Delete("/:id", r.adminDeleteExam)
}

// @Summary List exam configs
// @Tags Admin: Exams
// @Security AdminAuth
// @Produce json
// @Param exam query string false "Exam"
// @Success 200 {array} entity.ExamConfig
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/exams [get]
func (r *Routes) adminListExams(ctx *fiber.Ctx) error {
	var exam *entity.ExamCategory
	if value := ctx.Query("exam"); value != "" {
		parsed := entity.ExamCategory(value)
		exam = &parsed
	}

	configs, err := r.uc.Exam.AdminList(ctx.UserContext(), exam)
	if err != nil {
		r.l.Error(err, "http - v1 - adminListExams")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to list exams")
	}

	return ctx.Status(http.StatusOK).JSON(configs)
}

// @Summary Create exam config
// @Tags Admin: Exams
// @Security AdminAuth
// @Accept json
// @Produce json
// @Param request body entity.ExamConfigCreateRequest true "Exam payload"
// @Success 201 {object} entity.ExamConfig
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/exams [post]
func (r *Routes) adminCreateExam(ctx *fiber.Ctx) error {
	var payload entity.ExamConfigCreateRequest
	if err := ctx.BodyParser(&payload); err != nil {
		r.l.Error(err, "http - v1 - adminCreateExam - parse")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	if err := r.v.Struct(payload); err != nil {
		r.l.Error(err, "http - v1 - adminCreateExam - validation")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	config, err := r.uc.Exam.AdminCreate(ctx.UserContext(), payload)
	if err != nil {
		r.l.Error(err, "http - v1 - adminCreateExam - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to create exam")
	}

	return ctx.Status(http.StatusCreated).JSON(config)
}

// @Summary Get exam config
// @Tags Admin: Exams
// @Security AdminAuth
// @Produce json
// @Param id path string true "Exam ID"
// @Success 200 {object} entity.ExamConfig
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/exams/{id} [get]
func (r *Routes) adminGetExam(ctx *fiber.Ctx) error {
	id, err := parseUUID(ctx, "id")
	if err != nil {
		r.l.Error(err, "http - v1 - adminGetExam")
		return errorResponse(ctx, http.StatusBadRequest, "invalid id")
	}

	config, err := r.uc.Exam.AdminGet(ctx.UserContext(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - adminGetExam - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to load exam")
	}

	return ctx.Status(http.StatusOK).JSON(config)
}

// @Summary Update exam config
// @Tags Admin: Exams
// @Security AdminAuth
// @Accept json
// @Produce json
// @Param id path string true "Exam ID"
// @Param request body entity.ExamConfigUpdateRequest true "Exam update"
// @Success 200 {object} entity.ExamConfig
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/exams/{id} [patch]
func (r *Routes) adminUpdateExam(ctx *fiber.Ctx) error {
	id, err := parseUUID(ctx, "id")
	if err != nil {
		r.l.Error(err, "http - v1 - adminUpdateExam")
		return errorResponse(ctx, http.StatusBadRequest, "invalid id")
	}

	var payload entity.ExamConfigUpdateRequest
	if err := ctx.BodyParser(&payload); err != nil {
		r.l.Error(err, "http - v1 - adminUpdateExam - parse")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	updated, err := r.uc.Exam.AdminUpdate(ctx.UserContext(), id, payload)
	if err != nil {
		r.l.Error(err, "http - v1 - adminUpdateExam - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to update exam")
	}

	return ctx.Status(http.StatusOK).JSON(updated)
}

// @Summary Delete exam config
// @Tags Admin: Exams
// @Security AdminAuth
// @Param id path string true "Exam ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /admin/exams/{id} [delete]
func (r *Routes) adminDeleteExam(ctx *fiber.Ctx) error {
	id, err := parseUUID(ctx, "id")
	if err != nil {
		r.l.Error(err, "http - v1 - adminDeleteExam")
		return errorResponse(ctx, http.StatusBadRequest, "invalid id")
	}

	if err := r.uc.Exam.AdminDelete(ctx.UserContext(), id); err != nil {
		r.l.Error(err, "http - v1 - adminDeleteExam - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to delete exam")
	}

	return ctx.SendStatus(http.StatusNoContent)
}
