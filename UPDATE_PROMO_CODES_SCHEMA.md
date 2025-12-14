# Pembaruan Skema: Promo Codes untuk Membership & Token Standalone
## Dokumentasi Perluasan Fitur Promo Code

---

## 1. Analisis Kebutuhan

### 1.1. Limitasi Skema Current

Skema `promo_codes` saat ini hanya mendukung:
- ✅ Diskon untuk pembelian **membership** (Basic/Pro/Max)
- ✅ Filter berdasarkan plan dan duration
- ❌ **TIDAK mendukung** diskon untuk pembelian **token standalone**

### 1.2. Kebutuhan Baru

Platform REXTRA perlu mendukung promo code untuk:
1. **Membership purchase** (sudah ada)
2. **Token standalone purchase** (belum ada)
3. **Kombinasi membership + bonus token** (future)

**Contoh Use Case:**
- Promo "Beli 50 token dapat diskon 20%"
- Promo "Beli token minimal 100, gratis 20 token bonus"
- Promo "Beli membership Pro dapat bonus 50 token"

---

## 2. Skema Baru: Promo Codes Universal

### 2.1. Tabel `promo_codes` (Updated)

**Perubahan:**
- Tambah kolom `promo_type` untuk distinguish membership vs token promo
- Tambah kolom `min_token_purchase` untuk minimum pembelian token
- Tambah kolom `bonus_token` untuk bonus token (jika applicable)
- Rename beberapa kolom agar lebih generic

| Nama Kolom | Tipe | Deskripsi Lengkap |
|------------|------|-------------------|
| id | UUID | ID unik promo (PK) |
| code | VARCHAR(50) | Kode promo (unique, case-insensitive) |
| description | TEXT | Deskripsi promo |
| **promo_type** | VARCHAR(20) | **NEW**: 'membership' / 'token' / 'both' |
| discount_type | VARCHAR(20) | 'percentage' / 'fixed_amount' / 'free_100' / 'bonus_token' |
| discount_value | DECIMAL(10,2) | Nilai diskon (% atau Rp) - nullable jika type = bonus_token |
| **bonus_token** | INTEGER | **NEW**: Jumlah bonus token (nullable) |
| applicable_plans | JSONB | Array plan_code yang berlaku (null = semua) - untuk promo membership |
| applicable_durations | JSONB | Array duration_months yang berlaku (null = semua) - untuk promo membership |
| **min_token_purchase** | INTEGER | **NEW**: Minimal pembelian token (nullable) - untuk promo token |
| max_usage | INTEGER | Maksimal penggunaan total (null = unlimited) |
| **max_usage_per_user** | INTEGER | **NEW**: Maksimal penggunaan per user (null = unlimited) |
| current_usage | INTEGER | Jumlah penggunaan saat ini |
| valid_from | TIMESTAMP | Waktu mulai berlaku |
| valid_until | TIMESTAMP | Waktu berakhir |
| is_active | BOOLEAN | Status promo aktif |
| created_at | TIMESTAMP | Waktu pembuatan |
| updated_at | TIMESTAMP | Waktu update terakhir |

### 2.2. SQL DDL (PostgreSQL) - Updated

```sql
CREATE TABLE promo_codes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code VARCHAR(50) UNIQUE NOT NULL,
  description TEXT,
  
  -- Promo type (NEW)
  promo_type VARCHAR(20) NOT NULL CHECK (promo_type IN ('membership', 'token', 'both')),
  
  -- Discount configuration
  discount_type VARCHAR(20) NOT NULL CHECK (discount_type IN ('percentage', 'fixed_amount', 'free_100', 'bonus_token')),
  discount_value DECIMAL(10,2), -- nullable jika discount_type = 'bonus_token'
  bonus_token INTEGER, -- NEW: untuk bonus token
  
  -- Applicability rules
  applicable_plans JSONB, -- untuk membership promo
  applicable_durations JSONB, -- untuk membership promo
  min_token_purchase INTEGER, -- NEW: untuk token promo
  
  -- Usage limits
  max_usage INTEGER, -- total usage limit
  max_usage_per_user INTEGER, -- NEW: per-user limit
  current_usage INTEGER DEFAULT 0,
  
  -- Validity period
  valid_from TIMESTAMP NOT NULL,
  valid_until TIMESTAMP NOT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  
  -- Timestamps
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_promo_code ON promo_codes(code);
CREATE INDEX idx_promo_type ON promo_codes(promo_type);
CREATE INDEX idx_promo_active ON promo_codes(is_active);
CREATE INDEX idx_promo_valid ON promo_codes(valid_from, valid_until);
```

