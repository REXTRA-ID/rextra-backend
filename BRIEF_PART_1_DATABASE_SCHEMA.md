# 🔧 BRIEF BACKEND API REXTRA CLUB MEMBERSHIP
## PART 1: DATABASE SCHEMA & MODELS

---

## A. Gambaran Penugasan Backend API

**Nama Tugas**: Backend API REXTRA CLUB Membership System

**Penanggung Jawab**: [Nama PIC Backend Team]

**Status Pengerjaan**: Belum Dimulai

**Status Pengecekan**: Belum Diperiksa

**Tanggal Mulai**: [DD/MM/YYYY]

**Tanggal Selesai**: [DD/MM/YYYY]

---

## B. Brief Penugasan Backend

### 1. Tujuan / Objective

Backend API ini bertanggung jawab untuk:

- **Mengelola sistem membership REXTRA CLUB** dengan tiga paket langganan (Basic, Pro, Max) dan status Starter Plan untuk user baru atau non-member yang membeli token
- **Melacak status keanggotaan user** dari Starter → Member (Basic/Pro/Max) → Non-Member, termasuk durasi langganan dan perpanjangan otomatis
- **Mengelola sistem token** untuk akses fitur AI-based dengan tracking penggunaan per fitur dan history transaksi token
- **Mengelola sistem REXTRA POIN** sebagai reward dari pembelian membership dengan tracking transaksi dan balance
- **Integrasi payment gateway Xendit** untuk proses pembayaran membership dan token, termasuk webhook handling untuk auto-activation
- **Menyediakan endpoint untuk aplikasi mobile/web** mahasiswa dalam mengelola membership, melihat benefit, membeli paket, dan tracking token & poin
- **Mengelola akses fitur berbasis role membership** (Starter, Basic, Pro, Max, Non-Member) dengan pembatasan penggunaan harian dan token-based

---

## C. Desain Data & Entity

### 1. Daftar Entity Utama

1. **memberships** – Menyimpan data membership aktif user (status, paket, durasi, token balance, poin balance)
2. **membership_plans** – Menyimpan master data paket membership (Basic, Pro, Max) dengan harga dan benefit
3. **membership_durations** – Menyimpan pilihan durasi langganan (1, 3, 6, 12 bulan) dengan diskon dan bonus token
4. **token_transactions** – Menyimpan history transaksi token (purchase, usage, refund)
5. **token_usage_history** – Menyimpan detail penggunaan token per fitur AI
6. **poin_transactions** – Menyimpan history transaksi REXTRA POIN (earn, redeem)
7. **payment_transactions** – Menyimpan data transaksi pembayaran via Xendit
8. **promo_codes** – (Optional - fase admin) Menyimpan kode promo untuk diskon membership

**Catatan**: Tabel `users` sudah ada di sistem REXTRA yang lain, jadi tidak perlu dibuat ulang. Tabel membership akan reference ke `users.id` yang sudah ada.

---

### 2. Struktur Tabel Detail

---

#### Tabel: memberships

**Fungsi**: Menyimpan data membership aktif setiap user, termasuk status (Starter/Basic/Pro/Max/Non-Member), balance token & poin, dan masa berlaku.

| Nama Kolom | Tipe | Deskripsi singkat |
|------------|------|-------------------|
| id | UUID | ID unik membership (PK) |
| user_id | UUID | FK ke users.id (unique) |
| membership_status | VARCHAR(20) | starter/basic/pro/max/non_member |
| plan_id | UUID | FK ke membership_plans.id (nullable untuk starter/non_member) |
| duration_id | UUID | FK ke membership_durations.id (nullable untuk starter/non_member) |
| current_token_balance | INTEGER | Sisa token yang dimiliki user |
| current_poin_balance | INTEGER | Sisa REXTRA POIN yang dimiliki user |
| started_at | TIMESTAMP | Waktu mulai membership aktif |
| expired_at | TIMESTAMP | Waktu berakhir membership |
| is_active | BOOLEAN | Status aktif membership |
| auto_renew | BOOLEAN | Status auto-renewal (future feature) |
| created_at | TIMESTAMP | Waktu pembuatan record |
| updated_at | TIMESTAMP | Waktu update terakhir |

**SQL DDL (PostgreSQL)**:

```sql
CREATE TABLE memberships (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  membership_status VARCHAR(20) NOT NULL CHECK (membership_status IN ('starter', 'basic', 'pro', 'max', 'non_member')),
  plan_id UUID REFERENCES membership_plans(id) ON DELETE SET NULL,
  duration_id UUID REFERENCES membership_durations(id) ON DELETE SET NULL,
  current_token_balance INTEGER DEFAULT 0,
  current_poin_balance INTEGER DEFAULT 0,
  started_at TIMESTAMP,
  expired_at TIMESTAMP,
  is_active BOOLEAN DEFAULT TRUE,
  auto_renew BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_membership_user ON memberships(user_id);
CREATE INDEX idx_membership_status ON memberships(membership_status);
CREATE INDEX idx_membership_active ON memberships(is_active);
CREATE INDEX idx_membership_expired ON memberships(expired_at);
```

