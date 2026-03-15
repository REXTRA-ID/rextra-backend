package repository

import (
	"context"
	"rextra-backend/internal/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PengaturanRepository interface {
	GetOrCreateInvoiceSettings(ctx context.Context) (entity.InvoiceSettings, error)
	UpdateInvoiceSettings(ctx context.Context, model entity.InvoiceSettings) (entity.InvoiceSettings, error)
	GetOrCreateTrxIdSettings(ctx context.Context) (entity.TransactionIdSettings, error)
	UpdateTrxIdSettings(ctx context.Context, model entity.TransactionIdSettings) (entity.TransactionIdSettings, error)
	GetOrCreateNotifSettings(ctx context.Context) (entity.MembershipNotificationSettings, error)
	UpdateNotifSettings(ctx context.Context, model entity.MembershipNotificationSettings) (entity.MembershipNotificationSettings, error)
}

type pengaturanRepository struct {
	db *gorm.DB
}

func NewPengaturanRepository(db *gorm.DB) PengaturanRepository {
	return &pengaturanRepository{db: db}
}

func (r *pengaturanRepository) GetOrCreateInvoiceSettings(ctx context.Context) (entity.InvoiceSettings, error) {
	var result entity.InvoiceSettings
	defaults := entity.DefaultInvoiceSettings()
	err := r.db.WithContext(ctx).Where(entity.InvoiceSettings{ID: defaults.ID}).Attrs(defaults).FirstOrCreate(&result).Error
	return result, err
}

func (r *pengaturanRepository) UpdateInvoiceSettings(ctx context.Context, model entity.InvoiceSettings) (entity.InvoiceSettings, error) {
	err := r.db.WithContext(ctx).Omit(clause.Associations).Save(&model).Error
	return model, err
}

func (r *pengaturanRepository) GetOrCreateTrxIdSettings(ctx context.Context) (entity.TransactionIdSettings, error) {
	var result entity.TransactionIdSettings
	defaults := entity.DefaultTransactionIdSettings()
	err := r.db.WithContext(ctx).Where(entity.TransactionIdSettings{ID: defaults.ID}).Attrs(defaults).FirstOrCreate(&result).Error
	return result, err
}

func (r *pengaturanRepository) UpdateTrxIdSettings(ctx context.Context, model entity.TransactionIdSettings) (entity.TransactionIdSettings, error) {
	err := r.db.WithContext(ctx).Omit(clause.Associations).Save(&model).Error
	return model, err
}

func (r *pengaturanRepository) GetOrCreateNotifSettings(ctx context.Context) (entity.MembershipNotificationSettings, error) {
	var result entity.MembershipNotificationSettings
	defaults := entity.DefaultMembershipNotificationSettings()
	err := r.db.WithContext(ctx).Where(entity.MembershipNotificationSettings{ID: defaults.ID}).Attrs(defaults).FirstOrCreate(&result).Error
	return result, err
}

func (r *pengaturanRepository) UpdateNotifSettings(ctx context.Context, model entity.MembershipNotificationSettings) (entity.MembershipNotificationSettings, error) {
	err := r.db.WithContext(ctx).Omit(clause.Associations).Save(&model).Error
	return model, err
}