---

### 2.3. Tabel Baru: `promo_code_usage` (Tracking per User)

Untuk track penggunaan promo per user (enforce `max_usage_per_user`):

| Nama Kolom | Tipe | Deskripsi |
|------------|------|-----------|
| id | UUID | ID unik (PK) |
| promo_code_id | UUID | FK ke promo_codes.id |
| user_id | UUID | FK ke users.id |
| payment_transaction_id | UUID | FK ke payment_transactions.id |
| used_at | TIMESTAMP | Waktu penggunaan |

**SQL DDL:**

```sql
CREATE TABLE promo_code_usage (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  promo_code_id UUID NOT NULL REFERENCES promo_codes(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  payment_transaction_id UUID NOT NULL REFERENCES payment_transactions(id) ON DELETE CASCADE,
  used_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_promo_usage_promo ON promo_code_usage(promo_code_id);
CREATE INDEX idx_promo_usage_user ON promo_code_usage(user_id);
CREATE INDEX idx_promo_usage_payment ON promo_code_usage(payment_transaction_id);
CREATE UNIQUE INDEX idx_promo_usage_unique ON promo_code_usage(promo_code_id, payment_transaction_id);
```

---

## 3. Contoh Data Seed (Updated)

### 3.1. Promo untuk Membership

```sql
-- Promo 100% gratis membership untuk early adopters
INSERT INTO promo_codes (
  code, description, promo_type, discount_type, discount_value, bonus_token,
  applicable_plans, applicable_durations, max_usage, max_usage_per_user,
  valid_from, valid_until
) VALUES (
  'REXTRA2025',
  'Promo spesial 100% gratis membership untuk mahasiswa baru di tahun 2025',
  'membership', -- promo type
  'free_100',
  100,
  NULL, -- no bonus token
  '["BASIC", "PRO", "MAX"]'::jsonb,
  '[1, 3, 6, 12]'::jsonb,
  1000, -- max 1000 total usage
  1, -- max 1x per user
  '2025-01-01 00:00:00',
  '2025-12-31 23:59:59'
);

-- Promo diskon 50% untuk paket Pro + bonus 20 token
INSERT INTO promo_codes (
  code, description, promo_type, discount_type, discount_value, bonus_token,
  applicable_plans, applicable_durations, max_usage, max_usage_per_user,
  valid_from, valid_until
) VALUES (
  'PRO50BONUS',
  'Diskon 50% untuk paket Pro durasi 6 dan 12 bulan + Bonus 20 token',
  'membership',
  'percentage',
  50,
  20, -- bonus 20 token
  '["PRO"]'::jsonb,
  '[6, 12]'::jsonb,
  500,
  1,
  '2025-01-01 00:00:00',
  '2025-03-31 23:59:59'
);
```

### 3.2. Promo untuk Token Standalone

```sql
-- Promo diskon 20% untuk pembelian token standalone
INSERT INTO promo_codes (
  code, description, promo_type, discount_type, discount_value, bonus_token,
  applicable_plans, applicable_durations, min_token_purchase,
  max_usage, max_usage_per_user, valid_from, valid_until
) VALUES (
  'TOKEN20OFF',
  'Diskon 20% untuk pembelian token standalone (minimal 50 token)',
  'token', -- promo type untuk token
  'percentage',
  20,
  NULL,
  NULL, -- tidak applicable untuk membership
  NULL, -- tidak applicable untuk membership
  50, -- minimal beli 50 token
  NULL, -- unlimited total usage
  3, -- max 3x per user
  '2025-01-01 00:00:00',
  '2025-12-31 23:59:59'
);

-- Promo bonus token (beli 100 dapat bonus 30)
INSERT INTO promo_codes (
  code, description, promo_type, discount_type, discount_value, bonus_token,
  applicable_plans, applicable_durations, min_token_purchase,
  max_usage, max_usage_per_user, valid_from, valid_until
) VALUES (
  'BONUS30TOKEN',
  'Beli 100 token dapat bonus 30 token gratis!',
  'token',
  'bonus_token', -- type khusus untuk bonus token
  NULL, -- no price discount
  30, -- bonus 30 token
  NULL,
  NULL,
  100, -- minimal beli 100 token
  NULL, -- unlimited
  5, -- max 5x per user
  '2025-02-01 00:00:00',
  '2025-02-28 23:59:59'
);

-- Promo fixed amount Rp10.000 off untuk token
INSERT INTO promo_codes (
  code, description, promo_type, discount_type, discount_value, bonus_token,
  applicable_plans, applicable_durations, min_token_purchase,
  max_usage, max_usage_per_user, valid_from, valid_until
) VALUES (
  'TOKEN10K',
  'Potongan Rp10.000 untuk pembelian token minimal 30 token',
  'token',
  'fixed_amount',
  10000,
  NULL,
  NULL,
  NULL,
  30, -- minimal 30 token
  1000,
  2,
  '2025-01-15 00:00:00',
  '2025-01-31 23:59:59'
);
```

