# REXTRA TRANSACTION TESTER (PRO) - MOBILE SIMULATION VERSION (STABLE V4)
$baseUrl = "http://localhost:8000/api/v1"
$dbHost = "103.171.84.248"
$dbUser = "postgres"
$dbPass = "password"
$dbName = "rextra"
$dbPort = "5433"

# Helper to run DB queries - ULTRA STABLE VERSION 4
function Invoke-DBQuery($Query) {
    # Write query to temporary file
    $Query | Out-File -FilePath "query.sql" -Encoding utf8 -Force
    
    # Write Node.js loader
    $jsLoader = @"
const { Client } = require('pg');
const fs = require('fs');
const client = new Client({ host: '$dbHost', user: '$dbUser', password: '$dbPass', database: '$dbName', port: $dbPort });
async function run() {
    try {
        const q = fs.readFileSync('query.sql', 'utf8').trim();
        await client.connect();
        const res = await client.query(q);
        process.stdout.write(JSON.stringify(res.rows));
    } catch (e) {
        process.stderr.write('SQL_ERROR: ' + e.message);
    } finally {
        await client.end();
    }
}
run();
"@
    $jsLoader | Out-File -FilePath "db_exec.js" -Encoding utf8 -Force
    
    # Execute and capture
    $out = node db_exec.js 2>$null
    
    # Cleanup
    Remove-Item "db_exec.js" -ErrorAction SilentlyContinue
    Remove-Item "query.sql" -ErrorAction SilentlyContinue
    
    if ([string]::IsNullOrWhiteSpace($out) -or $out.Contains("SQL_ERROR")) {
        return @()
    }
    
    try { return ($out | ConvertFrom-Json) } catch { return @() }
}

function Show-Dashboard($data) {
    if (!$data) { return }
    $m = $data.membership; $u = $data.ui_state; $w = $data.wallet
    
    Clear-Host
    Write-Host "--- SIMULASI LAYAR MOBILE REXTRA ---" -ForegroundColor Gray
    Write-Host "=========================================="

    $bgColor = "Blue"
    if ($m.plan_name -eq "Pro") { $bgColor = "Magenta" }
    elseif ($m.plan_name -eq "Max") { $bgColor = "Yellow" }
    elseif ($m.plan_name -eq "Basic") { $bgColor = "Green" }
    
    Write-Host " [ SECTION A - MEMBERSHIP CARD ] " -BackgroundColor $bgColor -ForegroundColor Black
    Write-Host "  STATUS: $($m.status_label) " -ForegroundColor Gray
    Write-Host "  PAKET:  $($m.plan_name.ToUpper()) " -ForegroundColor White
    Write-Host "  Aktif hingga: $($m.active_until)"
    Write-Host "  Sisa Waktu:   $($m.remaining_days) Hari" -ForegroundColor Yellow
    Write-Host " "
    Write-Host "  TOKEN TERSEDIA: $($w.token_balance) Token" -ForegroundColor Cyan
    Write-Host " ------------------------------------------ " -ForegroundColor $bgColor

    if ($u.is_near_expiry) {
        Write-Host "`n [ ALERT: PENGINGAT ] " -BackgroundColor Yellow -ForegroundColor Black
        Write-Host " Membership kamu akan segera berakhir dalam $($m.remaining_days) hari." -ForegroundColor Yellow
    }

    if ($u.is_expired) {
        Write-Host "`n [ ALERT: MEMBERSHIP BERAKHIR ] " -BackgroundColor Red -ForegroundColor White
        Write-Host " Masa aktif membership kamu sudah habis." -ForegroundColor Red
    }

    Write-Host "`n [ QUICK ACTIONS ]"
    Write-Host " 1. Top Up Token | 2. Perbarui Membership | 3. Riwayat Transaksi"
}

function Show-Catalog($data) {
    Clear-Host
    Write-Host "--- SIMULASI LAYAR PENAWARAN CLUB ---" -ForegroundColor Gray
    Write-Host "=========================================="
    if ($data.current_context) {
        Write-Host " [ INFO BANNER ] " -BackgroundColor Gray -ForegroundColor White
        Write-Host " $($data.current_context.message)" -ForegroundColor Gray
        Write-Host " ------------------------------------------ "
    }
    foreach ($p in $data.plans) {
        $c = "White"
        if ($p.plan_name -eq "Basic") { $c = "Green" } 
        elseif ($p.plan_name -eq "Pro") { $c = "Magenta" } 
        elseif ($p.plan_name -eq "Max") { $c = "Yellow" }
        
        Write-Host "`n [$($p.plan_name.ToUpper()) PLAN] " -ForegroundColor $c
        if ($p.label) { Write-Host " Label: $($p.label)" -ForegroundColor Yellow }
        Write-Host " $($p.pricing_info) | Bonus: $($p.token_bonus_info)"
        Write-Host " Intro: $($p.marketing_intro)" -ForegroundColor Gray
        Write-Host " BUTTON -> [ $($p.cta.label) ]" -BackgroundColor $c -ForegroundColor Black
        Write-Host " ------------------------------------------ "
    }
}

