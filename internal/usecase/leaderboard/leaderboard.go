package leaderboard

import (
	"context"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
)

// UseCase groups leaderboard logic.
type UseCase struct {
	repo repo.LeaderboardRepository
}

// New creates Leaderboard usecase.
func New(repo repo.LeaderboardRepository) *UseCase {
	return &UseCase{repo: repo}
}

// Entries returns entries and stats.
func (uc *UseCase) Entries(ctx context.Context, limit int) ([]entity.LeaderboardEntry, entity.LeaderboardStats, error) {
	entries, err := uc.repo.List(ctx, limit)
	if err != nil {
		return nil, entity.LeaderboardStats{}, err
	}

	stats, err := uc.repo.Stats(ctx)
	if err != nil {
		return nil, entity.LeaderboardStats{}, err
	}

	return entries, stats, nil
}
