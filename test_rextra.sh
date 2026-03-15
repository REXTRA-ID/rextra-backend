#!/bin/bash

# REXTRA TRANSACTION TESTER (UBUNTU VERSION)
BASE_URL="http://localhost:8001/api/v1"

# Colors
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
GREEN='\033[0;32m'
RED='\033[0;31m'
GRAY='\033[0;90m'
NC='\033[0m' # No Color

clear
echo -e "${YELLOW}==========================================${NC}"
echo -e "${YELLOW}   REXTRA TRANSACTION TESTER (UBUNTU)     ${NC}"
echo -e "${YELLOW}==========================================${NC}"

# 1. Login
echo -e "\n${CYAN}Authenticating...${NC}"
LOGIN_DATA='{"email": "admin@email.com", "password": "admin123"}'
LOGIN_RESP=$(curl -s -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d "$LOGIN_DATA")

TOKEN=$(echo "$LOGIN_RESP" | jq -r '.data.access_token')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
    echo -e "${RED}Gagal Login: $(echo "$LOGIN_RESP" | jq -r '.message')${NC}"
    exit 1
fi

echo -e "${GREEN}Login Berhasil!${NC}"

while true; do
    echo -e "\n${YELLOW}=== MENU UTAMA REXTRA ===${NC}"
    echo "1. Buat Transaksi Baru"
    echo "2. Lihat Daftar Transaksi Saya"
    echo "3. Keluar"
    read -p "Pilih menu (1-3): " mainOpt

    if [ "$mainOpt" == "2" ]; then
        echo -e "\n${CYAN}[DAFTAR TRANSAKSI ANDA]${NC}"
        LIST_RESP=$(curl -s -X GET "$BASE_URL/my/payment" \
            -H "Authorization: Bearer $TOKEN")
        
        echo "$LIST_RESP" | jq -c '.data.data[]' | while read -r trx; do
            ID=$(echo "$trx" | jq -r '.transaction_id')
            STATUS=$(echo "$trx" | jq -r '.payment_status')
            PLAN=$(echo "$trx" | jq -r '.to_plan')
            DUR=$(echo "$trx" | jq -r '.to_duration_months')
            TOTAL=$(echo "$trx" | jq -r '.total_amount')
            URL=$(echo "$trx" | jq -r '.payment_url')
            CODE=$(echo "$trx" | jq -r '.pay_code')

            COLOR=$NC
            if [ "$STATUS" == "PAID" ]; then COLOR=$GREEN; fi
            if [ "$STATUS" == "pending" ]; then COLOR=$YELLOW; fi
            if [ "$STATUS" == "cancelled" ]; then COLOR=$RED; fi

            echo -e "[${ID}] ${COLOR}${STATUS}${NC} | ${PLAN} (${DUR} Bln) | Total: Rp ${TOTAL}"
            if [ "$STATUS" == "pending" ]; then
                echo -e "   ${GRAY}>> URL Bayar: ${URL}${NC}"
                if [ "$CODE" != "null" ]; then echo -e "   ${GRAY}>> Kode Bayar: ${CODE}${NC}"; fi
            fi
        done
        continue
    elif [ "$mainOpt" == "3" ]; then
        break
    fi

    # 2. Pilih Paket
    PLANS=$(curl -s -X GET "$BASE_URL/membership/plan" | jq -c '.data[]')
    echo -e "\n${CYAN}[PILIH PAKET]${NC}"
    i=1
    declare -a PLAN_IDS
    while read -r plan; do
        NAME=$(echo "$plan" | jq -r '.plan_name')
        TIER=$(echo "$plan" | jq -r '.tier_label')
        PLAN_IDS[$i]=$(echo "$plan" | jq -r '.id')
        echo "$i. $NAME ($TIER)"
        ((i++))
    done <<< "$PLANS"
    
    read -p "Pilih nomor paket: " planIdx
    SELECTED_PLAN_ID=${PLAN_IDS[$planIdx]}

    # Prepare Data
    PREPARE=$(curl -s -G "$BASE_URL/checkout/prepare" \
        -H "Authorization: Bearer $TOKEN" \
        --data-urlencode "plan_id=$SELECTED_PLAN_ID" \
        --data-urlencode "change_type=PEMBELIAN_BARU")

    # Pilih Durasi
    echo -e "\n--- PILIH DURASI ---"
    DURATIONS=$(echo "$PREPARE" | jq -c '.data.plan.durations[]')
    i=1
    declare -a DUR_IDS
    while read -r dur; do
        MONTHS=$(echo "$dur" | jq -r '.duration_months')
        PRICE=$(echo "$dur" | jq -r '.final_price')
        DUR_IDS[$i]=$(echo "$dur" | jq -r '.id')
        echo "$i. $MONTHS Bulan - Rp $PRICE"
        ((i++))
    done <<< "$DURATIONS"
    
    read -p "Pilih nomor durasi: " durIdx
    SELECTED_DUR_ID=${DUR_IDS[$durIdx]}

    # Pilih Bundle
    echo -e "\n--- ADD-ON TOKEN ---"
    BUNDLES=$(echo "$PREPARE" | jq -c '.data.token_bundles[]')
    echo "0. Tanpa Tambahan Token"
    i=1
    declare -a BUNDLE_IDS
    while read -r bun; do
        NAME=$(echo "$bun" | jq -r '.name')
        AMT=$(echo "$bun" | jq -r '.token_amount')
        PRICE=$(echo "$bun" | jq -r '.price_rp')
        BUNDLE_IDS[$i]=$(echo "$bun" | jq -r '.id')
        echo "$i. $NAME ($AMT Token) - Rp $PRICE"
        ((i++))
    done <<< "$BUNDLES"
    
    read -p "Pilih nomor bundle: " bunIdx
    SELECTED_BUNDLE_ID="null"
    if [ "$bunIdx" -gt 0 ]; then SELECTED_BUNDLE_ID="\"${BUNDLE_IDS[$bunIdx]}\""; fi

    read -p "Kode Promo (Kosongkan jika tidak ada): " promo

    # 3. Kalkulasi
    CALC_BODY="{\"plan_id\":\"$SELECTED_PLAN_ID\",\"duration_id\":\"$SELECTED_DUR_ID\",\"change_type\":\"PEMBELIAN_BARU\",\"token_bundle_package_id\":$SELECTED_BUNDLE_ID,\"promo_code\":\"$promo\"}"
    CALC=$(curl -s -X POST "$BASE_URL/checkout/calculate" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d "$CALC_BODY")
    
    TOTAL=$(echo "$CALC" | jq -r '.data.total_amount')
    echo -e "\n${GREEN}TOTAL HARGA: Rp $TOTAL${NC}"

    # 4. Pilih Kanal
    echo -e "\n--- PILIH KANAL BAYAR ---"
    CHANNELS=$(curl -s -X GET "$BASE_URL/payment/channels?amount=$TOTAL" | jq -c '.data[]')
    i=1
    declare -a CHAN_CODES
    while read -r chan; do
        NAME=$(echo "$chan" | jq -r '.name')
        CHAN_CODES[$i]=$(echo "$chan" | jq -r '.code')
        echo "$i. $NAME"
        ((i++))
    done <<< "$CHANNELS"
    
    read -p "Pilih nomor kanal: " chanIdx
    SELECTED_CHAN=${CHAN_CODES[$chanIdx]}

    # Initiate
    INIT_BODY="{\"plan_id\":\"$SELECTED_PLAN_ID\",\"duration_id\":\"$SELECTED_DUR_ID\",\"change_type\":\"PEMBELIAN_BARU\",\"token_bundle_package_id\":$SELECTED_BUNDLE_ID,\"promo_code\":\"$promo\",\"payment_method\":\"$SELECTED_CHAN\"}"
    INIT=$(curl -s -X POST "$BASE_URL/checkout/initiate" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d "$INIT_BODY")
    
    TRX_ID=$(echo "$INIT" | jq -r '.data.transaction_id')
    PAY_URL=$(echo "$INIT" | jq -r '.data.payment_url')

    if [ "$TRX_ID" != "null" ]; then
        echo -e "\n${GREEN}--- TRANSAKSI BERHASIL DIBUAT: $TRX_ID ---${NC}"
        while true; do
            echo -e "\n${YELLOW}[TINDAKAN]${NC}"
            echo "1. Tampilkan URL Bayar | 2. Batalkan | 3. Simulasikan LUNAS | 4. Kembali ke Menu"
            read -p "Pilih: " opt
            if [ "$opt" == "1" ]; then 
                echo -e "${CYAN}URL Bayar: $PAY_URL${NC}"
            elif [ "$opt" == "2" ]; then
                curl -s -X PUT "$BASE_URL/admin/payment/$TRX_ID/cancel" \
                    -H "Authorization: Bearer $TOKEN" \
                    -H "Content-Type: application/json" \
                    -d '{"note":"Batal dari script"}'
                echo -e "${RED}DIBATALKAN!${NC}"; break
            elif [ "$opt" == "3" ]; then
                curl -s -X POST "$BASE_URL/payment/simulate/$TRX_ID"
                echo -e "${GREEN}LUNAS (SIMULASI)!${NC}"; break
            else
                break
            fi
        done
    else
        echo -e "${RED}Gagal: $(echo "$INIT" | jq -r '.message')${NC}"
    fi
done
