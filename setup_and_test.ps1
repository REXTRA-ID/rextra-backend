# REXTRA TRANSACTION TESTER (PRO)
$baseUrl = "http://localhost:8000/api/v1"

Clear-Host
Write-Host "==========================================" -ForegroundColor Yellow
Write-Host "   REXTRA TRANSACTION TESTER (PRO)        " -ForegroundColor Yellow
Write-Host "==========================================" -ForegroundColor Yellow

# 1. Login
Write-Host "`nAuthenticating..." -ForegroundColor Cyan
$loginBody = @{ email = "admin@email.com"; password = "admin123" } | ConvertTo-Json
try {
    $loginResp = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method Post -Body $loginBody -ContentType "application/json"
    $token = $loginResp.data.access_token
    Write-Host "Login Berhasil!" -ForegroundColor Green
} catch {
    Write-Host "Gagal Login: $($_.Exception.Message)" -ForegroundColor Red
    exit
}

while ($true) {
    Write-Host "`n=== MENU UTAMA REXTRA ===" -ForegroundColor Yellow
    Write-Host "1. Buat Transaksi Baru"
    Write-Host "2. Lihat Daftar Transaksi Saya"
    Write-Host "3. Keluar"
    
    $mainOpt = Read-Host "Pilih menu (1-3)"

    if ($mainOpt -eq "2") {
        Write-Host "`n[DAFTAR TRANSAKSI ANDA]" -ForegroundColor Cyan
        try {
            $list = Invoke-RestMethod -Uri "$baseUrl/my/payment" -Method Get -Headers @{ Authorization = "Bearer $token" }
            $list.data.data | ForEach-Object {
                $statusColor = "White"
                if ($_.payment_status -eq "PAID") { $statusColor = "Green" }
                elseif ($_.payment_status -eq "pending") { $statusColor = "Yellow" }
                elseif ($_.payment_status -eq "cancelled") { $statusColor = "Red" }
                
                Write-Host "[$($_.transaction_id)] " -NoNewline
                Write-Host "$($_.payment_status)" -ForegroundColor $statusColor -NoNewline
                Write-Host " | $($_.to_plan) ($($_.to_duration_months) Bln) | Total: Rp $($_.total_amount)"
                
                if ($_.payment_status -eq "pending") {
                    Write-Host "   >> URL Bayar: $($_.payment_url)" -ForegroundColor Gray
                    if ($_.pay_code) { Write-Host "   >> Kode Bayar: $($_.pay_code)" -ForegroundColor Gray }
                }
            }
        } catch {
            Write-Host "Gagal mengambil daftar: $($_.Exception.Message)" -ForegroundColor Red
        }
        continue
    } elseif ($mainOpt -eq "3") {
        break
    }

    # PROSES BUAT TRANSAKSI (Menu 1)
    # 2. Ambil Plan
    $plans = (Invoke-RestMethod -Uri "$baseUrl/membership/plan" -Method Get).data
    Write-Host "`n[PILIH PAKET]" -ForegroundColor Cyan
    for ($i=0; $i -lt $plans.Count; $i++) {
        Write-Host "$($i+1). $($plans[$i].plan_name) ($($plans[$i].tier_label))"
    }
    $planIdx = [int](Read-Host "Pilih nomor paket") - 1
    if ($planIdx -lt 0) { $planIdx = 2 }
    $selectedPlan = $plans[$planIdx]

    $prepare = Invoke-RestMethod -Uri "$baseUrl/checkout/prepare?plan_id=$($selectedPlan.id)&change_type=PEMBELIAN_BARU" -Method Get -Headers @{ Authorization = "Bearer $token" }

    # Pilih Durasi
    Write-Host "`n--- PILIH DURASI ---"
    $durations = $prepare.data.plan.durations
    for ($i=0; $i -lt $durations.Count; $i++) {
        Write-Host "$($i+1). $($durations[$i].duration_months) Bulan - Rp $($durations[$i].final_price)"
    }
    $durIdx = [int](Read-Host "Pilih nomor durasi") - 1
    if ($durIdx -lt 0) { $durIdx = 0 }
    $selectedDur = $durations[$durIdx]

    # Pilih Bundle
    Write-Host "`n--- ADD-ON TOKEN ---"
    $bundles = $prepare.data.token_bundles
    Write-Host "0. Tanpa Tambahan Token"
    for ($i=0; $i -lt $bundles.Count; $i++) {
        Write-Host "$($i+1). $($bundles[$i].name) ($($bundles[$i].token_amount) Token) - Rp $($bundles[$i].price_rp)"
    }
    $bunIdx = [int](Read-Host "Pilih nomor bundle")
    $selectedBunId = $null
    if ($bunIdx -gt 0) { $selectedBunId = $bundles[$bunIdx-1].id }

    $promo = Read-Host "`nKode Promo"

    # 3. Kalkulasi
    $calcBody = @{ plan_id = $selectedPlan.id; duration_id = $selectedDur.id; change_type = "PEMBELIAN_BARU"; token_bundle_package_id = $selectedBunId; promo_code = $promo } | ConvertTo-Json
    $calc = Invoke-RestMethod -Uri "$baseUrl/checkout/calculate" -Method Post -Body $calcBody -ContentType "application/json" -Headers @{ Authorization = "Bearer $token" }
    $total = $calc.data.total_amount

    Write-Host "`nTOTAL HARGA: Rp $total" -ForegroundColor Green

    # 4. Pilih Kanal
    Write-Host "`n--- PILIH KANAL BAYAR ---"
    $channels = (Invoke-RestMethod -Uri "$baseUrl/payment/channels?amount=$total" -Method Get).data
    for ($i=0; $i -lt $channels.Count; $i++) {
        Write-Host "$($i+1). $($channels[$i].name) (Rp $($channels[$i].fee_flat) + $($channels[$i].fee_percent)%)"
    }
    $chanIdx = [int](Read-Host "Pilih nomor kanal") - 1
    $selectedChan = $channels[$chanIdx]

    # Initiate
    $initBody = @{ plan_id = $selectedPlan.id; duration_id = $selectedDur.id; change_type = "PEMBELIAN_BARU"; token_bundle_package_id = $selectedBunId; promo_code = $promo; payment_method = $selectedChan.code } | ConvertTo-Json
    try {
        $init = Invoke-RestMethod -Uri "$baseUrl/checkout/initiate" -Method Post -Body $initBody -ContentType "application/json" -Headers @{ Authorization = "Bearer $token" }
        $trxId = $init.data.transaction_id
        Write-Host "`n--- TRANSAKSI BERHASIL DIBUAT: $trxId ---" -ForegroundColor Green
        
        while ($true) {
            Write-Host "`n[TINDAKAN]" -ForegroundColor Yellow
            Write-Host "1. Buka Browser | 2. Batalkan | 3. Simulasikan LUNAS | 4. Kembali ke Menu"
            $opt = Read-Host "Pilih"
            if ($opt -eq "1") { Start-Process $init.data.payment_url }
            elseif ($opt -eq "2") {
                Invoke-RestMethod -Uri "$baseUrl/admin/payment/$trxId/cancel" -Method Put -Body (@{ note="Batal" }|ConvertTo-Json) -ContentType "application/json" -Headers @{ Authorization = "Bearer $token" }
                Write-Host "DIBATALKAN!" -ForegroundColor Red; break
            } elseif ($opt -eq "3") {
                Invoke-RestMethod -Uri "$baseUrl/payment/simulate/$trxId" -Method Post
                Write-Host "LUNAS (SIMULASI)!" -ForegroundColor Green; break
            } else { break }
        }
    } catch {
        Write-Host "Gagal: $($_.Exception.Message)" -ForegroundColor Red
    }
}
