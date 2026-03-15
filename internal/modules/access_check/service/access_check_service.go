package service

import (
	"context"
	"errors"
	"fmt"

	entity "rextra-backend/internal/entity"
	myerror "rextra-backend/internal/pkg/error"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CheckAccessRequest adalah parameter yang dikirim modul pemanggil.
type CheckAccessRequest struct {
	UserID         uuid.UUID
	EntitlementKey string
	// ReferenceID opsional — diisi pemanggil jika akses menghasilkan record baru.
	// Format: "{entity_type}:{uuid}", contoh: "cv_generate:abc-123"
	// Akan disimpan di UsageLog.ReferenceID untuk tracing.
	ReferenceID *string
}

// CheckAccessResponse adalah hasil dari CheckAccess.
type CheckAccessResponse struct {
	Granted         bool
	DenyReason      string // diisi jika Granted = false
	RestrictionType entity.RestrictionType
	// TokenCost hanya diisi jika RestrictionType = token_gated
	TokenCost int
	// QuotaRemaining hanya diisi jika RestrictionType = frequency_limited dan Granted = true.
	// Dipakai UI untuk tampilkan "sisa kuota kamu N" setelah akses berhasil — tanpa extra request.
	// Untuk saldo token saat ini, UI fetch terpisah dari modul token: GET /api/v1/token/wallet/balance.
	// JANGAN branch UI berdasarkan DenyReason string — gunakan RestrictionType enum yang stabil.
	QuotaRemaining *int
}

type accessCheckRepository interface {
	// Dari entitlement + membership module:
	// Cari DurationAccessMapping milik membership aktif user untuk entitlement yang diminta.
	// Harus Preload("Entitlement") karena RestrictionType dibaca dari sana.
	// Return gorm.ErrRecordNotFound jika plan user tidak punya akses ke entitlement ini.
	// Parameter planDurationID diambil dari membership.DurationID — field di entity Memberships.
	// BUKAN PlanDurationID (itu ada di SubscriptionCycle, bukan Memberships).
	// [E-5 FIX] Diubah ke *uuid.UUID karena Memberships.DurationID adalah nullable pointer.
	// Nil berarti user adalah plan Standard/Starter tanpa durasi berbayar — nil guard wajib
	// dilakukan di service sebelum memanggil method ini.
	GetMappingByMembershipAndKey(ctx context.Context, planDurationID *uuid.UUID, entitlementKey string) (entity.DurationAccessMapping, error)

	// Dari part4 — UserEntitlementQuota:
	// Cari quota record milik membership + entitlement tertentu.
	// Dipakai hanya untuk frequency_limited.
	GetQuotaByMembershipAndKey(ctx context.Context, membershipID uuid.UUID, entitlementKey string) (entity.UserEntitlementQuota, error)

	// Update QuotaUsed dan QuotaRemaining setelah ConsumeOne().
	// Wajib dipanggil dalam DB transaction yang sama dengan InsertUsageLog.
	UpdateQuota(ctx context.Context, tx *gorm.DB, quota entity.UserEntitlementQuota) error

	// Ambil membership aktif user.
	// Return gorm.ErrRecordNotFound jika user tidak punya membership aktif.
	GetActiveMembershipByUserID(ctx context.Context, userID uuid.UUID) (entity.Memberships, error)

	// Insert UsageLog — dipanggil untuk semua hasil (granted maupun denied).
	InsertUsageLog(ctx context.Context, tx *gorm.DB, log entity.UsageLog) error
}

type (
	AccessCheckService interface {
		// CheckAccess memutuskan apakah user boleh menggunakan suatu entitlement.
		// Selalu insert UsageLog — baik granted maupun denied.
		// Untuk token_gated: deduct token dalam satu DB transaction dengan UsageLog.
		// Untuk frequency_limited: ConsumeOne() dalam satu DB transaction dengan UsageLog.
		CheckAccess(ctx context.Context, req CheckAccessRequest) (CheckAccessResponse, error)
	}

	accessCheckService struct {
		repo      accessCheckRepository
		tokenRepo TokenRepository
		db        *gorm.DB
	}
)

func NewAccessCheckService(
	repo accessCheckRepository,
	tokenRepo TokenRepository,
	db *gorm.DB,
) AccessCheckService {
	return &accessCheckService{
		repo:      repo,
		tokenRepo: tokenRepo,
		db:        db,
	}
}

func (s *accessCheckService) CheckAccess(ctx context.Context, req CheckAccessRequest) (CheckAccessResponse, error) {
	// Step 1: Ambil membership aktif user
	membership, err := s.repo.GetActiveMembershipByUserID(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// User tidak punya membership aktif — deny, tapi tetap log
			_ = s.insertLog(ctx, nil, uuid.Nil, req, entity.UsageLogResultDenied, strPtr("tidak punya membership aktif"), nil)
			return CheckAccessResponse{
				Granted:    false,
				DenyReason: "tidak punya membership aktif",
			}, nil
		}
		return CheckAccessResponse{}, myerror.DatabaseError(err)
	}

	// Step 2: Cek apakah plan user punya akses ke entitlement ini
	// DurationAccessMapping di-preload dengan Entitlement untuk baca RestrictionType
	// DurationID dari entity Memberships — bukan PlanDurationID (field itu ada di SubscriptionCycle).
	//
	// [E-5 FIX] Nil guard: Standard/Starter tidak punya DurationID karena plan mereka tidak
	// terdaftar di DurationAccessMapping (plan tanpa durasi berbayar). Tanpa guard ini, memanggil
	// GetMappingByMembershipAndKey dengan *uuid.UUID nil akan panic (nil pointer dereference).
	// Langsung deny — plan tanpa durasi memang tidak punya konfigurasi entitlement.
	if membership.DurationID == nil {
		reason := fmt.Sprintf("plan %s tidak memiliki konfigurasi durasi aktif — akses fitur tidak tersedia", membership.PlanName)
		_ = s.insertLog(ctx, nil, membership.ID, req, entity.UsageLogResultDenied, strPtr(reason), nil)
		return CheckAccessResponse{
			Granted:    false,
			DenyReason: reason,
		}, nil
	}
	mapping, err := s.repo.GetMappingByMembershipAndKey(ctx, membership.DurationID, req.EntitlementKey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			reason := fmt.Sprintf("plan %s tidak punya akses ke fitur ini", membership.PlanName)
			_ = s.insertLog(ctx, nil, membership.ID, req, entity.UsageLogResultDenied, strPtr(reason), nil)
			return CheckAccessResponse{
				Granted:    false,
				DenyReason: reason,
			}, nil
		}
		return CheckAccessResponse{}, myerror.DatabaseError(err)
	}

	restrictionType := mapping.Entitlement.RestrictionType

	// Step 3: Handle berdasarkan RestrictionType
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
		return CheckAccessResponse{}, myerror.New(
			fmt.Sprintf("restriction type tidak dikenal: %s", restrictionType),
			myerror.SystemError,
		)
	}
}

