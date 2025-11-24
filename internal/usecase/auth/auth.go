package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
	"github.com/evrone/go-clean-template/pkg/jwt"
)

// UseCase for auth flows.
type UseCase struct {
	users      repo.UserRepository
	userJWT    *jwt.Service
	adminJWT   *jwt.Service
	adminCreds AdminCredentials
}

// AdminCredentials carries env-configured bootstrap admin identity.
type AdminCredentials struct {
	Username    string
	Password    string
	DisplayName string
	UserID      uuid.UUID
	PrimaryExam entity.ExamCategory
	Email       string
	Role        entity.UserRole
	Permissions []string
	CreatedAt   time.Time
}

// ErrInvalidAdminCredentials returned when username/password mismatch.
var ErrInvalidAdminCredentials = errors.New("invalid admin credentials")

// New constructs UseCase.
func New(users repo.UserRepository, userJWT, adminJWT *jwt.Service, creds AdminCredentials) *UseCase {
	if creds.PrimaryExam == "" {
		creds.PrimaryExam = entity.ExamCategoryNEETPG
	}
	if creds.Role == "" {
		creds.Role = entity.UserRoleSuperAdmin
	}
	return &UseCase{users: users, userJWT: userJWT, adminJWT: adminJWT, adminCreds: creds}
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

	token, err := uc.userJWT.Generate(jwt.Claims{
		UserID: user.ID.String(),
		Role:   user.Role,
		Exam:   user.PrimaryExam,
	})
	if err != nil {
		return entity.AuthResponse{}, fmt.Errorf("auth - Generate token: %w", err)
	}

	return entity.AuthResponse{AccessToken: token, User: user}, nil
}

// AdminLogin authenticates env-configured admin credentials.
func (uc *UseCase) AdminLogin(_ context.Context, req entity.AdminLoginRequest) (entity.AuthResponse, error) {
	if req.Username != uc.adminCreds.Username || req.Password != uc.adminCreds.Password {
		return entity.AuthResponse{}, ErrInvalidAdminCredentials
	}

	claims := jwt.Claims{
		UserID: uc.adminCreds.UserID.String(),
		Role:   uc.adminCreds.Role,
		Exam:   uc.adminCreds.PrimaryExam,
	}

	token, err := uc.adminJWT.Generate(claims)
	if err != nil {
		return entity.AuthResponse{}, fmt.Errorf("auth - AdminLogin - generate token: %w", err)
	}

	email := uc.adminCreds.Email
	user := entity.User{
		ID:          uc.adminCreds.UserID,
		DisplayName: uc.adminCreds.DisplayName,
		Email:       &email,
		PrimaryExam: uc.adminCreds.PrimaryExam,
		Role:        uc.adminCreds.Role,
		CreatedAt:   uc.adminCreds.CreatedAt,
	}

	return entity.AuthResponse{AccessToken: token, User: user}, nil
}
