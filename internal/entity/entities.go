package entity

import (
	"time"

	"github.com/google/uuid"
)

// ExamCategory defines supported exams.
type ExamCategory string

const (
	ExamCategoryNEETPG ExamCategory = "NEET_PG"
	ExamCategoryNEETUG ExamCategory = "NEET_UG"
	ExamCategoryJEE    ExamCategory = "JEE"
	ExamCategoryUPSC   ExamCategory = "UPSC"
)

// UserRole defines user roles.
type UserRole string

const (
	UserRoleUser       UserRole = "USER"
	UserRoleAdmin      UserRole = "ADMIN"
	UserRoleSuperAdmin UserRole = "SUPER_ADMIN"
)

// QuestionChoiceType determines single/multi choice.
type QuestionChoiceType string

const (
	ChoiceTypeSingle QuestionChoiceType = "single"
	ChoiceTypeMulti  QuestionChoiceType = "multi"
)

// PracticeMode describes session modes.
type PracticeMode string

const (
	PracticeModeSmart    PracticeMode = "smart"
	PracticeModeCustom   PracticeMode = "custom"
	PracticeModeRevision PracticeMode = "revision"
	PracticeModeExam     PracticeMode = "exam"
)

// PracticeSessionStatus defines session lifecycle.
type PracticeSessionStatus string

const (
	PracticeStatusInProgress PracticeSessionStatus = "in_progress"
	PracticeStatusCompleted  PracticeSessionStatus = "completed"
	PracticeStatusAbandoned  PracticeSessionStatus = "abandoned"
)

// ExamConfigType for different exam kinds.
type ExamConfigType string

const (
	ExamTypeMock        ExamConfigType = "MOCK"
	ExamTypeSubjectTest ExamConfigType = "SUBJECT_TEST"
	ExamTypeRewardEvent ExamConfigType = "REWARD_EVENT"
	ExamTypeDailyTest   ExamConfigType = "DAILY_TEST"
)

// ExamStatus represents scheduling state.
type ExamStatus string

const (
	ExamStatusDraft     ExamStatus = "DRAFT"
	ExamStatusScheduled ExamStatus = "SCHEDULED"
	ExamStatusOngoing   ExamStatus = "ONGOING"
	ExamStatusCompleted ExamStatus = "COMPLETED"
)

// WalletTxType constants.
type WalletTxType string

const (
	WalletTxReward     WalletTxType = "REWARD"
	WalletTxExamEntry  WalletTxType = "EXAM_ENTRY"
	WalletTxCoupon     WalletTxType = "COUPON"
	WalletTxAdjustment WalletTxType = "ADJUSTMENT"
	WalletTxReferral   WalletTxType = "REFERRAL"
	WalletTxSpin       WalletTxType = "SPIN"
	WalletTxBonus      WalletTxType = "BONUS"
)

// ReferralStatus defines states for referrals.
type ReferralStatus string

const (
	ReferralInvited   ReferralStatus = "INVITED"
	ReferralJoined    ReferralStatus = "JOINED"
	ReferralActivated ReferralStatus = "ACTIVATED"
)

// RevisionResult describes SRS outcome.
type RevisionResult string

const (
	RevisionResultCorrect   RevisionResult = "CORRECT"
	RevisionResultIncorrect RevisionResult = "INCORRECT"
)

// User represents a learner or admin.
type User struct {
	ID          uuid.UUID    `json:"id"`
	DisplayName string       `json:"displayName"`
	Email       *string      `json:"email,omitempty"`
	TelegramID  *string      `json:"telegramId,omitempty"`
	PrimaryExam ExamCategory `json:"primaryExam"`
	Role        UserRole     `json:"role"`
	CreatedAt   time.Time    `json:"createdAt"`
}

// ExamProfile tracks stats per exam.
type ExamProfile struct {
	ExamCategory      ExamCategory `json:"exam"`
	TotalQuestions    int          `json:"totalQuestions"`
	TotalCorrect      int          `json:"totalCorrect"`
	TotalTimeSeconds  int          `json:"totalTimeSeconds"`
	OverallLevel      int          `json:"overallLevel"`
	CurrentStreakDays int          `json:"currentStreakDays"`
	LongestStreakDays int          `json:"longestStreakDays"`
	LastLoginAt       *time.Time   `json:"lastLoginAt,omitempty"`
}

