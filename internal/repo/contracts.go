package repo

import (
	"context"

	"github.com/google/uuid"

	"github.com/evrone/go-clean-template/internal/entity"
)

type (
    // TranslationRepo -.
    TranslationRepo interface {
        Store(context.Context, entity.Translation) error
        GetHistory(context.Context) ([]entity.Translation, error)
    }

	// TranslationWebAPI -.
	TranslationWebAPI interface {
		Translate(entity.Translation) (entity.Translation, error)
	}
)

// QuestionFilter carries optional filters.
type QuestionFilter struct {
	Exam      *entity.ExamCategory
	SubjectID *uuid.UUID
	TopicID   *uuid.UUID
}

// PodcastFilter describes query args.
type PodcastFilter struct {
	SubjectID *uuid.UUID
	TopicID   *uuid.UUID
}

// AnalyticsFilter describes dashboard query.
type AnalyticsFilter struct {
	Exam  *entity.ExamCategory
	Range string
}

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

type (
	UserRepository interface {
		GetByID(ctx context.Context, userID uuid.UUID) (entity.User, error)
		GetByTelegramID(ctx context.Context, telegramID string) (entity.User, error)
		Create(ctx context.Context, user entity.User) (entity.User, error)
		GetExamProfile(ctx context.Context, userID uuid.UUID) (entity.ExamProfile, error)
	}

	SubjectRepository interface {
		ListByExam(ctx context.Context, exam *entity.ExamCategory) ([]entity.Subject, error)
		Create(ctx context.Context, subject entity.Subject) (entity.Subject, error)
	}

	TopicRepository interface {
		ListBySubject(ctx context.Context, subjectID uuid.UUID) ([]entity.Topic, error)
		Create(ctx context.Context, topic entity.Topic) (entity.Topic, error)
	}

	QuestionRepository interface {
		List(ctx context.Context, filter QuestionFilter) ([]entity.Question, error)
		Create(ctx context.Context, question entity.Question) (entity.Question, error)
		GetByID(ctx context.Context, id uuid.UUID) (entity.Question, error)
		Update(ctx context.Context, question entity.Question) (entity.Question, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	PracticeSessionRepository interface {
		CreateSession(ctx context.Context, session entity.PracticeSession) (entity.PracticeSession, error)
		ListSessions(ctx context.Context, userID uuid.UUID) ([]entity.PracticeSession, error)
		GetSession(ctx context.Context, id uuid.UUID) (entity.PracticeSession, error)
		ListSessionQuestions(ctx context.Context, sessionID uuid.UUID) ([]entity.PracticeSessionQuestion, error)
		GetSessionQuestion(ctx context.Context, id uuid.UUID) (entity.PracticeSessionQuestion, error)
		UpdateSessionQuestion(ctx context.Context, question entity.PracticeSessionQuestion) (entity.PracticeSessionQuestion, error)
	}

	RevisionRepository interface {
		ListDue(ctx context.Context, userID uuid.UUID) ([]entity.RevisionItem, error)
	}

	ExamRepository interface {
		ListConfigs(ctx context.Context) ([]entity.ExamConfig, error)
		CreateConfig(ctx context.Context, config entity.ExamConfig) (entity.ExamConfig, error)
		GetConfig(ctx context.Context, id uuid.UUID) (entity.ExamConfig, error)
		UpdateConfig(ctx context.Context, config entity.ExamConfig) (entity.ExamConfig, error)
		DeleteConfig(ctx context.Context, id uuid.UUID) error
		ListSummaries(ctx context.Context) ([]entity.ExamSummary, error)
	}

	PodcastRepository interface {
		List(ctx context.Context, filter PodcastFilter) ([]entity.PodcastEpisode, error)
		Get(ctx context.Context, id uuid.UUID) (entity.PodcastEpisode, error)
		Create(ctx context.Context, episode entity.PodcastEpisode) (entity.PodcastEpisode, error)
		Update(ctx context.Context, episode entity.PodcastEpisode) (entity.PodcastEpisode, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	WalletRepository interface {
		GetSummary(ctx context.Context, userID uuid.UUID) (entity.WalletSummary, error)
		ListTransactions(ctx context.Context, userID uuid.UUID) ([]entity.WalletTransaction, error)
	}

	CouponRepository interface {
		List(ctx context.Context) ([]entity.Coupon, error)
		Create(ctx context.Context, coupon entity.Coupon) (entity.Coupon, error)
		Get(ctx context.Context, id uuid.UUID) (entity.Coupon, error)
		Update(ctx context.Context, coupon entity.Coupon) (entity.Coupon, error)
		Delete(ctx context.Context, id uuid.UUID) error
		Redeem(ctx context.Context, code string, userID uuid.UUID) (entity.WalletSummary, error)
	}

	ReferralRepository interface {
		GetSummary(ctx context.Context, userID uuid.UUID) (entity.ReferralSummary, error)
	}

    AISettingsRepository interface {
        Get(ctx context.Context) (entity.AISettings, error)
        Update(ctx context.Context, settings entity.AISettings) (entity.AISettings, error)
    }

    AnalyticsRepository interface {
        Overview(ctx context.Context, filter AnalyticsFilter) (entity.AnalyticsOverview, error)
    }
 
    LeaderboardRepository interface {
        List(ctx context.Context, limit int) ([]entity.LeaderboardEntry, error)
        Stats(ctx context.Context) (entity.LeaderboardStats, error)
    }

    FeedRepository interface {
        List(ctx context.Context) ([]entity.FeedPost, error)
    }
)