while ($true) {
    Clear-Host
    Write-Host "==========================================" -ForegroundColor Yellow
    Write-Host "   REXTRA MOBILE SIMULATOR (STABLE V4)    " -ForegroundColor Yellow
    Write-Host "==========================================" -ForegroundColor Yellow
    Write-Host "1. Akun Baru (Standard - Perlu Klaim Starter)"
    Write-Host "2. Akun Pro (Hampir Expired H-3)"
    Write-Host "3. Akun Basic (Sudah Expired)"
    Write-Host "4. Akun Max (Aktif Normal Sisa 60 Hari)"
    Write-Host "5. Keluar"

    $opt = Read-Host "`nPilih Skenario"
    if ($opt -eq "5") { break }

    $mail = "sim_mobile_$($opt)@email.com"
    $hash = '$2a$04$jP9U2bPN8tfpHfuqMni6bOFz6xokF2nrdpZOQRGPqpTceGPpfO0Wq' # admin123

    Write-Host "`n[1/3] Menyiapkan State Database..." -ForegroundColor Gray
    
    # Comprehensive Cleanup to handle all Foreign Keys
    $cleanupSQL = @"
DO \$\$
DECLARE
    uid UUID;
BEGIN
    SELECT id INTO uid FROM users WHERE email = '$mail';
    IF uid IS NOT NULL THEN
        -- Delete from deep dependencies first
        DELETE FROM ikigai_total_scores WHERE test_session_id IN (SELECT id FROM careerprofile_test_sessions WHERE user_id = uid);
        DELETE FROM career_recommendations WHERE test_session_id IN (SELECT id FROM careerprofile_test_sessions WHERE user_id = uid);
        DELETE FROM fit_check_results WHERE user_id = uid;
        DELETE FROM careerprofile_test_sessions WHERE user_id = uid;
        DELETE FROM user_career_profiles WHERE user_id = uid;
        DELETE FROM usage_logs WHERE membership_id IN (SELECT id FROM memberships WHERE user_id = uid);
        DELETE FROM subscription_cycles WHERE membership_id IN (SELECT id FROM memberships WHERE user_id = uid);
        DELETE FROM user_entitlement_quotas WHERE membership_id IN (SELECT id FROM memberships WHERE user_id = uid);
        DELETE FROM points_ledger WHERE membership_id IN (SELECT id FROM memberships WHERE user_id = uid);
        DELETE FROM token_ledger WHERE wallet_id IN (SELECT id FROM token_wallet WHERE user_id = uid);
        DELETE FROM token_wallet WHERE user_id = uid;
        DELETE FROM payment_transactions WHERE user_id = uid;
        DELETE FROM token_usage_histories WHERE user_id = uid;
        DELETE FROM poin_transactions WHERE user_id = uid;
        DELETE FROM discount_redemptions WHERE user_id = uid;
        DELETE FROM sessions WHERE user_id = uid;
        DELETE FROM memberships WHERE user_id = uid;
        DELETE FROM users WHERE id = uid;
    END IF;
END \$\$;
"@
    $null = Invoke-DBQuery $cleanupSQL

    # Create User with ON CONFLICT
    $userSQL = "INSERT INTO users (id, fullname, email, password, is_verified, phone_number, role) VALUES (gen_random_uuid(), 'Mobile Tester', '$mail', '$hash', true, '081222333444', 'USER') ON CONFLICT (email) DO UPDATE SET is_verified = true RETURNING id;"
    $uRes = Invoke-DBQuery $userSQL
    if ($uRes.Count -eq 0) { Write-Host "Gagal Insert/Get User!"; pause; continue }
    $uid = $uRes[0].id

    # Create Membership State
    if ($opt -eq "1") {
        # Standard active, no claim yet
        $p = (Invoke-DBQuery "SELECT id FROM membership_plans WHERE plan_name = 'Standard'")[0].id
        $d = (Invoke-DBQuery "SELECT id FROM plan_durations WHERE plan_id = '$p' AND duration_months = 1")[0].id
        $null = Invoke-DBQuery "INSERT INTO memberships (id, user_id, plan_id, plan_name, duration_id, duration_months, is_active, has_claimed_starter) VALUES (gen_random_uuid(), '$uid', '$p', 'Standard', '$d', 1, true, false) ON CONFLICT (user_id) DO UPDATE SET plan_name = 'Standard', is_active = true, has_claimed_starter = false;"
    }
    elseif ($opt -eq "2") {
        # Pro Near Expiry
        $p = (Invoke-DBQuery "SELECT id FROM membership_plans WHERE plan_name = 'Pro'")[0].id
        $d = (Invoke-DBQuery "SELECT id FROM plan_durations WHERE plan_id = '$p' AND duration_months = 1")[0].id
        $null = Invoke-DBQuery "INSERT INTO memberships (id, user_id, plan_id, plan_name, duration_id, duration_months, is_active, started_at, expired_at, has_claimed_starter) VALUES (gen_random_uuid(), '$uid', '$p', 'Pro', '$d', 1, true, NOW() - interval '27 days', NOW() + interval '3 days', true) ON CONFLICT (user_id) DO UPDATE SET plan_name = 'Pro', is_active = true, expired_at = NOW() + interval '3 days';"
    }
    elseif ($opt -eq "3") {
        # Basic Expired
        $p = (Invoke-DBQuery "SELECT id FROM membership_plans WHERE plan_name = 'Basic'")[0].id
        $d = (Invoke-DBQuery "SELECT id FROM plan_durations WHERE plan_id = '$p' AND duration_months = 1")[0].id
        $null = Invoke-DBQuery "INSERT INTO memberships (id, user_id, plan_id, plan_name, duration_id, duration_months, is_active, started_at, expired_at, has_claimed_starter) VALUES (gen_random_uuid(), '$uid', '$p', 'Basic', '$d', 1, false, NOW() - interval '40 days', NOW() - interval '1 day', true) ON CONFLICT (user_id) DO UPDATE SET plan_name = 'Basic', is_active = false, expired_at = NOW() - interval '1 day';"
    }
    elseif ($opt -eq "4") {
        # Max Active
        $p = (Invoke-DBQuery "SELECT id FROM membership_plans WHERE plan_name = 'Max'")[0].id
        $d = (Invoke-DBQuery "SELECT id FROM plan_durations WHERE plan_id = '$p' AND duration_months = 3")[0].id
        $null = Invoke-DBQuery "INSERT INTO memberships (id, user_id, plan_id, plan_name, duration_id, duration_months, is_active, started_at, expired_at, has_claimed_starter) VALUES (gen_random_uuid(), '$uid', '$p', 'Max', '$d', 3, true, NOW() - interval '30 days', NOW() + interval '60 days', true) ON CONFLICT (user_id) DO UPDATE SET plan_name = 'Max', is_active = true, expired_at = NOW() + interval '60 days';"
    }

    Write-Host "[2/3] Login..." -ForegroundColor Cyan
    try {
        $loginBody = @{ email = $mail; password = "admin123" } | ConvertTo-Json
        $token = (Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method Post -Body $loginBody -ContentType "application/json").data.access_token

        if ($opt -eq "1") {
            Write-Host "`nSelamat Datang di REXTRA!" -ForegroundColor Green
            $claim = Read-Host "Klaim Trial Starter 30 Hari? (y/n)"
            if ($claim -eq "y") {
                $null = Invoke-RestMethod -Uri "$baseUrl/my/claim-starter" -Method Post -Headers @{ Authorization = "Bearer $token" }
                Write-Host "Klaim Berhasil!" -ForegroundColor Green
            }
        }

        while ($true) {
            Write-Host "`n[3/3] Memuat Dashboard..." -ForegroundColor Cyan
            $dash = (Invoke-RestMethod -Uri "$baseUrl/my/membership/dashboard" -Method Get -Headers @{ Authorization = "Bearer $token" }).data
            Show-Dashboard $dash
            
            $act = Read-Host "`nPilih Tindakan (1-3, atau Enter untuk kembali ke Menu Utama)"
            if ($act -eq "2") {
                $catalog = (Invoke-RestMethod -Uri "$baseUrl/membership/plan" -Method Get -Headers @{ Authorization = "Bearer $token" }).data
                Show-Catalog $catalog
                Read-Host "`nTekan Enter untuk kembali ke Dashboard"
            } elseif ([string]::IsNullOrWhiteSpace($act)) {
                break
            }
        }
    } catch {
        Write-Host "`nERROR: $($_.Exception.Message)" -ForegroundColor Red
        pause
    }
}
