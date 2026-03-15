package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserEntitlementQuota struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	MembershipID uuid.UUID `json:"membership_id" gorm:"type:uuid;not null;index:idx_ueq_unique,unique"`
	EntitlementID uuid.UUID `json:"entitlement_id" gorm:"type:uuid;not null;index:idx_ueq_unique,unique"`

	// Snapshot key untuk kemudahan query tanpa join — immutable setelah dibuat
	EntitlementKey string `json:"entitlement_key" gorm:"type:varchar(200);not null;index"`

	// QuotaGranted adalah snapshot dari DurationAccessMapping.UsageLimit saat user beli.
	// Nilai ini tidak berubah sepanjang siklus — mencerminkan berapa total jatah yang diberikan.
	// Contoh: user beli Pro 3 bulan → QuotaGranted = 6 (langsung di awal, bukan 2 per bulan).
	QuotaGranted int `json:"quota_granted" gorm:"not null;default:0"`

	// QuotaUsed bertambah 1 setiap user berhasil menggunakan fitur ini.
	// Di-update oleh access-check service dalam satu DB transaction yang sama
	// dengan pembuatan UsageLog — tidak boleh update salah satu tanpa yang lain.
	QuotaUsed int `json:"quota_used" gorm:"not null;default:0"`

	// QuotaRemaining adalah derived value: QuotaGranted - QuotaUsed.
	// Disimpan eksplisit (bukan computed column) agar query access check bisa
	// langsung filter WHERE quota_remaining > 0 tanpa kalkulasi di application layer.
	// Wajib di-update bersamaan dengan QuotaUsed — selalu QuotaGranted - QuotaUsed.
	QuotaRemaining int `json:"quota_remaining" gorm:"not null;default:0"`

	// Periode berlaku — mengikuti siklus membership yang bersangkutan.
	// Diisi saat create dari Memberships.StartedAt dan Memberships.ExpiredAt.
	// Dipakai access-check untuk double-check apakah siklus ini masih aktif
	// selain mengecek QuotaRemaining.
	CycleStartedAt time.Time  `json:"cycle_started_at" gorm:"not null"`
	CycleExpiredAt *time.Time `json:"cycle_expired_at,omitempty"`

	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`

	// BelongsTo Memberships — quota ini terikat ke satu siklus membership.
	// Load saat perlu invalidasi massal (misal: user upgrade di tengah siklus,
	// semua quota lama perlu digantikan quota baru sesuai plan baru).
	Membership Memberships `json:"membership,omitempty" gorm:"foreignKey:MembershipID;references:ID"`

	// BelongsTo Entitlement — load saat runtime access check untuk verifikasi
	// RestrictionType masih frequency_limited dan entitlement masih aktif.
	Entitlement Entitlement `json:"entitlement,omitempty" gorm:"foreignKey:EntitlementID;references:ID"`
}

func (UserEntitlementQuota) TableName() string {
	return "user_entitlement_quotas"
}

func (u *UserEntitlementQuota) BeforeCreate(tx *gorm.DB) error {
	now := time.Now().UTC()
	u.CreatedAt = now
	u.UpdatedAt = now
	return nil
}

// NewUserEntitlementQuota membuat quota record baru saat user berhasil beli membership.
// quotaGranted diambil dari DurationAccessMapping.UsageLimit untuk kombinasi plan+durasi yang dibeli.
// cycleExpiredAt diambil dari Memberships.ExpiredAt.
func NewUserEntitlementQuota(
	membershipID uuid.UUID,
	entitlementID uuid.UUID,
	entitlementKey string,
	quotaGranted int,
	cycleStartedAt time.Time,
	cycleExpiredAt *time.Time,
) UserEntitlementQuota {
	return UserEntitlementQuota{
		MembershipID:   membershipID,
		EntitlementID:  entitlementID,
		EntitlementKey: entitlementKey,
		QuotaGranted:   quotaGranted,
		QuotaUsed:      0,
		QuotaRemaining: quotaGranted,
		CycleStartedAt: cycleStartedAt,
		CycleExpiredAt: cycleExpiredAt,
	}
}

// ConsumeOne mengurangi sisa kuota sebesar 1.
// Dipanggil oleh access-check service setelah memverifikasi QuotaRemaining > 0.
// Selalu dipanggil dalam DB transaction yang sama dengan pembuatan UsageLog.
func (u *UserEntitlementQuota) ConsumeOne() {
	u.QuotaUsed++
	u.QuotaRemaining--
	u.UpdatedAt = time.Now().UTC()
}

// IsExhausted mengecek apakah kuota sudah habis.
func (u *UserEntitlementQuota) IsExhausted() bool {
	return u.QuotaRemaining <= 0
}

// IsActive mengecek apakah siklus quota ini masih berlaku saat ini.
func (u *UserEntitlementQuota) IsActive() bool {
	now := time.Now().UTC()
	if now.Before(u.CycleStartedAt) {
		return false
	}
	if u.CycleExpiredAt != nil && now.After(*u.CycleExpiredAt) {
		return false
	}
	return true
}