**Business Rules**:
- Setiap user hanya punya **1 record membership** (one-to-one relationship)
- Saat user baru registrasi → auto create membership dengan status `starter`
- Saat starter expired → update status jadi `non_member`
- Saat beli membership → update status jadi `basic/pro/max`
- Token dan Poin balance **tidak pernah dihapus**, hanya berkurang saat dipakai

---

#### Tabel: membership_plans

**Fungsi**: Master data untuk paket membership (Basic, Pro, Max) dengan detail harga dan benefit.

| Nama Kolom | Tipe | Deskripsi singkat |
|------------|------|-------------------|
| id | UUID | ID unik plan (PK) |
| plan_name | VARCHAR(50) | Nama paket (Basic/Pro/Max) |
| plan_code | VARCHAR(20) | Kode paket (BASIC/PRO/MAX) - unique |
| monthly_token | INTEGER | Jumlah token per bulan |
| base_monthly_price | DECIMAL(10,2) | Harga dasar bulanan (Rp) |
| description | TEXT | Deskripsi paket |
| benefits | JSONB | (OPTIONAL) List benefit untuk display UI katalog saja, BUKAN untuk control akses fitur |
| is_active | BOOLEAN | Status paket aktif |
| created_at | TIMESTAMP | Waktu pembuatan |
| updated_at | TIMESTAMP | Waktu update terakhir |

**SQL DDL (PostgreSQL)**:

```sql
CREATE TABLE membership_plans (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  plan_name VARCHAR(50) NOT NULL,
  plan_code VARCHAR(20) UNIQUE NOT NULL,
  monthly_token INTEGER NOT NULL,
  base_monthly_price DECIMAL(10,2) NOT NULL,
  description TEXT,
  benefits JSONB,
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_plan_code ON membership_plans(plan_code);
CREATE INDEX idx_plan_active ON membership_plans(is_active);
```

**Data Seed**:

```sql
INSERT INTO membership_plans (plan_name, plan_code, monthly_token, base_monthly_price, description, benefits) VALUES
(
  'Basic', 
  'BASIC', 
  10, 
  15000.00, 
  'Paket dasar untuk mahasiswa yang baru memulai persiapan karier', 
  '["Akses CRUD System unlimited", "10 token/bulan", "Akses REXTRA Academy kategori Basic", "Portofolio AI Check 5x/bulan"]'::jsonb
),
(
  'Pro', 
  'PRO', 
  20, 
  20000.00, 
  'Paket untuk mahasiswa aktif dengan kebutuhan AI moderat',
  '["Semua benefit Basic", "20 token/bulan", "Akses 6 modul pembelajaran karier", "Prioritas Rencana Karier AI", "Portofolio AI Check 15x/bulan"]'::jsonb
),
(
  'Max', 
  'MAX', 
  40, 
  30000.00, 
  'Paket premium untuk mahasiswa intensif persiapan karier',
  '["Semua benefit Pro", "40 token/bulan", "Akses penuh semua modul & fitur", "Unlimited Portofolio AI Check", "Priority support"]'::jsonb
);
```

**⚠️ CATATAN PENTING - Control Akses Fitur:**

Kolom `benefits` di tabel ini **HANYA untuk display UI** di halaman katalog membership (pricing page). Kolom ini **BUKAN** dipakai untuk control akses fitur.

**Control akses fitur diatur di middleware masing-masing API fitur** dengan logic sebagai berikut:

**Contoh Logic Control Akses per Fitur:**

