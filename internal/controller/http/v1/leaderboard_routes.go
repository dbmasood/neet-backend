package v1

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func registerLeaderboardRoutes(api fiber.Router, r *Routes) {
	api.Get("/leaderboard", r.leaderboard)
}

// @Summary Leaderboard
// @Tags App: Leaderboard
// @Security UserAuth
// @Produce json
// @Param limit query int false "Max entries" default(20)
// @Success 200 {object} entity.LeaderboardStats
// @Success 200 {array} entity.LeaderboardEntry
// @Failure 500 {object} response.Error
// @Router /leaderboard [get]
func (r *Routes) leaderboard(ctx *fiber.Ctx) error {
	limit := 20
	if q := ctx.Query("limit"); q != "" {
		if parsed, err := strconv.Atoi(q); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	entries, stats, err := r.uc.Leaderboard.Entries(ctx.UserContext(), limit)
	if err != nil {
		r.l.Error(err, "http - v1 - leaderboard")
		return errorResponse(ctx, http.StatusInternalServerError, "unable to load leaderboard")
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"stats":   stats,
		"entries": entries,
	})
}
