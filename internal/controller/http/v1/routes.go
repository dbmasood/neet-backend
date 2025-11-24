package v1

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/evrone/go-clean-template/internal/controller/http/middleware"
	"github.com/evrone/go-clean-template/internal/usecase"
	"github.com/evrone/go-clean-template/pkg/jwt"
	"github.com/evrone/go-clean-template/pkg/logger"
)

// Routes groups service handlers.
type Routes struct {
	uc usecase.UseCases
	l  logger.Interface
	v  *validator.Validate
}

// RegisterRoutes registers domain routes.
func RegisterRoutes(api fiber.Router, uc usecase.UseCases, userJWT, adminJWT *jwt.Service, l logger.Interface) {
	r := &Routes{uc: uc, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	authGroup := api.Group("/auth")
	registerAuthRoutes(authGroup, r)

	adminAuthGroup := authGroup.Group("/admin")
	adminAuthGroup.Use(middleware.AdminAuth(adminJWT))
	registerAdminAuthRoutes(adminAuthGroup, r)

	userGroup := api.Group("/")
	userGroup.Use(middleware.UserAuth(userJWT))
	registerUserRoutes(userGroup, r)
	registerLeaderboardRoutes(userGroup, r)
	registerFeedRoutes(userGroup, r)

	practiceGroup := api.Group("/practice")
	practiceGroup.Use(middleware.UserAuth(userJWT))
	registerPracticeRoutes(practiceGroup, r)

	revisionGroup := api.Group("/revision")
	revisionGroup.Use(middleware.UserAuth(userJWT))
	registerPracticeRevisionRoutes(revisionGroup, r)

	podcastGroup := api.Group("/podcasts")
	podcastGroup.Use(middleware.UserAuth(userJWT))
	registerPodcastRoutes(podcastGroup, r)

	eventsGroup := api.Group("/events")
	eventsGroup.Use(middleware.UserAuth(userJWT))
	registerEventsRoutes(eventsGroup, r)

	registerWalletRoutes(api, r, userJWT)

	adminGroup := api.Group("/admin")
	adminGroup.Use(middleware.AdminAuth(adminJWT))

	registerAdminQuestionRoutes(adminGroup.Group("/questions"), r)
	registerAdminSubjectsTopicsRoutes(adminGroup, r)
	registerAdminExamsRoutes(adminGroup.Group("/exams"), r)
	registerAdminPodcastsRoutes(adminGroup.Group("/podcasts"), r)
	registerAdminCouponsRoutes(adminGroup.Group("/coupons"), r)
	registerAdminAISettingsRoutes(adminGroup.Group("/ai-settings"), r)
	registerAdminAnalyticsRoutes(adminGroup.Group("/analytics"), r)
	registerAdminEventsRoutes(adminGroup.Group("/events"), r)
	registerAdminReferralRoutes(adminGroup.Group("/referrals"), r)
	registerAdminUsersRoutes(adminGroup.Group("/users"), r)
}

func (r *Routes) getUserID(ctx *fiber.Ctx) (uuid.UUID, error) {
	raw := ctx.Locals("userID")
	if raw == nil {
		return uuid.Nil, fmt.Errorf("user id missing")
	}

	str, ok := raw.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid user id")
	}

	return uuid.Parse(str)
}

func parseUUID(ctx *fiber.Ctx, key string) (uuid.UUID, error) {
	value := ctx.Params(key)
	if value == "" {
		return uuid.Nil, fmt.Errorf("%s is required", key)
	}

	return uuid.Parse(value)
}

func parseQueryUUID(ctx *fiber.Ctx, key string) (*uuid.UUID, error) {
	value := ctx.Query(key)
	if value == "" {
		return nil, nil
	}

	parsed, err := uuid.Parse(value)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}
