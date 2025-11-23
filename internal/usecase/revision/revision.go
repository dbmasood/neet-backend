package revision

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
)

// UseCase manages revision queue.
type UseCase struct {
	revisions repo.RevisionRepository
}

// New constructs UseCase.
func New(revisions repo.RevisionRepository) *UseCase {
	return &UseCase{revisions: revisions}
}

// GetQueue returns due items.
func (uc *UseCase) GetQueue(ctx context.Context, userID uuid.UUID) ([]entity.RevisionItem, error) {
	items, err := uc.revisions.ListDue(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("revision - ListDue: %w", err)
	}

	return items, nil
}
