package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/gofiber/fiber/v2"
)

func registerAdminPodcastsRoutes(api fiber.Router, r *Routes) {
	api.Get("", r.adminListPodcasts)
	api.Post("", r.adminCreatePodcast)
	api.Get("/:id", r.adminGetPodcast)
	api.Patch("/:id", r.adminUpdatePodcast)
	api.Delete("/:id", r.adminDeletePodcast)
}

// @Summary List podcasts
// @Tags Admin: Podcasts
// @Security AdminAuth
// @Produce json
// @Success 200 {array} entity.PodcastEpisode
// @Failure 401 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /admin/podcasts [get]
func (r *Routes) adminListPodcasts(ctx *fiber.Ctx) error {
	episodes, err := r.uc.Podcast.List(ctx.UserContext(), repo.PodcastFilter{})
	if err != nil {
		r.l.Error(err, "http - v1 - adminListPodcasts")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to list podcasts")
	}

	return ctx.Status(http.StatusOK).JSON(episodes)
}

// @Summary Create podcast episode
// @Tags Admin: Podcasts
// @Security AdminAuth
// @Accept json
// @Produce json
// @Param request body entity.PodcastCreateRequest true "Podcast payload"
// @Success 201 {object} entity.PodcastEpisode
// @Failure 400 {object} response.Error
// @Failure 401 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /admin/podcasts [post]
func (r *Routes) adminCreatePodcast(ctx *fiber.Ctx) error {
	var payload entity.PodcastCreateRequest
	if err := ctx.BodyParser(&payload); err != nil {
		r.l.Error(err, "http - v1 - adminCreatePodcast - parse")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	if err := r.v.Struct(payload); err != nil {
		r.l.Error(err, "http - v1 - adminCreatePodcast - validation")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	episode, err := r.uc.Podcast.AdminCreate(ctx.UserContext(), payload)
	if err != nil {
		r.l.Error(err, "http - v1 - adminCreatePodcast - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to create podcast")
	}

	return ctx.Status(http.StatusCreated).JSON(episode)
}

// @Summary Get podcast episode
// @Tags Admin: Podcasts
// @Security AdminAuth
// @Produce json
// @Param id path string true "Episode ID"
// @Success 200 {object} entity.PodcastEpisode
// @Failure 400 {object} response.Error
// @Failure 401 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /admin/podcasts/{id} [get]
func (r *Routes) adminGetPodcast(ctx *fiber.Ctx) error {
	id, err := parseUUID(ctx, "id")
	if err != nil {
		r.l.Error(err, "http - v1 - adminGetPodcast")
		return errorResponse(ctx, http.StatusBadRequest, "invalid id")
	}

	episode, err := r.uc.Podcast.AdminGet(ctx.UserContext(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - adminGetPodcast - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to load podcast")
	}

	return ctx.Status(http.StatusOK).JSON(episode)
}

// @Summary Update podcast episode
// @Tags Admin: Podcasts
// @Security AdminAuth
// @Accept json
// @Produce json
// @Param id path string true "Episode ID"
// @Param request body entity.PodcastCreateRequest true "Podcast payload"
// @Success 200 {object} entity.PodcastEpisode
// @Failure 400 {object} response.Error
// @Failure 401 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /admin/podcasts/{id} [patch]
func (r *Routes) adminUpdatePodcast(ctx *fiber.Ctx) error {
	id, err := parseUUID(ctx, "id")
	if err != nil {
		r.l.Error(err, "http - v1 - adminUpdatePodcast")
		return errorResponse(ctx, http.StatusBadRequest, "invalid id")
	}

	var payload entity.PodcastCreateRequest
	if err := ctx.BodyParser(&payload); err != nil {
		r.l.Error(err, "http - v1 - adminUpdatePodcast - parse")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	if err := r.v.Struct(payload); err != nil {
		r.l.Error(err, "http - v1 - adminUpdatePodcast - validation")
		return errorResponse(ctx, http.StatusBadRequest, "invalid body")
	}

	updated, err := r.uc.Podcast.AdminUpdate(ctx.UserContext(), id, payload)
	if err != nil {
		r.l.Error(err, "http - v1 - adminUpdatePodcast - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to update podcast")
	}

	return ctx.Status(http.StatusOK).JSON(updated)
}

// @Summary Delete podcast episode
// @Tags Admin: Podcasts
// @Security AdminAuth
// @Param id path string true "Episode ID"
// @Success 204
// @Failure 400 {object} response.Error
// @Failure 401 {object} response.Error
// @Router /admin/podcasts/{id} [delete]
func (r *Routes) adminDeletePodcast(ctx *fiber.Ctx) error {
	id, err := parseUUID(ctx, "id")
	if err != nil {
		r.l.Error(err, "http - v1 - adminDeletePodcast")
		return errorResponse(ctx, http.StatusBadRequest, "invalid id")
	}

	if err := r.uc.Podcast.AdminDelete(ctx.UserContext(), id); err != nil {
		r.l.Error(err, "http - v1 - adminDeletePodcast - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to delete podcast")
	}

	return ctx.SendStatus(http.StatusNoContent)
}