// ─── Handlers per RestrictionType ─────────────────────────────────────────────

func (s *accessCheckService) handleUnlimited(ctx context.Context, membership entity.Memberships, req CheckAccessRequest) (CheckAccessResponse, error) {
	if err := s.insertLog(ctx, nil, membership.ID, req, entity.UsageLogResultSuccess, nil, nil); err != nil {
		return CheckAccessResponse{}, myerror.DatabaseError(err)
	}
	return CheckAccessResponse{
		Granted:         true,
		RestrictionType: entity.RestrictionUnlimited,
	}, nil
}

func (s *accessCheckService) handleLocked(ctx context.Context, membership entity.Memberships, req CheckAccessRequest) (CheckAccessResponse, error) {
	reason := "fitur ini terkunci untuk plan kamu"
	// Best effort — kegagalan log tidak boleh block user. Konsekuensi: analytics "berapa kali
	// user coba akses fitur locked" tidak 100% akurat. Acceptable untuk operasional normal.
	_ = s.insertLog(ctx, nil, membership.ID, req, entity.UsageLogResultDenied, strPtr(reason), nil)
	return CheckAccessResponse{
		Granted:         false,
		DenyReason:      reason,
		RestrictionType: entity.RestrictionLocked,
	}, nil
}

func (s *accessCheckService) handleTokenGated(ctx context.Context, membership entity.Memberships, tokenCost int, req CheckAccessRequest) (CheckAccessResponse, error) {
	// Cek saldo token
	wallet, err := s.tokenRepo.GetWalletByUserID(ctx, nil, req.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			reason := fmt.Sprintf("saldo token tidak cukup, dibutuhkan %d token", tokenCost)
			_ = s.insertLog(ctx, nil, membership.ID, req, entity.UsageLogResultDenied, strPtr(reason), nil)
			return CheckAccessResponse{
				Granted:         false,
				DenyReason:      reason,
				RestrictionType: entity.RestrictionTokenGated,
				TokenCost:       tokenCost,
			}, nil
		}
		return CheckAccessResponse{}, myerror.DatabaseError(err)
	}

	if wallet.Balance < int64(tokenCost) {
		reason := fmt.Sprintf("saldo token tidak cukup, dibutuhkan %d token, saldo kamu %d", tokenCost, wallet.Balance)
		_ = s.insertLog(ctx, nil, membership.ID, req, entity.UsageLogResultDenied, strPtr(reason), nil)
		return CheckAccessResponse{
			Granted:         false,
			DenyReason:      reason,
			RestrictionType: entity.RestrictionTokenGated,
			TokenCost:       tokenCost,
		}, nil
	}

	// Saldo cukup — deduct token + insert UsageLog dalam satu DB transaction.
	// Kalau salah satu gagal, keduanya rollback: tidak ada token terpotong tanpa log,
	// dan tidak ada log tanpa token terpotong.
	var result CheckAccessResponse
	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		description := fmt.Sprintf("penggunaan fitur: %s", req.EntitlementKey)
		if err := s.tokenRepo.DeductToken(ctx, tx, req.UserID, tokenCost, description); err != nil {
			return err
		}
		if err := s.insertLog(ctx, tx, membership.ID, req, entity.UsageLogResultSuccess, nil, req.ReferenceID); err != nil {
			return err
		}
		result = CheckAccessResponse{
			Granted:         true,
			RestrictionType: entity.RestrictionTokenGated,
			TokenCost:       tokenCost,
		}
		return nil
	})
	if txErr != nil {
		return CheckAccessResponse{}, myerror.DatabaseError(txErr)
	}
	return result, nil
}

