package exam

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
)

// UseCase manages exam config.
type UseCase struct {
	repo repo.ExamRepository
}

// New constructs UseCase.
func New(repo repo.ExamRepository) *UseCase {
	return &UseCase{repo: repo}
}

// AdminList returns configs with optional exam filter.
func (uc *UseCase) AdminList(ctx context.Context, exam *entity.ExamCategory) ([]entity.ExamConfig, error) {
	configs, err := uc.repo.ListConfigs(ctx)
	if err != nil {
		return nil, fmt.Errorf("exam - ListConfigs: %w", err)
	}

	if exam == nil || *exam == "" {
		return configs, nil
	}

	filtered := make([]entity.ExamConfig, 0, len(configs))
	for _, cfg := range configs {
		if cfg.Exam == *exam {
			filtered = append(filtered, cfg)
		}
	}

	return filtered, nil
}

// AdminCreate stores a config.
func (uc *UseCase) AdminCreate(ctx context.Context, req entity.ExamConfigCreateRequest) (entity.ExamConfig, error) {
	config := entity.ExamConfig{
		ID:               uuid.New(),
		Exam:             req.Exam,
		Name:             req.Name,
		Type:             req.Type,
		Description:      req.Description,
		NumQuestions:     req.NumQuestions,
		TimeLimitMinutes: req.TimeLimitMinutes,
		MarksPerCorrect:  req.MarksPerCorrect,
		NegativePerWrong: req.NegativePerWrong,
		EntryFee:         req.EntryFee,
		ScheduleStartAt:  req.ScheduleStartAt,
		ScheduleEndAt:    req.ScheduleEndAt,
		Status:           entity.ExamStatusDraft,
	}

	created, err := uc.repo.CreateConfig(ctx, config)
	if err != nil {
		return entity.ExamConfig{}, fmt.Errorf("exam - CreateConfig: %w", err)
	}

	return created, nil
}

// AdminGet retrieves config.
func (uc *UseCase) AdminGet(ctx context.Context, id uuid.UUID) (entity.ExamConfig, error) {
	config, err := uc.repo.GetConfig(ctx, id)
	if err != nil {
		return entity.ExamConfig{}, fmt.Errorf("exam - GetConfig: %w", err)
	}

	return config, nil
}

// AdminUpdate modifies config.
func (uc *UseCase) AdminUpdate(ctx context.Context, id uuid.UUID, req entity.ExamConfigUpdateRequest) (entity.ExamConfig, error) {
	config, err := uc.repo.GetConfig(ctx, id)
	if err != nil {
		return entity.ExamConfig{}, fmt.Errorf("exam - GetConfig: %w", err)
	}

	if req.Name != nil {
		config.Name = *req.Name
	}
	if req.Type != nil {
		config.Type = *req.Type
	}
	if req.Description != nil {
		config.Description = *req.Description
	}
	if req.NumQuestions != nil {
		config.NumQuestions = *req.NumQuestions
	}
	if req.TimeLimitMinutes != nil {
		config.TimeLimitMinutes = *req.TimeLimitMinutes
	}
	if req.MarksPerCorrect != nil {
		config.MarksPerCorrect = *req.MarksPerCorrect
	}
	if req.NegativePerWrong != nil {
		config.NegativePerWrong = *req.NegativePerWrong
	}
	if req.EntryFee != nil {
		config.EntryFee = *req.EntryFee
	}
	if req.ScheduleStartAt != nil {
		config.ScheduleStartAt = req.ScheduleStartAt
	}
	if req.ScheduleEndAt != nil {
		config.ScheduleEndAt = req.ScheduleEndAt
	}
	if req.Status != nil {
		config.Status = *req.Status
	}

	updated, err := uc.repo.UpdateConfig(ctx, config)
	if err != nil {
		return entity.ExamConfig{}, fmt.Errorf("exam - UpdateConfig: %w", err)
	}

	return updated, nil
}

// AdminDelete removes config.
func (uc *UseCase) AdminDelete(ctx context.Context, id uuid.UUID) error {
	if err := uc.repo.DeleteConfig(ctx, id); err != nil {
		return fmt.Errorf("exam - DeleteConfig: %w", err)
	}

	return nil
}

// ListEvents returns exam summaries.
func (uc *UseCase) ListEvents(ctx context.Context) ([]entity.ExamSummary, error) {
	summaries, err := uc.repo.ListSummaries(ctx)
	if err != nil {
		return nil, fmt.Errorf("exam - ListSummaries: %w", err)
	}

	return summaries, nil
}
