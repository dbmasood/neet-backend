package referral

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
)

// UseCase handles referral stats.
type UseCase struct {
	repo repo.ReferralRepository
}

// New constructs UseCase.
func New(repo repo.ReferralRepository) *UseCase {
	return &UseCase{repo: repo}
}

// Summary returns referral snapshot.
func (uc *UseCase) Summary(ctx context.Context, userID uuid.UUID) (entity.ReferralSummary, error) {
	summary, err := uc.repo.GetSummary(ctx, userID)
	if err != nil {
		return entity.ReferralSummary{}, fmt.Errorf("referral - Summary: %w", err)
	}

	return summary, nil
}
