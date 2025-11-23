package coupon

import (
	"context"
	"fmt"
	"github.com/google/uuid"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
)

// UseCase handles coupons.
type UseCase struct {
	repo repo.CouponRepository
}

// New constructs UseCase.
func New(repo repo.CouponRepository) *UseCase {
	return &UseCase{repo: repo}
}

// Redeem applies a coupon code.
func (uc *UseCase) Redeem(ctx context.Context, userID uuid.UUID, req entity.CouponRedeemRequest) (entity.WalletSummary, error) {
	summary, err := uc.repo.Redeem(ctx, req.Code, userID)
	if err != nil {
		return entity.WalletSummary{}, fmt.Errorf("coupon - Redeem: %w", err)
	}

	return summary, nil
}

// AdminList returns coupons.
func (uc *UseCase) AdminList(ctx context.Context) ([]entity.Coupon, error) {
	coupons, err := uc.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("coupon - List: %w", err)
	}

	return coupons, nil
}

// AdminCreate stores coupon.
func (uc *UseCase) AdminCreate(ctx context.Context, req entity.CouponCreateRequest) (entity.Coupon, error) {
	coupon := entity.Coupon{
		ID:             uuid.New(),
		Code:           req.Code,
		Description:    req.Description,
		Type:           req.Type,
		Amount:         req.Amount,
		MaxUsesTotal:   req.MaxUsesTotal,
		MaxUsesPerUser: req.MaxUsesPerUser,
		ExpiresAt:      req.ExpiresAt,
		IsActive:       req.IsActive,
	}

	created, err := uc.repo.Create(ctx, coupon)
	if err != nil {
		return entity.Coupon{}, fmt.Errorf("coupon - Create: %w", err)
	}

	return created, nil
}

// AdminGet returns coupon by id.
func (uc *UseCase) AdminGet(ctx context.Context, id uuid.UUID) (entity.Coupon, error) {
	coupon, err := uc.repo.Get(ctx, id)
	if err != nil {
		return entity.Coupon{}, fmt.Errorf("coupon - Get: %w", err)
	}

	return coupon, nil
}

// AdminUpdate modifies coupon details.
func (uc *UseCase) AdminUpdate(ctx context.Context, id uuid.UUID, req entity.CouponCreateRequest) (entity.Coupon, error) {
	coupon, err := uc.repo.Get(ctx, id)
	if err != nil {
		return entity.Coupon{}, fmt.Errorf("coupon - Get: %w", err)
	}

	coupon.Code = req.Code
	coupon.Description = req.Description
	coupon.Type = req.Type
	coupon.Amount = req.Amount
	coupon.MaxUsesTotal = req.MaxUsesTotal
	coupon.MaxUsesPerUser = req.MaxUsesPerUser
	coupon.ExpiresAt = req.ExpiresAt
	coupon.IsActive = req.IsActive

	updated, err := uc.repo.Update(ctx, coupon)
	if err != nil {
		return entity.Coupon{}, fmt.Errorf("coupon - Update: %w", err)
	}

	return updated, nil
}

// AdminDelete removes a coupon.
func (uc *UseCase) AdminDelete(ctx context.Context, id uuid.UUID) error {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("coupon - Delete: %w", err)
	}

	return nil
}
