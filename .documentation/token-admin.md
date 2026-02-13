# Rextra Token Admin Documentation

Dokumen ini menjelaskan fungsionalitas, arsitektur, dan logika dari modul token di Rextra.

## Logika bisnis

- Nama dari **`Repository`** berkaitan dengan nama entity. Jika nama entity adalah `token_wallet`, maka nama repositorynya adalah `token_wallet_repository`.
- Logika bisnis sepenuhnya berada pada lapisan **`service`**. Setiap fungsi menggambarkan penggunaannya.


### Bundle Service
Bundle service adalah layanan yang digunakan untuk mengelola bundle token. Layanan ini berfungsi untuk mengelola data bundle token, seperti pembuatan, penghapusan, dan pembaruan data bundle token.

### Custom Pricing Service
Custom pricing service adalah layanan yang digunakan untuk mengelola harga kustom token. Dimana service ini berguna untuk mengartur harga dan diskon jika user melakukan pembelian secara custom(tidak melewati bundle).

- hanya terdapat **1 config custom pricing** yang **aktif**
- fungsi pada service berguna sesuai apa nama fungsi tersebut. Contoh fungsi `Getcurrent` digunakan untuk mengambil config custom pricing yang aktif.
- untuk fungsi `CreateNewConfig` akan mengaktifkan config custom pricing yang baru. Dan akan menonaktifkan config custom pricing yang lama. **sejauh ini belum ada service untuk memilih config mana yang berlaku, jadi selalu yang terbaru**

### Ledger Service
Ledger service adalah layanan yang digunakan untuk mengelola ledger token. Layanan ini berfungsi untuk mencatan semua aktivitas user yang berkaitan dengan token. Seperti pembelian, penjualan, dan pengembalian token.


## Summary Service
Summary service adalah layanan yang digunakan untuk mengelola summary token. Seperti data KPI dan statistik berdasarkan hari. GetGraphTrendByDirection dan GetGraphTrendBySourceType masing - masing berguna untuk mengambil statistik berdasarkan kondisi tertentu seperti Direction (IN / OUT) atau Source Type (Usage, Membership, dan Lain-lain).

- Query statistik pada database dijalankan **secara atomic**(satu persatu) dan **tidak** langsung **big query** di repository layer.

## Top up Transaction Service
Top up transaction service adalah layanan yang digunakan untuk mengelola top up token. Layanan ini berfungsi untuk mencatat semua aktivitas user yang berkaitan dengan top up token. Seperti pembelian, penjualan, dan pengembalian token.

## Wallet Service
Wallet service hanya digunakan oleh user untuk melihat pemakaian dan saldo token.
