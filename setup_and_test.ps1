# REXTRA TRANSACTION TESTER (PRO) - MOBILE SIMULATION VERSION
$baseUrl = "http://localhost:8000/api/v1"
$dbHost = "103.171.84.248"
$dbUser = "postgres"
$dbPass = "password"
$dbName = "rextra"
$dbPort = "5433"

# Helper to run DB queries with robust logic
function Invoke-DBQuery($Query) {
    $jsCode = "const { Client } = require('pg'); const client = new Client({ host: '$dbHost', user: '$dbUser', password: '$dbPass', database: '$dbName', port: $dbPort }); async function run() { try { await client.connect(); const res = await client.query(\`$Query\`); process.stdout.write(JSON.stringify(res.rows)); } catch (e) { process.stderr.write(e.message); } finally { await client.end(); } } run();"
    $jsCode | Out-File -FilePath "temp_db.js" -Encoding utf8 -Force
    $out = node temp_db.js 2>$null
    Remove-Item "temp_db.js" -ErrorAction SilentlyContinue
    if ([string]::IsNullOrWhiteSpace($out)) { return @() }
    try { return ($out | ConvertFrom-Json) } catch { return @() }
}

function Show-MobileDashboard($DashboardData) {
    if (!$DashboardData) { Write-Host "Error: Data Dashboard Tidak Ditemukan" -ForegroundColor Red; return }
    $m = $DashboardData.membership
    $u = $DashboardData.ui_state
    $w = $DashboardData.wallet
    
    Clear-Host
    Write-Host "--- SIMULASI LAYAR MOBILE REXTRA ---" -ForegroundColor Gray
    Write-Host "=========================================="

    $bgColor = "Blue"
    if ($m.plan_name -eq "Pro") { $bgColor = "Magenta" }
    elseif ($m.plan_name -eq "Max") { $bgColor = "Yellow" }
    elseif ($m.plan_name -eq "Basic") { $bgColor = "Green" }
    
    Write-Host " [ SECTION A - MEMBERSHIP CARD ] " -BackgroundColor $bgColor -ForegroundColor Black
    Write-Host "  $($m.status_label) " -ForegroundColor Gray
    Write-Host "  $($m.plan_name.ToUpper()) " -ForegroundColor White
    Write-Host " "
    Write-Host "  Aktif hingga $($m.active_until)"
    Write-Host "  Sisa $($m.remaining_days) hari" -ForegroundColor Yellow
    Write-Host " "
    Write-Host "  Token Tersedia:"
    Write-Host "  $($w.token_balance) Token" -ForegroundColor Cyan
    Write-Host " ------------------------------------------ " -ForegroundColor $bgColor

    if ($u.is_near_expiry) {
        Write-Host "`n [ ALERT: PENGINGAT ] " -BackgroundColor Yellow -ForegroundColor Black
        Write-Host " Membership kamu akan segera berakhir" -ForegroundColor Yellow
        Write-Host " Langganan REXTRA CLUB kamu akan berakhir dalam $($m.remaining_days) hari."
        Write-Host " Perpanjang sekarang agar akses tetap aktif tanpa jeda."
    }

    if ($u.is_expired) {
        Write-Host "`n [ ALERT: MEMBERSHIP BERAKHIR ] " -BackgroundColor Red -ForegroundColor White
        Write-Host " Masa aktif membership kamu sudah habis." -ForegroundColor Red
        Write-Host " Segera perbarui untuk mendapatkan kembali akses fitur."
    }

    Write-Host "`n [ QUICK ACTIONS ]" -ForegroundColor Gray
    Write-Host " Top Up | Perpanjang | Riwayat"

    Write-Host "`n [ SECTION B - BENEFIT MEMBERSHIP ] " -ForegroundColor White
    Write-Host " ------------------------------------------ "
    if ($DashboardData.benefits) {
        foreach ($group in $DashboardData.benefits) {
            Write-Host " * $($group.feature_name): $($group.access_label)" -ForegroundColor Cyan
            if ($group.sub_features) {
                foreach ($sub in $group.sub_features) {
                    Write-Host "   - $($sub.name): $($sub.access_label)"
                }
            }
            Write-Host " ------------------------------------------ "
        }
    }
}

