package usecase

import (
	"github.com/evrone/go-clean-template/internal/usecase/admin"
	"github.com/evrone/go-clean-template/internal/usecase/ai"
	"github.com/evrone/go-clean-template/internal/usecase/analytics"
	"github.com/evrone/go-clean-template/internal/usecase/auth"
	"github.com/evrone/go-clean-template/internal/usecase/coupon"
	"github.com/evrone/go-clean-template/internal/usecase/exam"
	"github.com/evrone/go-clean-template/internal/usecase/feed"
	"github.com/evrone/go-clean-template/internal/usecase/leaderboard"
	"github.com/evrone/go-clean-template/internal/usecase/podcast"
	"github.com/evrone/go-clean-template/internal/usecase/practice"
	"github.com/evrone/go-clean-template/internal/usecase/question"
	"github.com/evrone/go-clean-template/internal/usecase/referral"
	"github.com/evrone/go-clean-template/internal/usecase/revision"
	"github.com/evrone/go-clean-template/internal/usecase/user"
	"github.com/evrone/go-clean-template/internal/usecase/wallet"
)

// UseCases groups all domain usecases.
type UseCases struct {
	Admin       *admin.UseCase
	Auth        *auth.UseCase
	User        *user.UseCase
	Practice    *practice.UseCase
	Revision    *revision.UseCase
	Question    *question.UseCase
	Exam        *exam.UseCase
	Podcast     *podcast.UseCase
	Wallet      *wallet.UseCase
	Coupon      *coupon.UseCase
	Referral    *referral.UseCase
	AI          *ai.UseCase
	Analytics   *analytics.UseCase
	Leaderboard *leaderboard.UseCase
	Feed        *feed.UseCase
}