// MeResponse returned by /me.
type MeResponse struct {
	User        User        `json:"user"`
	ExamProfile ExamProfile `json:"examProfile"`
}

// TelegramAuthRequest payload for Telegram login.
type TelegramAuthRequest struct {
	TelegramID  string       `json:"telegramId" validate:"required"`
	DisplayName string       `json:"displayName"`
	Exam        ExamCategory `json:"exam" validate:"required"`
}

// AdminLoginRequest payload for admin login.
type AdminLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse returns token and user info.
type AuthResponse struct {
	AccessToken string `json:"accessToken"`
	User        User   `json:"user"`
}

// AdminProfile represents the bootstrap admin identity.
type AdminProfile struct {
	ID          uuid.UUID    `json:"id"`
	DisplayName string       `json:"displayName"`
	Email       string       `json:"email"`
	Role        string       `json:"role"`
	PrimaryExam ExamCategory `json:"primaryExam"`
	CreatedAt   time.Time    `json:"createdAt"`
	Permissions []string     `json:"permissions"`
}

// AnalyticsPoint represents a metric value for a date.
type AnalyticsPoint struct {
	Date  string `json:"date"`
	Value int    `json:"value"`
}

// AnalyticsTimeSeries response for dashboard charts.
type AnalyticsTimeSeries struct {
	Metric string           `json:"metric"`
	Exam   ExamCategory     `json:"exam"`
	Range  string           `json:"range"`
	Points []AnalyticsPoint `json:"points"`
}

// SubjectAccuracyItem describes subject level accuracy.
type SubjectAccuracyItem struct {
	SubjectID   string  `json:"subjectId"`
	SubjectName string  `json:"subjectName"`
	Accuracy    float64 `json:"accuracy"`
}

// SubjectAccuracyResponse wraps accuracy per subject.
type SubjectAccuracyResponse struct {
	Exam     ExamCategory          `json:"exam"`
	Subjects []SubjectAccuracyItem `json:"subjects"`
}

// WeakTopicItem describes a struggling topic.
type WeakTopicItem struct {
	SubjectID   string  `json:"subjectId"`
	SubjectName string  `json:"subjectName"`
	TopicID     string  `json:"topicId"`
	TopicName   string  `json:"topicName"`
	Accuracy    float64 `json:"accuracy"`
	Attempts    int     `json:"attempts"`
}

// WeakTopicsResponse wraps weakest topics.
type WeakTopicsResponse struct {
	Items []WeakTopicItem `json:"items"`
}

// AdminEventSummary describes admin view of upcoming events.
type AdminEventSummary struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Exam            ExamCategory   `json:"exam"`
	Type            ExamConfigType `json:"type"`
	StartAt         time.Time      `json:"startAt"`
	RegisteredCount int            `json:"registeredCount"`
	Status          ExamStatus     `json:"status"`
}

// AdminEventsResponse lists upcoming events.
type AdminEventsResponse struct {
	Items []AdminEventSummary `json:"items"`
}

// AdminReferralSummary exposes referral KPIs.
type AdminReferralSummary struct {
	Range          string `json:"range"`
	TotalReferrals int    `json:"totalReferrals"`
	RewardsPaid    int    `json:"rewardsPaid"`
	NewUsers       int    `json:"newUsers"`
}

// AdminUserStatus enumerates admin user states.
type AdminUserStatus string

const (
	AdminUserStatusActive    AdminUserStatus = "active"
	AdminUserStatusInactive  AdminUserStatus = "inactive"
	AdminUserStatusInvited   AdminUserStatus = "invited"
	AdminUserStatusSuspended AdminUserStatus = "suspended"
)

// AdminUserRole enumerates console roles.
type AdminUserRole string

const (
	AdminUserRoleSuperAdmin AdminUserRole = "superadmin"
	AdminUserRoleAdmin      AdminUserRole = "admin"
	AdminUserRoleManager    AdminUserRole = "manager"
	AdminUserRoleCashier    AdminUserRole = "cashier"
)

