# REXTRA TRANSACTION TESTER (PRO) - DYNAMIC VERSION
$baseUrl = "http://localhost:8000/api/v1"

function Show-Error {
    param($Msg, $Exception)
    Write-Host "`nERROR: $Msg" -ForegroundColor Red
    if ($Exception.Response) {
        try {
            $stream = $Exception.Response.GetResponseStream()
            $reader = New-Object System.IO.StreamReader($stream)
            $body = $reader.ReadToEnd() | ConvertFrom-Json
            Write-Host "Pesan Server: $($body.error)" -ForegroundColor Yellow
        } catch {}
    } else {
        Write-Host "Detail: $($Exception.Message)" -ForegroundColor Gray
    }
}

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
    Show-Error "Gagal Login" $_
    exit
}

while ($true) {
    # 1.5 Cek Membership Aktif
    Write-Host "`nMengecek status membership Anda..." -ForegroundColor Gray
    $myMem = $null
    try {
        $memResp = Invoke-RestMethod -Uri "$baseUrl/my/membership" -Method Get -Headers @{ Authorization = "Bearer $token" }
        $myMem = $memResp.data
        Write-Host "Status: $($myMem.plan_name) ($($myMem.remaining_days) Hari Tersisa)" -ForegroundColor Cyan
    } catch {
        Write-Host "Status: Belum Memiliki Membership Aktif" -ForegroundColor Gray
    }

    Write-Host "`n=== MENU UTAMA REXTRA ===" -ForegroundColor Yellow
    Write-Host "1. Buat Transaksi Baru / Perpanjang / Upgrade"
    Write-Host "2. Lihat Daftar Transaksi Saya"
    Write-Host "3. Keluar"
    
    $mainOpt = Read-Host "Pilih menu (1-3)"

    if ($mainOpt -eq "2") {
        Write-Host "`n[DAFTAR TRANSAKSI ANDA]" -ForegroundColor Cyan
        try {
            $list = Invoke-RestMethod -Uri "$baseUrl/checkout/transactions" -Method Get -Headers @{ Authorization = "Bearer $token" }
            if ($list.data.Count -eq 0) { Write-Host "(Belum ada riwayat transaksi)" -ForegroundColor Gray }
            $list.data | ForEach-Object {
                $statusColor = "White"
                if ($_.payment_status -eq "paid") { $statusColor = "Green" }
                elseif ($_.payment_status -eq "pending") { $statusColor = "Yellow" }
                elseif ($_.payment_status -eq "cancelled") { $statusColor = "Red" }
                
                Write-Host "[$($_.transaction_id)] " -NoNewline
                Write-Host "$($_.payment_status)" -ForegroundColor $statusColor -NoNewline
                Write-Host " | $($_.change_type) | Total: Rp $($_.total_amount)"
                if ($_.items) { $itemNames = $_.items | ForEach-Object { $_.name }; Write-Host "   Items: $($itemNames -join ', ')" -ForegroundColor Gray }
            }
        } catch {
            Show-Error "Gagal mengambil daftar" $_
        }
        continue
    } elseif ($mainOpt -eq "3") {
        break
    }

    # PROSES BUAT TRANSAKSI (Menu 1)
    # 2. Ambil Plan
    $plans = (Invoke-RestMethod -Uri "$baseUrl/membership/plan" -Method Get).data
    Write-Host "`n[PILIH PAKET TUJUAN]" -ForegroundColor Cyan
    for ($i=0; $i -lt $plans.Count; $i++) {
        Write-Host "$($i+1). $($plans[$i].plan_name) ($($plans[$i].tier_label))"
    }
    $planChoice = Read-Host "Pilih nomor paket"
    $planIdx = [int]$planChoice - 1
    if ($planIdx -lt 0 -or $planIdx -ge $plans.Count) { $planIdx = 0 }
    $selectedPlan = $plans[$planIdx]

    # Tentukan Change Type
    $changeType = "PEMBELIAN_BARU"
    if ($myMem -and $myMem.plan_name -ne "Starter") {
        Write-Host "`nAnda memiliki paket aktif. Pilih jenis transaksi:" -ForegroundColor Yellow
        Write-Host "1. RENEWAL (Perpanjang Paket yang Sama)"
        Write-Host "2. UPGRADE (Pindah ke Paket Lebih Tinggi)"
        Write-Host "3. DOWNGRADE (Pindah ke Paket Lebih Rendah)"
        Write-Host "4. PEMBELIAN_BARU (Timpa Total)"
        $ctOpt = Read-Host "Pilih (1-4, Default: 1)"
        if ($ctOpt -eq "2") { $changeType = "UPGRADE" }
        elseif ($ctOpt -eq "3") { $changeType = "DOWNGRADE" }
        elseif ($ctOpt -eq "4") { $changeType = "PEMBELIAN_BARU" }
        else { $changeType = "RENEWAL" }
    }

    try {
        $prepare = Invoke-RestMethod -Uri "$baseUrl/checkout/prepare?plan_id=$($selectedPlan.id)&change_type=$changeType" -Method Get -Headers @{ Authorization = "Bearer $token" }
    } catch {
        Show-Error "Gagal menyiapkan checkout" $_
        continue
    }

    # Pilih Durasi
    Write-Host "`n--- PILIH DURASI ---"
    $durations = $prepare.data.plan.durations
    for ($i=0; $i -lt $durations.Count; $i++) {
        Write-Host "$($i+1). $($durations[$i].duration_months) Bulan - Rp $($durations[$i].final_price)"
    }
    $durChoice = Read-Host "Pilih nomor durasi"
    $durIdx = [int]$durChoice - 1
    if ($durIdx -lt 0 -or $durIdx -ge $durations.Count) { $durIdx = 0 }
    $selectedDur = $durations[$durIdx]

    # Pilih Bundle
    Write-Host "`n--- ADD-ON TOKEN ---"
    $bundles = $prepare.data.token_bundles
    Write-Host "0. Tanpa Tambahan Token"
    for ($i=0; $i -lt $bundles.Count; $i++) {
        Write-Host "$($i+1). $($bundles[$i].name) ($($bundles[$i].token_amount) Token) - Rp $($bundles[$i].price_rp)"
    }
    $bunChoice = Read-Host "Pilih nomor bundle"
    $bunIdx = [int]$bunChoice
    $selectedBunId = $null
    if ($bunIdx -gt 0 -and $bunIdx -le $bundles.Count) { $selectedBunId = $bundles[$bunIdx-1].id }

    # PROMO SECTION
    $promo = ""
    while ($true) {
        Write-Host "`n--- PROMO & DISKON ---" -ForegroundColor Cyan
        $pubPromos = $prepare.data.eligible_discounts
        if ($pubPromos.Count -gt 0) {
            Write-Host "Diskon Publik Tersedia:"
            $pubPromos | ForEach-Object { Write-Host "- [$($_.code)] $($_.name): $($_.description)" -ForegroundColor Green }
        }
        
        $promoInput = Read-Host "Masukkan Kode Promo (Kosongkan jika tidak ada)"
        if ([string]::IsNullOrWhiteSpace($promoInput)) { $promo = $null; break }
        
        try {
            $valBody = @{ code = $promoInput; plan_id = $selectedPlan.id; duration_id = $selectedDur.id } | ConvertTo-Json
            $valResp = Invoke-RestMethod -Uri "$baseUrl/promo/validate" -Method Post -Body $valBody -ContentType "application/json" -Headers @{ Authorization = "Bearer $token" }
            Write-Host "PROMO BERHASIL: $($valResp.data.message)" -ForegroundColor Green
            $promo = $promoInput
            break
        } catch {
            Show-Error "Promo Tidak Valid" $_
            $retry = Read-Host "Coba kode lain? (y/n)"
            if ($retry -ne "y") { $promo = $null; break }
        }
    }

    # 3. Kalkulasi
    $calcBody = @{ 
        plan_id = $selectedPlan.id; duration_id = $selectedDur.id; change_type = $changeType; 
        token_bundle_package_id = $selectedBunId; promo_code = $promo; use_credit = $true
    } | ConvertTo-Json
    
    try {
        $calc = Invoke-RestMethod -Uri "$baseUrl/checkout/calculate" -Method Post -Body $calcBody -ContentType "application/json" -Headers @{ Authorization = "Bearer $token" }
        $total = $calc.data.total_amount
        Write-Host "`n--- RINGKASAN PEMBAYARAN ---" -ForegroundColor Yellow
        Write-Host "Harga Paket: Rp $($calc.data.membership_price)"
        if ($calc.data.token_price -gt 0) { Write-Host "Add-on Token: Rp $($calc.data.token_price)" }
        if ($calc.data.discount_amount -gt 0) { Write-Host "Diskon: -Rp $($calc.data.discount_amount)" -ForegroundColor Green }
        if ($calc.data.duration_credit -gt 0) { Write-Host "Kredit Sisa Paket Lama: -Rp $($calc.data.duration_credit)" -ForegroundColor Cyan }
        Write-Host "TOTAL AKHIR: Rp $total" -ForegroundColor Green
    } catch {
        Show-Error "Gagal kalkulasi harga" $_
        continue
    }

    # 4. Pilih Kanal
    Write-Host "`n--- PILIH KANAL BAYAR ---"
    try {
        $channels = (Invoke-RestMethod -Uri "$baseUrl/checkout/channels" -Method Get -Headers @{ Authorization = "Bearer $token" }).data
        for ($i=0; $i -lt $channels.Count; $i++) { Write-Host "$($i+1). $($channels[$i].name) (Admin: Rp $($channels[$i].admin_fee))" }
        $chanChoice = Read-Host "Pilih nomor kanal"
        $chanIdx = [int]$chanChoice - 1
        if ($chanIdx -lt 0 -or $chanIdx -ge $channels.Count) { $chanIdx = 0 }
        $selectedChan = $channels[$chanIdx]
    } catch {
        Show-Error "Gagal mengambil daftar channel" $_
        continue
    }

    # Initiate
    $initBody = @{ 
        plan_id = $selectedPlan.id; duration_id = $selectedDur.id; change_type = $changeType; 
        token_bundle_package_id = $selectedBunId; promo_code = $promo; 
        payment_method = $selectedChan.code; use_credit = $true
    } | ConvertTo-Json
    
    try {
        $init = Invoke-RestMethod -Uri "$baseUrl/checkout/initiate" -Method Post -Body $initBody -ContentType "application/json" -Headers @{ Authorization = "Bearer $token" }
        $trxId = $init.data.transaction_id
        Write-Host "`n--- TRANSAKSI BERHASIL DIBUAT: $trxId ---" -ForegroundColor Green
        
        while ($true) {
            Write-Host "`n[TINDAKAN]" -ForegroundColor Yellow
            Write-Host "1. Buka Browser Pembayaran | 2. Batalkan | 3. Simulasikan LUNAS | 4. Kembali"
            $opt = Read-Host "Pilih"
            if ($opt -eq "1") { Start-Process $init.data.payment_url }
            elseif ($opt -eq "2") {
                $cancelBody = @{ cancel_reason = "CHANGE_ORDER"; cancel_note = "Dibatalkan via tester" } | ConvertTo-Json
                Invoke-RestMethod -Uri "$baseUrl/checkout/cancel/$trxId" -Method Post -Body $cancelBody -ContentType "application/json" -Headers @{ Authorization = "Bearer $token" }
                Write-Host "TRANSAKSI $trxId TELAH DIBATALKAN!" -ForegroundColor Red; break
            } elseif ($opt -eq "3") {
                Invoke-RestMethod -Uri "$baseUrl/payment/simulate/$trxId" -Method Post
                Write-Host "PEMBAYARAN BERHASIL DISIMULASIKAN!" -ForegroundColor Green; break
            } else { break }
        }
    } catch {
        Show-Error "Gagal membuat transaksi" $_
    }
}