### 3.3. Promo Universal (Both)

```sql
-- Promo universal untuk membership atau token
INSERT INTO promo_codes (
  code, description, promo_type, discount_type, discount_value, bonus_token,
  applicable_plans, applicable_durations, min_token_purchase,
  max_usage, max_usage_per_user, valid_from, valid_until
) VALUES (
  'NEWYEAR15',
  'Diskon 15% untuk pembelian membership atau token (minimal 20 token)',
  'both', -- berlaku untuk membership DAN token
  'percentage',
  15,
  NULL,
  '["BASIC", "PRO", "MAX"]'::jsonb, -- applicable plans
  '[3, 6, 12]'::jsonb, -- applicable durations
  20, -- minimal token jika beli token
  2000,
  1,
  '2025-01-01 00:00:00',
  '2025-01-15 23:59:59'
);
```

---

## 4. Relasi Antar Tabel (Updated)

```
users (existing table)
  ↓ (1:1)
memberships ────→ membership_plans (N:1)
            └───→ membership_durations (N:1)

users
  ↓ (1:N)
token_transactions ────→ token_usage_history (1:1 for usage type)

users
  ↓ (1:N)
poin_transactions

users
  ↓ (1:N)
payment_transactions ───→ membership_plans (N:1)
                     ├──→ membership_durations (N:1)
                     └──→ promo_codes (N:1, optional) ✅ UPDATED: mendukung membership & token

promo_codes
  ↓ (1:N)
promo_code_usage ───→ users (N:1)
                 └──→ payment_transactions (N:1)
```

**Penjelasan Relasi Baru:**
- `payment_transactions.promo_code` → Simpan kode promo yang digunakan (VARCHAR)
- `promo_code_usage` → Track usage per user untuk enforce `max_usage_per_user`
- Satu promo code bisa dipakai banyak user (1:N)
- Satu payment transaction hanya pakai 1 promo code (N:1)

---

## 5. Validasi Logic (Updated)

### 5.1. Validasi Promo Code untuk Membership

```go
func ValidatePromoCodeForMembership(
  promoCode string, 
  planCode string, 
  durationMonths int, 
  userID uuid.UUID,
) (*PromoCode, error) {
  
  // 1. Find promo code (case-insensitive)
  promo, err := promoRepo.GetByCode(strings.ToUpper(promoCode))
  if err != nil {
    return nil, errors.New("Promo code tidak ditemukan")
  }
  
  // 2. Check promo_type
  if promo.PromoType != "membership" && promo.PromoType != "both" {
    return nil, errors.New("Promo code tidak berlaku untuk membership")
  }
  
  // 3. Check is_active
  if !promo.IsActive {
    return nil, errors.New("Promo code tidak aktif")
  }
  
  // 4. Check valid date
  now := time.Now()
  if now.Before(promo.ValidFrom) || now.After(promo.ValidUntil) {
    return nil, errors.New("Promo code sudah expired atau belum berlaku")
  }
  
  // 5. Check max_usage (total)
  if promo.MaxUsage != nil && promo.CurrentUsage >= *promo.MaxUsage {
    return nil, errors.New("Promo code sudah mencapai batas penggunaan maksimal")
  }
  
  // 6. Check max_usage_per_user (NEW)
  if promo.MaxUsagePerUser != nil {
    usageCount := promoUsageRepo.CountByPromoAndUser(promo.ID, userID)
    if usageCount >= *promo.MaxUsagePerUser {
      return nil, errors.New(fmt.Sprintf("Anda sudah menggunakan promo ini %d kali (maksimal %d kali)", usageCount, *promo.MaxUsagePerUser))
    }
  }
  
  // 7. Check applicable plans
  if promo.ApplicablePlans != nil {
    var applicablePlans []string
    json.Unmarshal(promo.ApplicablePlans, &applicablePlans)
    
    if len(applicablePlans) > 0 && !contains(applicablePlans, strings.ToUpper(planCode)) {
      return nil, errors.New(fmt.Sprintf("Promo code tidak berlaku untuk paket %s", planCode))
    }
  }
  
  // 8. Check applicable durations
  if promo.ApplicableDurations != nil {
    var applicableDurations []int
    json.Unmarshal(promo.ApplicableDurations, &applicableDurations)
    
    if len(applicableDurations) > 0 && !containsInt(applicableDurations, durationMonths) {
      return nil, errors.New(fmt.Sprintf("Promo code tidak berlaku untuk durasi %d bulan", durationMonths))
    }
  }
  
  // 9. Promo valid
  return promo, nil
}
```

