package feed

import (
	"context"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
)

// UseCase provides feed content.
type UseCase struct {
	repo repo.FeedRepository
}

// New creates feed usecase.
func New(repo repo.FeedRepository) *UseCase {
	return &UseCase{repo: repo}
}

// List returns posts.
func (uc *UseCase) List(ctx context.Context) ([]entity.FeedPost, error) {
	return uc.repo.List(ctx)
}
