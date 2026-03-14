package service

import (
	"context"
	"errors"
	"fmt"

	"rextra-backend/internal/entity"
	myerror "rextra-backend/internal/pkg/error"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	CheckAccessRequest struct {
		UserID         uuid.UUID
		EntitlementKey string
		ReferenceID    *string
	}

	CheckAccessResponse struct {
		Granted         bool
		DenyReason      string
		RestrictionType entity.RestrictionType
		TokenCost       int
		QuotaRemaining  *int
	}

	AccessCheckService interface {
		CheckAccess(ctx context.Context, req CheckAccessRequest) (CheckAccessResponse, error)
	}

	accessCheckRepository interface {
		GetMappingByMembershipAndKey(ctx context.Context, planDurationID *uuid.UUID, entitlementKey string) (entity.DurationAccessMapping, error)
		GetQuotaByMembershipAndKey(ctx context.Context, membershipID uuid.UUID, entitlementKey string) (entity.UserEntitlementQuota, error)
		UpdateQuota(ctx context.Context, tx *gorm.DB, quota entity.UserEntitlementQuota) error
		GetActiveMembershipByUserID(ctx context.Context, userID uuid.UUID) (entity.Memberships, error)
		InsertUsageLog(ctx context.Context, tx *gorm.DB, log entity.UsageLog) error
	}

	accessCheckService struct {
		repo      accessCheckRepository
		tokenRepo TokenRepository
		db        *gorm.DB
	}

	accessCheckRepositoryImpl struct {
		db *gorm.DB
	}
)

func NewAccessCheckService(repo accessCheckRepository, tokenRepo TokenRepository, db *gorm.DB) AccessCheckService {
	return &accessCheckService{repo: repo, tokenRepo: tokenRepo, db: db}
}

func NewAccessCheckRepository(db *gorm.DB) accessCheckRepository {
	return &accessCheckRepositoryImpl{db: db}
}

func (s *accessCheckService) CheckAccess(ctx context.Context, req CheckAccessRequest) (CheckAccessResponse, error) {
	membership, err := s.repo.GetActiveMembershipByUserID(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = s.insertLog(ctx, nil, uuid.Nil, req, entity.UsageLogResultDenied, strPtr("tidak punya membership aktif"), nil)
			return CheckAccessResponse{Granted: false, DenyReason: "tidak punya membership aktif"}, nil
		}
		return CheckAccessResponse{}, myerror.DatabaseError(err)
	}

	if membership.DurationID == nil {
		reason := fmt.Sprintf("plan %s tidak memiliki konfigurasi durasi aktif — akses fitur tidak tersedia", membership.PlanName)
		_ = s.insertLog(ctx, nil, membership.ID, req, entity.UsageLogResultDenied, strPtr(reason), nil)
		return CheckAccessResponse{Granted: false, DenyReason: reason}, nil
	}

	mapping, err := s.repo.GetMappingByMembershipAndKey(ctx, membership.DurationID, req.EntitlementKey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			reason := fmt.Sprintf("plan %s tidak punya akses ke fitur ini", membership.PlanName)
			_ = s.insertLog(ctx, nil, membership.ID, req, entity.UsageLogResultDenied, strPtr(reason), nil)
			return CheckAccessResponse{Granted: false, DenyReason: reason}, nil
		}
		return CheckAccessResponse{}, myerror.DatabaseError(err)
	}

	restrictionType := mapping.Entitlement.RestrictionType
	switch restrictionType {
	case entity.RestrictionUnlimited:
		return s.handleUnlimited(ctx, membership, req)
	case entity.RestrictionLocked:
		return s.handleLocked(ctx, membership, req)
	case entity.RestrictionTokenGated:
		return s.handleTokenGated(ctx, membership, mapping.Entitlement.TokenCost, req)
	case entity.RestrictionFrequencyLimited:
		return s.handleFrequencyLimited(ctx, membership, req)
	default:
		return CheckAccessResponse{}, myerror.New(fmt.Sprintf("restriction type tidak dikenal: %s", restrictionType), myerror.SystemError)
	}
}

func (s *accessCheckService) handleUnlimited(ctx context.Context, membership entity.Memberships, req CheckAccessRequest) (CheckAccessResponse, error) {
	if err := s.insertLog(ctx, nil, membership.ID, req, entity.UsageLogResultSuccess, nil, nil); err != nil { return CheckAccessResponse{}, myerror.DatabaseError(err) }
	return CheckAccessResponse{Granted: true, RestrictionType: entity.RestrictionUnlimited}, nil
}

func (s *accessCheckService) handleLocked(ctx context.Context, membership entity.Memberships, req CheckAccessRequest) (CheckAccessResponse, error) {
	reason := "fitur ini terkunci untuk plan kamu"
	_ = s.insertLog(ctx, nil, membership.ID, req, entity.UsageLogResultDenied, strPtr(reason), nil)
	return CheckAccessResponse{Granted: false, DenyReason: reason, RestrictionType: entity.RestrictionLocked}, nil
}

