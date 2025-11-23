package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/gofiber/fiber/v2"
)

func registerPodcastRoutes(api fiber.Router, r *Routes) {
	api.Get("", r.listPodcasts)
	api.Get("/:id", r.getPodcast)
}

// @Summary List podcast episodes
// @Tags App: Podcasts
// @Security UserAuth
// @Produce json
// @Param subjectId query string false "Subject ID"
// @Param topicId query string false "Topic ID"
// @Success 200 {array} entity.PodcastEpisode
// @Failure 401 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /podcasts [get]
func (r *Routes) listPodcasts(ctx *fiber.Ctx) error {
	filter := repo.PodcastFilter{}
	if s, err := parseQueryUUID(ctx, "subjectId"); err != nil {
		r.l.Error(err, "http - v1 - listPodcasts - subject")
		return errorResponse(ctx, http.StatusBadRequest, "invalid subjectId")
	} else {
		filter.SubjectID = s
	}

	if t, err := parseQueryUUID(ctx, "topicId"); err != nil {
		r.l.Error(err, "http - v1 - listPodcasts - topic")
		return errorResponse(ctx, http.StatusBadRequest, "invalid topicId")
	} else {
		filter.TopicID = t
	}

	episodes, err := r.uc.Podcast.List(ctx.UserContext(), filter)
	if err != nil {
		r.l.Error(err, "http - v1 - listPodcasts - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to load podcasts")
	}

	return ctx.Status(http.StatusOK).JSON(episodes)
}

// @Summary Get podcast episode
// @Tags App: Podcasts
// @Security UserAuth
// @Produce json
// @Param id path string true "Episode ID"
// @Success 200 {object} entity.PodcastEpisode
// @Failure 400 {object} response.Error
// @Failure 401 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /podcasts/{id} [get]
func (r *Routes) getPodcast(ctx *fiber.Ctx) error {
	id, err := parseUUID(ctx, "id")
	if err != nil {
		r.l.Error(err, "http - v1 - getPodcast")
		return errorResponse(ctx, http.StatusBadRequest, "invalid id")
	}

	episode, err := r.uc.Podcast.Get(ctx.UserContext(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - getPodcast - usecase")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to load episode")
	}

	return ctx.Status(http.StatusOK).JSON(episode)
}
