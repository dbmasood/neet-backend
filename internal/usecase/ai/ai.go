package ai

import (
	"context"
	"fmt"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
)

// UseCase manages AI settings.
type UseCase struct {
	repo repo.AISettingsRepository
}

// New constructs UseCase.
func New(repo repo.AISettingsRepository) *UseCase {
	return &UseCase{repo: repo}
}

// Get returns current settings.
func (uc *UseCase) Get(ctx context.Context) (entity.AISettings, error) {
	settings, err := uc.repo.Get(ctx)
	if err != nil {
		return entity.AISettings{}, fmt.Errorf("ai - Get: %w", err)
	}

	return settings, nil
}

// Update persists new settings.
func (uc *UseCase) Update(ctx context.Context, settings entity.AISettings) (entity.AISettings, error) {
	updated, err := uc.repo.Update(ctx, settings)
	if err != nil {
		return entity.AISettings{}, fmt.Errorf("ai - Update: %w", err)
	}

	return updated, nil
}