```go
// 1. Fitur berbasis TOKEN (Kenali Diri, CV Generator, AI Interviewer)
// Cek: current_token_balance >= token_required
// TIDAK peduli membership_status apa (bahkan non_member bisa asal punya token)
func KenaliDiriMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {
    userID := c.GetString("user_id")
    membership := getMembership(userID)
    
    tokenRequired := 3 // Kenali Diri butuh 3 token
    if membership.CurrentTokenBalance < tokenRequired {
      c.JSON(400, response.Error("Insufficient token", "INSUFFICIENT_TOKEN"))
      c.Abort()
      return
    }
    
    c.Next()
  }
}

// 2. Fitur CRUD System (Portofolio, Rencana Karier, Persona)
// Cek: membership_status IN ('starter', 'basic', 'pro', 'max')
// Non-Member TIDAK bisa akses
func PortofolioMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {
    userID := c.GetString("user_id")
    membership := getMembership(userID)
    
    allowedStatus := []string{"starter", "basic", "pro", "max"}
    if !contains(allowedStatus, membership.MembershipStatus) {
      c.JSON(403, response.Error("Akses ditolak. Bergabung dengan REXTRA CLUB untuk akses fitur ini.", "FORBIDDEN"))
      c.Abort()
      return
    }
    
    c.Next()
  }
}

// 3. Fitur dengan pembatasan HARIAN (Portofolio AI Check)
// Cek: membership_status + hitung usage hari ini
func PortofolioAICheckMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {
    userID := c.GetString("user_id")
    membership := getMembership(userID)
    
    // Define daily limit per membership status
    dailyLimits := map[string]int{
      "starter":    0,  // Starter tidak bisa pakai AI Check
      "basic":      5,  // Basic: 5x per bulan (anggap ~0.16/hari)
      "pro":        15, // Pro: 15x per bulan (anggap ~0.5/hari)
      "max":        -1, // Max: unlimited
      "non_member": 0,  // Non-member tidak bisa
    }
    
    limit := dailyLimits[membership.MembershipStatus]
    
    if limit == 0 {
      c.JSON(403, response.Error("Fitur tidak tersedia untuk membership Anda", "FORBIDDEN"))
      c.Abort()
      return
    }
    
    if limit > 0 {
      // Hitung usage hari ini
      usageToday := countUsageToday(userID, "portofolio_ai_check")
      if usageToday >= limit {
        c.JSON(400, response.Error("Daily limit reached", "DAILY_LIMIT_REACHED"))
        c.Abort()
        return
      }
    }
    
    // limit == -1 berarti unlimited (Max plan)
    c.Next()
  }
}

// 4. Fitur PREMIUM khusus plan tertentu (misalnya: Priority Support)
// Cek: membership_status harus 'max'
func PrioritySupportMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {
    userID := c.GetString("user_id")
    membership := getMembership(userID)
    
    if membership.MembershipStatus != "max" {
      c.JSON(403, response.Error("Fitur Priority Support hanya untuk paket Max", "FORBIDDEN"))
      c.Abort()
      return
    }
    
    c.Next()
  }
}
```

**Kesimpulan:**
- `benefits` JSONB → Pure untuk display di frontend (pricing page)
- Control akses fitur → Logic di middleware per API, bukan dari kolom `benefits`
- Frontend bisa hardcode list benefits di code, atau fetch dari kolom `benefits` ini (optional)



---

#### Tabel: membership_durations

**Fungsi**: Master data untuk pilihan durasi langganan dengan diskon dan bonus token yang berbeda-beda.

| Nama Kolom | Tipe | Deskripsi singkat |
|------------|------|-------------------|
| id | UUID | ID unik duration (PK) |
| duration_months | INTEGER | Durasi langganan dalam bulan (1/3/6/12) |
| discount_percentage | DECIMAL(5,2) | Persentase diskon (0/8/15/30) |
| token_bonus_percentage | DECIMAL(5,2) | Persentase bonus token (0/10/20/40) |
| rextra_poin_multiplier | INTEGER | Multiplier untuk REXTRA POIN |
| is_active | BOOLEAN | Status durasi aktif |
| created_at | TIMESTAMP | Waktu pembuatan |
| updated_at | TIMESTAMP | Waktu update terakhir |

**SQL DDL (PostgreSQL)**:

```sql
CREATE TABLE membership_durations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  duration_months INTEGER NOT NULL,
  discount_percentage DECIMAL(5,2) DEFAULT 0,
  token_bonus_percentage DECIMAL(5,2) DEFAULT 0,
  rextra_poin_multiplier INTEGER DEFAULT 1,
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_duration_months ON membership_durations(duration_months);
```

**Data Seed**:

```sql
INSERT INTO membership_durations (duration_months, discount_percentage, token_bonus_percentage, rextra_poin_multiplier) VALUES
(1, 0, 0, 1),    -- 1 bulan: no discount, no bonus
(3, 8, 10, 3),   -- 3 bulan: 8% discount, 10% bonus token, 3x poin
(6, 15, 20, 6),  -- 6 bulan: 15% discount, 20% bonus token, 6x poin
(12, 30, 40, 12); -- 12 bulan: 30% discount, 40% bonus token, 12x poin
```

**Business Rules - Perhitungan Harga**:

```
Total Harga (HT) = (Harga Bulanan × Lama Langganan) × (1 − Diskon)
Total Token (TT) = (Token Bulanan × Lama Langganan) × (1 + Bonus Token)
Efektif Harga per Bulan (EHB) = HT ÷ Lama Langganan
Efektif Harga per Token (EPT) = HT ÷ TT
REXTRA POIN = (Token Bulanan × Lama Langganan) × Multiplier
```

**Contoh Perhitungan - Basic 6 Bulan**:
```
Base Price = Rp15.000 × 6 = Rp90.000
Discount = Rp90.000 × 15% = Rp13.500
Total Harga = Rp90.000 - Rp13.500 = Rp76.500 (tapi di spec Rp51.000)

Base Token = 10 × 6 = 60 token
Bonus Token = 60 × 20% = 12 token
Total Token = 60 + 12 = 72 token

REXTRA POIN = (10 × 6) × 6 = 360 poin (tapi di spec 60 poin)
```