---

### 5.2. Validasi Promo Code untuk Token Standalone (NEW)

```go
func ValidatePromoCodeForToken(
  promoCode string,
  tokenQuantity int,
  userID uuid.UUID,
) (*PromoCode, error) {
  
  // 1. Find promo code
  promo, err := promoRepo.GetByCode(strings.ToUpper(promoCode))
  if err != nil {
    return nil, errors.New("Promo code tidak ditemukan")
  }
  
  // 2. Check promo_type (NEW)
  if promo.PromoType != "token" && promo.PromoType != "both" {
    return nil, errors.New("Promo code tidak berlaku untuk pembelian token")
  }
  
  // 3. Check is_active
  if !promo.IsActive {
    return nil, errors.New("Promo code tidak aktif")
  }
  
  // 4. Check valid date
  now := time.Now()
  if now.Before(promo.ValidFrom) || now.After(promo.ValidUntil) {
    return nil, errors.New("Promo code sudah expired atau belum berlaku")
  }
  
  // 5. Check max_usage (total)
  if promo.MaxUsage != nil && promo.CurrentUsage >= *promo.MaxUsage {
    return nil, errors.New("Promo code sudah mencapai batas penggunaan maksimal")
  }
  
  // 6. Check max_usage_per_user (NEW)
  if promo.MaxUsagePerUser != nil {
    usageCount := promoUsageRepo.CountByPromoAndUser(promo.ID, userID)
    if usageCount >= *promo.MaxUsagePerUser {
      return nil, errors.New(fmt.Sprintf("Anda sudah menggunakan promo ini %d kali (maksimal %d kali)", usageCount, *promo.MaxUsagePerUser))
    }
  }
  
  // 7. Check min_token_purchase (NEW)
  if promo.MinTokenPurchase != nil && tokenQuantity < *promo.MinTokenPurchase {
    return nil, errors.New(fmt.Sprintf("Minimal pembelian %d token untuk menggunakan promo ini", *promo.MinTokenPurchase))
  }
  
  // 8. Promo valid
  return promo, nil
}
```

---

### 5.3. Calculate Discount & Bonus (NEW)

```go
type PromoResult struct {
  PriceDiscount float64
  BonusToken    int
}

func CalculatePromoDiscount(promo *PromoCode, basePrice float64, tokenQuantity int) PromoResult {
  result := PromoResult{}
  
  switch promo.DiscountType {
  case "percentage":
    // Percentage discount
    result.PriceDiscount = basePrice * (promo.DiscountValue / 100)
    
  case "fixed_amount":
    // Fixed amount discount
    discount := promo.DiscountValue
    if discount > basePrice {
      discount = basePrice // tidak boleh lebih dari harga
    }
    result.PriceDiscount = discount
    
  case "free_100":
    // 100% free
    result.PriceDiscount = basePrice
    
  case "bonus_token":
    // No price discount, only bonus token
    result.PriceDiscount = 0
    if promo.BonusToken != nil {
      result.BonusToken = *promo.BonusToken
    }
  }
  
  // Add bonus token if applicable (even with price discount)
  if promo.BonusToken != nil && promo.DiscountType != "bonus_token" {
    result.BonusToken = *promo.BonusToken
  }
  
  return result
}
```

---

### 5.4. Record Promo Usage (NEW)

Setelah payment sukses, record promo usage:

```go
func RecordPromoUsage(promo *PromoCode, userID uuid.UUID, paymentID uuid.UUID) error {
  // 1. Create promo_code_usage record
  usage := &PromoCodeUsage{
    PromoCodeID:         promo.ID,
    UserID:              userID,
    PaymentTransactionID: paymentID,
  }
  
  if err := promoUsageRepo.Create(usage); err != nil {
    return err
  }
  
  // 2. Increment current_usage
  promo.CurrentUsage++
  if err := promoRepo.Update(promo); err != nil {
    return err
  }
  
  return nil
}
```

---

## 6. Contoh Use Case

### 6.1. User Beli Membership dengan Promo

```
User pilih: Pro 6 bulan (Rp102.000)
User input promo: "PRO50BONUS"

Backend:
1. Validate promo for membership ✓
2. Calculate discount:
   - Base: Rp102.000
   - Discount 50%: Rp51.000
   - Final: Rp51.000
   - Bonus token: 20
3. Create payment transaction
4. After payment success:
   - Activate membership Pro
   - Give 144 token (dari durasi) + 20 bonus token = 164 token
   - Record promo usage
```

### 6.2. User Beli Token Standalone dengan Promo

```
User pilih: 50 token (Rp50.000)
User input promo: "TOKEN20OFF"

Backend:
1. Validate promo for token ✓
   - Check min_token_purchase: 50 >= 50 ✓
2. Calculate discount:
   - Base: Rp50.000
   - Discount 20%: Rp10.000
   - Final: Rp40.000
3. Create payment transaction
4. After payment success:
   - Give 50 token
   - Activate Starter (if non-member)
   - Record promo usage
```

### 6.3. User Beli Token dengan Promo Bonus Token

```
User pilih: 100 token (Rp100.000)
User input promo: "BONUS30TOKEN"

Backend:
1. Validate promo for token ✓
   - Check min_token_purchase: 100 >= 100 ✓
2. Calculate:
   - Base: Rp100.000
   - Discount: Rp0 (no price discount)
   - Final: Rp100.000
   - Bonus token: 30
3. Create payment transaction
4. After payment success:
   - Give 100 token + 30 bonus = 130 token
   - Activate Starter (if non-member)
   - Record promo usage
```

---

## 7. Update Backend Handler

### 7.1. Purchase Membership (Updated)

```go
func (h *MembershipHandler) PurchaseMembership(c *gin.Context) {
  // ... (existing code)
  
  var promoDiscount float64
  var bonusToken int
  var promoApplied *PromoCode
  
  // Validate promo code if provided
  if req.PromoCode != "" {
    promo, err := h.promoService.ValidatePromoCodeForMembership(
      req.PromoCode, 
      req.PlanCode, 
      req.DurationMonths,
      userID,
    )
    if err != nil {
      c.JSON(400, response.Error(err.Error(), "INVALID_PROMO_CODE"))
      return
    }
    
    // Calculate promo discount & bonus
    pricing := h.calculatorService.CalculatePricing(plan, duration)
    promoResult := h.calculatorService.CalculatePromoDiscount(
      promo, 
      pricing.PriceAfterPlanDiscount,
      0, // not for token purchase
    )
    
    promoDiscount = promoResult.PriceDiscount
    bonusToken = promoResult.BonusToken
    promoApplied = promo
  }
  
  // ... (create payment, save promo_code, bonus_token)
  
  // In webhook after payment success:
  // - Activate membership dengan total_token + bonus_token
  // - Record promo usage
}
```

### 7.2. Purchase Token Standalone (Updated)

```go
func (h *TokenHandler) PurchaseTokenStandalone(c *gin.Context) {
  // ... (existing code)
  
  var promoDiscount float64
  var bonusToken int
  var promoApplied *PromoCode
  
  // Validate promo code if provided (NEW)
  if req.PromoCode != "" {
    promo, err := h.promoService.ValidatePromoCodeForToken(
      req.PromoCode,
      req.TokenQuantity,
      userID,
    )
    if err != nil {
      c.JSON(400, response.Error(err.Error(), "INVALID_PROMO_CODE"))
      return
    }
    
    // Calculate promo discount & bonus
    basePrice := float64(req.TokenQuantity) * tokenPricePerUnit
    promoResult := h.calculatorService.CalculatePromoDiscount(
      promo,
      basePrice,
      req.TokenQuantity,
    )
    
    promoDiscount = promoResult.PriceDiscount
    bonusToken = promoResult.BonusToken
    promoApplied = promo
  }
  
  // Calculate final price
  finalPrice := (float64(req.TokenQuantity) * tokenPricePerUnit) - promoDiscount
  
  // Create payment transaction with promo
  payment := &PaymentTransaction{
    // ... (existing fields)
    PromoCode: req.PromoCode,
    GrossAmount: float64(req.TokenQuantity) * tokenPricePerUnit,
    DiscountAmount: promoDiscount,
    FinalAmount: finalPrice,
  }
  
  // ... (rest of the code)
  
  // In webhook after payment success:
  // - Give req.TokenQuantity + bonusToken
  // - Record promo usage
}
```

