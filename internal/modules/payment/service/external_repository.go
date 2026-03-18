package service

import (
	"context"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PromoRepository — subset dari promo repository
type PromoRepository interface {
	CountEligibleVouchers(ctx context.Context, planID string) (int, error)
	ValidateAndCalculateDiscount(ctx context.Context, code string, planID string, subtotal int64) (int64, error)
	RecordRedemption(ctx context.Context, code string, userID uuid.UUID, transactionID uuid.UUID) error
}

// UserRepository — subset dari user repository untuk tripay customer
type UserRepository interface {
	GetById(ctx context.Context, tx *gorm.DB, id string) (entity.User, error)
}

// PromoFacade — implementasi dari PromoRepository untuk membungkus module promo_code
type PromoFacade struct {}

func NewPromoFacade() PromoRepository {
	return &PromoFacade{}
}

func (f *PromoFacade) CountEligibleVouchers(ctx context.Context, planID string) (int, error) {
	// TODO: implementasi sesungguhnya karena promo belum ada usecase count_vouchers
	return 0, nil
}

func (f *PromoFacade) ValidateAndCalculateDiscount(ctx context.Context, code string, planID string, subtotal int64) (int64, error) {
	// TODO: implementasi sesungguhnya karena validasi promo module butuh refaktor untuk checkout baru
	return 0, nil
}

func (f *PromoFacade) RecordRedemption(ctx context.Context, code string, userID uuid.UUID, transactionID uuid.UUID) error {
	// TODO: implementasi sesungguhnya karena redemption promo butuh sinkronisasi DTO
	return nil
}