**⚠️ CATATAN PENTING**: Ada inkonsistensi antara rumus dan tabel harga di document. Tolong konfirmasi dengan tim product:
1. Apakah pakai rumus atau pakai fixed price dari tabel?
2. Kalau pakai fixed price, perlu tambah kolom `fixed_price` di tabel ini

---

#### Tabel: token_transactions

**Fungsi**: Menyimpan semua transaksi token (pembelian, penggunaan, refund) untuk audit trail.

| Nama Kolom | Tipe | Deskripsi singkat |
|------------|------|-------------------|
| id | UUID | ID unik transaksi (PK) |
| user_id | UUID | FK ke users.id |
| transaction_type | VARCHAR(30) | purchase_membership/purchase_standalone/usage/refund/monthly_refill |
| token_amount | INTEGER | Jumlah token (+ untuk in, - untuk out) |
| token_balance_before | INTEGER | Saldo token sebelum transaksi |
| token_balance_after | INTEGER | Saldo token setelah transaksi |
| reference_type | VARCHAR(50) | Tipe referensi (payment/feature_usage/membership_renewal/monthly_refill) |
| reference_id | UUID | ID referensi ke tabel lain |
| description | TEXT | Deskripsi transaksi |
| created_at | TIMESTAMP | Waktu transaksi |

**SQL DDL (PostgreSQL)**:

```sql
CREATE TABLE token_transactions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  transaction_type VARCHAR(30) NOT NULL CHECK (transaction_type IN ('purchase_membership', 'purchase_standalone', 'usage', 'refund', 'monthly_refill')),
  token_amount INTEGER NOT NULL,
  token_balance_before INTEGER NOT NULL,
  token_balance_after INTEGER NOT NULL,
  reference_type VARCHAR(50),
  reference_id UUID,
  description TEXT,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_token_tx_user ON token_transactions(user_id);
CREATE INDEX idx_token_tx_type ON token_transactions(transaction_type);
CREATE INDEX idx_token_tx_created ON token_transactions(created_at DESC);
```

**Business Rules**:
- Setiap perubahan token balance **HARUS** create record di tabel ini
- `token_balance_before` dan `token_balance_after` untuk validasi consistency
- `transaction_type`:
  - `purchase_membership`: Token dari pembelian membership plan
  - `purchase_standalone`: Token dari pembelian token terpisah
  - `monthly_refill`: Token refill otomatis setiap bulan untuk member aktif
  - `usage`: Penggunaan token untuk fitur AI
  - `refund`: Refund token (misalnya error atau cancel)

---

#### Tabel: token_usage_history

**Fungsi**: Menyimpan detail penggunaan token per fitur AI untuk analytics dan monitoring usage pattern.

| Nama Kolom | Tipe | Deskripsi singkat |
|------------|------|-------------------|
| id | UUID | ID unik usage (PK) |
| user_id | UUID | FK ke users.id |
| token_transaction_id | UUID | FK ke token_transactions.id |
| feature_name | VARCHAR(100) | Nama fitur yang digunakan (kenali_diri/cv_generator/ai_interviewer) |
| token_used | INTEGER | Jumlah token yang dipakai |
| usage_metadata | JSONB | Metadata tambahan (feature-specific data) |
| created_at | TIMESTAMP | Waktu penggunaan |

**SQL DDL (PostgreSQL)**:

```sql
CREATE TABLE token_usage_history (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_transaction_id UUID REFERENCES token_transactions(id) ON DELETE SET NULL,
  feature_name VARCHAR(100) NOT NULL,
  token_used INTEGER NOT NULL,
  usage_metadata JSONB,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_token_usage_user ON token_usage_history(user_id);
CREATE INDEX idx_token_usage_feature ON token_usage_history(feature_name);
CREATE INDEX idx_token_usage_created ON token_usage_history(created_at DESC);
```

**Contoh usage_metadata untuk berbagai fitur**:

```json
// Kenali Diri (Assessment)
{
  "assessment_id": "uuid-assessment",
  "question_count": 15,
  "completion_time_seconds": 420,
  "result_persona": "builder"
}

// CV Generator
{
  "cv_template_id": "template-modern-01",
  "sections_generated": ["profile", "experience", "skills", "education"],
  "total_pages": 2,
  "format": "pdf"
}

// AI Interviewer
{
  "interview_session_id": "uuid-session",
  "job_position": "Frontend Developer",
  "question_count": 10,
  "duration_minutes": 25,
  "average_score": 85
}
```

**Business Rules**:
- Setiap kali user pakai fitur AI → create record di `token_transactions` (type: usage) DAN di tabel ini
- Relasi one-to-one: 1 token_transaction (usage) = 1 token_usage_history
- Feature name harus konsisten, pakai constant:
  - `kenali_diri`
  - `cv_generator`
  - `ai_interviewer`