# --- MAIN FLOW ---
while ($true) {
    Clear-Host
    Write-Host "==========================================" -ForegroundColor Yellow
    Write-Host "   REXTRA SYSTEM TESTER (MOBILE FLOW)     " -ForegroundColor Yellow
    Write-Host "==========================================" -ForegroundColor Yellow

    Write-Host "Pilih Skenario Pengetesan:"
    Write-Host "1. Akun Baru (Manual Klaim Starter)"
    Write-Host "2. Akun Lama - Hampir Expired (H-3)"
    Write-Host "3. Akun Lama - Sudah Expired"
    Write-Host "4. Akun Lama - Masih Aktif (Sisa 60 Hari)"
    Write-Host "5. Keluar"

    $scenChoice = Read-Host "`nPilih nomor"
    if ($scenChoice -eq "5") { break }

    $testMail = "tester_mobile_$($scenChoice)@email.com"
    $testHash = '$2a$04$jP9U2bPN8tfpHfuqMni6bOFz6xokF2nrdpZOQRGPqpTceGPpfO0Wq' # admin123

    Write-Host "`n[1/3] Menyiapkan State Database..." -ForegroundColor Gray
    $null = Invoke-DBQuery "DELETE FROM subscription_cycles WHERE membership_id IN (SELECT id FROM memberships WHERE user_id IN (SELECT id FROM users WHERE email = '$testMail'))"
    $null = Invoke-DBQuery "DELETE FROM payment_transactions WHERE user_id IN (SELECT id FROM users WHERE email = '$testMail')"
    $null = Invoke-DBQuery "DELETE FROM memberships WHERE user_id IN (SELECT id FROM users WHERE email = '$testMail')"
    $null = Invoke-DBQuery "DELETE FROM users WHERE email = '$testMail'"

    $userRows = Invoke-DBQuery "INSERT INTO users (id, fullname, email, password, is_verified, phone_number, role) VALUES (gen_random_uuid(), 'Mobile Tester', '$testMail', '$testHash', true, '081222333444', 'USER') RETURNING id"
    if ($userRows.Count -eq 0) { Write-Host "Gagal menyiapkan data user!"; pause; continue }
    $testUid = $userRows[0].id

    if ($scenChoice -eq "1") {
        $null = Invoke-DBQuery "INSERT INTO memberships (id, user_id, plan_name, is_active, has_claimed_starter) VALUES (gen_random_uuid(), '$testUid', 'Starter', false, false)"
    }
    elseif ($scenChoice -eq "2") {
        $pRes = Invoke-DBQuery "SELECT id FROM membership_plans WHERE plan_name = 'Pro'"
        $pId = $pRes[0].id
        $dRes = Invoke-DBQuery "SELECT id FROM plan_durations WHERE plan_id = '$pId' AND duration_months = 1"
        $dId = $dRes[0].id
        $st = (Get-Date).AddDays(-27).ToString("yyyy-MM-dd HH:mm:ss")
        $en = (Get-Date).AddDays(3).ToString("yyyy-MM-dd HH:mm:ss")
        $null = Invoke-DBQuery "INSERT INTO memberships (id, user_id, plan_id, plan_name, duration_id, duration_months, is_active, started_at, expired_at, has_claimed_starter) VALUES (gen_random_uuid(), '$testUid', '$pId', 'Pro', '$dId', 1, true, '$st', '$en', true)"
    }
    elseif ($scenChoice -eq "3") {
        $pRes = Invoke-DBQuery "SELECT id FROM membership_plans WHERE plan_name = 'Basic'"
        $pId = $pRes[0].id
        $dRes = Invoke-DBQuery "SELECT id FROM plan_durations WHERE plan_id = '$pId' AND duration_months = 1"
        $dId = $dRes[0].id
        $st = (Get-Date).AddDays(-40).ToString("yyyy-MM-dd HH:mm:ss")
        $en = (Get-Date).AddDays(-1).ToString("yyyy-MM-dd HH:mm:ss")
        $null = Invoke-DBQuery "INSERT INTO memberships (id, user_id, plan_id, plan_name, duration_id, duration_months, is_active, started_at, expired_at, has_claimed_starter) VALUES (gen_random_uuid(), '$testUid', '$pId', 'Basic', '$dId', 1, false, '$st', '$en', true)"
    }
    elseif ($scenChoice -eq "4") {
        $pRes = Invoke-DBQuery "SELECT id FROM membership_plans WHERE plan_name = 'Max'"
        $pId = $pRes[0].id
        $dRes = Invoke-DBQuery "SELECT id FROM plan_durations WHERE plan_id = '$pId' AND duration_months = 3"
        $dId = $dRes[0].id
        $st = (Get-Date).AddDays(-30).ToString("yyyy-MM-dd HH:mm:ss")
        $en = (Get-Date).AddDays(60).ToString("yyyy-MM-dd HH:mm:ss")
        $null = Invoke-DBQuery "INSERT INTO memberships (id, user_id, plan_id, plan_name, duration_id, duration_months, is_active, started_at, expired_at, has_claimed_starter) VALUES (gen_random_uuid(), '$testUid', '$pId', 'Max', '$dId', 3, true, '$st', '$en', true)"
    }

    Write-Host "[2/3] Authenticating..." -ForegroundColor Cyan
    try {
        $loginBody = @{ email = $testMail; password = "admin123" } | ConvertTo-Json
        $loginResp = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method Post -Body $loginBody -ContentType "application/json"
        $authToken = $loginResp.data.access_token

        if ($scenChoice -eq "1") {
            Write-Host "`nSelamat Datang di REXTRA!" -ForegroundColor Green
            Write-Host "Kamu berhak mendapatkan akses STARTER PLAN GRATIS selama 30 hari."
            $claimPrompt = Read-Host "Klaim Sekarang? (y/n)"
            if ($claimPrompt -eq "y") {
                $null = Invoke-RestMethod -Uri "$baseUrl/my/claim-starter" -Method Post -Headers @{ Authorization = "Bearer $authToken" }
                Write-Host "Klaim Starter Sukses! Menuju Dashboard..." -ForegroundColor Green
                Start-Sleep -Seconds 1
            }
        }

        Write-Host "[3/3] Memuat Halaman Dashboard..." -ForegroundColor Cyan
        $dashboardData = Invoke-RestMethod -Uri "$baseUrl/my/membership/dashboard" -Method Get -Headers @{ Authorization = "Bearer $authToken" }
        Show-MobileDashboard $dashboardData.data
    } catch {
        Write-Host "Gagal eksekusi skenario: $($_.Exception.Message)" -ForegroundColor Red
    }
    
    Write-Host "`nSkenario Selesai. Tekan tombol apa saja untuk kembali ke Menu."
    $null = [Console]::ReadKey()
}
