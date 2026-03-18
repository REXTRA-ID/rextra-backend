package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Memberships struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null;uniqueIndex"` // 1 user = 1 row

	// Status plan saat ini
	PlanID   *uuid.UUID   `json:"plan_id,omitempty" gorm:"type:uuid;index"`
	PlanName EnumPlanName `json:"plan_name" gorm:"type:varchar(20);not null"` // snapshot nama plan

	// Durasi aktif saat ini
	DurationID     *uuid.UUID `json:"duration_id,omitempty" gorm:"type:uuid;index"`
	DurationMonths *int       `json:"duration_months,omitempty"` // snapshot durasi bulan

	// Periode aktif
	StartedAt *time.Time `json:"started_at,omitempty"`
	ExpiredAt *time.Time `json:"expired_at,omitempty"`

	// Status
	IsActive  bool `json:"is_active" gorm:"not null;default:true"`
	AutoRenew bool `json:"auto_renew" gorm:"not null;default:false"`

	// Saldo snapshot — source of truth tetap di TokenWallet (Hamzah) dan PointsLedger,
	// ini hanya untuk kemudahan query tanpa join ke tabel lain di setiap request
	CurrentTokenBalance int `json:"current_token_balance" gorm:"not null;default:0"`
	CurrentPoinBalance  int `json:"current_poin_balance" gorm:"not null;default:0"`

	// Statistik
	PaidCycleCount   int `json:"paid_cycle_count" gorm:"not null;default:0"`  // berapa kali bayar
	EntitlementCount int `json:"entitlement_count" gorm:"not null;default:0"` // berapa entitlement aktif

	Timestamp

	// BelongsTo MembershipPlans — membership user selalu merujuk ke satu plan aktif.
	// Load relasi ini saat menampilkan halaman profil membership user atau saat admin
	// melihat detail user. Juga dibutuhkan saat kalkulasi upgrade/downgrade untuk
	// mengetahui kategori plan (paid/unpaid) dan konfigurasi lainnya.
	Plan *MembershipPlans `json:"plan,omitempty" gorm:"foreignKey:PlanID;references:ID"`

	// BelongsTo MembershipDuration — membership user merujuk ke konfigurasi durasi yang aktif.
	// Dibutuhkan saat user ingin upgrade/downgrade: sistem perlu memanggil
	// duration.CalculateCredit(remainingDays) untuk menghitung kredit sisa durasi.
	Duration *MembershipDuration `json:"duration,omitempty" gorm:"foreignKey:DurationID;references:ID"`

	// HasMany SubscriptionCycle — histori semua siklus berlangganan user.
	// Load relasi ini saat menampilkan tab "Riwayat Berlangganan" di halaman profil user.
	SubscriptionCycles []SubscriptionCycle `json:"subscription_cycles,omitempty" gorm:"foreignKey:MembershipID;references:ID"`

	// HasMany PaymentTransactions — histori semua transaksi pembayaran user.
	// Load relasi ini saat menampilkan tab "Riwayat Tagihan & Invoice" di halaman profil user.
	PaymentTransactions []PaymentTransactions `json:"payment_transactions,omitempty" gorm:"foreignKey:MembershipID;references:ID"`

	// HasMany PoinTransactions — histori mutasi poin user (earn dan spend).
	// Load relasi ini saat menampilkan tab "Riwayat Poin" di halaman profil user.
	PoinTransactions []PoinTransactions `json:"poin_transactions,omitempty" gorm:"foreignKey:MembershipID;references:ID"`

	// HasMany PointsLedger — ledger poin detail untuk rekonsiliasi saldo.
	// Digunakan oleh background job untuk memastikan CurrentPoinBalance konsisten
	// dengan total sum entry di ledger ini.
	// PointsLedger []PointsLedger `json:"points_ledger,omitempty" gorm:"foreignKey:MembershipID;references:ID"`

	// HasMany UsageLog — log semua aktivitas akses fitur user.
	// Load relasi ini saat admin ingin audit aktivitas user tertentu
	// atau debugging kenapa akses user ditolak.
	// UsageLogs []UsageLog `json:"usage_logs,omitempty" gorm:"foreignKey:MembershipID;references:ID"`

	// HasMany UserEntitlementQuota — semua quota record milik membership ini.
	// Load saat payment callback untuk buat quota baru, atau saat upgrade/downgrade
	// untuk invalidasi quota lama sebelum digantikan quota baru.
	// Jangan preload ini secara default — hanya load saat dibutuhkan secara eksplisit.
	UserEntitlementQuotas []UserEntitlementQuota `json:"user_entitlement_quotas,omitempty" gorm:"foreignKey:MembershipID;references:ID"`
}

func (Memberships) TableName() string {
	return "memberships"
}

func (m *Memberships) BeforeCreate(tx *gorm.DB) error {
	now := time.Now().UTC()
	m.CreatedAt = now
	m.UpdatedAt = now
	return nil
}

func NewMembership(userID uuid.UUID, planName EnumPlanName) Memberships {
	return Memberships{
		UserID:   userID,
		PlanName: planName,
		IsActive: true,
	}
}