func (s *accessCheckService) handleFrequencyLimited(ctx context.Context, membership entity.Memberships, req CheckAccessRequest) (CheckAccessResponse, error) {
	quota, err := s.repo.GetQuotaByMembershipAndKey(ctx, membership.ID, req.EntitlementKey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Edge case: quota record belum dibuat karena payment callback gagal.
			// Deny dengan pesan berbeda dari "kuota habis" agar UI bisa tampilkan
			// state "akses sementara tidak tersedia" yang berbeda dari state kuota habis.
			reason := "kuota belum tersedia, hubungi support"
			_ = s.insertLog(ctx, nil, membership.ID, req, entity.UsageLogResultDenied, strPtr(reason), nil)
			return CheckAccessResponse{
				Granted:    false,
				DenyReason: reason,
			}, nil
		}
		return CheckAccessResponse{}, myerror.DatabaseError(err)
	}

	// Cek apakah siklus quota masih aktif
	if !quota.IsActive() {
		reason := "kuota kamu sudah expired, perbarui membership untuk mendapatkan kuota baru"
		_ = s.insertLog(ctx, nil, membership.ID, req, entity.UsageLogResultDenied, strPtr(reason), nil)
		return CheckAccessResponse{
			Granted:         false,
			DenyReason:      reason,
			RestrictionType: entity.RestrictionFrequencyLimited,
		}, nil
	}

	// Cek apakah kuota habis
	if quota.IsExhausted() {
		reason := fmt.Sprintf("kuota kamu sudah habis (%d/%d)", quota.QuotaUsed, quota.QuotaGranted)
		_ = s.insertLog(ctx, nil, membership.ID, req, entity.UsageLogResultDenied, strPtr(reason), nil)
		return CheckAccessResponse{
			Granted:         false,
			DenyReason:      reason,
			RestrictionType: entity.RestrictionFrequencyLimited,
		}, nil
	}

	// Kuota tersedia — ConsumeOne() + UpdateQuota + insert UsageLog dalam satu DB transaction.
	// Kalau salah satu gagal, keduanya rollback: kuota tidak berkurang tanpa log.
	var result CheckAccessResponse
	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		quota.ConsumeOne()
		if err := s.repo.UpdateQuota(ctx, tx, quota); err != nil {
			return err
		}
		if err := s.insertLog(ctx, tx, membership.ID, req, entity.UsageLogResultSuccess, nil, req.ReferenceID); err != nil {
			return err
		}
		remaining := quota.QuotaRemaining
		result = CheckAccessResponse{
			Granted:         true,
			RestrictionType: entity.RestrictionFrequencyLimited,
			QuotaRemaining:  &remaining,
		}
		return nil
	})
	if txErr != nil {
		return CheckAccessResponse{}, myerror.DatabaseError(txErr)
	}
	return result, nil
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func (s *accessCheckService) insertLog(
	ctx context.Context,
	tx *gorm.DB,
	membershipID uuid.UUID,
	req CheckAccessRequest,
	result entity.UsageLogResult,
	errorMessage *string,
	referenceID *string,
) error {
	// Ambil entitlement name untuk UsageLog — best effort, tidak fatal kalau gagal
	// karena entitlement key saja sudah cukup untuk tracing
	entitlementName := req.EntitlementKey

	log := entity.NewUsageLog(membershipID, req.EntitlementKey, entitlementName, result, errorMessage, referenceID)
	return s.repo.InsertUsageLog(ctx, tx, log)
}