func (s *accessCheckService) handleTokenGated(ctx context.Context, membership entity.Memberships, tokenCost int, req CheckAccessRequest) (CheckAccessResponse, error) {
	wallet, err := s.tokenRepo.GetWalletByUserID(ctx, nil, req.UserID)
	if err != nil || wallet.Balance < int64(tokenCost) {
		reason := fmt.Sprintf("saldo token tidak cukup, dibutuhkan %d token", tokenCost)
		_ = s.insertLog(ctx, nil, membership.ID, req, entity.UsageLogResultDenied, strPtr(reason), nil)
		return CheckAccessResponse{Granted: false, DenyReason: reason, RestrictionType: entity.RestrictionTokenGated, TokenCost: tokenCost}, nil
	}

	var result CheckAccessResponse
	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.tokenRepo.DeductToken(ctx, tx, req.UserID, tokenCost, fmt.Sprintf("penggunaan fitur: %s", req.EntitlementKey)); err != nil { return err }
		if err := s.insertLog(ctx, tx, membership.ID, req, entity.UsageLogResultSuccess, nil, req.ReferenceID); err != nil { return err }
		result = CheckAccessResponse{Granted: true, RestrictionType: entity.RestrictionTokenGated, TokenCost: tokenCost}
		return nil
	})
	if txErr != nil { return CheckAccessResponse{}, myerror.DatabaseError(txErr) }
	return result, nil
}

func (s *accessCheckService) handleFrequencyLimited(ctx context.Context, membership entity.Memberships, req CheckAccessRequest) (CheckAccessResponse, error) {
	quota, err := s.repo.GetQuotaByMembershipAndKey(ctx, membership.ID, req.EntitlementKey)
	if err != nil || !quota.IsActive() || quota.IsExhausted() {
		reason := "kuota kamu sudah habis atau expired"
		_ = s.insertLog(ctx, nil, membership.ID, req, entity.UsageLogResultDenied, strPtr(reason), nil)
		return CheckAccessResponse{Granted: false, DenyReason: reason, RestrictionType: entity.RestrictionFrequencyLimited}, nil
	}

	var result CheckAccessResponse
	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		quota.ConsumeOne()
		if err := s.repo.UpdateQuota(ctx, tx, quota); err != nil { return err }
		if err := s.insertLog(ctx, tx, membership.ID, req, entity.UsageLogResultSuccess, nil, req.ReferenceID); err != nil { return err }
		remaining := quota.QuotaRemaining
		result = CheckAccessResponse{Granted: true, RestrictionType: entity.RestrictionFrequencyLimited, QuotaRemaining: &remaining}
		return nil
	})
	if txErr != nil { return CheckAccessResponse{}, myerror.DatabaseError(txErr) }
	return result, nil
}

func (s *accessCheckService) insertLog(ctx context.Context, tx *gorm.DB, membershipID uuid.UUID, req CheckAccessRequest, result entity.UsageLogResult, errorMessage *string, referenceID *string) error {
	log := entity.NewUsageLog(membershipID, req.EntitlementKey, req.EntitlementKey, result, errorMessage, referenceID)
	return s.repo.InsertUsageLog(ctx, tx, log)
}

func strPtr(s string) *string { return &s }

// ─── Repository Implementation ──────────────────────────────────────────────

func (r *accessCheckRepositoryImpl) GetActiveMembershipByUserID(ctx context.Context, userID uuid.UUID) (entity.Memberships, error) {
	var membership entity.Memberships
	err := r.db.WithContext(ctx).Where("user_id = ? AND is_active = true AND (expired_at > NOW() OR expired_at IS NULL)", userID).First(&membership).Error
	return membership, err
}

func (r *accessCheckRepositoryImpl) GetMappingByMembershipAndKey(ctx context.Context, planDurationID *uuid.UUID, entitlementKey string) (entity.DurationAccessMapping, error) {
	var mapping entity.DurationAccessMapping
	err := r.db.WithContext(ctx).Preload("Entitlement").Joins("JOIN entitlements ON entitlements.id = duration_access_mappings.entitlement_id").Where("duration_access_mappings.plan_duration_id = ? AND duration_access_mappings.status = 'aktif' AND entitlements.key = ? AND entitlements.status = 'active'", *planDurationID, entitlementKey).First(&mapping).Error
	return mapping, err
}

func (r *accessCheckRepositoryImpl) GetQuotaByMembershipAndKey(ctx context.Context, membershipID uuid.UUID, entitlementKey string) (entity.UserEntitlementQuota, error) {
	var quota entity.UserEntitlementQuota
	err := r.db.WithContext(ctx).Where("membership_id = ? AND entitlement_key = ?", membershipID, entitlementKey).First(&quota).Error
	return quota, err
}

func (r *accessCheckRepositoryImpl) UpdateQuota(ctx context.Context, tx *gorm.DB, quota entity.UserEntitlementQuota) error {
	if tx == nil { tx = r.db }
	return tx.WithContext(ctx).Model(&quota).Updates(map[string]any{"quota_used": quota.QuotaUsed, "quota_remaining": quota.QuotaRemaining}).Error
}

func (r *accessCheckRepositoryImpl) InsertUsageLog(ctx context.Context, tx *gorm.DB, log entity.UsageLog) error {
	if tx == nil { tx = r.db }
	return tx.WithContext(ctx).Create(&log).Error
}