---

## 8. Migration Script

```sql
-- Step 1: Backup existing promo_codes table
CREATE TABLE promo_codes_backup AS SELECT * FROM promo_codes;

-- Step 2: Add new columns to existing table
ALTER TABLE promo_codes 
ADD COLUMN promo_type VARCHAR(20) DEFAULT 'membership' NOT NULL,
ADD COLUMN bonus_token INTEGER,
ADD COLUMN min_token_purchase INTEGER,
ADD COLUMN max_usage_per_user INTEGER;

-- Step 3: Add constraint for promo_type
ALTER TABLE promo_codes
ADD CONSTRAINT check_promo_type CHECK (promo_type IN ('membership', 'token', 'both'));

-- Step 4: Update discount_type constraint to include 'bonus_token'
ALTER TABLE promo_codes
DROP CONSTRAINT IF EXISTS promo_codes_discount_type_check;

ALTER TABLE promo_codes
ADD CONSTRAINT promo_codes_discount_type_check 
CHECK (discount_type IN ('percentage', 'fixed_amount', 'free_100', 'bonus_token'));

-- Step 5: Create index for promo_type
CREATE INDEX idx_promo_type ON promo_codes(promo_type);

-- Step 6: Create promo_code_usage table
CREATE TABLE promo_code_usage (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  promo_code_id UUID NOT NULL REFERENCES promo_codes(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  payment_transaction_id UUID NOT NULL REFERENCES payment_transactions(id) ON DELETE CASCADE,
  used_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_promo_usage_promo ON promo_code_usage(promo_code_id);
CREATE INDEX idx_promo_usage_user ON promo_code_usage(user_id);
CREATE INDEX idx_promo_usage_payment ON promo_code_usage(payment_transaction_id);
CREATE UNIQUE INDEX idx_promo_usage_unique ON promo_code_usage(promo_code_id, payment_transaction_id);

-- Step 7: Verify migration
SELECT * FROM promo_codes LIMIT 5;
```

---

## 9. Summary Perubahan

### 9.1. Tabel Changes

| Tabel | Perubahan |
|-------|-----------|
| **promo_codes** | ✅ Tambah kolom: `promo_type`, `bonus_token`, `min_token_purchase`, `max_usage_per_user` |
| **promo_code_usage** | ✅ Tabel baru untuk tracking usage per user |
| **payment_transactions** | ⚠️ Sudah ada kolom `promo_code` (VARCHAR) - no change needed |

### 9.2. Feature Support

| Fitur | Before | After |
|-------|--------|-------|
| Promo untuk membership | ✅ | ✅ |
| Promo untuk token standalone | ❌ | ✅ |
| Promo universal (both) | ❌ | ✅ |
| Bonus token | ❌ | ✅ |
| Per-user usage limit | ❌ | ✅ |
| Minimum token purchase | ❌ | ✅ |

### 9.3. Validation Logic

| Validation | Before | After |
|------------|--------|-------|
| Check promo type | ❌ | ✅ |
| Check min token purchase | ❌ | ✅ |
| Check per-user limit | ❌ | ✅ |
| Calculate bonus token | ❌ | ✅ |

---

## 10. Kesimpulan

Dengan skema baru ini, sistem promo code REXTRA sekarang:

✅ **Mendukung 3 tipe promo:**
- Membership only
- Token standalone only  
- Universal (both)

✅ **Mendukung 4 tipe diskon:**
- Percentage discount
- Fixed amount discount
- Free 100%
- Bonus token

✅ **Fitur tambahan:**
- Minimum token purchase requirement
- Per-user usage limit
- Tracking usage per user
- Bonus token untuk membership purchase

✅ **Backward compatible:**
- Existing membership promo tetap work
- Migration script provided
- No breaking changes

**Skema ini siap untuk production!** 🚀

---

**Prepared by**: Tim Backend REXTRA  
**Document Version**: 1.0  
**Last Updated**: November 2025