---

#### Tabel: poin_transactions

**Fungsi**: Menyimpan semua transaksi REXTRA POIN (earn, redeem, adjustment) untuk gamification system.

| Nama Kolom | Tipe | Deskripsi singkat |
|------------|------|-------------------|
| id | UUID | ID unik transaksi (PK) |
| user_id | UUID | FK ke users.id |
| transaction_type | VARCHAR(20) | earn/redeem/adjustment |
| poin_amount | INTEGER | Jumlah poin (+ untuk earn, - untuk redeem) |
| poin_balance_before | INTEGER | Saldo poin sebelum transaksi |
| poin_balance_after | INTEGER | Saldo poin setelah transaksi |
| source | VARCHAR(100) | Sumber poin (membership_purchase/mission_complete/event_reward/redeem_reward) |
| reference_id | UUID | ID referensi ke tabel lain |
| description | TEXT | Deskripsi transaksi |
| created_at | TIMESTAMP | Waktu transaksi |

**SQL DDL (PostgreSQL)**:

```sql
CREATE TABLE poin_transactions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  transaction_type VARCHAR(20) NOT NULL CHECK (transaction_type IN ('earn', 'redeem', 'adjustment')),
  poin_amount INTEGER NOT NULL,
  poin_balance_before INTEGER NOT NULL,
  poin_balance_after INTEGER NOT NULL,
  source VARCHAR(100),
  reference_id UUID,
  description TEXT,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_poin_tx_user ON poin_transactions(user_id);
CREATE INDEX idx_poin_tx_type ON poin_transactions(transaction_type);
CREATE INDEX idx_poin_tx_created ON poin_transactions(created_at DESC);
```

**Business Rules**:
- Setiap perubahan poin balance **HARUS** create record di tabel ini
- `transaction_type`:
  - `earn`: User dapat poin (dari membership purchase, mission complete, event reward)
  - `redeem`: User tukar poin untuk reward
  - `adjustment`: Admin manual adjustment (fix error, bonus special)
- `source`:
  - `membership_purchase`: Poin dari beli membership
  - `mission_complete`: Poin dari selesaikan mission di Persona REXTRA
  - `event_reward`: Poin dari partisipasi event REXTRA
  - `redeem_reward`: Penukaran poin untuk voucher/merchandise

---

#### Tabel: payment_transactions

**Fungsi**: Menyimpan data transaksi pembayaran via Xendit untuk membership dan token standalone.

| Nama Kolom | Tipe | Deskripsi singkat |
|------------|------|-------------------|
| id | UUID | ID unik payment (PK) |
| user_id | UUID | FK ke users.id |
| payment_type | VARCHAR(30) | membership/token_standalone |
| plan_id | UUID | FK ke membership_plans.id (nullable untuk token) |
| duration_id | UUID | FK ke membership_durations.id (nullable untuk token) |
| token_quantity | INTEGER | Jumlah token yang dibeli (untuk standalone) |
| gross_amount | DECIMAL(12,2) | Total harga sebelum diskon |
| discount_amount | DECIMAL(12,2) | Jumlah diskon |
| final_amount | DECIMAL(12,2) | Total harga akhir yang dibayar |
| promo_code | VARCHAR(50) | Kode promo yang digunakan (optional) |
| xendit_invoice_id | VARCHAR(255) | Invoice ID dari Xendit (unique) |
| xendit_external_id | VARCHAR(255) | External ID untuk tracking (unique) |
| payment_method | VARCHAR(50) | Metode pembayaran (va/ewallet/qris/credit_card) |
| payment_status | VARCHAR(20) | pending/paid/failed/expired |
| paid_at | TIMESTAMP | Waktu pembayaran berhasil |
| expired_at | TIMESTAMP | Waktu invoice expired |
| xendit_callback_data | JSONB | Raw callback data dari Xendit |
| created_at | TIMESTAMP | Waktu pembuatan transaksi |
| updated_at | TIMESTAMP | Waktu update terakhir |

**SQL DDL (PostgreSQL)**:

```sql
CREATE TABLE payment_transactions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  payment_type VARCHAR(30) NOT NULL CHECK (payment_type IN ('membership', 'token_standalone')),
  plan_id UUID REFERENCES membership_plans(id) ON DELETE SET NULL,
  duration_id UUID REFERENCES membership_durations(id) ON DELETE SET NULL,
  token_quantity INTEGER,
  gross_amount DECIMAL(12,2) NOT NULL,
  discount_amount DECIMAL(12,2) DEFAULT 0,
  final_amount DECIMAL(12,2) NOT NULL,
  promo_code VARCHAR(50),
  xendit_invoice_id VARCHAR(255) UNIQUE,
  xendit_external_id VARCHAR(255) UNIQUE,
  payment_method VARCHAR(50),
  payment_status VARCHAR(20) NOT NULL CHECK (payment_status IN ('pending', 'paid', 'failed', 'expired')),
  paid_at TIMESTAMP,
  expired_at TIMESTAMP,
  xendit_callback_data JSONB,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_payment_user ON payment_transactions(user_id);
CREATE INDEX idx_payment_status ON payment_transactions(payment_status);
CREATE INDEX idx_payment_xendit_invoice ON payment_transactions(xendit_invoice_id);
CREATE INDEX idx_payment_xendit_external ON payment_transactions(xendit_external_id);
CREATE INDEX idx_payment_created ON payment_transactions(created_at DESC);
```

