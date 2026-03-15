# REXTRA MISSING API SPECIFICATION

Berdasarkan analisis antara UI Brief dan Folder `module_fiks_banget`, berikut adalah daftar API yang belum ada atau perlu penyesuaian detail untuk mendukung UI Screen.

## 1. Endpoint: Unduh Bukti Pembayaran (Receipt)
**Screen Terkait:** 14 (Sukses New User), 20A (Sukses Upgrade), 24 (Detail Transaksi Berhasil)

| Requirement | Detail |
|:---|:---|
| **Endpoint** | `GET /api/v1/my/payment/{transaction_id}/receipt` |
| **Auth** | User (Owner of Transaction) |
| **Data Output** | File PDF atau Image (Faktur/Kwitansi) |
| **Logic BE** | Mengambil data dari `payment_transactions`, menggabungkannya dengan template HTML, lalu dikonversi ke PDF (misal pakai library `gofpdf` atau `wkhtmltopdf`). |
| **Status** | ❌ Belum Ada |

## 2. Endpoint: Riwayat Transaksi Gabungan (Sync Check)
**Screen Terkait:** 23 (Halaman Riwayat Transaksi)

| Requirement | Detail |
|:---|:---|
| **Endpoint** | `GET /api/v1/transaction-history/transactions` |
| **Logic BE** | Modul ini sudah ada briefnya di `module_fiks_banget`, tapi perlu dipastikan apakah sudah menggabungkan tabel `payment_transactions` (milik kita) dengan `topup_transaction` (milik modul token/Hamzah). |
| **Status** | ⚠️ Perlu Integrasi antar Modul |

## 3. Penyesuaian: Detail Transaksi (VA & Instruksi Bayar)
**Screen Terkait:** 10 (Detail Menunggu Pembayaran), 18 (Detail Upgrade), 19 (Detail Downgrade)

| Requirement | Detail |
|:---|:---|
| **Endpoint** | `GET /api/v1/my/payment/{id}` |
| **Adjustment** | Saat ini hanya return data basic. Harus dipastikan menyertakan data Tripay lengkap: `pay_code`, `payment_name`, `checkout_url`, dan list `instructions` (cara bayar). |
| **Status** | 🟡 Partial (Sudah ada tapi butuh diperkaya datanya) |

---

## 💡 Usulan Brief Baru: API Receipt
Untuk melancarkan UI Screen "Pembayaran Berhasil", kita perlu membuat modul `receipt` sederhana:

### Draft Spec API Receipt:
- **Tujuan:** Memberikan bukti resmi pembelian membership/token.
- **Payload Response:** Stream File (Content-Type: `application/pdf`).
- **Data yang ditampilkan di PDF:**
  - ID Transaksi & Tanggal Lunas.
  - Nama User & Email.
  - Nama Paket & Durasi.
  - Rincian Biaya (Subtotal, Diskon, Biaya Admin, Total).
  - Status: LUNAS.

---
*Dokumen ini dibuat untuk melengkapi REXTRA_API_MAPPING.md agar tim Backend & Frontend sinkron.*
