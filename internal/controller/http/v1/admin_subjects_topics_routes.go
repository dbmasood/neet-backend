package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/controller/http/v1/request"
	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var (
	_ = entity.Subject{}
	_ = entity.Topic{}
)

func registerAdminSubjectsTopicsRoutes(api fiber.Router, r *Routes) {
	subjects := api.Group("/subjects")
	subjects.Get("", r.adminListSubjects)
	subjects.Post("", r.adminCreateSubject)

	topics := api.Group("/topics")
	topics.Get("", r.adminListTopics)
	topics.Post("", r.adminCreateTopic)
}

// @Summary List subjects
// @Tags Admin: Subjects
// @Security AdminAuth
// @Produce json
// @Success 200 {array} entity.Subject
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/subjects [get]
func (r *Routes) adminListSubjects(ctx *fiber.Ctx) error {
	subjects, err := r.uc.User.ListSubjects(ctx.UserContext(), nil)
	if err != nil {
		r.l.Error(err, "http - v1 - adminListSubjects")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to list subjects")
	}

	return ctx.Status(http.StatusOK).JSON(subjects)
}

// @Summary Create subject
// @Tags Admin: Subjects
// @Security AdminAuth
// @Accept json
// @Produce json
// @Param request body request.AdminSubjectCreateRequest true "Subject payload"
// @Success 201 {object} entity.Subject
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/subjects [post]
func (r *Routes) adminCreateSubject(ctx *fiber.Ctx) error {
	var payload request.AdminSubjectCreateRequest
	if err := ctx.BodyParser(&payload); err != nil {
		r.l.Error(err, "http - v1 - adminCreateSubject - parse")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	if err := r.v.Struct(payload); err != nil {
		r.l.Error(err, "http - v1 - adminCreateSubject - validation")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	subject, err := r.uc.User.CreateSubject(ctx.UserContext(), payload.Exam, payload.Name)
	if err != nil {
		r.l.Error(err, "http - v1 - adminCreateSubject - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to create subject")
	}

	return ctx.Status(http.StatusCreated).JSON(subject)
}

// @Summary List topics
// @Tags Admin: Topics
// @Security AdminAuth
// @Produce json
// @Param subjectId query string false "Subject ID"
// @Success 200 {array} entity.Topic
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/topics [get]
func (r *Routes) adminListTopics(ctx *fiber.Ctx) error {
	subjectID, err := parseQueryUUID(ctx, "subjectId")
	if err != nil {
		r.l.Error(err, "http - v1 - adminListTopics")
		return errorResponse(ctx, http.StatusBadRequest, "invalid subjectId")
	}

	var id uuid.UUID
	if subjectID != nil {
		id = *subjectID
	}

	topics, err := r.uc.User.ListTopics(ctx.UserContext(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - adminListTopics - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to list topics")
	}

	return ctx.Status(http.StatusOK).JSON(topics)
}

// @Summary Create topic
// @Tags Admin: Topics
// @Security AdminAuth
// @Accept json
// @Produce json
// @Param request body request.AdminTopicCreateRequest true "Topic payload"
// @Success 201 {object} entity.Topic
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/topics [post]
func (r *Routes) adminCreateTopic(ctx *fiber.Ctx) error {
	var payload request.AdminTopicCreateRequest
	if err := ctx.BodyParser(&payload); err != nil {
		r.l.Error(err, "http - v1 - adminCreateTopic - parse")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	if err := r.v.Struct(payload); err != nil {
		r.l.Error(err, "http - v1 - adminCreateTopic - validation")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	topic, err := r.uc.User.CreateTopic(ctx.UserContext(), payload.SubjectID, payload.Name)
	if err != nil {
		r.l.Error(err, "http - v1 - adminCreateTopic - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to create topic")
	}

	return ctx.Status(http.StatusCreated).JSON(topic)
}
