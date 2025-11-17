# 🔧 BRIEF BACKEND API REXTRA CLUB MEMBERSHIP
## PART 3: BUSINESS LOGIC & AUTO-TRANSITION

---

## Daftar Isi Part 3

1. [Alur Membership & Status User](#1-alur-membership--status-user)
2. [Membership Activation Logic](#2-membership-activation-logic-after-payment-success)
3. [Xendit Webhook Handler (Detail Implementation)](#3-xendit-webhook-handler-detail-implementation)
4. [Token Refill Logic (Monthly Scheduler)](#4-token-refill-logic-monthly-scheduler)
5. [Membership Expiration Logic](#5-membership-expiration-logic)
6. [Starter Plan Expiration Logic](#6-starter-plan-expiration-logic)
7. [Upgrade/Downgrade Calculation Logic](#7-upgradedowngrade-calculation-logic)
8. [Promo Code Validation Logic](#8-promo-code-validation-logic)
9. [Price Calculator Utility](#9-price-calculator-utility)

---

## 1. Alur Membership & Status User

### 1.1. Status Membership User

REXTRA CLUB memiliki **5 status membership** yang berbeda:

#### Status 1: **Starter** (Free Trial Plan)

**Fokus**: Masa percobaan 30 hari untuk user baru atau non-member yang membeli token standalone

**Cara Mendapatkan**:
1. User baru selesai registrasi + profiling Persona REXTRA → **Auto-activated Starter**
2. User dengan status Non-Member membeli token standalone → **Auto-activated Starter (30 hari baru)**

**Benefit**:
- Akses CRUD System unlimited (Persona, Rencana Karier, Kamus Karier, Eksplorasi, Portofolio tanpa AI, Event REXTRA)
- Akses terbatas ke fitur berbasis AI (tidak dapat token bulanan, harus beli token standalone)
- Dapat mendaftar & akses recording event REXTRA
- Akses katalog event REXTRA selamanya (bahkan setelah expired)

**Durasi**: 30 hari

**Setelah Expired**: Status berubah menjadi **Non-Member**

---

#### Status 2: **Basic** (REXTRA CLUB Paid Membership)

**Fokus**: Paket dasar untuk mahasiswa yang baru memulai persiapan karier

**Token per Bulan**: 10 token

**Harga**:
- 1 bulan: Rp15.000
- 3 bulan: Rp27.600 (efektif Rp9.200/bln)
- 6 bulan: Rp51.000 (efektif Rp8.500/bln)
- 12 bulan: Rp90.000 (efektif Rp7.500/bln)

**REXTRA POIN**:
- 1 bulan: 10 poin
- 3 bulan: 30 poin
- 6 bulan: 60 poin
- 12 bulan: 120 poin

**Benefit**:
- Semua benefit Starter
- 10 token/bulan untuk fitur AI (Kenali Diri, CV Generator, AI Interviewer)
- Akses REXTRA Academy kategori Basic
- Portofolio AI Check: 5 kegiatan/bulan
- Analisis penulisan data kegiatan: 5 kegiatan/bulan
- Analisis performa Rencana Karier: per semester

---

#### Status 3: **Pro** (REXTRA CLUB Paid Membership)

**Fokus**: Paket untuk mahasiswa aktif dengan kebutuhan AI moderat

**Token per Bulan**: 20 token

**Harga**:
- 1 bulan: Rp20.000
- 3 bulan: Rp55.200 (efektif Rp18.400/bln)
- 6 bulan: Rp102.000 (efektif Rp17.000/bln)
- 12 bulan: Rp180.000 (efektif Rp15.000/bln)

**REXTRA POIN**:
- 1 bulan: 20 poin
- 3 bulan: 60 poin
- 6 bulan: 120 poin
- 12 bulan: 240 poin

**Benefit**:
- Semua benefit Basic
- 20 token/bulan untuk fitur AI
- Akses penuh ke 6 modul pembelajaran karier digital
- Prioritas dalam penggunaan fitur Rencana Karier AI dan Portofolio AI Check
- Portofolio AI Check: 15 kegiatan/bulan
- Analisis penulisan data kegiatan: 15 kegiatan/bulan

---

#### Status 4: **Max** (REXTRA CLUB Paid Membership)

**Fokus**: Paket premium untuk mahasiswa intensif persiapan karier digital

**Token per Bulan**: 40 token

**Harga**:
- 1 bulan: Rp30.000
- 3 bulan: Rp82.800 (efektif Rp27.600/bln)
- 6 bulan: Rp153.000 (efektif Rp25.500/bln)
- 12 bulan: Rp252.000 (efektif Rp21.000/bln)

**REXTRA POIN**:
- 1 bulan: 30 poin
- 3 bulan: 90 poin
- 6 bulan: 180 poin
- 12 bulan: 360 poin

**Benefit**:
- Semua benefit Pro
- 40 token/bulan untuk fitur AI
- Akses penuh ke semua kategori modul edukasi
- Unlimited Portofolio AI Check
- Unlimited analisis penulisan data kegiatan
- Priority support dari tim REXTRA
- Akses ke semua fitur dan modul tanpa batasan

---

#### Status 5: **Non-Member** (No Active Plan)

**Fokus**: User yang tidak memiliki langganan aktif

**Cara Mendapatkan**:
- Status Starter expired (30 hari) dan tidak membeli membership
- Status membership (Basic/Pro/Max) expired dan tidak renewal

**Benefit**:
- Hanya dapat mengakses Katalog Event REXTRA (lihat event, detail, benefit)
- Tidak dapat mengakses fitur CRUD System dan fitur AI
- Dapat mendaftar event REXTRA yang dibuka untuk umum

**Cara Keluar dari Status Ini**:
- Membeli token standalone → dapat Starter 30 hari
- Membeli membership (Basic/Pro/Max) → langsung dapat status sesuai paket

---

### 1.2. Aturan Transisi Status Membership

```
[USER BARU]
    ↓
Registrasi + Profiling Persona
    ↓
[STARTER] (30 hari)
    ↓
    ├── Tidak beli apapun → [NON-MEMBER]
    ├── Beli Token Standalone → [STARTER] (30 hari baru)
    └── Beli Membership → [BASIC/PRO/MAX]

[NON-MEMBER]
    ↓
    ├── Beli Token Standalone → [STARTER] (30 hari)
    └── Beli Membership → [BASIC/PRO/MAX]

[BASIC/PRO/MAX] (Active)
    ↓
    ├── Expired & tidak renewal → [NON-MEMBER]
    ├── Renewal paket sama → [BASIC/PRO/MAX] (durasi bertambah)
    ├── Upgrade H-7 → [PRO/MAX] (dapat diskon 15%)
    └── Downgrade H-7 → [BASIC/PRO] (no diskon)

[BASIC/PRO/MAX] (Active) + Beli Token Standalone
    ↓
Token bertambah, status tetap [BASIC/PRO/MAX]
(TIDAK dapat Starter extension)
```

**Catatan Penting:**
- **Starter hanya diberikan** jika status user adalah **Non-Member**
- **Member aktif yang beli token** → token bertambah tapi **TIDAK dapat Starter**
- **Upgrade/Downgrade** hanya bisa dilakukan dalam periode **H-7** (7 hari sebelum expired)
- **Upgrade** → dapat diskon khusus 15%
- **Downgrade** → tidak dapat diskon, bayar full price

---

## 2. Membership Activation Logic (After Payment Success)

Setelah pembayaran berhasil via Xendit webhook, sistem akan melakukan aktivasi membership dengan langkah berikut:

### 2.1. Untuk Purchase Membership (Basic/Pro/Max)

**Step 1: Hitung Total Token dengan Bonus**

```go
func CalculateTotalToken(plan *MembershipPlan, duration *MembershipDuration) int {
  baseToken := plan.MonthlyToken * duration.DurationMonths
  bonusToken := float64(baseToken) * (duration.TokenBonusPercentage / 100)
  totalToken := baseToken + int(bonusToken)
  
  return totalToken
}
```

**Contoh Perhitungan:**
```
Plan: Pro (20 token/bulan)
Duration: 6 bulan (bonus 20%)

Base Token = 20 × 6 = 120 token
Bonus Token = 120 × 20% = 24 token
Total Token = 120 + 24 = 144 token
```

---

**Step 2: Hitung REXTRA POIN**

```go
func CalculateRextraPoin(plan *MembershipPlan, duration *MembershipDuration) int {
  basePoin := plan.MonthlyToken * duration.DurationMonths
  rextraPoin := basePoin * duration.RextraPoinMultiplier
  
  return rextraPoin
}
```

**Contoh Perhitungan:**
```
Plan: Pro (20 token/bulan)
Duration: 6 bulan (multiplier 6x)

Base Poin = 20 × 6 = 120
REXTRA POIN = 120 × 1 = 120 poin

Note: Multiplier sepertinya 1x untuk semua, bukan duration_months
Kalau mau pakai duration_months sebagai multiplier:
REXTRA POIN = (20 × 6) × 6 = 720 poin (terlalu banyak)

Sesuai spec document:
Basic 6 bulan = 60 poin → berarti (10 × 6) × 1 = 60 ✓
Pro 6 bulan = 120 poin → berarti (20 × 6) × 1 = 120 ✓
```

**⚠️ KLARIFIKASI DIPERLUKAN**: Apakah `rextra_poin_multiplier` di tabel `membership_durations` seharusnya selalu 1, atau pakai duration_months?

---

**Step 3: Update/Create Membership Record**

```go
func ActivateMembership(userID uuid.UUID, payment *PaymentTransaction) error {
  // 1. Get plan and duration from payment
  plan := payment.Plan
  duration := payment.Duration
  
  // 2. Calculate token and poin
  totalToken := CalculateTotalToken(plan, duration)
  rextraPoin := CalculateRextraPoin(plan, duration)
  
  // 3. Get or create membership
  membership, err := membershipRepo.GetByUserID(userID)
  if err != nil {
    // First time purchase - create new membership
    membership = &Membership{
      UserID: userID,
    }
  }
  
  // 4. Update membership fields
  now := time.Now()
  expiredAt := now.AddDate(0, duration.DurationMonths, 0)
  
  membership.MembershipStatus = strings.ToLower(plan.PlanCode) // basic/pro/max
  membership.PlanID = &plan.ID
  membership.DurationID = &duration.ID
  membership.CurrentTokenBalance += totalToken
  membership.CurrentPoinBalance += rextraPoin
  membership.StartedAt = &now
  membership.ExpiredAt = &expiredAt
  membership.IsActive = true
  
  // 5. Save membership
  if err := membershipRepo.Save(membership); err != nil {
    return err
  }
  
  // 6. Create token transaction record
  tokenTx := &TokenTransaction{
    UserID:             userID,
    TransactionType:    "purchase_membership",
    TokenAmount:        totalToken,
    TokenBalanceBefore: membership.CurrentTokenBalance - totalToken,
    TokenBalanceAfter:  membership.CurrentTokenBalance,
    ReferenceType:      "payment",
    ReferenceID:        &payment.ID,
    Description:        fmt.Sprintf("Token dari pembelian membership %s %d bulan", plan.PlanName, duration.DurationMonths),
  }
  tokenRepo.Create(tokenTx)
  
  // 7. Create poin transaction record
  poinTx := &PoinTransaction{
    UserID:            userID,
    TransactionType:   "earn",
    PoinAmount:        rextraPoin,
    PoinBalanceBefore: membership.CurrentPoinBalance - rextraPoin,
    PoinBalanceAfter:  membership.CurrentPoinBalance,
    Source:            "membership_purchase",
    ReferenceID:       &payment.ID,
    Description:       fmt.Sprintf("REXTRA POIN dari pembelian membership %s %d bulan", plan.PlanName, duration.DurationMonths),
  }
  poinRepo.Create(poinTx)
  
  // 8. Send notification to user
  notificationService.Send(userID, NotificationPayload{
    Title:   "Membership Aktif! 🎉",
    Message: fmt.Sprintf("Selamat! Membership %s Anda sudah aktif hingga %s", plan.PlanName, expiredAt.Format("02 Jan 2006")),
    Type:    "membership_activated",
  })
  
  return nil
}
```

---

### 2.2. Untuk Purchase Token Standalone

**Step 1: Add Token ke Membership**

```go
func ActivateTokenStandalone(userID uuid.UUID, payment *PaymentTransaction) error {
  // 1. Get membership
  membership, err := membershipRepo.GetByUserID(userID)
  if err != nil {
    return err
  }
  
  // 2. Add token
  tokenQuantity := *payment.TokenQuantity
  membership.CurrentTokenBalance += tokenQuantity
  
  // 3. Save membership
  membershipRepo.Save(membership)
  
  // 4. Create token transaction
  tokenTx := &TokenTransaction{
    UserID:             userID,
    TransactionType:    "purchase_standalone",
    TokenAmount:        tokenQuantity,
    TokenBalanceBefore: membership.CurrentTokenBalance - tokenQuantity,
    TokenBalanceAfter:  membership.CurrentTokenBalance,
    ReferenceType:      "payment",
    ReferenceID:        &payment.ID,
    Description:        fmt.Sprintf("Pembelian %d token standalone", tokenQuantity),
  }
  tokenRepo.Create(tokenTx)
  
  // 5. CEK membership_status user
  isNonMember := membership.MembershipStatus == "non_member"
  
  if isNonMember {
    // 6. Activate Starter Plan (30 hari)
    now := time.Now()
    expiredAt := now.AddDate(0, 0, 30) // +30 hari
    
    membership.MembershipStatus = "starter"
    membership.PlanID = nil
    membership.DurationID = nil
    membership.StartedAt = &now
    membership.ExpiredAt = &expiredAt
    membership.IsActive = true
    
    membershipRepo.Save(membership)
    
    // 7. Send notification with Starter info
    notificationService.Send(userID, NotificationPayload{
      Title:   "Token Berhasil Dibeli! 🎁",
      Message: fmt.Sprintf("Anda mendapat %d token + Bonus Starter Plan 30 hari!", tokenQuantity),
      Type:    "token_purchased_with_starter",
    })
  } else {
    // 8. User sudah Member - skip Starter, hanya tambah token
    notificationService.Send(userID, NotificationPayload{
      Title:   "Token Berhasil Ditambahkan!",
      Message: fmt.Sprintf("%d token telah ditambahkan ke saldo Anda", tokenQuantity),
      Type:    "token_purchased",
    })
  }
  
  return nil
}
```

**Flow Decision:**

```
User beli token standalone
    ↓
Cek membership_status
    ↓
    ├── Non-Member
    │       ↓
    │   Add token + Activate Starter (30 hari)
    │
    └── Starter/Basic/Pro/Max (sudah punya membership aktif)
            ↓
        Add token saja (NO Starter extension)
```

---

## 3. Xendit Webhook Handler (Detail Implementation)

Webhook ini adalah **jantung** dari sistem payment. Saat Xendit mengirim notifikasi pembayaran sukses, sistem harus auto-activate membership/token.

### 3.1. Webhook Security

```go
func (h *WebhookHandler) XenditCallback(c *gin.Context) {
  // 1. Verify Xendit callback token untuk security
  callbackToken := c.GetHeader("X-CALLBACK-TOKEN")
  expectedToken := config.Get("XENDIT_CALLBACK_TOKEN")
  
  if callbackToken != expectedToken {
    log.Error("Invalid Xendit callback token")
    c.JSON(401, gin.H{"error": "Unauthorized"})
    return
  }
  
  // Continue processing...
}
```

---

### 3.2. Parse Xendit Payload

```go
type XenditInvoiceCallback struct {
  ID             string    `json:"id"`
  ExternalID     string    `json:"external_id"`
  UserID         string    `json:"user_id"`
  Status         string    `json:"status"` // PAID, EXPIRED, FAILED
  PaidAmount     float64   `json:"paid_amount"`
  PaidAt         time.Time `json:"paid_at"`
  PaymentMethod  string    `json:"payment_method"`
  PaymentChannel string    `json:"payment_channel"`
  Description    string    `json:"description"`
  InvoiceURL     string    `json:"invoice_url"`
  Updated        time.Time `json:"updated"`
  Created        time.Time `json:"created"`
  Currency       string    `json:"currency"`
  PaymentID      string    `json:"payment_id"`
}

func (h *WebhookHandler) XenditCallback(c *gin.Context) {
  // ... (after security check)
  
  // 2. Bind request body
  var callback XenditInvoiceCallback
  if err := c.ShouldBindJSON(&callback); err != nil {
    log.Error("Failed to parse Xendit callback", err)
    c.JSON(400, gin.H{"error": "Invalid payload"})
    return
  }
  
  log.Info("Received Xendit callback", map[string]interface{}{
    "invoice_id":  callback.ID,
    "external_id": callback.ExternalID,
    "status":      callback.Status,
    "amount":      callback.PaidAmount,
  })
  
  // Continue processing...
}
```

---

### 3.3. Find Payment Transaction

```go
func (h *WebhookHandler) XenditCallback(c *gin.Context) {
  // ... (after parsing)
  
  // 3. Find payment_transaction by xendit_invoice_id atau xendit_external_id
  payment, err := h.paymentService.GetByXenditInvoiceID(callback.ID)
  if err != nil {
    // Try by external_id
    payment, err = h.paymentService.GetByXenditExternalID(callback.ExternalID)
    if err != nil {
      log.Error("Payment transaction not found", map[string]interface{}{
        "invoice_id":  callback.ID,
        "external_id": callback.ExternalID,
      })
      c.JSON(404, gin.H{"error": "Transaction not found"})
      return
    }
  }
  
  // Continue processing...
}
```

---

### 3.4. Handle Payment Status

```go
func (h *WebhookHandler) XenditCallback(c *gin.Context) {
  // ... (after finding payment)
  
  // 4. Handle berdasarkan status
  switch callback.Status {
  case "PAID":
    h.handlePaidStatus(payment, &callback)
  case "EXPIRED":
    h.handleExpiredStatus(payment, &callback)
  case "FAILED":
    h.handleFailedStatus(payment, &callback)
  default:
    log.Info("Unhandled payment status", callback.Status)
  }
  
  // 5. Always return 200 OK to Xendit
  c.JSON(200, gin.H{
    "status":  "success",
    "message": "Webhook processed successfully",
  })
}
```

---

### 3.5. Handle PAID Status (Most Important)

```go
func (h *WebhookHandler) handlePaidStatus(payment *PaymentTransaction, callback *XenditInvoiceCallback) {
  // 1. Check idempotency - jangan double process
  if payment.PaymentStatus == "paid" {
    log.Warn("Payment already processed", payment.ID)
    return
  }
  
  // 2. Update payment_transaction
  payment.PaymentStatus = "paid"
  payment.PaidAt = &callback.PaidAt
  payment.PaymentMethod = fmt.Sprintf("%s_%s", callback.PaymentMethod, callback.PaymentChannel)
  
  // Save raw callback data untuk audit
  callbackJSON, _ := json.Marshal(callback)
  payment.XenditCallbackData = datatypes.JSON(callbackJSON)
  
  if err := h.paymentService.UpdatePaymentTransaction(payment); err != nil {
    log.Error("Failed to update payment transaction", err)
    return
  }
  
  // 3. Activate Membership/Token berdasarkan payment_type
  if payment.PaymentType == "membership" {
    h.activateMembershipPayment(payment)
  } else if payment.PaymentType == "token_standalone" {
    h.activateTokenStandalonePayment(payment)
  }
}

func (h *WebhookHandler) activateMembershipPayment(payment *PaymentTransaction) {
  // Use the ActivateMembership function from section 2.1
  if err := h.membershipService.ActivateMembership(payment.UserID, payment); err != nil {
    log.Error("Failed to activate membership", err)
    
    // CRITICAL: Rollback payment status jika gagal activate
    payment.PaymentStatus = "pending"
    h.paymentService.UpdatePaymentTransaction(payment)
    
    // Send alert to admin
    h.alertService.SendCriticalAlert("Membership activation failed", map[string]interface{}{
      "payment_id": payment.ID,
      "user_id":    payment.UserID,
      "error":      err.Error(),
    })
    
    return
  }
  
  log.Info("Membership activated successfully", map[string]interface{}{
    "user_id":  payment.UserID,
    "plan":     payment.Plan.PlanCode,
    "duration": payment.Duration.DurationMonths,
  })
}

func (h *WebhookHandler) activateTokenStandalonePayment(payment *PaymentTransaction) {
  // Use the ActivateTokenStandalone function from section 2.2
  if err := h.tokenService.ActivateTokenStandalone(payment.UserID, payment); err != nil {
    log.Error("Failed to activate token standalone", err)
    
    // CRITICAL: Rollback payment status
    payment.PaymentStatus = "pending"
    h.paymentService.UpdatePaymentTransaction(payment)
    
    // Send alert to admin
    h.alertService.SendCriticalAlert("Token standalone activation failed", map[string]interface{}{
      "payment_id":     payment.ID,
      "user_id":        payment.UserID,
      "token_quantity": *payment.TokenQuantity,
      "error":          err.Error(),
    })
    
    return
  }
  
  log.Info("Token standalone activated successfully", map[string]interface{}{
    "user_id":        payment.UserID,
    "token_quantity": *payment.TokenQuantity,
  })
}
```

---

### 3.6. Handle EXPIRED/FAILED Status

```go
func (h *WebhookHandler) handleExpiredStatus(payment *PaymentTransaction, callback *XenditInvoiceCallback) {
  // Update payment status
  payment.PaymentStatus = "expired"
  callbackJSON, _ := json.Marshal(callback)
  payment.XenditCallbackData = datatypes.JSON(callbackJSON)
  h.paymentService.UpdatePaymentTransaction(payment)
  
  // Send notification to user
  h.notificationService.Send(payment.UserID, NotificationPayload{
    Title:   "Pembayaran Expired",
    Message: "Invoice pembayaran Anda telah expired. Silakan coba lagi untuk melanjutkan.",
    Type:    "payment_expired",
  })
  
  log.Info("Payment expired", payment.ID)
}

func (h *WebhookHandler) handleFailedStatus(payment *PaymentTransaction, callback *XenditInvoiceCallback) {
  // Update payment status
  payment.PaymentStatus = "failed"
  callbackJSON, _ := json.Marshal(callback)
  payment.XenditCallbackData = datatypes.JSON(callbackJSON)
  h.paymentService.UpdatePaymentTransaction(payment)
  
  // Send notification to user
  h.notificationService.Send(payment.UserID, NotificationPayload{
    Title:   "Pembayaran Gagal",
    Message: "Pembayaran Anda gagal diproses. Silakan coba lagi atau hubungi customer support.",
    Type:    "payment_failed",
  })
  
  log.Info("Payment failed", payment.ID)
}
```

---

### 3.7. Complete Webhook Handler Flow

```go
func (h *WebhookHandler) XenditCallback(c *gin.Context) {
  // 1. Verify callback token
  callbackToken := c.GetHeader("X-CALLBACK-TOKEN")
  if callbackToken != config.Get("XENDIT_CALLBACK_TOKEN") {
    log.Error("Invalid Xendit callback token")
    c.JSON(401, gin.H{"error": "Unauthorized"})
    return
  }
  
  // 2. Parse callback payload
  var callback XenditInvoiceCallback
  if err := c.ShouldBindJSON(&callback); err != nil {
    log.Error("Failed to parse Xendit callback", err)
    c.JSON(400, gin.H{"error": "Invalid payload"})
    return
  }
  
  // 3. Log received callback
  log.Info("Received Xendit callback", map[string]interface{}{
    "invoice_id":  callback.ID,
    "external_id": callback.ExternalID,
    "status":      callback.Status,
    "amount":      callback.PaidAmount,
  })
  
  // 4. Find payment transaction
  payment, err := h.paymentService.GetByXenditInvoiceID(callback.ID)
  if err != nil {
    payment, err = h.paymentService.GetByXenditExternalID(callback.ExternalID)
    if err != nil {
      log.Error("Payment transaction not found", err)
      c.JSON(404, gin.H{"error": "Transaction not found"})
      return
    }
  }
  
  // 5. Handle status
  switch callback.Status {
  case "PAID":
    h.handlePaidStatus(payment, &callback)
  case "EXPIRED":
    h.handleExpiredStatus(payment, &callback)
  case "FAILED":
    h.handleFailedStatus(payment, &callback)
  default:
    log.Info("Unhandled payment status", callback.Status)
  }
  
  // 6. Always return 200 OK
  c.JSON(200, gin.H{
    "status":  "success",
    "message": "Webhook processed successfully",
  })
}
```

**Important Notes:**
- **Idempotency**: Selalu cek `payment.PaymentStatus` sebelum process untuk avoid double activation
- **Error Handling**: Jika activation gagal, rollback payment status ke "pending"
- **Logging**: Log semua webhook events untuk audit trail
- **Alert**: Send critical alert ke admin jika activation gagal
- **Return 200**: Selalu return 200 OK ke Xendit, even if error, supaya Xendit stop retry

---

## 4. Token Refill Logic (Monthly Scheduler)

Setiap awal bulan (tanggal 1), sistem akan auto-refill token untuk semua member aktif.

### 4.1. Scheduled Job Setup

```go
// Use cron library: github.com/robfig/cron/v3

func InitScheduledJobs() {
  c := cron.New()
  
  // Run every 1st of month at 00:00
  c.AddFunc("0 0 1 * *", MonthlyTokenRefill)
  
  // Run every day at 00:00 for expiration checks
  c.AddFunc("0 0 * * *", CheckExpiredMemberships)
  c.AddFunc("0 0 * * *", CheckExpiredStarterPlan)
  
  // Run every day at 09:00 for reminders
  c.AddFunc("0 9 * * *", SendExpirationReminders)
  
  c.Start()
  
  log.Info("Scheduled jobs initialized")
}
```

---

### 4.2. Monthly Token Refill Implementation

```go
func MonthlyTokenRefill() {
  log.Info("Starting monthly token refill job")
  
  // 1. Query semua membership aktif dengan paid plan
  memberships, err := membershipRepo.GetActiveMembers([]string{"basic", "pro", "max"})
  if err != nil {
    log.Error("Failed to get active members", err)
    return
  }
  
  log.Info(fmt.Sprintf("Found %d active members for token refill", len(memberships)))
  
  successCount := 0
  failCount := 0
  
  // 2. Loop setiap membership
  for _, membership := range memberships {
    // 3. Get monthly_token dari plan
    plan := membership.Plan
    if plan == nil {
      log.Warn("Membership has no plan", membership.ID)
      failCount++
      continue
    }
    
    monthlyToken := plan.MonthlyToken
    
    // 4. Check if membership expired (safety check)
    if membership.ExpiredAt != nil && time.Now().After(*membership.ExpiredAt) {
      log.Warn("Skipping expired membership", membership.ID)
      failCount++
      continue
    }
    
    // 5. Add token
    tokenBalanceBefore := membership.CurrentTokenBalance
    membership.CurrentTokenBalance += monthlyToken
    
    if err := membershipRepo.Update(membership); err != nil {
      log.Error("Failed to update membership token balance", err)
      failCount++
      continue
    }
    
    // 6. Create token_transaction record
    tokenTx := &TokenTransaction{
      UserID:             membership.UserID,
      TransactionType:    "monthly_refill",
      TokenAmount:        monthlyToken,
      TokenBalanceBefore: tokenBalanceBefore,
      TokenBalanceAfter:  membership.CurrentTokenBalance,
      ReferenceType:      "monthly_refill",
      ReferenceID:        nil,
      Description:        fmt.Sprintf("Refill token bulanan untuk paket %s", plan.PlanName),
    }
    
    if err := tokenRepo.Create(tokenTx); err != nil {
      log.Error("Failed to create token transaction", err)
      failCount++
      continue
    }
    
    // 7. Send notification to user
    notificationService.Send(membership.UserID, NotificationPayload{
      Title:   "Token Bulanan Telah Di-refill! 🎁",
      Message: fmt.Sprintf("Anda mendapat %d token untuk bulan ini. Saldo token Anda sekarang: %d", monthlyToken, membership.CurrentTokenBalance),
      Type:    "token_refill",
    })
    
    successCount++
    
    log.Info("Token refilled successfully", map[string]interface{}{
      "user_id":      membership.UserID,
      "plan":         plan.PlanCode,
      "token_added":  monthlyToken,
      "total_tokens": membership.CurrentTokenBalance,
    })
  }
  
  log.Info("Monthly token refill completed", map[string]interface{}{
    "total_members":   len(memberships),
    "success_count":   successCount,
    "fail_count":      failCount,
  })
  
  // 8. Send summary report to admin
  alertService.SendInfoAlert("Monthly Token Refill Completed", map[string]interface{}{
    "total_members": len(memberships),
    "success":       successCount,
    "failed":        failCount,
    "timestamp":     time.Now(),
  })
}
```

**Key Points:**
- Run setiap **tanggal 1 pukul 00:00**
- Hanya refill untuk membership dengan status `basic/pro/max` yang `is_active = true`
- Skip membership yang sudah expired (safety check)
- Create `token_transaction` record dengan type `monthly_refill`
- Send notification ke user
- Send summary report ke admin

---

## 5. Membership Expiration Logic

Sistem akan check membership yang expired setiap hari dan update statusnya.

### 5.1. Check Expired Memberships

```go
func CheckExpiredMemberships() {
  log.Info("Starting expired membership check")
  
  // 1. Query semua membership yang expired
  now := time.Now()
  expiredMemberships, err := membershipRepo.GetExpiredMemberships(now, []string{"basic", "pro", "max"})
  if err != nil {
    log.Error("Failed to get expired memberships", err)
    return
  }
  
  if len(expiredMemberships) == 0 {
    log.Info("No expired memberships found")
    return
  }
  
  log.Info(fmt.Sprintf("Found %d expired memberships", len(expiredMemberships)))
  
  // 2. Loop setiap expired membership
  for _, membership := range expiredMemberships {
    // 3. Update membership ke Non-Member
    membership.IsActive = false
    membership.MembershipStatus = "non_member"
    membership.PlanID = nil
    membership.DurationID = nil
    // Note: Token dan POIN TIDAK dihapus, tetap ada
    
    if err := membershipRepo.Update(membership); err != nil {
      log.Error("Failed to update expired membership", err)
      continue
    }
    
    // 4. Send notification ke user
    notificationService.Send(membership.UserID, NotificationPayload{
      Title:   "Membership Berakhir",
      Message: "Membership Anda telah berakhir. Perpanjang sekarang untuk tetap menikmati benefit REXTRA CLUB!",
      Type:    "membership_expired",
      CTA: map[string]string{
        "text": "Perpanjang Sekarang",
        "url":  "/membership/plans",
      },
    })
    
    log.Info("Membership expired and updated to non-member", map[string]interface{}{
      "user_id":         membership.UserID,
      "expired_at":      membership.ExpiredAt,
      "previous_status": membership.MembershipStatus,
    })
  }
  
  log.Info("Expired membership check completed", map[string]interface{}{
    "expired_count": len(expiredMemberships),
  })
}
```

**Query for GetExpiredMemberships:**

```go
func (r *MembershipRepository) GetExpiredMemberships(now time.Time, statuses []string) ([]*Membership, error) {
  var memberships []*Membership
  
  err := r.db.
    Preload("Plan").
    Preload("Duration").
    Where("is_active = ?", true).
    Where("membership_status IN ?", statuses).
    Where("expired_at <= ?", now).
    Find(&memberships).Error
  
  return memberships, err
}
```

---

### 5.2. Send Expiration Reminders

Send reminder ke user sebelum membership expired (H-7, H-3, H-1).

```go
func SendExpirationReminders() {
  log.Info("Starting expiration reminder job")
  
  now := time.Now()
  
  // 1. Reminder H-7
  send7DayReminders(now)
  
  // 2. Reminder H-3
  send3DayReminders(now)
  
  // 3. Reminder H-1
  send1DayReminders(now)
  
  log.Info("Expiration reminder job completed")
}

func send7DayReminders(now time.Time) {
  from := now.Add(7 * 24 * time.Hour)
  to := from.Add(24 * time.Hour)
  
  memberships, _ := membershipRepo.GetExpiringBetween(from, to, []string{"basic", "pro", "max"})
  
  for _, membership := range memberships {
    notificationService.Send(membership.UserID, NotificationPayload{
      Title:   "Membership akan Berakhir dalam 7 Hari",
      Message: fmt.Sprintf("Membership %s Anda akan berakhir pada %s. Perpanjang sekarang untuk tetap mendapat benefit!", membership.Plan.PlanName, membership.ExpiredAt.Format("02 Jan 2006")),
      Type:    "membership_expiring_7d",
      CTA: map[string]string{
        "text": "Perpanjang Sekarang",
        "url":  "/membership/renew",
      },
    })
    
    log.Info("Sent 7-day reminder", membership.UserID)
  }
}

func send3DayReminders(now time.Time) {
  from := now.Add(3 * 24 * time.Hour)
  to := from.Add(24 * time.Hour)
  
  memberships, _ := membershipRepo.GetExpiringBetween(from, to, []string{"basic", "pro", "max"})
  
  for _, membership := range memberships {
    notificationService.Send(membership.UserID, NotificationPayload{
      Title:   "⚠️ Membership akan Berakhir dalam 3 Hari!",
      Message: fmt.Sprintf("Jangan lewatkan benefit eksklusif! Membership Anda berakhir %s", membership.ExpiredAt.Format("02 Jan 2006")),
      Type:    "membership_expiring_3d",
      CTA: map[string]string{
        "text": "Perpanjang Sekarang",
        "url":  "/membership/renew",
      },
    })
    
    log.Info("Sent 3-day reminder", membership.UserID)
  }
}

func send1DayReminders(now time.Time) {
  from := now.Add(24 * time.Hour)
  to := from.Add(24 * time.Hour)
  
  memberships, _ := membershipRepo.GetExpiringBetween(from, to, []string{"basic", "pro", "max"})
  
  for _, membership := range memberships {
    notificationService.Send(membership.UserID, NotificationPayload{
      Title:   "🚨 TERAKHIR! Membership Berakhir Besok",
      Message: "Membership Anda berakhir besok! Perpanjang sekarang untuk tidak kehilangan akses ke semua fitur premium.",
      Type:    "membership_expiring_1d",
      CTA: map[string]string{
        "text": "Perpanjang Sekarang",
        "url":  "/membership/renew",
      },
    })
    
    log.Info("Sent 1-day reminder", membership.UserID)
  }
}
```

**Query for GetExpiringBetween:**

```go
func (r *MembershipRepository) GetExpiringBetween(from, to time.Time, statuses []string) ([]*Membership, error) {
  var memberships []*Membership
  
  err := r.db.
    Preload("Plan").
    Where("is_active = ?", true).
    Where("membership_status IN ?", statuses).
    Where("expired_at >= ? AND expired_at < ?", from, to).
    Find(&memberships).Error
  
  return memberships, err
}
```

---

## 6. Starter Plan Expiration Logic

Similar dengan membership expiration, tapi khusus untuk Starter Plan.

```go
func CheckExpiredStarterPlan() {
  log.Info("Starting expired starter plan check")
  
  // 1. Query semua Starter yang expired
  now := time.Now()
  expiredStarters, err := membershipRepo.GetExpiredMemberships(now, []string{"starter"})
  if err != nil {
    log.Error("Failed to get expired starters", err)
    return
  }
  
  if len(expiredStarters) == 0 {
    log.Info("No expired starter plans found")
    return
  }
  
  log.Info(fmt.Sprintf("Found %d expired starter plans", len(expiredStarters)))
  
  // 2. Loop setiap expired starter
  for _, membership := range expiredStarters {
    // 3. Update membership ke Non-Member
    membership.IsActive = false
    membership.MembershipStatus = "non_member"
    membership.ExpiredAt = nil
    membership.StartedAt = nil
    // Token balance tetap (tidak dihapus)
    
    if err := membershipRepo.Update(membership); err != nil {
      log.Error("Failed to update expired starter", err)
      continue
    }
    
    // 4. Send notification ke user
    notificationService.Send(membership.UserID, NotificationPayload{
      Title:   "Starter Plan Berakhir",
      Message: "Periode Starter Plan Anda telah berakhir. Bergabung dengan REXTRA CLUB untuk akses penuh ke semua fitur!",
      Type:    "starter_expired",
      CTA: map[string]string{
        "text": "Lihat Paket Membership",
        "url":  "/membership/plans",
      },
    })
    
    log.Info("Starter plan expired", membership.UserID)
  }
  
  log.Info("Expired starter plan check completed", map[string]interface{}{
    "expired_count": len(expiredStarters),
  })
}
```

**Key Difference dari Membership Expiration:**
- Starter Plan: `expired_at` dan `started_at` di-set NULL setelah expired
- Membership: Hanya update status, tapi `expired_at` tetap ada (untuk history)

---

## 7. Upgrade/Downgrade Calculation Logic

### 7.1. Upgrade Price Calculation (dengan 15% Discount)

```go
func CalculateUpgradePrice(currentPlan, newPlan *MembershipPlan, duration *MembershipDuration) *PricingResult {
  // 1. Hitung base price
  basePrice := newPlan.BaseMonthlyPrice * float64(duration.DurationMonths)
  
  // 2. Apply plan discount
  planDiscount := basePrice * (duration.DiscountPercentage / 100)
  priceAfterPlanDiscount := basePrice - planDiscount
  
  // 3. Apply upgrade discount (15% extra discount)
  upgradeDiscountPercentage := 15.0
  upgradeDiscount := priceAfterPlanDiscount * (upgradeDiscountPercentage / 100)
  finalPrice := priceAfterPlanDiscount - upgradeDiscount
  
  // 4. Calculate tokens
  baseToken := newPlan.MonthlyToken * duration.DurationMonths
  totalToken := int(float64(baseToken) * (1 + duration.TokenBonusPercentage/100))
  
  // 5. Calculate REXTRA POIN
  rextraPoin := baseToken * duration.RextraPoinMultiplier
  
  return &PricingResult{
    BasePrice:              basePrice,
    PlanDiscount:           planDiscount,
    UpgradeDiscount:        upgradeDiscount,
    TotalDiscount:          planDiscount + upgradeDiscount,
    FinalPrice:             finalPrice,
    EffectiveMonthlyPrice:  finalPrice / float64(duration.DurationMonths),
    TotalToken:             totalToken,
    RextraPoin:             rextraPoin,
  }
}
```

**Contoh Perhitungan Upgrade:**
```
Current Plan: Pro
New Plan: Max
Duration: 6 bulan

Base Price = Rp30.000 × 6 = Rp180.000
Plan Discount (15%) = Rp180.000 × 15% = Rp27.000
Price After Plan Discount = Rp180.000 - Rp27.000 = Rp153.000
Upgrade Discount (15%) = Rp153.000 × 15% = Rp22.950
Final Price = Rp153.000 - Rp22.950 = Rp130.050

Total Discount = Rp27.000 + Rp22.950 = Rp49.950 (27.75%)
```

---

### 7.2. Downgrade Price Calculation (NO Discount)

```go
func CalculateDowngradePrice(currentPlan, newPlan *MembershipPlan, duration *MembershipDuration) *PricingResult {
  // 1. Hitung base price
  basePrice := newPlan.BaseMonthlyPrice * float64(duration.DurationMonths)
  
  // 2. Apply plan discount SAJA (NO upgrade discount)
  planDiscount := basePrice * (duration.DiscountPercentage / 100)
  finalPrice := basePrice - planDiscount
  
  // 3. Calculate tokens
  baseToken := newPlan.MonthlyToken * duration.DurationMonths
  totalToken := int(float64(baseToken) * (1 + duration.TokenBonusPercentage/100))
  
  // 4. Calculate REXTRA POIN
  rextraPoin := baseToken * duration.RextraPoinMultiplier
  
  return &PricingResult{
    BasePrice:              basePrice,
    PlanDiscount:           planDiscount,
    UpgradeDiscount:        0, // NO UPGRADE DISCOUNT FOR DOWNGRADE
    TotalDiscount:          planDiscount,
    FinalPrice:             finalPrice,
    EffectiveMonthlyPrice:  finalPrice / float64(duration.DurationMonths),
    TotalToken:             totalToken,
    RextraPoin:             rextraPoin,
  }
}
```

**Contoh Perhitungan Downgrade:**
```
Current Plan: Pro
New Plan: Basic
Duration: 3 bulan

Base Price = Rp15.000 × 3 = Rp45.000
Plan Discount (8%) = Rp45.000 × 8% = Rp3.600
Final Price = Rp45.000 - Rp3.600 = Rp41.400

NO UPGRADE DISCOUNT
Total Discount = Rp3.600 (8% only)
```

---

### 7.3. Validation Logic for Upgrade/Downgrade

```go
func ValidateUpgradeDowngrade(membership *Membership, newPlanCode string) error {
  // 1. Check if in H-7 period
  daysRemaining := int(time.Until(*membership.ExpiredAt).Hours() / 24)
  if daysRemaining > 7 {
    return errors.New("Upgrade/Downgrade hanya bisa dilakukan dalam 7 hari sebelum membership expired")
  }
  
  // 2. Define plan hierarchy
  planHierarchy := map[string]int{
    "basic": 1,
    "pro":   2,
    "max":   3,
  }
  
  currentLevel := planHierarchy[membership.MembershipStatus]
  newLevel := planHierarchy[strings.ToLower(newPlanCode)]
  
  // 3. Validate upgrade
  if newLevel > currentLevel {
    // This is upgrade - OK
    return nil
  }
  
  // 4. Validate downgrade
  if newLevel < currentLevel {
    // This is downgrade
    if membership.MembershipStatus == "basic" {
      return errors.New("Tidak bisa downgrade dari paket Basic (sudah paket terendah)")
    }
    // Downgrade OK
    return nil
  }
  
  // 5. Same level - not allowed
  return errors.New("Paket baru harus berbeda dari paket saat ini")
}
```

---

## 8. Promo Code Validation Logic

```go
func ValidatePromoCode(promoCode string, planCode string, durationMonths int) (*PromoCode, error) {
  // 1. Find promo code (case-insensitive)
  promo, err := promoRepo.GetByCode(strings.ToUpper(promoCode))
  if err != nil {
    return nil, errors.New("Promo code tidak ditemukan")
  }
  
  // 2. Check is_active
  if !promo.IsActive {
    return nil, errors.New("Promo code tidak aktif")
  }
  
  // 3. Check valid date
  now := time.Now()
  if now.Before(promo.ValidFrom) {
    return nil, errors.New("Promo code belum berlaku")
  }
  if now.After(promo.ValidUntil) {
    return nil, errors.New("Promo code sudah expired")
  }
  
  // 4. Check max usage
  if promo.MaxUsage != nil && promo.CurrentUsage >= *promo.MaxUsage {
    return nil, errors.New("Promo code sudah mencapai batas penggunaan maksimal")
  }
  
  // 5. Check applicable plans
  if promo.ApplicablePlans != nil {
    var applicablePlans []string
    json.Unmarshal(promo.ApplicablePlans, &applicablePlans)
    
    if len(applicablePlans) > 0 && !contains(applicablePlans, strings.ToUpper(planCode)) {
      return nil, errors.New(fmt.Sprintf("Promo code tidak berlaku untuk paket %s", planCode))
    }
  }
  
  // 6. Check applicable durations
  if promo.ApplicableDurations != nil {
    var applicableDurations []int
    json.Unmarshal(promo.ApplicableDurations, &applicableDurations)
    
    if len(applicableDurations) > 0 && !containsInt(applicableDurations, durationMonths) {
      return nil, errors.New(fmt.Sprintf("Promo code tidak berlaku untuk durasi %d bulan", durationMonths))
    }
  }
  
  // 7. Promo valid - return
  return promo, nil
}

func CalculatePromoDiscount(promo *PromoCode, priceAfterPlanDiscount float64) float64 {
  switch promo.DiscountType {
  case "percentage":
    // Percentage discount (e.g. 50% off)
    return priceAfterPlanDiscount * (promo.DiscountValue / 100)
    
  case "fixed_amount":
    // Fixed amount discount (e.g. Rp20.000 off)
    discount := promo.DiscountValue
    if discount > priceAfterPlanDiscount {
      // Discount tidak boleh lebih besar dari harga
      return priceAfterPlanDiscount
    }
    return discount
    
  case "free_100":
    // 100% free
    return priceAfterPlanDiscount
    
  default:
    return 0
  }
}

func IncrementPromoUsage(promoCode string) error {
  promo, err := promoRepo.GetByCode(promoCode)
  if err != nil {
    return err
  }
  
  promo.CurrentUsage++
  return promoRepo.Update(promo)
}
```

**Usage in Purchase Flow:**

```go
func (s *MembershipService) PurchaseMembership(req *PurchaseRequest) (*PaymentTransaction, error) {
  // ... (validate plan, duration)
  
  var promoDiscount float64
  var promoApplied *PromoCode
  
  // Validate promo code if provided
  if req.PromoCode != "" {
    promo, err := s.promoService.ValidatePromoCode(req.PromoCode, req.PlanCode, req.DurationMonths)
    if err != nil {
      return nil, err
    }
    
    // Calculate promo discount
    pricing := CalculatePricing(plan, duration)
    promoDiscount = CalculatePromoDiscount(promo, pricing.PriceAfterPlanDiscount)
    promoApplied = promo
  }
  
  // ... (create payment transaction)
  
  // Increment promo usage ONLY after payment success
  // This is done in webhook handler after PAID status
  
  return payment, nil
}

// In Webhook Handler - After Payment Success
func (h *WebhookHandler) handlePaidStatus(payment *PaymentTransaction, callback *XenditCallback) {
  // ... (update payment, activate membership)
  
  // Increment promo usage
  if payment.PromoCode != "" {
    if err := h.promoService.IncrementPromoUsage(payment.PromoCode); err != nil {
      log.Error("Failed to increment promo usage", err)
      // Not critical - don't rollback payment
    }
  }
}
```

---

## 9. Price Calculator Utility

Central utility untuk semua perhitungan harga, token, dan poin.

```go
package calculator

type PricingResult struct {
  BasePrice              float64
  PlanDiscount           float64
  PromoDiscount          float64
  UpgradeDiscount        float64
  TotalDiscount          float64
  FinalPrice             float64
  EffectiveMonthlyPrice  float64
  PriceAfterPlanDiscount float64
  TotalToken             int
  RextraPoin             int
}

type PriceCalculator struct {
  // Dependencies if needed
}

func NewPriceCalculator() *PriceCalculator {
  return &PriceCalculator{}
}

// Calculate standard membership pricing
func (c *PriceCalculator) CalculatePricing(plan *MembershipPlan, duration *MembershipDuration) *PricingResult {
  // 1. Base price
  basePrice := plan.BaseMonthlyPrice * float64(duration.DurationMonths)
  
  // 2. Plan discount
  planDiscount := basePrice * (duration.DiscountPercentage / 100)
  priceAfterPlanDiscount := basePrice - planDiscount
  
  // 3. Final price (no promo, no upgrade discount)
  finalPrice := priceAfterPlanDiscount
  
  // 4. Calculate tokens
  baseToken := plan.MonthlyToken * duration.DurationMonths
  bonusToken := float64(baseToken) * (duration.TokenBonusPercentage / 100)
  totalToken := baseToken + int(bonusToken)
  
  // 5. Calculate REXTRA POIN
  rextraPoin := baseToken * duration.RextraPoinMultiplier
  
  return &PricingResult{
    BasePrice:              basePrice,
    PlanDiscount:           planDiscount,
    PromoDiscount:          0,
    UpgradeDiscount:        0,
    TotalDiscount:          planDiscount,
    FinalPrice:             finalPrice,
    EffectiveMonthlyPrice:  finalPrice / float64(duration.DurationMonths),
    PriceAfterPlanDiscount: priceAfterPlanDiscount,
    TotalToken:             totalToken,
    RextraPoin:             rextraPoin,
  }
}

// Calculate pricing with promo code
func (c *PriceCalculator) CalculatePricingWithPromo(plan *MembershipPlan, duration *MembershipDuration, promo *PromoCode) *PricingResult {
  // 1. Get base pricing
  result := c.CalculatePricing(plan, duration)
  
  // 2. Apply promo discount
  if promo != nil {
    promoDiscount := c.calculatePromoDiscount(promo, result.PriceAfterPlanDiscount)
    result.PromoDiscount = promoDiscount
    result.TotalDiscount = result.PlanDiscount + promoDiscount
    result.FinalPrice = result.PriceAfterPlanDiscount - promoDiscount
    
    // Ensure final price not negative
    if result.FinalPrice < 0 {
      result.FinalPrice = 0
    }
    
    result.EffectiveMonthlyPrice = result.FinalPrice / float64(duration.DurationMonths)
  }
  
  return result
}

// Calculate upgrade pricing (with 15% upgrade discount)
func (c *PriceCalculator) CalculateUpgradePricing(plan *MembershipPlan, duration *MembershipDuration) *PricingResult {
  // 1. Get base pricing
  result := c.CalculatePricing(plan, duration)
  
  // 2. Apply 15% upgrade discount
  upgradeDiscountPercentage := 15.0
  upgradeDiscount := result.PriceAfterPlanDiscount * (upgradeDiscountPercentage / 100)
  
  result.UpgradeDiscount = upgradeDiscount
  result.TotalDiscount = result.PlanDiscount + upgradeDiscount
  result.FinalPrice = result.PriceAfterPlanDiscount - upgradeDiscount
  result.EffectiveMonthlyPrice = result.FinalPrice / float64(duration.DurationMonths)
  
  return result
}

// Calculate downgrade pricing (NO upgrade discount)
func (c *PriceCalculator) CalculateDowngradePricing(plan *MembershipPlan, duration *MembershipDuration) *PricingResult {
  // Just return standard pricing (no upgrade discount)
  return c.CalculatePricing(plan, duration)
}

// Helper: Calculate promo discount
func (c *PriceCalculator) calculatePromoDiscount(promo *PromoCode, priceAfterPlanDiscount float64) float64 {
  switch promo.DiscountType {
  case "percentage":
    return priceAfterPlanDiscount * (promo.DiscountValue / 100)
  case "fixed_amount":
    discount := promo.DiscountValue
    if discount > priceAfterPlanDiscount {
      return priceAfterPlanDiscount
    }
    return discount
  case "free_100":
    return priceAfterPlanDiscount
  default:
    return 0
  }
}

// Helper: Calculate savings (for display)
func (c *PriceCalculator) CalculateSavings(plan *MembershipPlan, duration *MembershipDuration) float64 {
  // Savings = (base price without discount) - (final price)
  basePrice := plan.BaseMonthlyPrice * float64(duration.DurationMonths)
  result := c.CalculatePricing(plan, duration)
  return basePrice - result.FinalPrice
}
```

**Usage Example:**

```go
calculator := NewPriceCalculator()

// Standard pricing
pricing := calculator.CalculatePricing(proPlan, duration6Months)
fmt.Printf("Final Price: Rp%.0f\n", pricing.FinalPrice)
fmt.Printf("Total Token: %d\n", pricing.TotalToken)

// Pricing with promo
pricingWithPromo := calculator.CalculatePricingWithPromo(proPlan, duration6Months, promoCode)
fmt.Printf("Final Price with Promo: Rp%.0f\n", pricingWithPromo.FinalPrice)

// Upgrade pricing
upgradePricing := calculator.CalculateUpgradePricing(maxPlan, duration6Months)
fmt.Printf("Upgrade Price (15% off): Rp%.0f\n", upgradePricing.FinalPrice)
```

---

## PART 3 SELESAI ✅

**Part 3 ini sudah mencakup:**
- ✅ Alur Membership & 5 Status User detail
- ✅ Membership Activation Logic (purchase membership & token standalone)
- ✅ Xendit Webhook Handler lengkap dengan security, idempotency, error handling
- ✅ Token Refill Logic (monthly scheduler)
- ✅ Membership Expiration Logic
- ✅ Starter Plan Expiration Logic
- ✅ Expiration Reminders (H-7, H-3, H-1)
- ✅ Upgrade/Downgrade Calculation Logic dengan validation
- ✅ Promo Code Validation Logic lengkap
- ✅ Price Calculator Utility (central pricing calculator)

**Total Brief:**
- Part 1: Database Schema (8 tables)
- Part 2: API Endpoints (14 endpoints)
- Part 3: Business Logic & Auto-Transition (9 sections)

**Belum ada di brief:**
- Error Code Standards (bisa jadi Part 4 optional)
- Deployment Guide (bisa jadi Part 4 optional)
- Testing Strategy (bisa jadi Part 4 optional)

**Mau lanjut bikin Part 4 (optional) atau brief udah cukup lengkap?**