// AdminUser describes an operator.
type AdminUser struct {
	ID          uuid.UUID       `json:"id"`
	FirstName   string          `json:"firstName"`
	LastName    string          `json:"lastName"`
	Username    string          `json:"username"`
	Email       string          `json:"email"`
	PhoneNumber string          `json:"phoneNumber"`
	Status      AdminUserStatus `json:"status"`
	Role        AdminUserRole   `json:"role"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

// AdminUsersMeta contains pagination info.
type AdminUsersMeta struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

// AdminUserList envelope for list response.
type AdminUserList struct {
	Items []AdminUser    `json:"items"`
	Meta  AdminUsersMeta `json:"meta"`
}

// AdminUserFilter for querying users.
type AdminUserFilter struct {
	Page     int
	PageSize int
	Statuses []AdminUserStatus
	Role     *AdminUserRole
	Username string
}

// AdminUserCreateRequest payload.
type AdminUserCreateRequest struct {
	FirstName   string          `json:"firstName" validate:"required"`
	LastName    string          `json:"lastName" validate:"required"`
	Username    string          `json:"username" validate:"required"`
	Email       string          `json:"email" validate:"required,email"`
	PhoneNumber string          `json:"phoneNumber" validate:"required"`
	Role        AdminUserRole   `json:"role" validate:"required"`
	Status      AdminUserStatus `json:"status" validate:"required"`
	Password    string          `json:"password" validate:"required"`
}

// AdminUserUpdateRequest payload.
type AdminUserUpdateRequest struct {
	FirstName   *string          `json:"firstName,omitempty"`
	LastName    *string          `json:"lastName,omitempty"`
	Username    *string          `json:"username,omitempty"`
	Email       *string          `json:"email,omitempty"`
	PhoneNumber *string          `json:"phoneNumber,omitempty"`
	Role        *AdminUserRole   `json:"role,omitempty"`
	Status      *AdminUserStatus `json:"status,omitempty"`
	Password    *string          `json:"password,omitempty"`
}

// AdminBulkStatusRequest payload.
type AdminBulkStatusRequest struct {
	UserIDs []uuid.UUID     `json:"userIds" validate:"required,min=1,dive,required"`
	Status  AdminUserStatus `json:"status" validate:"required"`
}

// AdminBulkStatusResponse result.
type AdminBulkStatusResponse struct {
	Updated int `json:"updated"`
}

// AdminBulkDeleteRequest payload.
type AdminBulkDeleteRequest struct {
	UserIDs []uuid.UUID `json:"userIds" validate:"required,min=1,dive,required"`
}

// AdminBulkDeleteResponse result.
type AdminBulkDeleteResponse struct {
	Deleted int `json:"deleted"`
}

// AdminInviteRequest payload.
type AdminInviteRequest struct {
	Email   string        `json:"email" validate:"required,email"`
	Role    AdminUserRole `json:"role" validate:"required"`
	Message string        `json:"message"`
}

// AdminInviteResponse contains invite metadata.
type AdminInviteResponse struct {
	Invited   bool      `json:"invited"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// Subject describes a subject.
type Subject struct {
	ID       uuid.UUID    `json:"id"`
	Exam     ExamCategory `json:"exam"`
	Name     string       `json:"name"`
	IsActive bool         `json:"isActive"`
}

// Topic describes a topic under a subject.
type Topic struct {
	ID        uuid.UUID `json:"id"`
	SubjectID uuid.UUID `json:"subjectId"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"isActive"`
}

// LeaderboardEntry represents ranked user data.
type LeaderboardEntry struct {
	ID            uuid.UUID `json:"id"`
	DisplayName   string    `json:"displayName"`
	Score         int       `json:"score"`
	TotalCorrect  int       `json:"totalCorrect"`
	TotalAttempt  int       `json:"totalAttempt"`
	DayStreak     int       `json:"dayStreak"`
	EarnedRewards int       `json:"earnedRewards"`
	AvatarURL     string    `json:"avatarUrl,omitempty"`
}

// LeaderboardStats contains aggregate data for the leaderboard.
type LeaderboardStats struct {
	TotalUsers       int     `json:"totalUsers"`
	AverageAccuracy  float64 `json:"averageAccuracy"`
	LeaderboardRange string  `json:"leaderboardRange"`
}

// FeedPost represents a post in the public feed.
type FeedPost struct {
	ID        uuid.UUID `json:"id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	ImageURL  string    `json:"imageUrl,omitempty"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
	Author    string    `json:"author"`
	CTA       string    `json:"cta"`
	Likes     int       `json:"likes"`
	Comments  int       `json:"comments"`
	ReadTime  string    `json:"readTime"`
}

// Question represents a practice question.
type Question struct {
	ID              uuid.UUID          `json:"id"`
	Exam            ExamCategory       `json:"exam"`
	SubjectID       uuid.UUID          `json:"subjectId"`
	TopicID         uuid.UUID          `json:"topicId"`
	QuestionText    string             `json:"questionText"`
	OptionA         string             `json:"optionA"`
	OptionB         string             `json:"optionB"`
	OptionC         string             `json:"optionC"`
	OptionD         string             `json:"optionD"`
	CorrectOption   int                `json:"correctOption"`
	Explanation     *string            `json:"explanation,omitempty"`
	DifficultyLevel int                `json:"difficultyLevel"`
	ChoiceType      QuestionChoiceType `json:"choiceType"`
	IsClinical      bool               `json:"isClinical"`
	IsImageBased    bool               `json:"isImageBased"`
	IsHighYield     bool               `json:"isHighYield"`
	IsActive        bool               `json:"isActive"`
}

// QuestionCreateRequest body.
type QuestionCreateRequest struct {
	Exam            ExamCategory       `json:"exam" validate:"required"`
	SubjectID       uuid.UUID          `json:"subjectId" validate:"required"`
	TopicID         uuid.UUID          `json:"topicId" validate:"required"`
	QuestionText    string             `json:"questionText" validate:"required"`
	OptionA         string             `json:"optionA" validate:"required"`
	OptionB         string             `json:"optionB" validate:"required"`
	OptionC         string             `json:"optionC" validate:"required"`
	OptionD         string             `json:"optionD" validate:"required"`
	CorrectOption   int                `json:"correctOption" validate:"required,min=1,max=4"`
	Explanation     string             `json:"explanation"`
	DifficultyLevel int                `json:"difficultyLevel"`
	ChoiceType      QuestionChoiceType `json:"choiceType"`
	IsClinical      bool               `json:"isClinical"`
	IsImageBased    bool               `json:"isImageBased"`
	IsHighYield     bool               `json:"isHighYield"`
	IsActive        bool               `json:"isActive"`
}

// QuestionUpdateRequest body.
type QuestionUpdateRequest struct {
	SubjectID       *uuid.UUID          `json:"subjectId,omitempty"`
	TopicID         *uuid.UUID          `json:"topicId,omitempty"`
	QuestionText    *string             `json:"questionText,omitempty"`
	OptionA         *string             `json:"optionA,omitempty"`
	OptionB         *string             `json:"optionB,omitempty"`
	OptionC         *string             `json:"optionC,omitempty"`
	OptionD         *string             `json:"optionD,omitempty"`
	CorrectOption   *int                `json:"correctOption,omitempty"`
	Explanation     *string             `json:"explanation,omitempty"`
	DifficultyLevel *int                `json:"difficultyLevel,omitempty"`
	ChoiceType      *QuestionChoiceType `json:"choiceType,omitempty"`
	IsClinical      *bool               `json:"isClinical,omitempty"`
	IsImageBased    *bool               `json:"isImageBased,omitempty"`
	IsHighYield     *bool               `json:"isHighYield,omitempty"`
	IsActive        *bool               `json:"isActive,omitempty"`
}

// PracticeSession tracks a session.
type PracticeSession struct {
	ID                    uuid.UUID             `json:"id"`
	Mode                  PracticeMode          `json:"mode"`
	Exam                  ExamCategory          `json:"exam"`
	Status                PracticeSessionStatus `json:"status"`
	TotalQuestionsPlanned *int                  `json:"totalQuestionsPlanned,omitempty"`
	StartedAt             time.Time             `json:"startedAt"`
	CompletedAt           *time.Time            `json:"completedAt,omitempty"`
}

// PracticeSessionCreateRequest body.
type PracticeSessionCreateRequest struct {
	Mode             PracticeMode `json:"mode" validate:"required"`
	Exam             ExamCategory `json:"exam"`
	SubjectIDs       []uuid.UUID  `json:"subjectIds"`
	TopicIDs         []uuid.UUID  `json:"topicIds"`
	DifficultyLevels []int        `json:"difficultyLevels"`
	NumQuestions     int          `json:"numQuestions"`
	TimeLimitMinutes *int         `json:"timeLimitMinutes,omitempty"`
}

// PracticeSessionQuestion holds question within a session.
type PracticeSessionQuestion struct {
	ID             uuid.UUID  `json:"id"`
	SequenceIndex  int        `json:"sequenceIndex"`
	Question       Question   `json:"question"`
	SelectedOption *int       `json:"selectedOption,omitempty"`
	IsCorrect      *bool      `json:"isCorrect,omitempty"`
	TimeTakenMs    *int       `json:"timeTakenMs,omitempty"`
	AnsweredAt     *time.Time `json:"answeredAt,omitempty"`
}

// PracticeSessionDetail includes questions.
type PracticeSessionDetail struct {
	Session   PracticeSession           `json:"session"`
	Questions []PracticeSessionQuestion `json:"questions"`
}

// PracticeAnswerRequest payload.
type PracticeAnswerRequest struct {
	SessionQuestionID uuid.UUID `json:"sessionQuestionId" validate:"required"`
	SelectedOption    int       `json:"selectedOption" validate:"required"`
	TimeTakenMs       *int      `json:"timeTakenMs,omitempty"`
}

// RevisionItem for SRS queue.
type RevisionItem struct {
	ID            uuid.UUID `json:"id"`
	Question      Question  `json:"question"`
	NextReviewAt  time.Time `json:"nextReviewAt"`
	IntervalIndex int       `json:"intervalIndex"`
	TimesReviewed int       `json:"timesReviewed"`
}

// ExamConfig describes an exam event.
type ExamConfig struct {
	ID               uuid.UUID      `json:"id"`
	Exam             ExamCategory   `json:"exam"`
	Name             string         `json:"name"`
	Type             ExamConfigType `json:"type"`
	Description      string         `json:"description"`
	NumQuestions     int            `json:"numQuestions"`
	TimeLimitMinutes int            `json:"timeLimitMinutes"`
	MarksPerCorrect  float64        `json:"marksPerCorrect"`
	NegativePerWrong float64        `json:"negativePerWrong"`
	EntryFee         int            `json:"entryFee"`
	ScheduleStartAt  *time.Time     `json:"scheduleStartAt,omitempty"`
	ScheduleEndAt    *time.Time     `json:"scheduleEndAt,omitempty"`
	Status           ExamStatus     `json:"status"`
}

// ExamConfigCreateRequest body.
type ExamConfigCreateRequest struct {
	Exam             ExamCategory   `json:"exam" validate:"required"`
	Name             string         `json:"name" validate:"required"`
	Type             ExamConfigType `json:"type" validate:"required"`
	Description      string         `json:"description"`
	NumQuestions     int            `json:"numQuestions" validate:"required"`
	TimeLimitMinutes int            `json:"timeLimitMinutes" validate:"required"`
	MarksPerCorrect  float64        `json:"marksPerCorrect"`
	NegativePerWrong float64        `json:"negativePerWrong"`
	EntryFee         int            `json:"entryFee"`
	ScheduleStartAt  *time.Time     `json:"scheduleStartAt,omitempty"`
	ScheduleEndAt    *time.Time     `json:"scheduleEndAt,omitempty"`
}

// ExamConfigUpdateRequest body.
type ExamConfigUpdateRequest struct {
	Name             *string         `json:"name,omitempty"`
	Type             *ExamConfigType `json:"type,omitempty"`
	Description      *string         `json:"description,omitempty"`
	NumQuestions     *int            `json:"numQuestions,omitempty"`
	TimeLimitMinutes *int            `json:"timeLimitMinutes,omitempty"`
	MarksPerCorrect  *float64        `json:"marksPerCorrect,omitempty"`
	NegativePerWrong *float64        `json:"negativePerWrong,omitempty"`
	EntryFee         *int            `json:"entryFee,omitempty"`
	ScheduleStartAt  *time.Time      `json:"scheduleStartAt,omitempty"`
	ScheduleEndAt    *time.Time      `json:"scheduleEndAt,omitempty"`
	Status           *ExamStatus     `json:"status,omitempty"`
}

// ExamSummary returned by events list.
type ExamSummary struct {
	Config       ExamConfig `json:"config"`
	IsRegistered bool       `json:"isRegistered"`
	IsCompleted  bool       `json:"isCompleted"`
	BestScore    *float64   `json:"bestScore,omitempty"`
}

// PodcastEpisode describes audio content.
type PodcastEpisode struct {
	ID              uuid.UUID    `json:"id"`
	Exam            ExamCategory `json:"exam"`
	SubjectID       *uuid.UUID   `json:"subjectId,omitempty"`
	TopicID         *uuid.UUID   `json:"topicId,omitempty"`
	Title           string       `json:"title"`
	Description     string       `json:"description"`
	AudioURL        string       `json:"audioUrl"`
	DurationSeconds int          `json:"durationSeconds"`
	Tags            []string     `json:"tags"`
	IsActive        bool         `json:"isActive"`
}

// PodcastCreateRequest body.
type PodcastCreateRequest struct {
	Exam            ExamCategory `json:"exam" validate:"required"`
	SubjectID       *uuid.UUID   `json:"subjectId,omitempty"`
	TopicID         *uuid.UUID   `json:"topicId,omitempty"`
	Title           string       `json:"title" validate:"required"`
	Description     string       `json:"description"`
	AudioURL        string       `json:"audioUrl" validate:"required"`
	DurationSeconds int          `json:"durationSeconds"`
	Tags            []string     `json:"tags"`
	IsActive        bool         `json:"isActive"`
}

// WalletSummary describes balances.
type WalletSummary struct {
	Balance        int `json:"balance"`
	LifetimeEarned int `json:"lifetimeEarned"`
	LifetimeSpent  int `json:"lifetimeSpent"`
}

// WalletTransaction is a ledger entry.
type WalletTransaction struct {
	ID          uuid.UUID    `json:"id"`
	Amount      int          `json:"amount"`
	Type        WalletTxType `json:"type"`
	Description string       `json:"description"`
	CreatedAt   time.Time    `json:"createdAt"`
}

// Coupon describes a coupon.
type Coupon struct {
	ID             uuid.UUID  `json:"id"`
	Code           string     `json:"code"`
	Description    string     `json:"description"`
	Type           string     `json:"type"`
	Amount         int        `json:"amount"`
	MaxUsesTotal   int        `json:"maxUsesTotal"`
	MaxUsesPerUser int        `json:"maxUsesPerUser"`
	ExpiresAt      *time.Time `json:"expiresAt,omitempty"`
	IsActive       bool       `json:"isActive"`
}

// CouponCreateRequest body.
type CouponCreateRequest struct {
	Code           string     `json:"code" validate:"required"`
	Description    string     `json:"description"`
	Type           string     `json:"type" validate:"required"`
	Amount         int        `json:"amount"`
	MaxUsesTotal   int        `json:"maxUsesTotal"`
	MaxUsesPerUser int        `json:"maxUsesPerUser"`
	ExpiresAt      *time.Time `json:"expiresAt,omitempty"`
	IsActive       bool       `json:"isActive"`
}

// CouponRedeemRequest body.
type CouponRedeemRequest struct {
	Code string `json:"code" validate:"required"`
}

// ReferralSummary tracks referral stats.
type ReferralSummary struct {
	ReferralCode string `json:"referralCode"`
	TotalInvited int    `json:"totalInvited"`
	Joined       int    `json:"joined"`
	Activated    int    `json:"activated"`
	TotalEarned  int    `json:"totalEarned"`
}

// AISettings drives adaptive algorithms.
type AISettings struct {
	WeaknessMinAttempts      int   `json:"weaknessMinAttempts"`
	WeaknessThresholdPercent int   `json:"weaknessThresholdPercent"`
	StrongThresholdPercent   int   `json:"strongThresholdPercent"`
	RevisionIntervalsDays    []int `json:"revisionIntervalsDays"`
	IncludeGuessedCorrect    bool  `json:"includeGuessedCorrect"`
	RevisionEnabled          bool  `json:"revisionEnabled"`
}

// AnalyticsOverview returns dashboard metrics.
type AnalyticsOverview struct {
	TotalUsers          int     `json:"totalUsers"`
	ActiveUsers         int     `json:"activeUsers"`
	QuestionsAnswered   int     `json:"questionsAnswered"`
	AverageAccuracy     float64 `json:"averageAccuracy"`
	AverageStudyMinutes float64 `json:"averageStudyMinutes"`
	TotalRewards        int     `json:"totalRewards"`
}