**Business Rules**:
- Setiap payment request ke Xendit → create record dengan status `pending`
- Saat Xendit webhook callback → update status jadi `paid`/`failed`/`expired`
- Jika `payment_status = 'paid'` → trigger activation logic (membership atau token)
- `xendit_external_id` format: `rextra_{payment_type}_{user_id}_{timestamp}`
- `xendit_callback_data` simpan raw JSON dari webhook untuk audit

---

#### Tabel: promo_codes (Optional - Fase Admin)

**Fungsi**: Menyimpan kode promo untuk diskon membership yang bisa di-manage oleh admin.

| Nama Kolom | Tipe | Deskripsi singkat |
|------------|------|-------------------|
| id | UUID | ID unik promo (PK) |
| code | VARCHAR(50) | Kode promo (unique, case-insensitive) |
| description | TEXT | Deskripsi promo |
| discount_type | VARCHAR(20) | percentage/fixed_amount/free_100 |
| discount_value | DECIMAL(10,2) | Nilai diskon (% atau Rp) |
| applicable_plans | JSONB | Array plan_code yang berlaku (null = semua) |
| applicable_durations | JSONB | Array duration_months yang berlaku (null = semua) |
| max_usage | INTEGER | Maksimal penggunaan total (null = unlimited) |
| current_usage | INTEGER | Jumlah penggunaan saat ini |
| valid_from | TIMESTAMP | Waktu mulai berlaku |
| valid_until | TIMESTAMP | Waktu berakhir |
| is_active | BOOLEAN | Status promo aktif |
| created_at | TIMESTAMP | Waktu pembuatan |
| updated_at | TIMESTAMP | Waktu update terakhir |

**SQL DDL (PostgreSQL)**:

```sql
CREATE TABLE promo_codes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code VARCHAR(50) UNIQUE NOT NULL,
  description TEXT,
  discount_type VARCHAR(20) NOT NULL CHECK (discount_type IN ('percentage', 'fixed_amount', 'free_100')),
  discount_value DECIMAL(10,2),
  applicable_plans JSONB,
  applicable_durations JSONB,
  max_usage INTEGER,
  current_usage INTEGER DEFAULT 0,
  valid_from TIMESTAMP NOT NULL,
  valid_until TIMESTAMP NOT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_promo_code ON promo_codes(code);
CREATE INDEX idx_promo_active ON promo_codes(is_active);
CREATE INDEX idx_promo_valid ON promo_codes(valid_from, valid_until);
```

**Contoh Data Seed**:

```sql
-- Promo 100% gratis untuk early adopters
INSERT INTO promo_codes (code, description, discount_type, discount_value, applicable_plans, applicable_durations, max_usage, valid_from, valid_until) VALUES
(
  'REXTRA2025',
  'Promo spesial 100% gratis untuk mahasiswa baru di tahun 2025',
  'free_100',
  100,
  '["BASIC", "PRO", "MAX"]'::jsonb,
  '[1, 3, 6, 12]'::jsonb,
  1000,
  '2025-01-01 00:00:00',
  '2025-12-31 23:59:59'
);

-- Promo diskon 50% untuk paket Pro
INSERT INTO promo_codes (code, description, discount_type, discount_value, applicable_plans, applicable_durations, max_usage, valid_from, valid_until) VALUES
(
  'PRO50OFF',
  'Diskon 50% untuk paket Pro durasi 6 dan 12 bulan',
  'percentage',
  50,
  '["PRO"]'::jsonb,
  '[6, 12]'::jsonb,
  500,
  '2025-01-01 00:00:00',
  '2025-03-31 23:59:59'
);

-- Promo fixed amount Rp20.000 off
INSERT INTO promo_codes (code, description, discount_type, discount_value, applicable_plans, applicable_durations, max_usage, valid_from, valid_until) VALUES
(
  'DISKON20K',
  'Potongan harga Rp20.000 untuk semua paket',
  'fixed_amount',
  20000,
  NULL, -- berlaku untuk semua plan
  NULL, -- berlaku untuk semua duration
  NULL, -- unlimited usage
  '2025-01-01 00:00:00',
  '2025-12-31 23:59:59'
);
```

---

### 3. Relasi Antar Tabel

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
                     └──→ membership_durations (N:1)
                     └──→ promo_codes (N:1, optional)
