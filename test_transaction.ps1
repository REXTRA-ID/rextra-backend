# REXTRA API Interactive Test Script
$baseUrl = "http://localhost:8000/api/v1"

Clear-Host
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "   REXTRA MEMBERSHIP API TESTER (FIXED)   " -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan

# 1. Login Step (Hardcoded to avoid shell mangling @ symbol)
Write-Host "`n[STEP 1] Autentikasi User" -ForegroundColor Yellow
$userPart = "admin"
$domainPart = "email.com"
$email = $userPart + "@" + $domainPart
$password = "admin123"

Write-Host "Mencoba login otomatis sebagai: $email" -ForegroundColor Gray
try {
    $loginBody = @{ email = $email; password = $password } | ConvertTo-Json
    $loginResp = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method Post -Body $loginBody -ContentType "application/json"
    $token = $loginResp.data.access_token
    Write-Host "Login Berhasil!" -ForegroundColor Green
} catch {
    Write-Host "ERROR: Login Gagal." -ForegroundColor Red
    Write-Host "Penyebab: Server mungkin belum siap atau database belum terisi seeder."
    exit
}

# 2. Pick Plan
Write-Host "`n[STEP 2] Pilih Paket Membership" -ForegroundColor Yellow
$plansResp = Invoke-RestMethod -Uri "$baseUrl/membership/plan" -Method Get
$plans = $plansResp.data
for ($i=0; $i -lt $plans.Count; $i++) {
    Write-Host "$($i+1). $($plans[$i].plan_name) ($($plans[$i].tier_label))"
}
$planChoice = Read-Host "Pilih nomor paket (Default: 3)"
if ([string]::IsNullOrWhiteSpace($planChoice)) { $planChoice = 3 }
$selectedPlan = $plans[[int]$planChoice - 1]

# 3. Prepare & Pick Duration
Write-Host "`n[STEP 3] Pilih Durasi untuk $($selectedPlan.plan_name)" -ForegroundColor Yellow
$prepareUrl = "$baseUrl/checkout/prepare?plan_id=$($selectedPlan.id)&change_type=PEMBELIAN_BARU"
$prepareResp = Invoke-RestMethod -Uri $prepareUrl -Method Get -Headers @{ Authorization = "Bearer $token" }
$durations = $prepareResp.data.plan.durations

if ($null -eq $durations -or $durations.Count -eq 0) {
    Write-Host "Error: Durasi paket belum diset di database." -ForegroundColor Red
    exit
}

for ($i=0; $i -lt $durations.Count; $i++) {
    Write-Host "$($i+1). $($durations[$i].duration_months) Bulan - Rp $($durations[$i].final_price)"
}
$durChoice = Read-Host "Pilih nomor durasi (Default: 1)"
if ([string]::IsNullOrWhiteSpace($durChoice)) { $durChoice = 1 }
$selectedDuration = $durations[[int]$durChoice - 1]

# 4. Calculate Summary
Write-Host "`n[STEP 4] Preview Pembayaran" -ForegroundColor Yellow
$calcBody = @{
    plan_id = $selectedPlan.id
    duration_id = $selectedDuration.id
    change_type = "PEMBELIAN_BARU"
    use_credit = $true
} | ConvertTo-Json

$calcResp = Invoke-RestMethod -Uri "$baseUrl/checkout/calculate" -Method Post -Body $calcBody -ContentType "application/json" -Headers @{ Authorization = "Bearer $token" }

Write-Host "------------------------------------------"
Write-Host "Plan: $($selectedPlan.plan_name)"
Write-Host "Durasi: $($selectedDuration.duration_months) Bulan"
Write-Host "Harga: Rp $($calcResp.data.membership_price)"
Write-Host "TOTAL: Rp $($calcResp.data.total_amount)"
Write-Host "------------------------------------------"

# 5. Initiate Payment
$confirm = Read-Host "Buat Transaksi & Buka Browser Pembayaran? (y/n)"
if ($confirm -ne "y") { exit }

$initBody = @{
    plan_id = $selectedPlan.id
    duration_id = $selectedDuration.id
    change_type = "PEMBELIAN_BARU"
    payment_method = "BRIVA" 
} | ConvertTo-Json

try {
    $initResp = Invoke-RestMethod -Uri "$baseUrl/checkout/initiate" -Method Post -Body $initBody -ContentType "application/json" -Headers @{ Authorization = "Bearer $token" }
    Write-Host "`nSUKSES! Transaksi Created." -ForegroundColor Green
    Write-Host "URL: $($initResp.data.payment_url)" -ForegroundColor Cyan
    Start-Process $initResp.data.payment_url
} catch {
    Write-Host "Gagal membuat transaksi." -ForegroundColor Red
}
