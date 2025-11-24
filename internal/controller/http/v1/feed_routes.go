package v1

import (
	"net/http"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/gofiber/fiber/v2"
)

var _ = entity.FeedPost{}

func registerFeedRoutes(api fiber.Router, r *Routes) {
	api.Get("/feed", r.feed)
}

// @Summary Feed posts
// @Tags App: Feed
// @Security UserAuth
// @Produce json
// @Success 200 {array} entity.FeedPost
// @Failure 500 {object} ErrorResponse
// @Router /feed [get]
func (r *Routes) feed(ctx *fiber.Ctx) error {
	posts, err := r.uc.Feed.List(ctx.UserContext())
	if err != nil {
		r.l.Error(err, "http - v1 - feed")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to load feed")
	}

	return ctx.Status(http.StatusOK).JSON(posts)
}