```

**Penjelasan Relasi**:
1. **users ↔ memberships**: One-to-One (satu user hanya punya satu membership record)
2. **memberships → plans/durations**: Many-to-One (banyak membership bisa pakai plan/duration yang sama)
3. **users → token_transactions**: One-to-Many (satu user bisa punya banyak transaksi token)
4. **token_transactions → token_usage_history**: One-to-One untuk type='usage' (setiap usage punya detail)
5. **users → poin_transactions**: One-to-Many (satu user bisa punya banyak transaksi poin)
6. **users → payment_transactions**: One-to-Many (satu user bisa punya banyak payment)

---

### 4. Definisi Model (Golang dengan GORM)

```go
package domain

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/datatypes"
)

// Membership Model
type Membership struct {
    ID                   uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    UserID               uuid.UUID  `gorm:"type:uuid;unique;not null" json:"user_id"`
    MembershipStatus     string     `gorm:"type:varchar(20);not null" json:"membership_status"`
    PlanID               *uuid.UUID `gorm:"type:uuid" json:"plan_id"`
    DurationID           *uuid.UUID `gorm:"type:uuid" json:"duration_id"`
    CurrentTokenBalance  int        `gorm:"default:0" json:"current_token_balance"`
    CurrentPoinBalance   int        `gorm:"default:0" json:"current_poin_balance"`
    StartedAt            *time.Time `json:"started_at"`
    ExpiredAt            *time.Time `json:"expired_at"`
    IsActive             bool       `gorm:"default:true" json:"is_active"`
    AutoRenew            bool       `gorm:"default:false" json:"auto_renew"`
    CreatedAt            time.Time  `json:"created_at"`
    UpdatedAt            time.Time  `json:"updated_at"`
    
    // Relations
    Plan     *MembershipPlan     `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
    Duration *MembershipDuration `gorm:"foreignKey:DurationID" json:"duration,omitempty"`
}

// MembershipPlan Model
type MembershipPlan struct {
    ID                uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    PlanName          string         `gorm:"type:varchar(50);not null" json:"plan_name"`
    PlanCode          string         `gorm:"type:varchar(20);unique;not null" json:"plan_code"`
    MonthlyToken      int            `gorm:"not null" json:"monthly_token"`
    BaseMonthlyPrice  float64        `gorm:"type:decimal(10,2);not null" json:"base_monthly_price"`
    Description       string         `gorm:"type:text" json:"description"`
    Benefits          datatypes.JSON `gorm:"type:jsonb" json:"benefits,omitempty"` // OPTIONAL: untuk display UI saja, bukan control akses
    IsActive          bool           `gorm:"default:true" json:"is_active"`
    CreatedAt         time.Time      `json:"created_at"`
    UpdatedAt         time.Time      `json:"updated_at"`
}

// MembershipDuration Model
type MembershipDuration struct {
    ID                    uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    DurationMonths        int       `gorm:"not null" json:"duration_months"`
    DiscountPercentage    float64   `gorm:"type:decimal(5,2);default:0" json:"discount_percentage"`
    TokenBonusPercentage  float64   `gorm:"type:decimal(5,2);default:0" json:"token_bonus_percentage"`
    RextraPoinMultiplier  int       `gorm:"default:1" json:"rextra_poin_multiplier"`
    IsActive              bool      `gorm:"default:true" json:"is_active"`
    CreatedAt             time.Time `json:"created_at"`
    UpdatedAt             time.Time `json:"updated_at"`
}

// TokenTransaction Model
type TokenTransaction struct {
    ID                  uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    UserID              uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
    TransactionType     string     `gorm:"type:varchar(30);not null" json:"transaction_type"`
    TokenAmount         int        `gorm:"not null" json:"token_amount"`
    TokenBalanceBefore  int        `gorm:"not null" json:"token_balance_before"`
    TokenBalanceAfter   int        `gorm:"not null" json:"token_balance_after"`
    ReferenceType       string     `gorm:"type:varchar(50)" json:"reference_type"`
    ReferenceID         *uuid.UUID `gorm:"type:uuid" json:"reference_id"`
    Description         string     `gorm:"type:text" json:"description"`
    CreatedAt           time.Time  `json:"created_at"`
    
    // Relations
    Usage *TokenUsageHistory `gorm:"foreignKey:TokenTransactionID" json:"usage,omitempty"`
}

// TokenUsageHistory Model
type TokenUsageHistory struct {
    ID                   uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    UserID               uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
    TokenTransactionID   *uuid.UUID     `gorm:"type:uuid" json:"token_transaction_id"`
    FeatureName          string         `gorm:"type:varchar(100);not null" json:"feature_name"`
    TokenUsed            int            `gorm:"not null" json:"token_used"`
    UsageMetadata        datatypes.JSON `gorm:"type:jsonb" json:"usage_metadata"`
    CreatedAt            time.Time      `json:"created_at"`
    
    // Relations
    Transaction *TokenTransaction `gorm:"foreignKey:TokenTransactionID" json:"transaction,omitempty"`
}

// PoinTransaction Model
type PoinTransaction struct {
    ID                 uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    UserID             uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
    TransactionType    string     `gorm:"type:varchar(20);not null" json:"transaction_type"`
    PoinAmount         int        `gorm:"not null" json:"poin_amount"`
    PoinBalanceBefore  int        `gorm:"not null" json:"poin_balance_before"`
    PoinBalanceAfter   int        `gorm:"not null" json:"poin_balance_after"`
    Source             string     `gorm:"type:varchar(100)" json:"source"`
    ReferenceID        *uuid.UUID `gorm:"type:uuid" json:"reference_id"`
    Description        string     `gorm:"type:text" json:"description"`
    CreatedAt          time.Time  `json:"created_at"`
}

// PaymentTransaction Model
type PaymentTransaction struct {
    ID                  uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    UserID              uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
    PaymentType         string         `gorm:"type:varchar(30);not null" json:"payment_type"`
    PlanID              *uuid.UUID     `gorm:"type:uuid" json:"plan_id"`
    DurationID          *uuid.UUID     `gorm:"type:uuid" json:"duration_id"`
    TokenQuantity       *int           `json:"token_quantity"`
    GrossAmount         float64        `gorm:"type:decimal(12,2);not null" json:"gross_amount"`
    DiscountAmount      float64        `gorm:"type:decimal(12,2);default:0" json:"discount_amount"`
    FinalAmount         float64        `gorm:"type:decimal(12,2);not null" json:"final_amount"`
    PromoCode           string         `gorm:"type:varchar(50)" json:"promo_code"`
    XenditInvoiceID     string         `gorm:"type:varchar(255);unique" json:"xendit_invoice_id"`
    XenditExternalID    string         `gorm:"type:varchar(255);unique" json:"xendit_external_id"`
    PaymentMethod       string         `gorm:"type:varchar(50)" json:"payment_method"`
    PaymentStatus       string         `gorm:"type:varchar(20);not null" json:"payment_status"`
    PaidAt              *time.Time     `json:"paid_at"`
    ExpiredAt           *time.Time     `json:"expired_at"`
    XenditCallbackData  datatypes.JSON `gorm:"type:jsonb" json:"xendit_callback_data"`
    CreatedAt           time.Time      `json:"created_at"`
    UpdatedAt           time.Time      `json:"updated_at"`
    
    // Relations
    Plan     *MembershipPlan     `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
    Duration *MembershipDuration `gorm:"foreignKey:DurationID" json:"duration,omitempty"`
}

