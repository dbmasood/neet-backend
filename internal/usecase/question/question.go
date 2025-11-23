package question

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
)

// UseCase for admin question management.
type UseCase struct {
	repo repo.QuestionRepository
}

// New constructs UseCase.
func New(repo repo.QuestionRepository) *UseCase {
	return &UseCase{repo: repo}
}

// AdminList returns questions filtered.
func (uc *UseCase) AdminList(ctx context.Context, filter repo.QuestionFilter) ([]entity.Question, error) {
	questions, err := uc.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("question - List: %w", err)
	}

	return questions, nil
}

// AdminCreate persists a new question.
func (uc *UseCase) AdminCreate(ctx context.Context, req entity.QuestionCreateRequest) (entity.Question, error) {
	question := entity.Question{
		ID:              uuid.New(),
		Exam:            req.Exam,
		SubjectID:       req.SubjectID,
		TopicID:         req.TopicID,
		QuestionText:    req.QuestionText,
		OptionA:         req.OptionA,
		OptionB:         req.OptionB,
		OptionC:         req.OptionC,
		OptionD:         req.OptionD,
		CorrectOption:   req.CorrectOption,
		Explanation:     stringPointer(req.Explanation),
		DifficultyLevel: req.DifficultyLevel,
		ChoiceType:      req.ChoiceType,
		IsClinical:      req.IsClinical,
		IsImageBased:    req.IsImageBased,
		IsHighYield:     req.IsHighYield,
		IsActive:        req.IsActive,
	}

	created, err := uc.repo.Create(ctx, question)
	if err != nil {
		return entity.Question{}, fmt.Errorf("question - Create: %w", err)
	}

	return created, nil
}

// AdminGet returns a question by ID.
func (uc *UseCase) AdminGet(ctx context.Context, id uuid.UUID) (entity.Question, error) {
	question, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return entity.Question{}, fmt.Errorf("question - GetByID: %w", err)
	}

	return question, nil
}

// AdminUpdate updates question fields.
func (uc *UseCase) AdminUpdate(ctx context.Context, id uuid.UUID, req entity.QuestionUpdateRequest) (entity.Question, error) {
	question, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return entity.Question{}, fmt.Errorf("question - GetByID: %w", err)
	}

	if req.SubjectID != nil {
		question.SubjectID = *req.SubjectID
	}
	if req.TopicID != nil {
		question.TopicID = *req.TopicID
	}
	if req.QuestionText != nil {
		question.QuestionText = *req.QuestionText
	}
	if req.OptionA != nil {
		question.OptionA = *req.OptionA
	}
	if req.OptionB != nil {
		question.OptionB = *req.OptionB
	}
	if req.OptionC != nil {
		question.OptionC = *req.OptionC
	}
	if req.OptionD != nil {
		question.OptionD = *req.OptionD
	}
	if req.CorrectOption != nil {
		question.CorrectOption = *req.CorrectOption
	}
	if req.Explanation != nil {
		question.Explanation = req.Explanation
	}
	if req.DifficultyLevel != nil {
		question.DifficultyLevel = *req.DifficultyLevel
	}
	if req.ChoiceType != nil {
		question.ChoiceType = *req.ChoiceType
	}
	if req.IsClinical != nil {
		question.IsClinical = *req.IsClinical
	}
	if req.IsImageBased != nil {
		question.IsImageBased = *req.IsImageBased
	}
	if req.IsHighYield != nil {
		question.IsHighYield = *req.IsHighYield
	}
	if req.IsActive != nil {
		question.IsActive = *req.IsActive
	}

	updated, err := uc.repo.Update(ctx, question)
	if err != nil {
		return entity.Question{}, fmt.Errorf("question - Update: %w", err)
	}

	return updated, nil
}

// AdminDelete removes a question.
func (uc *UseCase) AdminDelete(ctx context.Context, id uuid.UUID) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("question - Delete: %w", err)
	}

	return nil
}

func stringPointer(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}
