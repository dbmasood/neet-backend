package analytics

import (
	"context"
	"fmt"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
)

// UseCase provides metrics for dashboards.
type UseCase struct {
	repo repo.AnalyticsRepository
}

// New constructs UseCase.
func New(repo repo.AnalyticsRepository) *UseCase {
	return &UseCase{repo: repo}
}

// Overview returns aggregated metrics.
func (uc *UseCase) Overview(ctx context.Context, filter repo.AnalyticsFilter) (entity.AnalyticsOverview, error) {
	overview, err := uc.repo.Overview(ctx, filter)
	if err != nil {
		return entity.AnalyticsOverview{}, fmt.Errorf("analytics - Overview: %w", err)
	}

	return overview, nil
}
