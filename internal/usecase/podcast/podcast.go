package podcast

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
)

// UseCase handles podcast flows.
type UseCase struct {
	repo repo.PodcastRepository
}

// New constructs UseCase.
func New(repo repo.PodcastRepository) *UseCase {
	return &UseCase{repo: repo}
}

// List returns episodes for filters.
func (uc *UseCase) List(ctx context.Context, filter repo.PodcastFilter) ([]entity.PodcastEpisode, error) {
	episodes, err := uc.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("podcast - List: %w", err)
	}

	return episodes, nil
}

// Get returns a single episode.
func (uc *UseCase) Get(ctx context.Context, id uuid.UUID) (entity.PodcastEpisode, error) {
	episode, err := uc.repo.Get(ctx, id)
	if err != nil {
		return entity.PodcastEpisode{}, fmt.Errorf("podcast - Get: %w", err)
	}

	return episode, nil
}

// AdminCreate creates an episode.
func (uc *UseCase) AdminCreate(ctx context.Context, req entity.PodcastCreateRequest) (entity.PodcastEpisode, error) {
	episode := entity.PodcastEpisode{
		ID:              uuid.New(),
		Exam:            req.Exam,
		SubjectID:       req.SubjectID,
		TopicID:         req.TopicID,
		Title:           req.Title,
		Description:     req.Description,
		AudioURL:        req.AudioURL,
		DurationSeconds: req.DurationSeconds,
		Tags:            req.Tags,
		IsActive:        req.IsActive,
	}

	created, err := uc.repo.Create(ctx, episode)
	if err != nil {
		return entity.PodcastEpisode{}, fmt.Errorf("podcast - Create: %w", err)
	}

	return created, nil
}

// AdminUpdate modifies an episode.
func (uc *UseCase) AdminUpdate(ctx context.Context, id uuid.UUID, req entity.PodcastCreateRequest) (entity.PodcastEpisode, error) {
	episode, err := uc.repo.Get(ctx, id)
	if err != nil {
		return entity.PodcastEpisode{}, fmt.Errorf("podcast - Get: %w", err)
	}

	episode.Exam = req.Exam
	episode.SubjectID = req.SubjectID
	episode.TopicID = req.TopicID
	episode.Title = req.Title
	episode.Description = req.Description
	episode.AudioURL = req.AudioURL
	episode.DurationSeconds = req.DurationSeconds
	episode.Tags = req.Tags
	episode.IsActive = req.IsActive

	updated, err := uc.repo.Update(ctx, episode)
	if err != nil {
		return entity.PodcastEpisode{}, fmt.Errorf("podcast - Update: %w", err)
	}

	return updated, nil
}

// AdminGet returns an episode for admins.
func (uc *UseCase) AdminGet(ctx context.Context, id uuid.UUID) (entity.PodcastEpisode, error) {
	return uc.repo.Get(ctx, id)
}

// AdminDelete removes an episode.
func (uc *UseCase) AdminDelete(ctx context.Context, id uuid.UUID) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("podcast - Delete: %w", err)
	}

	return nil
}