// PromoCode Model (Optional)
type PromoCode struct {
    ID                   uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    Code                 string         `gorm:"type:varchar(50);unique;not null" json:"code"`
    Description          string         `gorm:"type:text" json:"description"`
    DiscountType         string         `gorm:"type:varchar(20);not null" json:"discount_type"`
    DiscountValue        float64        `gorm:"type:decimal(10,2)" json:"discount_value"`
    ApplicablePlans      datatypes.JSON `gorm:"type:jsonb" json:"applicable_plans"`
    ApplicableDurations  datatypes.JSON `gorm:"type:jsonb" json:"applicable_durations"`
    MaxUsage             *int           `json:"max_usage"`
    CurrentUsage         int            `gorm:"default:0" json:"current_usage"`
    ValidFrom            time.Time      `gorm:"not null" json:"valid_from"`
    ValidUntil           time.Time      `gorm:"not null" json:"valid_until"`
    IsActive             bool           `gorm:"default:true" json:"is_active"`
    CreatedAt            time.Time      `json:"created_at"`
    UpdatedAt            time.Time      `json:"updated_at"`
}

// Table names
func (Membership) TableName() string          { return "memberships" }
func (MembershipPlan) TableName() string      { return "membership_plans" }
func (MembershipDuration) TableName() string  { return "membership_durations" }
func (TokenTransaction) TableName() string    { return "token_transactions" }
func (TokenUsageHistory) TableName() string   { return "token_usage_history" }
func (PoinTransaction) TableName() string     { return "poin_transactions" }
func (PaymentTransaction) TableName() string  { return "payment_transactions" }
func (PromoCode) TableName() string           { return "promo_codes" }
```

---

## PART 1 SELESAI ✅

**Next**: Part 2 akan cover API Endpoints Specification dengan request/response lengkap!

File ini sudah mencakup:
- ✅ Database schema 8 tabel (tanpa users)
- ✅ DDL SQL lengkap dengan indexes
- ✅ Data seed untuk master tables
- ✅ Business rules per tabel
- ✅ Relasi antar tabel
- ✅ Golang models dengan GORM

**Questions?** Konfirmasi dulu apakah ada yang perlu diperbaiki sebelum lanjut ke Part 2!