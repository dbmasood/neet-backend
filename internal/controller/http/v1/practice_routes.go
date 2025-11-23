package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

func registerPracticeRoutes(api fiber.Router, r *Routes) {
	api.Post("/sessions", r.createPracticeSession)
	api.Get("/sessions", r.listPracticeSessions)
	api.Get("/sessions/:id", r.getPracticeSession)
	api.Post("/sessions/:id/answers", r.answerPracticeQuestion)
}

func registerPracticeRevisionRoutes(api fiber.Router, r *Routes) {
	api.Get("/queue", r.revisionQueue)
}

// @Summary Create new practice session
// @Tags App: Practice
// @Security UserAuth
// @Accept json
// @Produce json
// @Param request body entity.PracticeSessionCreateRequest true "Session payload"
// @Success 201 {object} entity.PracticeSession
// @Failure 400 {object} response.Error
// @Failure 401 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /practice/sessions [post]
func (r *Routes) createPracticeSession(ctx *fiber.Ctx) error {
	var payload entity.PracticeSessionCreateRequest

	if err := ctx.BodyParser(&payload); err != nil {
		r.l.Error(err, "http - v1 - createPracticeSession")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	if err := r.v.Struct(payload); err != nil {
		r.l.Error(err, "http - v1 - createPracticeSession - validation")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	userID, err := r.getUserID(ctx)
	if err != nil {
		r.l.Error(err, "http - v1 - createPracticeSession - user")
		return errorResponse(ctx, http.StatusUnauthorized, "missing user")
	}

	session, err := r.uc.Practice.CreateSession(ctx.UserContext(), userID, payload)
	if err != nil {
		r.l.Error(err, "http - v1 - createPracticeSession - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to create session")
	}

	return ctx.Status(http.StatusCreated).JSON(session)
}

// @Summary List practice sessions
// @Tags App: Practice
// @Security UserAuth
// @Produce json
// @Success 200 {array} entity.PracticeSession
// @Failure 401 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /practice/sessions [get]
func (r *Routes) listPracticeSessions(ctx *fiber.Ctx) error {
	userID, err := r.getUserID(ctx)
	if err != nil {
		r.l.Error(err, "http - v1 - listPracticeSessions")
		return errorResponse(ctx, http.StatusUnauthorized, "missing user")
	}

	sessions, err := r.uc.Practice.ListSessions(ctx.UserContext(), userID)
	if err != nil {
		r.l.Error(err, "http - v1 - listPracticeSessions - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to fetch sessions")
	}

	return ctx.Status(http.StatusOK).JSON(sessions)
}

// @Summary Get practice session detail
// @Tags App: Practice
// @Security UserAuth
// @Produce json
// @Param id path string true "Session ID"
// @Success 200 {object} entity.PracticeSessionDetail
// @Failure 400 {object} response.Error
// @Failure 401 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /practice/sessions/{id} [get]
func (r *Routes) getPracticeSession(ctx *fiber.Ctx) error {
	sessionID, err := parseUUID(ctx, "id")
	if err != nil {
		r.l.Error(err, "http - v1 - getPracticeSession")
		return errorResponse(ctx, http.StatusBadRequest, "invalid session id")
	}

	userID, err := r.getUserID(ctx)
	if err != nil {
		r.l.Error(err, "http - v1 - getPracticeSession - user")
		return errorResponse(ctx, http.StatusUnauthorized, "missing user")
	}

	detail, err := r.uc.Practice.GetSessionDetail(ctx.UserContext(), sessionID, userID)
	if err != nil {
		r.l.Error(err, "http - v1 - getPracticeSession - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to load session")
	}

	return ctx.Status(http.StatusOK).JSON(detail)
}

// @Summary Submit practice answer
// @Tags App: Practice
// @Security UserAuth
// @Accept json
// @Produce json
// @Param id path string true "Session ID"
// @Param request body entity.PracticeAnswerRequest true "Answer payload"
// @Success 200 {object} entity.PracticeSessionQuestion
// @Failure 400 {object} response.Error
// @Failure 401 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /practice/sessions/{id}/answers [post]
func (r *Routes) answerPracticeQuestion(ctx *fiber.Ctx) error {
	sessionID, err := parseUUID(ctx, "id")
	if err != nil {
		r.l.Error(err, "http - v1 - answerPracticeQuestion")
		return errorResponse(ctx, http.StatusBadRequest, "invalid session id")
	}

	var payload entity.PracticeAnswerRequest
	if err := ctx.BodyParser(&payload); err != nil {
		r.l.Error(err, "http - v1 - answerPracticeQuestion - parse")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	if err := r.v.Struct(payload); err != nil {
		r.l.Error(err, "http - v1 - answerPracticeQuestion - validation")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	userID, err := r.getUserID(ctx)
	if err != nil {
		r.l.Error(err, "http - v1 - answerPracticeQuestion - user")
		return errorResponse(ctx, http.StatusUnauthorized, "missing user")
	}

	question, err := r.uc.Practice.AnswerQuestion(ctx.UserContext(), sessionID, payload, userID)
	if err != nil {
		r.l.Error(err, "http - v1 - answerPracticeQuestion - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to record answer")
	}

	return ctx.Status(http.StatusOK).JSON(question)
}

// @Summary Get revision queue
// @Tags App: Revision
// @Security UserAuth
// @Produce json
// @Success 200 {array} entity.RevisionItem
// @Failure 401 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /revision/queue [get]
func (r *Routes) revisionQueue(ctx *fiber.Ctx) error {
	userID, err := r.getUserID(ctx)
	if err != nil {
		r.l.Error(err, "http - v1 - revisionQueue")
		return errorResponse(ctx, http.StatusUnauthorized, "missing user")
	}

	items, err := r.uc.Revision.GetQueue(ctx.UserContext(), userID)
	if err != nil {
		r.l.Error(err, "http - v1 - revisionQueue - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to load revision queue")
	}

	return ctx.Status(http.StatusOK).JSON(items)
}