func strPtr(s string) *string {
	return &s
}

type accessCheckRepositoryImpl struct {
	db *gorm.DB
}

func NewAccessCheckRepository(db *gorm.DB) accessCheckRepository {
	return &accessCheckRepositoryImpl{db: db}
}

func (r *accessCheckRepositoryImpl) GetActiveMembershipByUserID(ctx context.Context, userID uuid.UUID) (entity.Memberships, error) {
	var membership entity.Memberships
	// Double check: is_active dikelola manual saat expired/downgrade, tapi kalau background job
	// expire belum jalan, expired_at > NOW() tetap jadi safety net. Kedua kondisi wajib ada.
	err := r.db.WithContext(ctx).
		// expired_at nullable: plan Standard (seumur hidup) punya expired_at = nil.
		// NULL > NOW() = false di SQL → user Standard akan selalu ter-exclude kalau tidak pakai OR.
		// Tambah OR expired_at IS NULL agar Standard tetap dapat akses jika punya entitlement.
		Where("user_id = ? AND is_active = true AND (expired_at > NOW() OR expired_at IS NULL)", userID).
		First(&membership).Error
	return membership, err
}

// [E-5 FIX] Parameter diubah ke *uuid.UUID agar konsisten dengan Memberships.DurationID.
// Nil check sudah dilakukan di service layer sebelum memanggil method ini,
// sehingga di sini aman untuk dereference (*planDurationID).
func (r *accessCheckRepositoryImpl) GetMappingByMembershipAndKey(ctx context.Context, planDurationID *uuid.UUID, entitlementKey string) (entity.DurationAccessMapping, error) {
	var mapping entity.DurationAccessMapping
	err := r.db.WithContext(ctx).
		Preload("Entitlement").
		Joins("JOIN entitlements ON entitlements.id = duration_access_mappings.entitlement_id").
		// is_active tidak ada di entity Entitlement — field yang benar adalah Status (varchar).
		// Tambah juga filter status mapping supaya mapping yang dinonaktifkan admin tidak memberikan akses.
		// *planDurationID aman di-dereference karena nil guard sudah dilewati di service.
		Where("duration_access_mappings.plan_duration_id = ? AND duration_access_mappings.status = 'aktif' AND entitlements.key = ? AND entitlements.status = 'active'", *planDurationID, entitlementKey).
		First(&mapping).Error
	return mapping, err
}

func (r *accessCheckRepositoryImpl) GetQuotaByMembershipAndKey(ctx context.Context, membershipID uuid.UUID, entitlementKey string) (entity.UserEntitlementQuota, error) {
	var quota entity.UserEntitlementQuota
	err := r.db.WithContext(ctx).
		Where("membership_id = ? AND entitlement_key = ?", membershipID, entitlementKey).
		First(&quota).Error
	return quota, err
}

func (r *accessCheckRepositoryImpl) UpdateQuota(ctx context.Context, tx *gorm.DB, quota entity.UserEntitlementQuota) error {
	if tx == nil {
		tx = r.db
	}
	return tx.WithContext(ctx).
		Model(&quota).
		Updates(map[string]any{
			"quota_used":      quota.QuotaUsed,
			"quota_remaining": quota.QuotaRemaining,
		}).Error
}

func (r *accessCheckRepositoryImpl) InsertUsageLog(ctx context.Context, tx *gorm.DB, log entity.UsageLog) error {
	if tx == nil {
		tx = r.db
	}
	return tx.WithContext(ctx).Create(&log).Error
}
