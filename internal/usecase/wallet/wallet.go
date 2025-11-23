package wallet

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
)

// UseCase handles wallet flows.
type UseCase struct {
	repo repo.WalletRepository
}

// New constructs UseCase.
func New(repo repo.WalletRepository) *UseCase {
	return &UseCase{repo: repo}
}

// Summary returns wallet snapshot.
func (uc *UseCase) Summary(ctx context.Context, userID uuid.UUID) (entity.WalletSummary, error) {
	summary, err := uc.repo.GetSummary(ctx, userID)
	if err != nil {
		return entity.WalletSummary{}, fmt.Errorf("wallet - GetSummary: %w", err)
	}

	return summary, nil
}

// ListTransactions returns history.
func (uc *UseCase) ListTransactions(ctx context.Context, userID uuid.UUID) ([]entity.WalletTransaction, error) {
	txs, err := uc.repo.ListTransactions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("wallet - ListTransactions: %w", err)
	}

	return txs, nil
}
