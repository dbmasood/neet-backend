package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/pkg/jwt"
)

// UseCase for auth flows.
type UseCase struct {
	users repo.UserRepository
	jwt   *jwt.Service
}

// New constructs UseCase.
func New(users repo.UserRepository, jwtService *jwt.Service) *UseCase {
	return &UseCase{users: users, jwt: jwtService}
}

// TelegramAuth authenticates via Telegram.
func (uc *UseCase) TelegramAuth(ctx context.Context, req entity.TelegramAuthRequest) (entity.AuthResponse, error) {
	user, err := uc.users.GetByTelegramID(ctx, req.TelegramID)
	if err != nil {
		return entity.AuthResponse{}, fmt.Errorf("auth - GetByTelegramID: %w", err)
	}

	if user.ID == uuid.Nil {
		user = entity.User{
			ID:          uuid.New(),
			DisplayName: req.DisplayName,
			PrimaryExam: req.Exam,
			Role:        entity.UserRoleUser,
			CreatedAt:   time.Now().UTC(),
		}

		user, err = uc.users.Create(ctx, user)
		if err != nil {
			return entity.AuthResponse{}, fmt.Errorf("auth - Create: %w", err)
		}
	}

	token, err := uc.jwt.Generate(jwt.Claims{
		UserID: user.ID.String(),
		Role:   user.Role,
		Exam:   user.PrimaryExam,
	})
	if err != nil {
		return entity.AuthResponse{}, fmt.Errorf("auth - Generate token: %w", err)
	}

	return entity.AuthResponse{AccessToken: token, User: user}, nil
}
