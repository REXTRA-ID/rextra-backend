# REXTRA API MAPPING (UI/UX BRIEF COMPLIANCE)

Dokumen ini memetakan antara kebutuhan Halaman UI (berdasarkan file brief) dengan Endpoint API yang tersedia atau yang perlu dibuat.

## 📊 Mapping Screen & API

| No | Nama Screen / State | Endpoint Terkait | Keterangan Data & Kondisi | Status |
|:---|:---|:---|:---|:---|
| 1 | **Dashboard Membership Tracker** | `GET /api/v1/my/membership` | Mengambil data membership aktif, sisa hari, dan emblem warna. | ✅ Tested |
| 2 | **Dashboard - Quota/Token Tracker** | `GET /api/v1/entitlement/my-quota` | Mengambil sisa frekuensi fitur (unlimited/limited) dan saldo token. | ✅ Tested |
| 3 | **State: Membership Akan Berakhir** | `GET /api/v1/my/membership` | Muncul jika `expired_at` sisa ≤ 7 hari. Backend memberikan flag/tanggal. | ✅ Tested |
| 4 | **Halaman Penawaran R-CLUB (Landing)** | `GET /api/v1/membership/plan` | Menampilkan list paket (Basic, Pro, Max) beserta benefitnya. | ✅ Tested |
| 5 | **Penawaran - State User Berlangganan** | `GET /api/v1/membership/plan` | Backend perlu memberi info mana "Plan Saat Ini" dan tombol (Upgrade/Downgrade). | ⚠️ Partial (Frontend Logic) |
| 6 | **Checkout - (New/Upgrade/Downgrade)** | `GET /api/v1/checkout/prepare` | Menampilkan detail paket terpilih, list durasi, dan **list bundle token**. | ✅ Tested |
| 7 | **Modal Pilih Durasi** | `GET /api/v1/checkout/prepare` | Data durasi (1, 3, 6, 12 bln) beserta harga coret & harga final. | ✅ Tested |
| 8 | **Modal Pilih Paket Topup Token** | `GET /api/v1/checkout/prepare` | Mengambil data `token_bundles` (Starter, Basic 25, Standard, dst). | ✅ Tested |
| 9 | **Halaman Promo REXTRA (List)** | `GET /api/v1/promo` | Menampilkan list voucher aktif. Filter Kategori (Semua, Club, Token) dilakukan di FE/BE. | ⚠️ Belum di-test detail |
| 10 | **Bottom Sheet Input Kode Promo** | `POST /api/v1/promo/validate` | Validasi kode manual. Menghasilkan error jika tidak ditemukan/tidak valid. | ✅ Tested |
| 11 | **Kalkulasi Harga (Real-time)** | `POST /api/v1/checkout/calculate` | Menghitung Subtotal, Diskon Promo, Kredit Sisa (jika upgrade), dan Total Akhir. | ✅ Tested |
| 12 | **Pilih Kanal Pembayaran & Biaya Admin** | `GET /api/v1/payment/channels?amount={total}` | Menampilkan list Tripay (BNI, QRIS, dll) + Biaya Admin (Flat & Persen). | ✅ Tested |
| 13 | **Initiate Checkout (Klik Bayar)** | `POST /api/v1/checkout/initiate` | Membuat transaksi di DB (Pending) & Generate URL Tripay. | ✅ Tested |
| 14 | **Detail Transaksi - Menunggu Bayar** | `GET /api/v1/my/payment/{id}` | Detail instruksi bayar, VA, URL Tripay, dan countdown waktu expired. | ✅ Tested |
| 15 | **Bottom Sheet Pembatalan** | `PUT /api/v1/admin/payment/{id}/cancel` | Mengirim alasan pembatalan (radio button/teks lain). | ✅ Tested |
| 16 | **Detail Transaksi - Dibatalkan** | `GET /api/v1/my/payment/{id}` | Menampilkan ilustrasi batal dan alasan yang tadi diinput. | ✅ Tested |
| 17 | **Detail Transaksi - Gagal (Expired)** | `GET /api/v1/my/payment/{id}` | Menampilkan status expired jika waktu bayar lewat. | ✅ Tested |
| 18 | **Detail Transaksi - Berhasil** | `GET /api/v1/my/payment/{id}` | Menampilkan info sukses, nominal lunas, dan tombol unduh bukti. | ✅ Tested |
| 19 | **Riwayat Transaksi (List)** | `GET /api/v1/my/payment` | List semua transaksi. Mendukung search ID & filter status. | ✅ Tested |
| 20 | **Unduh Bukti Pembayaran (Receipt)** | `GET /api/v1/my/payment/{id}/receipt` | **[BELUM ADA]** Endpoint untuk generate PDF/Image bukti bayar. | ❌ Belum Ada |

## 🛠️ Rencana Pengembangan Endpoint Baru

| Endpoint | Method | Tujuan | Payload / Requirement |
|:---|:---|:---|:---|
| `/api/v1/my/payment/{id}/receipt` | GET | Generate Bukti Bayar | Mengambil data transaksi sukses dan menghasilkan file PDF. |
| `/api/v1/promo/redeem` | POST | Pakai Voucher | Memastikan voucher "terkunci" untuk transaksi ini (sudah ada di logic initiate). |

## 📝 Catatan Penting untuk Frontend
- **State Kosong (Empty State):** Jika API List (Promo/Payment) mengembalikan array kosong `[]`, FE wajib menampilkan ilustrasi sesuai brief.
- **Filter Riwayat:** Backend sudah mendukung filter status di `GET /api/v1/my/payment?status=paid`.
- **Admin Fee:** Pastikan selalu mengirim query param `?amount=xxx` saat memanggil `/payment/channels` agar biaya admin Tripay akurat.

---
*Mapping ini disusun berdasarkan update terbaru sistem per 14 Maret 2026.*
