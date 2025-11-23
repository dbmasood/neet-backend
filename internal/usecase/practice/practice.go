package practice

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
)

// UseCase orchestrates practice flows.
type UseCase struct {
	sessions repo.PracticeSessionRepository
}

// New constructs UseCase.
func New(sessions repo.PracticeSessionRepository) *UseCase {
	return &UseCase{sessions: sessions}
}

// CreateSession starts a new session.
func (uc *UseCase) CreateSession(ctx context.Context, userID uuid.UUID, req entity.PracticeSessionCreateRequest) (entity.PracticeSession, error) {
	session := entity.PracticeSession{
		ID:        uuid.New(),
		Mode:      req.Mode,
		Exam:      req.Exam,
		Status:    entity.PracticeStatusInProgress,
		StartedAt: time.Now().UTC(),
	}

	created, err := uc.sessions.CreateSession(ctx, session)
	if err != nil {
		return entity.PracticeSession{}, fmt.Errorf("practice - CreateSession: %w", err)
	}

	return created, nil
}

// ListSessions returns practice history.
func (uc *UseCase) ListSessions(ctx context.Context, userID uuid.UUID) ([]entity.PracticeSession, error) {
	sessions, err := uc.sessions.ListSessions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("practice - ListSessions: %w", err)
	}

	return sessions, nil
}

// GetSessionDetail returns session questions.
func (uc *UseCase) GetSessionDetail(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) (entity.PracticeSessionDetail, error) {
	session, err := uc.sessions.GetSession(ctx, sessionID)
	if err != nil {
		return entity.PracticeSessionDetail{}, fmt.Errorf("practice - GetSession: %w", err)
	}

	questions, err := uc.sessions.ListSessionQuestions(ctx, sessionID)
	if err != nil {
		return entity.PracticeSessionDetail{}, fmt.Errorf("practice - ListSessionQuestions: %w", err)
	}

	return entity.PracticeSessionDetail{Session: session, Questions: questions}, nil
}

// AnswerQuestion records response.
func (uc *UseCase) AnswerQuestion(ctx context.Context, sessionID uuid.UUID, req entity.PracticeAnswerRequest, userID uuid.UUID) (entity.PracticeSessionQuestion, error) {
	question, err := uc.sessions.GetSessionQuestion(ctx, req.SessionQuestionID)
	if err != nil {
		return entity.PracticeSessionQuestion{}, fmt.Errorf("practice - GetSessionQuestion: %w", err)
	}

	correct := question.Question.CorrectOption == req.SelectedOption
	question.SelectedOption = &req.SelectedOption
	question.IsCorrect = &correct
	question.TimeTakenMs = req.TimeTakenMs

	updated, err := uc.sessions.UpdateSessionQuestion(ctx, question)
	if err != nil {
		return entity.PracticeSessionQuestion{}, fmt.Errorf("practice - UpdateSessionQuestion: %w", err)
	}

	return updated, nil
}