// Activate mengaktifkan membership dengan plan dan durasi baru
func (m *Memberships) Activate(plan *MembershipPlans, duration *MembershipDuration) {
	now := time.Now().UTC()
	expired := now.AddDate(0, duration.DurationMonth, 0)

	m.PlanID = &plan.ID
	m.PlanName = plan.PlanName
	m.DurationID = &duration.ID
	m.DurationMonths = &duration.DurationMonth
	m.StartedAt = &now
	m.ExpiredAt = &expired
	m.IsActive = true
	m.PaidCycleCount++
	m.UpdatedAt = now
}

// Deactivate menurunkan membership ke free plan
func (m *Memberships) Deactivate(fallbackPlanID uuid.UUID, fallbackPlanName EnumPlanName) {
	now := time.Now().UTC()
	m.PlanID = &fallbackPlanID
	m.PlanName = fallbackPlanName
	m.DurationID = nil
	m.DurationMonths = nil
	m.StartedAt = nil
	m.ExpiredAt = nil
	m.IsActive = false
	m.UpdatedAt = now
}

// ExpiringMembership mengubah status membership yang sudah expired menjadi plan non-member
func (m *Memberships) ExpiringMembership(plan *MembershipPlans) {
	now := time.Now().UTC()
	m.PlanID = &plan.ID
	m.PlanName = plan.PlanName
	m.DurationID = nil
	m.DurationMonths = nil
	m.StartedAt = nil
	m.ExpiredAt = nil
	m.IsActive = false
	m.UpdatedAt = now
}

// RemainingDays menghitung sisa hari sebelum expired
func (m *Memberships) RemainingDays() int {
	if m.ExpiredAt == nil {
		return 0
	}
	remaining := time.Until(*m.ExpiredAt).Hours() / 24
	if remaining < 0 {
		return 0
	}
	return int(remaining)
}

// IsEligibleForChangeOfPlan mengecek apakah membership masih aktif dan bisa diubah
func (m *Memberships) IsEligibleForChangeOfPlan() bool {
	if !m.IsActive || m.ExpiredAt == nil {
		return false
	}
	return time.Now().Before(*m.ExpiredAt)
}

// AddPoinBalance menambah saldo poin snapshot
func (m *Memberships) AddPoinBalance(amount int) {
	m.CurrentPoinBalance += amount
	m.UpdatedAt = time.Now().UTC()
}

// DeductPoinBalance mengurangi saldo poin snapshot
func (m *Memberships) DeductPoinBalance(amount int) {
	m.CurrentPoinBalance -= amount
	if m.CurrentPoinBalance < 0 {
		m.CurrentPoinBalance = 0
	}
	m.UpdatedAt = time.Now().UTC()
}

// UpdateMembership memperbarui status plan membership dan durasi
func (m *Memberships) UpdateMembership(plan *MembershipPlans, duration *MembershipDuration) {
	now := time.Now().UTC()
	expired := now.AddDate(0, duration.DurationMonth, 0)

	m.PlanID = &plan.ID
	m.PlanName = plan.PlanName
	m.DurationID = &duration.ID
	m.DurationMonths = &duration.DurationMonth
	m.StartedAt = &now
	m.ExpiredAt = &expired
	m.IsActive = true
	m.UpdatedAt = now
}

// CalculateTotalToken mengalkulasi jumlah token yang didapat berdasarkan plan dan bonus
func (m *Memberships) CalculateTotalToken(bonusToken *int) int {
	if m.Plan == nil || m.Duration == nil {
		return 0 // fallback if relations not loaded
	}
	baseToken := m.Plan.MonthlyToken * m.Duration.DurationMonth
	bonusFromDuration := float64(baseToken) * m.Duration.TokenBonusPercentage
	totalToken := baseToken + int(bonusFromDuration)

	if bonusToken != nil {
		totalToken += *bonusToken
	}

	return totalToken
}

// CalculateRextraPoin mengalkulasi tambahan rextra poin
func (m *Memberships) CalculateRextraPoin() int {
	if m.Duration == nil {
		return 0
	}
	// Assumption: base point per successful checkout is 10, times multiplier
	basePoint := 10
	return basePoint * m.Duration.RextraPoinMultiplier
}

// AddMembershipToken menambah balance token pada membership
func (m *Memberships) AddMembershipToken(amount int) {
	m.CurrentTokenBalance += amount
	m.UpdatedAt = time.Now().UTC()
}

// RefillToken menambah balance token sesuai dengan monthly token plan
func (m *Memberships) RefillToken() int {
	refillAmount := 0
	if m.Plan != nil {
		refillAmount = m.Plan.MonthlyToken
	}
	m.CurrentTokenBalance += refillAmount
	m.UpdatedAt = time.Now().UTC()
	return refillAmount
}

// UseMembershipToken mengurangi balance token pada membership
func (m *Memberships) UseMembershipToken(amount int) {
	m.CurrentTokenBalance -= amount
	m.UpdatedAt = time.Now().UTC()
}
