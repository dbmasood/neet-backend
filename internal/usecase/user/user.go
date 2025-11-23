package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
)

// UseCase for user-facing flows.
type UseCase struct {
	users    repo.UserRepository
	subjects repo.SubjectRepository
	topics   repo.TopicRepository
}

// New constructs UseCase.
func New(users repo.UserRepository, subjects repo.SubjectRepository, topics repo.TopicRepository) *UseCase {
	return &UseCase{users: users, subjects: subjects, topics: topics}
}

// Me returns current user payload.
func (uc *UseCase) Me(ctx context.Context, userID uuid.UUID) (entity.MeResponse, error) {
	user, err := uc.users.GetByID(ctx, userID)
	if err != nil {
		return entity.MeResponse{}, fmt.Errorf("user - GetByID: %w", err)
	}

	profile, err := uc.users.GetExamProfile(ctx, userID)
	if err != nil {
		return entity.MeResponse{}, fmt.Errorf("user - GetExamProfile: %w", err)
	}

	return entity.MeResponse{User: user, ExamProfile: profile}, nil
}

// ListSubjects returns subjects.
func (uc *UseCase) ListSubjects(ctx context.Context, exam *entity.ExamCategory) ([]entity.Subject, error) {
	subjects, err := uc.subjects.ListByExam(ctx, exam)
	if err != nil {
		return nil, fmt.Errorf("user - ListSubjects: %w", err)
	}

	return subjects, nil
}

// ListTopics returns topics for a subject.
func (uc *UseCase) ListTopics(ctx context.Context, subjectID uuid.UUID) ([]entity.Topic, error) {
	topics, err := uc.topics.ListBySubject(ctx, subjectID)
	if err != nil {
		return nil, fmt.Errorf("user - ListTopics: %w", err)
	}

	return topics, nil
}

// CreateSubject registers a new subject.
func (uc *UseCase) CreateSubject(ctx context.Context, exam entity.ExamCategory, name string) (entity.Subject, error) {
	subject := entity.Subject{
		ID:       uuid.New(),
		Exam:     exam,
		Name:     name,
		IsActive: true,
	}

	created, err := uc.subjects.Create(ctx, subject)
	if err != nil {
		return entity.Subject{}, fmt.Errorf("user - CreateSubject: %w", err)
	}

	return created, nil
}

// CreateTopic registers a new topic.
func (uc *UseCase) CreateTopic(ctx context.Context, subjectID uuid.UUID, name string) (entity.Topic, error) {
	topic := entity.Topic{
		ID:        uuid.New(),
		SubjectID: subjectID,
		Name:      name,
		IsActive:  true,
	}

	created, err := uc.topics.Create(ctx, topic)
	if err != nil {
		return entity.Topic{}, fmt.Errorf("user - CreateTopic: %w", err)
	}

	return created, nil
}
