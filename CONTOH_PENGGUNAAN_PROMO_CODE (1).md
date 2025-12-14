# Contoh Penggunaan Promo Code: All Scenarios
## Panduan Lengkap Implementasi Promo Code di REXTRA

---

## Daftar Isi

1. [Promo Membership](#1-promo-membership)
2. [Promo Token Standalone](#2-promo-token-standalone)
3. [Promo Universal (Both)](#3-promo-universal-both)
4. [Kasus Edge Cases & Validasi](#4-kasus-edge-cases--validasi)
5. [Flow Diagram](#5-flow-diagram)

---

## 1. Promo Membership

### Scenario 1.1: Diskon Percentage untuk Membership

**Promo Data:**
```sql
INSERT INTO promo_codes (
  code, description, promo_type, discount_type, discount_value, bonus_token,
  applicable_plans, applicable_durations, max_usage, max_usage_per_user,
  valid_from, valid_until
) VALUES (
  'PRO50OFF',
  'Diskon 50% untuk paket Pro durasi 6 dan 12 bulan',
  'membership',
  'percentage',
  50,
  NULL,
  '["PRO"]'::jsonb,
  '[6, 12]'::jsonb,
  500,
  1,
  '2025-01-01 00:00:00',
  '2025-03-31 23:59:59'
);
```

**User Action:**
```
User memilih:
- Plan: Pro
- Duration: 6 bulan
- Promo code: PRO50OFF
```

**Backend Calculation:**
```go
// Step 1: Get plan & duration
plan := getPlanByCode("PRO")        // monthly_token: 20, base_price: 20000
duration := getDurationByMonths(6)  // discount: 15%, token_bonus: 20%

// Step 2: Calculate base pricing
basePrice = 20000 × 6 = 120,000
planDiscount = 120,000 × 15% = 18,000
priceAfterPlanDiscount = 120,000 - 18,000 = 102,000

// Step 3: Apply promo
promoDiscount = 102,000 × 50% = 51,000
finalPrice = 102,000 - 51,000 = 51,000

// Step 4: Calculate tokens
baseToken = 20 × 6 = 120
bonusToken = 120 × 20% = 24
totalToken = 120 + 24 = 144

// Step 5: Calculate REXTRA POIN
rextraPoin = (20 × 6) × 1 = 120
```

**Result:**
```json
{
  "plan": "Pro",
  "duration": "6 bulan",
  "base_price": 120000,
  "plan_discount": 18000,
  "promo_discount": 51000,
  "final_price": 51000,
  "total_token": 144,
  "rextra_poin": 120,
  "savings": 69000,
  "savings_percentage": 57.5
}
```

**Database Records Created:**
```sql
-- payment_transactions
INSERT INTO payment_transactions (
  user_id, payment_type, plan_id, duration_id,
  gross_amount, discount_amount, final_amount, promo_code
) VALUES (
  'user-uuid', 'membership', 'pro-plan-id', '6month-duration-id',
  120000, 69000, 51000, 'PRO50OFF'
);

-- After payment success:
-- memberships: status=pro, token_balance+=144, poin_balance+=120
-- token_transactions: type=purchase_membership, amount=+144
-- poin_transactions: type=earn, amount=+120
-- promo_code_usage: record usage
-- promo_codes: current_usage += 1
```

---

### Scenario 1.2: Fixed Amount Discount untuk Membership

**Promo Data:**
```sql
INSERT INTO promo_codes (
  code, description, promo_type, discount_type, discount_value, bonus_token,
  applicable_plans, applicable_durations, max_usage, max_usage_per_user,
  valid_from, valid_until
) VALUES (
  'DISKON20K',
  'Potongan harga Rp20.000 untuk semua paket',
  'membership',
  'fixed_amount',
  20000,
  NULL,
  NULL, -- berlaku untuk semua plan
  NULL, -- berlaku untuk semua duration
  NULL, -- unlimited usage
  2,    -- max 2x per user
  '2025-01-01 00:00:00',
  '2025-12-31 23:59:59'
);
```

**User Action:**
```
User memilih:
- Plan: Basic
- Duration: 3 bulan
- Promo code: DISKON20K
```

**Backend Calculation:**
```go
// Step 1: Get plan & duration
plan := getPlanByCode("BASIC")      // monthly_token: 10, base_price: 15000
duration := getDurationByMonths(3)  // discount: 8%, token_bonus: 10%

// Step 2: Calculate base pricing
basePrice = 15000 × 3 = 45,000
planDiscount = 45,000 × 8% = 3,600
priceAfterPlanDiscount = 45,000 - 3,600 = 41,400

// Step 3: Apply promo (fixed amount)
promoDiscount = 20,000
if (promoDiscount > priceAfterPlanDiscount) {
  promoDiscount = priceAfterPlanDiscount  // max = price itself
}
finalPrice = 41,400 - 20,000 = 21,400

// Step 4: Calculate tokens
baseToken = 10 × 3 = 30
bonusToken = 30 × 10% = 3
totalToken = 30 + 3 = 33

// Step 5: Calculate REXTRA POIN
rextraPoin = (10 × 3) × 1 = 30
```

**Result:**
```json
{
  "plan": "Basic",
  "duration": "3 bulan",
  "base_price": 45000,
  "plan_discount": 3600,
  "promo_discount": 20000,
  "final_price": 21400,
  "total_token": 33,
  "rextra_poin": 30,
  "savings": 23600,
  "savings_percentage": 52.4
}
```

---

### Scenario 1.3: Free 100% untuk Membership

**Promo Data:**
```sql
INSERT INTO promo_codes (
  code, description, promo_type, discount_type, discount_value, bonus_token,
  applicable_plans, applicable_durations, max_usage, max_usage_per_user,
  valid_from, valid_until
) VALUES (
  'GRATIS2025',
  'Promo 100% GRATIS membership untuk early adopters',
  'membership',
  'free_100',
  100,
  NULL,
  '["BASIC", "PRO", "MAX"]'::jsonb,
  '[1, 3, 6, 12]'::jsonb,
  1000,
  1,
  '2025-01-01 00:00:00',
  '2025-12-31 23:59:59'
);
```

**User Action:**
```
User memilih:
- Plan: Max
- Duration: 1 bulan
- Promo code: GRATIS2025
```

**Backend Calculation:**
```go
// Step 1: Get plan & duration
plan := getPlanByCode("MAX")        // monthly_token: 40, base_price: 30000
duration := getDurationByMonths(1)  // no discount, no bonus

// Step 2: Calculate base pricing
basePrice = 30000 × 1 = 30,000
planDiscount = 0
priceAfterPlanDiscount = 30,000

// Step 3: Apply promo (free 100%)
promoDiscount = 30,000  // 100% of price
finalPrice = 30,000 - 30,000 = 0

// Step 4: Calculate tokens
baseToken = 40 × 1 = 40
bonusToken = 0
totalToken = 40

// Step 5: Calculate REXTRA POIN
rextraPoin = (40 × 1) × 1 = 40
```

**Result:**
```json
{
  "plan": "Max",
  "duration": "1 bulan",
  "base_price": 30000,
  "plan_discount": 0,
  "promo_discount": 30000,
  "final_price": 0,
  "total_token": 40,
  "rextra_poin": 40,
  "savings": 30000,
  "savings_percentage": 100,
  "message": "🎉 Selamat! Anda mendapat membership GRATIS!"
}
```

**Important Note:**
- Xendit tetap perlu create invoice meskipun amount = 0
- Atau bisa skip Xendit dan langsung activate membership
- Untuk 100% free, recommended skip payment gateway

---

### Scenario 1.4: Membership dengan Bonus Token

**Promo Data:**
```sql
INSERT INTO promo_codes (
  code, description, promo_type, discount_type, discount_value, bonus_token,
  applicable_plans, applicable_durations, max_usage, max_usage_per_user,
  valid_from, valid_until
) VALUES (
  'PRO30BONUS',
  'Paket Pro 12 bulan diskon 30% + Bonus 50 token',
  'membership',
  'percentage',
  30,
  50, -- bonus 50 token
  '["PRO"]'::jsonb,
  '[12]'::jsonb,
  200,
  1,
  '2025-02-01 00:00:00',
  '2025-02-28 23:59:59'
);
```

**User Action:**
```
User memilih:
- Plan: Pro
- Duration: 12 bulan
- Promo code: PRO30BONUS
```

**Backend Calculation:**
```go
// Step 1: Get plan & duration
plan := getPlanByCode("PRO")         // monthly_token: 20, base_price: 20000
duration := getDurationByMonths(12)  // discount: 30%, token_bonus: 40%

// Step 2: Calculate base pricing
basePrice = 20000 × 12 = 240,000
planDiscount = 240,000 × 30% = 72,000
priceAfterPlanDiscount = 240,000 - 72,000 = 168,000

// Step 3: Apply promo (30% discount + 50 bonus token)
promoDiscount = 168,000 × 30% = 50,400
finalPrice = 168,000 - 50,400 = 117,600

// Step 4: Calculate tokens
baseToken = 20 × 12 = 240
bonusFromDuration = 240 × 40% = 96
bonusFromPromo = 50
totalToken = 240 + 96 + 50 = 386

// Step 5: Calculate REXTRA POIN
rextraPoin = (20 × 12) × 1 = 240
```

**Result:**
```json
{
  "plan": "Pro",
  "duration": "12 bulan",
  "base_price": 240000,
  "plan_discount": 72000,
  "promo_discount": 50400,
  "final_price": 117600,
  "total_token": 386,
  "bonus_token_from_promo": 50,
  "rextra_poin": 240,
  "savings": 122400,
  "savings_percentage": 51,
  "special_bonus": "🎁 Bonus 50 token dari promo!"
}
```

---

## 2. Promo Token Standalone

### Scenario 2.1: Diskon Percentage untuk Token

**Promo Data:**
```sql
INSERT INTO promo_codes (
  code, description, promo_type, discount_type, discount_value, bonus_token,
  applicable_plans, applicable_durations, min_token_purchase,
  max_usage, max_usage_per_user, valid_from, valid_until
) VALUES (
  'TOKEN20OFF',
  'Diskon 20% untuk pembelian token standalone (minimal 50 token)',
  'token',
  'percentage',
  20,
  NULL,
  NULL, -- not for membership
  NULL, -- not for membership
  50,   -- minimal 50 token
  NULL, -- unlimited total
  3,    -- max 3x per user
  '2025-01-01 00:00:00',
  '2025-12-31 23:59:59'
);
```

**User Action:**
```
User memilih:
- Token quantity: 50 token
- Promo code: TOKEN20OFF
```

**Backend Calculation:**
```go
const tokenPricePerUnit = 1000 // Rp1.000 per token

// Step 1: Validate
if (tokenQuantity < 50) {
  return error("Minimal pembelian 50 token untuk promo ini")
}

// Step 2: Calculate base price
basePrice = 50 × 1000 = 50,000

// Step 3: Apply promo
promoDiscount = 50,000 × 20% = 10,000
finalPrice = 50,000 - 10,000 = 40,000

// Step 4: Tokens
totalToken = 50
bonusToken = 0 // no bonus token for this promo
```

**Result:**
```json
{
  "token_quantity": 50,
  "token_price_per_unit": 1000,
  "base_price": 50000,
  "promo_discount": 10000,
  "final_price": 40000,
  "total_token": 50,
  "bonus_token": 0,
  "savings": 10000,
  "savings_percentage": 20,
  "extra_benefit": "✅ Dapat Starter Plan 30 hari (jika Non-Member)"
}
```

**Database Records Created:**
```sql
-- payment_transactions
INSERT INTO payment_transactions (
  user_id, payment_type, token_quantity,
  gross_amount, discount_amount, final_amount, promo_code
) VALUES (
  'user-uuid', 'token_standalone', 50,
  50000, 10000, 40000, 'TOKEN20OFF'
);

-- After payment success:
-- memberships: token_balance += 50
-- IF user is non_member: membership_status = 'starter', expired_at = NOW() + 30 days
-- token_transactions: type=purchase_standalone, amount=+50
-- promo_code_usage: record usage
```

---

### Scenario 2.2: Fixed Amount Discount untuk Token

**Promo Data:**
```sql
INSERT INTO promo_codes (
  code, description, promo_type, discount_type, discount_value, bonus_token,
  applicable_plans, applicable_durations, min_token_purchase,
  max_usage, max_usage_per_user, valid_from, valid_until
) VALUES (
  'TOKEN15K',
  'Potongan Rp15.000 untuk pembelian token minimal 30 token',
  'token',
  'fixed_amount',
  15000,
  NULL,
  NULL,
  NULL,
  30,
  1000,
  2,
  '2025-01-15 00:00:00',
  '2025-01-31 23:59:59'
);
```

**User Action:**
```
User memilih:
- Token quantity: 100 token
- Promo code: TOKEN15K
```

**Backend Calculation:**
```go
// Step 1: Validate
if (tokenQuantity < 30) {
  return error("Minimal pembelian 30 token untuk promo ini")
}

// Step 2: Calculate base price
basePrice = 100 × 1000 = 100,000

// Step 3: Apply promo (fixed amount)
promoDiscount = 15,000
if (promoDiscount > basePrice) {
  promoDiscount = basePrice
}
finalPrice = 100,000 - 15,000 = 85,000

// Step 4: Tokens
totalToken = 100
bonusToken = 0
```

**Result:**
```json
{
  "token_quantity": 100,
  "token_price_per_unit": 1000,
  "base_price": 100000,
  "promo_discount": 15000,
  "final_price": 85000,
  "total_token": 100,
  "bonus_token": 0,
  "savings": 15000,
  "savings_percentage": 15
}
```

---

### Scenario 2.3: Bonus Token (No Price Discount)

**Promo Data:**
```sql
INSERT INTO promo_codes (
  code, description, promo_type, discount_type, discount_value, bonus_token,
  applicable_plans, applicable_durations, min_token_purchase,
  max_usage, max_usage_per_user, valid_from, valid_until
) VALUES (
  'BONUS30TOKEN',
  'Beli 100 token dapat bonus 30 token gratis!',
  'token',
  'bonus_token',
  NULL, -- no price discount
  30,   -- bonus 30 token
  NULL,
  NULL,
  100,  -- minimal 100 token
  NULL,
  5,
  '2025-02-01 00:00:00',
  '2025-02-28 23:59:59'
);
```

**User Action:**
```
User memilih:
- Token quantity: 100 token
- Promo code: BONUS30TOKEN
```

**Backend Calculation:**
```go
// Step 1: Validate
if (tokenQuantity < 100) {
  return error("Minimal pembelian 100 token untuk bonus ini")
}

// Step 2: Calculate base price
basePrice = 100 × 1000 = 100,000

// Step 3: Apply promo (bonus_token type - NO price discount)
promoDiscount = 0
finalPrice = 100,000

// Step 4: Tokens
purchasedToken = 100
bonusToken = 30
totalToken = 100 + 30 = 130
```

**Result:**
```json
{
  "token_quantity": 100,
  "token_price_per_unit": 1000,
  "base_price": 100000,
  "promo_discount": 0,
  "final_price": 100000,
  "purchased_token": 100,
  "bonus_token": 30,
  "total_token": 130,
  "savings": 0,
  "special_bonus": "🎁 Bonus 30 token GRATIS!",
  "effective_token_price": 769,
  "message": "Bayar 100 token, dapat 130 token!"
}
```

---

### Scenario 2.4: Token dengan Price Discount + Bonus Token

**Promo Data:**
```sql
INSERT INTO promo_codes (
  code, description, promo_type, discount_type, discount_value, bonus_token,
  applicable_plans, applicable_durations, min_token_purchase,
  max_usage, max_usage_per_user, valid_from, valid_until
) VALUES (
  'TOKENSUPER',
  'SUPER DEAL! Diskon 25% + Bonus 20 token untuk pembelian 80 token',
  'token',
  'percentage',
  25,   -- 25% discount
  20,   -- bonus 20 token
  NULL,
  NULL,
  80,
  100,
  1,
  '2025-03-01 00:00:00',
  '2025-03-31 23:59:59'
);
```

**User Action:**
```
User memilih:
- Token quantity: 80 token
- Promo code: TOKENSUPER
```

**Backend Calculation:**
```go
// Step 1: Validate
if (tokenQuantity < 80) {
  return error("Minimal pembelian 80 token untuk promo ini")
}

// Step 2: Calculate base price
basePrice = 80 × 1000 = 80,000

// Step 3: Apply promo (25% discount + 20 bonus token)
promoDiscount = 80,000 × 25% = 20,000
finalPrice = 80,000 - 20,000 = 60,000

// Step 4: Tokens
purchasedToken = 80
bonusToken = 20
totalToken = 80 + 20 = 100
```

**Result:**
```json
{
  "token_quantity": 80,
  "token_price_per_unit": 1000,
  "base_price": 80000,
  "promo_discount": 20000,
  "final_price": 60000,
  "purchased_token": 80,
  "bonus_token": 20,
  "total_token": 100,
  "savings": 20000,
  "savings_percentage": 25,
  "special_bonus": "🎁 Bonus 20 token + Diskon 25%!",
  "effective_token_price": 600,
  "message": "Bayar Rp60.000, dapat 100 token!"
}
```

---

## 3. Promo Universal (Both)

### Scenario 3.1: Universal Promo untuk Membership

**Promo Data:**
```sql
INSERT INTO promo_codes (
  code, description, promo_type, discount_type, discount_value, bonus_token,
  applicable_plans, applicable_durations, min_token_purchase,
  max_usage, max_usage_per_user, valid_from, valid_until
) VALUES (
  'NEWYEAR15',
  'Diskon 15% untuk pembelian membership atau token (minimal 20 token)',
  'both',
  'percentage',
  15,
  NULL,
  '["BASIC", "PRO", "MAX"]'::jsonb,
  '[3, 6, 12]'::jsonb,
  20, -- minimal token jika beli token
  2000,
  1,
  '2025-01-01 00:00:00',
  '2025-01-15 23:59:59'
);
```

**User Action (Membership):**
```
User memilih:
- Plan: Basic
- Duration: 6 bulan
- Promo code: NEWYEAR15
```

**Backend Calculation:**
```go
// Step 1: Validate promo for membership
promo.PromoType == "both" ✓
plan == "BASIC" in applicable_plans ✓
duration == 6 in applicable_durations ✓

// Step 2: Calculate
basePrice = 15000 × 6 = 90,000
planDiscount = 90,000 × 15% = 13,500
priceAfterPlanDiscount = 76,500

promoDiscount = 76,500 × 15% = 11,475
finalPrice = 76,500 - 11,475 = 65,025

totalToken = 10 × 6 × 1.2 = 72
rextraPoin = 60
```

**Result:**
```json
{
  "plan": "Basic",
  "duration": "6 bulan",
  "base_price": 90000,
  "plan_discount": 13500,
  "promo_discount": 11475,
  "final_price": 65025,
  "total_token": 72,
  "rextra_poin": 60,
  "savings": 24975
}
```

---

### Scenario 3.2: Universal Promo untuk Token

**User Action (Token):**
```
User memilih:
- Token quantity: 50 token
- Promo code: NEWYEAR15
```

**Backend Calculation:**
```go
// Step 1: Validate promo for token
promo.PromoType == "both" ✓
tokenQuantity >= min_token_purchase (50 >= 20) ✓

// Step 2: Calculate
basePrice = 50 × 1000 = 50,000
promoDiscount = 50,000 × 15% = 7,500
finalPrice = 50,000 - 7,500 = 42,500

totalToken = 50
bonusToken = 0
```

**Result:**
```json
{
  "token_quantity": 50,
  "token_price_per_unit": 1000,
  "base_price": 50000,
  "promo_discount": 7500,
  "final_price": 42500,
  "total_token": 50,
  "bonus_token": 0,
  "savings": 7500,
  "savings_percentage": 15
}
```

---

## 4. Kasus Edge Cases & Validasi

### Case 4.1: Promo Code Sudah Expired

**Scenario:**
```
User input: PROMO2024
Valid until: 2024-12-31 23:59:59
Current date: 2025-01-15
```

**Backend Validation:**
```go
now := time.Now() // 2025-01-15
if now.After(promo.ValidUntil) {
  return error("Promo code sudah expired pada 31 Dec 2024")
}
```

**Response:**
```json
{
  "status": "error",
  "message": "Promo code sudah expired",
  "error_code": "PROMO_EXPIRED",
  "data": {
    "code": "PROMO2024",
    "expired_at": "2024-12-31T23:59:59Z",
    "days_ago": 15
  }
}
```

---

### Case 4.2: Promo Code Belum Aktif

**Scenario:**
```
User input: MARCH2025
Valid from: 2025-03-01 00:00:00
Current date: 2025-02-20
```

**Backend Validation:**
```go
now := time.Now() // 2025-02-20
if now.Before(promo.ValidFrom) {
  return error(fmt.Sprintf("Promo code akan aktif mulai %s", promo.ValidFrom.Format("02 Jan 2006")))
}
```

**Response:**
```json
{
  "status": "error",
  "message": "Promo code belum aktif",
  "error_code": "PROMO_NOT_STARTED",
  "data": {
    "code": "MARCH2025",
    "valid_from": "2025-03-01T00:00:00Z",
    "days_remaining": 9
  }
}
```

---

### Case 4.3: Max Usage Tercapai (Total)

**Scenario:**
```
Promo: max_usage = 1000
Current usage: 1000
User mencoba pakai: PRO50OFF
```

**Backend Validation:**
```go
if promo.MaxUsage != nil && promo.CurrentUsage >= *promo.MaxUsage {
  return error("Promo code sudah mencapai batas penggunaan maksimal")
}
```

**Response:**
```json
{
  "status": "error",
  "message": "Promo code sudah habis",
  "error_code": "PROMO_MAX_USAGE_REACHED",
  "data": {
    "code": "PRO50OFF",
    "max_usage": 1000,
    "current_usage": 1000,
    "message": "Promo ini sudah digunakan oleh 1000 pengguna"
  }
}
```

---

### Case 4.4: Max Usage Per User Tercapai

**Scenario:**
```
Promo: max_usage_per_user = 2
User sudah pakai: 2 kali
User mencoba pakai lagi: TOKEN20OFF
```

**Backend Validation:**
```go
if promo.MaxUsagePerUser != nil {
  usageCount := promoUsageRepo.CountByPromoAndUser(promo.ID, userID)
  if usageCount >= *promo.MaxUsagePerUser {
    return error(fmt.Sprintf("Anda sudah menggunakan promo ini %d kali (maksimal %d kali)", usageCount, *promo.MaxUsagePerUser))
  }
}
```

**Response:**
```json
{
  "status": "error",
  "message": "Anda sudah mencapai batas penggunaan",
  "error_code": "USER_MAX_USAGE_REACHED",
  "data": {
    "code": "TOKEN20OFF",
    "user_usage_count": 2,
    "max_usage_per_user": 2,
    "message": "Promo ini hanya bisa digunakan 2 kali per user"
  }
}
```

---

### Case 4.5: Promo Tidak Berlaku untuk Plan Tertentu

**Scenario:**
```
User memilih: Plan Basic
Promo: applicable_plans = ["PRO", "MAX"]
User input: PRO50OFF
```

**Backend Validation:**
```go
if promo.ApplicablePlans != nil {
  var applicablePlans []string
  json.Unmarshal(promo.ApplicablePlans, &applicablePlans)
  
  if !contains(applicablePlans, "BASIC") {
    return error("Promo code hanya berlaku untuk paket Pro dan Max")
  }
}
```

**Response:**
```json
{
  "status": "error",
  "message": "Promo tidak berlaku untuk paket yang dipilih",
  "error_code": "INVALID_PLAN_FOR_PROMO",
  "data": {
    "code": "PRO50OFF",
    "selected_plan": "Basic",
    "applicable_plans": ["Pro", "Max"],
    "message": "Promo ini hanya berlaku untuk paket Pro dan Max"
  }
}
```

---

### Case 4.6: Minimal Token Purchase Tidak Terpenuhi

**Scenario:**
```
User beli: 30 token
Promo: min_token_purchase = 50
User input: TOKEN20OFF
```

**Backend Validation:**
```go
if promo.MinTokenPurchase != nil && tokenQuantity < *promo.MinTokenPurchase {
  return error(fmt.Sprintf("Minimal pembelian %d token untuk menggunakan promo ini", *promo.MinTokenPurchase))
}
```

**Response:**
```json
{
  "status": "error",
  "message": "Pembelian token kurang dari minimal",
  "error_code": "MIN_TOKEN_NOT_MET",
  "data": {
    "code": "TOKEN20OFF",
    "token_quantity": 30,
    "min_token_purchase": 50,
    "shortage": 20,
    "message": "Tambah 20 token lagi untuk bisa pakai promo ini"
  }
}
```

---

### Case 4.7: Promo Type Salah

**Scenario:**
```
User beli: Membership Pro
Promo type: "token"
User input: TOKEN20OFF
```

**Backend Validation:**
```go
if promo.PromoType == "token" && paymentType == "membership" {
  return error("Promo code hanya berlaku untuk pembelian token standalone")
}
```

**Response:**
```json
{
  "status": "error",
  "message": "Promo tidak berlaku untuk membership",
  "error_code": "INVALID_PROMO_TYPE",
  "data": {
    "code": "TOKEN20OFF",
    "promo_type": "token",
    "payment_type": "membership",
    "message": "Promo ini khusus untuk pembelian token"
  }
}
```

---

### Case 4.8: Promo Code Tidak Aktif

**Scenario:**
```
Promo: is_active = false
User input: OLDPROMO
```

**Backend Validation:**
```go
if !promo.IsActive {
  return error("Promo code tidak aktif")
}
```

**Response:**
```json
{
  "status": "error",
  "message": "Promo code tidak aktif",
  "error_code": "PROMO_INACTIVE",
  "data": {
    "code": "OLDPROMO",
    "message": "Promo ini sudah tidak berlaku lagi"
  }
}
```

---

## 5. Flow Diagram

### 5.1. Complete Promo Validation Flow

```
┌─────────────────────┐
│ User Input Promo    │
│ Code                │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ Find Promo by Code  │
│ (case-insensitive)  │
└──────────┬──────────┘
           │
           ▼
    ┌──────┴──────┐
    │ Promo Found?│
    └──────┬──────┘
           │
    ┌──────┴──────┐
    │     NO      │──────> ❌ PROMO_NOT_FOUND
    └─────────────┘
    ┌──────┴──────┐
    │     YES     │
    └──────┬──────┘
           │
           ▼
┌─────────────────────┐
│ Check is_active     │
└──────────┬──────────┘
           │
    ┌──────┴──────┐
    │   FALSE     │──────> ❌ PROMO_INACTIVE
    └─────────────┘
    ┌──────┴──────┐
    │   TRUE      │
    └──────┬──────┘
           │
           ▼
┌─────────────────────┐
│ Check Valid Date    │
│ (valid_from/until)  │
└──────────┬──────────┘
           │
    ┌──────┴──────┐
    │ EXPIRED?    │──────> ❌ PROMO_EXPIRED
    └─────────────┘
    ┌──────┴──────┐
    │ NOT STARTED?│──────> ❌ PROMO_NOT_STARTED
    └─────────────┘
    ┌──────┴──────┐
    │   VALID     │
    └──────┬──────┘
           │
           ▼
┌─────────────────────┐
│ Check Promo Type    │
└──────────┬──────────┘
           │
    ┌──────┴──────┐
    │ membership  │───┐
    └─────────────┘   │
    ┌──────┴──────┐   │
    │   token     │───┼───> Validate based on payment_type
    └─────────────┘   │
    ┌──────┴──────┐   │
    │    both     │───┘
    └──────┬──────┘
           │
           ▼
┌─────────────────────┐
│ Check Max Usage     │
│ (total)             │
└──────────┬──────────┘
           │
    ┌──────┴──────┐
    │ REACHED?    │──────> ❌ PROMO_MAX_USAGE_REACHED
    └─────────────┘
    ┌──────┴──────┐
    │   OK        │
    └──────┬──────┘
           │
           ▼
┌─────────────────────┐
│ Check Max Usage     │
│ Per User            │
└──────────┬──────────┘
           │
    ┌──────┴──────┐
    │ REACHED?    │──────> ❌ USER_MAX_USAGE_REACHED
    └─────────────┘
    ┌──────┴──────┐
    │   OK        │
    └──────┬──────┘
           │
           ▼
┌─────────────────────┐
│ IF membership:      │
│ - Check plans       │
│ - Check durations   │
└──────────┬──────────┘
           │
    ┌──────┴──────┐
    │ NOT MATCH?  │──────> ❌ INVALID_PLAN/DURATION
    └─────────────┘
    ┌──────┴──────┐
    │   MATCH     │
    └──────┬──────┘
           │
           ▼
┌─────────────────────┐
│ IF token:           │
│ - Check min_token   │
└──────────┬──────────┘
           │
    ┌──────┴──────┐
    │ NOT MET?    │──────> ❌ MIN_TOKEN_NOT_MET
    └─────────────┘
    ┌──────┴──────┐
    │   MET       │
    └──────┬──────┘
           │
           ▼
┌─────────────────────┐
│ ✅ PROMO VALID      │
│ Calculate Discount  │
│ & Bonus Token       │
└─────────────────────┘
```

---

### 5.2. After Payment Success Flow

```
┌─────────────────────┐
│ Xendit Callback     │
│ Status: PAID        │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ Get Payment Data    │
│ - promo_code        │
│ - payment_type      │
└──────────┬──────────┘
           │
           ▼
    ┌──────┴──────────┐
    │ Has Promo Code? │
    └──────┬──────────┘
           │
    ┌──────┴──────┐
    │     NO      │────┐
    └─────────────┘    │
    ┌──────┴──────┐    │
    │    YES      │    │
    └──────┬──────┘    │
           │           │
           ▼           │
┌─────────────────────┐│
│ Activate Membership ││
│ or Add Token        ││
└──────────┬──────────┘│
           │           │
           ▼           │
┌─────────────────────┐│
│ Add Bonus Token     ││
│ (if any)            ││
└──────────┬──────────┘│
           │           │
           ▼           │
┌─────────────────────┐│
│ Record Promo Usage  ││
│ in promo_code_usage ││
└──────────┬──────────┘│
           │           │
           ▼           │
┌─────────────────────┐│
│ Increment           ││
│ promo.current_usage ││
└──────────┬──────────┘│
           │           │
           └───────────┴───> ✅ Done
```

---

## 6. Testing Checklist

### 6.1. Membership Promo Testing

- [ ] Percentage discount (10%, 50%, 100%)
- [ ] Fixed amount discount
- [ ] Free 100% (skip payment)
- [ ] With bonus token
- [ ] Applicable to specific plans only
- [ ] Applicable to specific durations only
- [ ] Universal promo (all plans, all durations)

### 6.2. Token Promo Testing

- [ ] Percentage discount
- [ ] Fixed amount discount
- [ ] Bonus token only (no price discount)
- [ ] Discount + bonus token combined
- [ ] Min token purchase validation
- [ ] Universal promo for token

### 6.3. Validation Testing

- [ ] Expired promo code
- [ ] Not yet started promo
- [ ] Inactive promo
- [ ] Max usage reached (total)
- [ ] Max usage per user reached
- [ ] Wrong promo type (token promo for membership purchase)
- [ ] Plan not applicable
- [ ] Duration not applicable
- [ ] Min token not met
- [ ] Case-insensitive code matching

### 6.4. Edge Cases Testing

- [ ] Multiple promo attempts (should reject after max usage)
- [ ] Concurrent usage (race condition)
- [ ] Payment failed after promo validation (no usage increment)
- [ ] User tries to use same promo for different purchases
- [ ] Very large discount (> price) - should cap at price
- [ ] Promo becomes expired during payment process
- [ ] Starter Plan activation with token promo (non-member only)

---

## 7. Kesimpulan

Dokumen ini mencakup **semua kemungkinan skenario** penggunaan promo code di REXTRA:

✅ **3 Tipe Promo**: membership, token, both
✅ **4 Discount Types**: percentage, fixed_amount, free_100, bonus_token
✅ **8 Edge Cases**: expired, inactive, max usage, wrong type, dll
✅ **2 Flow Diagrams**: validation & post-payment
✅ **Complete Testing Checklist**

**Total Scenarios Covered**: 16 scenarios + 8 edge cases = **24 test cases**

---

**Prepared by**: Tim Backend REXTRA  
**Document Version**: 1.0  
**Last Updated**: November 2025  
**Status**: Ready for Implementation