package admin

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/evrone/go-clean-template/internal/entity"
)

var (
	// ErrInvalidMetric returned for unsupported metric names.
	ErrInvalidMetric = errors.New("invalid metric")
	// ErrInvalidRange returned for unsupported ranges.
	ErrInvalidRange = errors.New("invalid range")
	// ErrUserNotFound when user id missing.
	ErrUserNotFound = errors.New("user not found")
	// ErrDuplicateUsername when username exists.
	ErrDuplicateUsername = errors.New("username already exists")
	// ErrDuplicateEmail when email exists.
	ErrDuplicateEmail = errors.New("email already exists")
)

// UseCase orchestrates admin specific flows.
type UseCase struct {
	profile entity.AdminProfile

	mu    sync.RWMutex
	users map[uuid.UUID]entity.AdminUser
}

// New constructs UseCase with bootstrap profile.
func New(profile entity.AdminProfile) *UseCase {
	uc := &UseCase{
		profile: profile,
		users:   make(map[uuid.UUID]entity.AdminUser),
	}

	uc.seedUsers()
	return uc
}

func (uc *UseCase) seedUsers() {
	first, last := splitName(uc.profile.DisplayName)
	now := uc.profile.CreatedAt
	if now.IsZero() {
		now = time.Now().UTC()
	}

	uc.users[uc.profile.ID] = entity.AdminUser{
		ID:          uc.profile.ID,
		FirstName:   first,
		LastName:    last,
		Username:    slugifyName(uc.profile.DisplayName),
		Email:       uc.profile.Email,
		PhoneNumber: "+91 90000 00000",
		Status:      entity.AdminUserStatusActive,
		Role:        normalizeRole(uc.profile.Role),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Profile returns admin bootstrap identity.
func (uc *UseCase) Profile(context.Context) (entity.AdminProfile, error) {
	return uc.profile, nil
}

// TimeSeries returns chart data for metric.
func (uc *UseCase) TimeSeries(_ context.Context, metric string, exam *entity.ExamCategory, window string) (entity.AnalyticsTimeSeries, error) {
	metric = strings.ToLower(metric)
	if metric != "active_users" && metric != "questions_answered" {
		return entity.AnalyticsTimeSeries{}, ErrInvalidMetric
	}

	if window == "" {
		window = "7d"
	}

	days, err := rangeToDays(window)
	if err != nil {
		return entity.AnalyticsTimeSeries{}, err
	}

	targetExam := uc.profile.PrimaryExam
	if exam != nil && *exam != "" {
		targetExam = *exam
	}

	points := make([]entity.AnalyticsPoint, 0, days)
	now := time.Now().UTC()
	for i := days - 1; i >= 0; i-- {
		day := now.AddDate(0, 0, -i)
		points = append(points, entity.AnalyticsPoint{
			Date:  day.Format("2006-01-02"),
			Value: metricValue(metric, day),
		})
	}

	return entity.AnalyticsTimeSeries{
		Metric: metric,
		Exam:   targetExam,
		Range:  window,
		Points: points,
	}, nil
}

// SubjectAccuracy returns performance per subject.
func (uc *UseCase) SubjectAccuracy(_ context.Context, exam *entity.ExamCategory) (entity.SubjectAccuracyResponse, error) {
	targetExam := uc.profile.PrimaryExam
	if exam != nil && *exam != "" {
		targetExam = *exam
	}

	subjects := []entity.SubjectAccuracyItem{
		{SubjectID: "subj_anat", SubjectName: "Anatomy", Accuracy: 0.72},
		{SubjectID: "subj_biochem", SubjectName: "Biochemistry", Accuracy: 0.64},
		{SubjectID: "subj_path", SubjectName: "Pathology", Accuracy: 0.58},
		{SubjectID: "subj_pharma", SubjectName: "Pharmacology", Accuracy: 0.61},
	}

	return entity.SubjectAccuracyResponse{
		Exam:     targetExam,
		Subjects: subjects,
	}, nil
}

// WeakTopics returns weakest topics list.
func (uc *UseCase) WeakTopics(_ context.Context, exam *entity.ExamCategory, limit int) (entity.WeakTopicsResponse, error) {
	targetExam := uc.profile.PrimaryExam
	if exam != nil && *exam != "" {
		targetExam = *exam
	}
	prefix := strings.ReplaceAll(string(targetExam), "_", " ")

	all := []entity.WeakTopicItem{
		{SubjectID: "subj_pharma", SubjectName: fmt.Sprintf("%s Pharmacology", prefix), TopicID: "topic_autonomic", TopicName: "Autonomic Drugs", Accuracy: 0.42, Attempts: 1240},
		{SubjectID: "subj_path", SubjectName: fmt.Sprintf("%s Pathology", prefix), TopicID: "topic_neoplasia", TopicName: "Neoplasia", Accuracy: 0.48, Attempts: 980},
		{SubjectID: "subj_micro", SubjectName: fmt.Sprintf("%s Microbiology", prefix), TopicID: "topic_virology", TopicName: "Virology", Accuracy: 0.45, Attempts: 1110},
		{SubjectID: "subj_anat", SubjectName: fmt.Sprintf("%s Anatomy", prefix), TopicID: "topic_neuro", TopicName: "Neuro Anatomy", Accuracy: 0.41, Attempts: 890},
	}

	if limit <= 0 || limit > len(all) {
		limit = len(all)
	}

	return entity.WeakTopicsResponse{Items: all[:limit]}, nil
}

// UpcomingEvents returns scheduled mocks.
func (uc *UseCase) UpcomingEvents(_ context.Context, exam *entity.ExamCategory) (entity.AdminEventsResponse, error) {
	targetExam := uc.profile.PrimaryExam
	if exam != nil && *exam != "" {
		targetExam = *exam
	}

	items := []entity.AdminEventSummary{
		{
			ID:              "mock-bio-1",
			Name:            fmt.Sprintf("%s - High Yield Bio Mock", strings.ReplaceAll(string(targetExam), "_", " ")),
			Exam:            targetExam,
			Type:            entity.ExamTypeMock,
			StartAt:         time.Now().Add(72 * time.Hour),
			RegisteredCount: 2420,
			Status:          entity.ExamStatusScheduled,
		},
		{
			ID:              "daily-test-2",
			Name:            "Daily Rapid Fire",
			Exam:            targetExam,
			Type:            entity.ExamTypeDailyTest,
			StartAt:         time.Now().Add(24 * time.Hour),
			RegisteredCount: 1340,
			Status:          entity.ExamStatusScheduled,
		},
	}

	return entity.AdminEventsResponse{Items: items}, nil
}

// ReferralSummary returns referral KPIs.
func (uc *UseCase) ReferralSummary(_ context.Context, window string) (entity.AdminReferralSummary, error) {
	if window == "" {
		window = "30d"
	}

	summary := entity.AdminReferralSummary{
		Range:          window,
		TotalReferrals: 980,
		RewardsPaid:    18200,
		NewUsers:       320,
	}
	if window == "7d" {
		summary.TotalReferrals = 210
		summary.RewardsPaid = 4200
		summary.NewUsers = 75
	} else if window == "today" {
		summary.TotalReferrals = 28
		summary.RewardsPaid = 980
		summary.NewUsers = 12
	}

	return summary, nil
}

// ListUsers returns paginated admin users.
func (uc *UseCase) ListUsers(_ context.Context, filter entity.AdminUserFilter) (entity.AdminUserList, error) {
	uc.mu.RLock()
	defer uc.mu.RUnlock()

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	var filtered []entity.AdminUser
	for _, user := range uc.users {
		if len(filter.Statuses) > 0 && !statusMatch(user.Status, filter.Statuses) {
			continue
		}
		if filter.Role != nil && user.Role != *filter.Role {
			continue
		}
		if filter.Username != "" && !strings.Contains(strings.ToLower(user.Username), strings.ToLower(filter.Username)) {
			continue
		}
		filtered = append(filtered, user)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})

	total := len(filtered)
	start := (filter.Page - 1) * filter.PageSize
	if start > total {
		start = total
	}
	end := start + filter.PageSize
	if end > total {
		end = total
	}

	return entity.AdminUserList{
		Items: filtered[start:end],
		Meta: entity.AdminUsersMeta{
			Page:     filter.Page,
			PageSize: filter.PageSize,
			Total:    total,
		},
	}, nil
}

// CreateUser registers a new admin user.
func (uc *UseCase) CreateUser(_ context.Context, req entity.AdminUserCreateRequest) (entity.AdminUser, error) {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if uc.usernameExists(req.Username, uuid.Nil) {
		return entity.AdminUser{}, ErrDuplicateUsername
	}
	if uc.emailExists(req.Email, uuid.Nil) {
		return entity.AdminUser{}, ErrDuplicateEmail
	}

	now := time.Now().UTC()
	user := entity.AdminUser{
		ID:          uuid.New(),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Username:    req.Username,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Status:      req.Status,
		Role:        req.Role,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	uc.users[user.ID] = user
	return user, nil
}

// UpdateUser mutates existing admin user.
func (uc *UseCase) UpdateUser(_ context.Context, id uuid.UUID, req entity.AdminUserUpdateRequest) (entity.AdminUser, error) {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	user, ok := uc.users[id]
	if !ok {
		return entity.AdminUser{}, ErrUserNotFound
	}

	if req.Username != nil && *req.Username != user.Username && uc.usernameExists(*req.Username, id) {
		return entity.AdminUser{}, ErrDuplicateUsername
	}
	if req.Email != nil && *req.Email != user.Email && uc.emailExists(*req.Email, id) {
		return entity.AdminUser{}, ErrDuplicateEmail
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.PhoneNumber != nil {
		user.PhoneNumber = *req.PhoneNumber
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.Status != nil {
		user.Status = *req.Status
	}
	user.UpdatedAt = time.Now().UTC()

	uc.users[id] = user
	return user, nil
}

// DeleteUser removes an admin user.
func (uc *UseCase) DeleteUser(_ context.Context, id uuid.UUID) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	if _, ok := uc.users[id]; !ok {
		return ErrUserNotFound
	}
	delete(uc.users, id)
	return nil
}

// BulkStatus updates statuses for multiple users.
func (uc *UseCase) BulkStatus(_ context.Context, req entity.AdminBulkStatusRequest) (int, error) {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	updated := 0
	for _, id := range req.UserIDs {
		if user, ok := uc.users[id]; ok {
			user.Status = req.Status
			user.UpdatedAt = time.Now().UTC()
			uc.users[id] = user
			updated++
		}
	}

	return updated, nil
}

// BulkDelete removes users in batch.
func (uc *UseCase) BulkDelete(_ context.Context, req entity.AdminBulkDeleteRequest) (int, error) {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	deleted := 0
	for _, id := range req.UserIDs {
		if _, ok := uc.users[id]; ok {
			delete(uc.users, id)
			deleted++
		}
	}
	return deleted, nil
}

// InviteUser records an invite command.
func (uc *UseCase) InviteUser(_ context.Context, req entity.AdminInviteRequest) (entity.AdminInviteResponse, error) {
	return entity.AdminInviteResponse{
		Invited:   true,
		ExpiresAt: time.Now().Add(72 * time.Hour).UTC(),
	}, nil
}

func (uc *UseCase) usernameExists(username string, exclude uuid.UUID) bool {
	for id, u := range uc.users {
		if id == exclude {
			continue
		}
		if strings.EqualFold(u.Username, username) {
			return true
		}
	}
	return false
}

func (uc *UseCase) emailExists(email string, exclude uuid.UUID) bool {
	for id, u := range uc.users {
		if id == exclude {
			continue
		}
		if strings.EqualFold(u.Email, email) {
			return true
		}
	}
	return false
}

func statusMatch(candidate entity.AdminUserStatus, statuses []entity.AdminUserStatus) bool {
	for _, status := range statuses {
		if candidate == status {
			return true
		}
	}
	return false
}

func rangeToDays(window string) (int, error) {
	switch strings.ToLower(window) {
	case "today":
		return 1, nil
	case "7d", "7day", "7days":
		return 7, nil
	case "30d", "30day", "30days":
		return 30, nil
	default:
		return 0, ErrInvalidRange
	}
}

func metricValue(metric string, t time.Time) int {
	base := 750
	if metric == "questions_answered" {
		base = 1800
	}
	offset := int(t.Unix()/86400)%200 + 100
	return base + offset
}

func splitName(full string) (string, string) {
	parts := strings.Fields(full)
	if len(parts) == 0 {
		return "Admin", "User"
	}
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], strings.Join(parts[1:], " ")
}

func slugifyName(full string) string {
	slug := strings.ToLower(strings.ReplaceAll(full, " ", "."))
	return strings.ReplaceAll(slug, "__", ".")
}

func normalizeRole(role string) entity.AdminUserRole {
	switch strings.ToLower(role) {
	case string(entity.AdminUserRoleSuperAdmin):
		return entity.AdminUserRoleSuperAdmin
	case string(entity.AdminUserRoleManager):
		return entity.AdminUserRoleManager
	case string(entity.AdminUserRoleCashier):
		return entity.AdminUserRoleCashier
	default:
		return entity.AdminUserRoleAdmin
	}
}
