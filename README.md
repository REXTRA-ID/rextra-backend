# Rextra Backend

Selamat datang di Rextra Backend! Proyek ini merupakan layanan backend yang dibangun menggunakan Go, dengan arsitektur modular yang modern dan skalabel. Dokumentasi ini akan memandu Anda dalam melakukan setup, memahami arsitektur, dan cara berkontribusi pada proyek.

##  Daftar Isi
0.  [Database Schema](https://dbdiagram.io/d/rextra-backend-697de7d0bd82f5fce234f212)
1.  [Cara Menjalankan Proyek](#1-cara-menjalankan-proyek)
    *   [Prasyarat](#prasyarat)
    *   [Konfigurasi](#konfigurasi)
    *   [Instalasi & Menjalankan (Docker - Direkomendasikan)](#instalasi--menjalankan-docker---direkomendasikan)
    *   [Instalasi & Menjalankan (Lokal)](#instalasi--menjalankan-lokal)
    *   [Perintah Penting](#perintah-penting)
2.  [Arsitektur Perangkat Lunak](#2-arsitektur-perangkat-lunak)
    *   [Filosofi Desain](#filosofi-desain)
    *   [Struktur Direktori](#struktur-direktori)
    *   [Diagram Arsitektur](#diagram-arsitektur)
3.  [Menambahkan Fitur Baru](#3-menambahkan-fitur-baru)
    *   [Langkah-langkah](#langkah-langkah)

---

## 0. Database Schema
Klik [disini](https://dbdiagram.io/d/rextra-backend-697de7d0bd82f5fce234f212) untuk melihat skema database.

## 1. Cara Menjalankan Proyek

Bagian ini menjelaskan cara menyiapkan dan menjalankan proyek di lingkungan pengembangan Anda.

### Prasyarat

Pastikan perangkat Anda telah terinstal:
*   [Go](https://golang.org/doc/install) (versi 1.21 atau lebih baru)
*   [Docker](https://www.docker.com/get-started) dan Docker Compose
*   [Make](https://www.gnu.org/software/make/)

### Konfigurasi

Proyek ini dikonfigurasi menggunakan *environment variables*. Untuk memulai, salin file konfigurasi contoh:

```bash
# Untuk lingkungan pengembangan dengan Docker
cp .env.example .env.dev
```
```bash
# Untuk lingkungan lokal tanpa Docker
cp .env.example .env.local
```

Kemudian, sesuaikan isi file `.env.dev` atau `.env.local` dengan konfigurasi lokal Anda (misalnya, kredensial database, key API, dll).

**Variabel Penting:**
*   `APP_MODE`: Mode aplikasi (`development`, `production`).
*   `APP_HOST`, `APP_PORT`: Host dan port untuk server.
*   `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`: Kredensial koneksi database PostgreSQL.
*   `REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`: Kredensial koneksi Redis.
*   `FIREBASE_CREDENTIALS_BASE64`: Kredensial Firebase dalam format base64.
*   `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`, `AWS_S3_BUCKET`: Konfigurasi untuk penyimpanan file S3.

### Instalasi & Menjalankan (Docker - Direkomendasikan)

Ini adalah cara yang paling mudah dan konsisten untuk menjalankan proyek, karena semua dependensi (Go, Postgres, Redis) sudah diatur di dalam Docker.

1.  **Build dan Jalankan Container:**
    Perintah ini akan membuat image Docker dan menjalankan semua layanan (aplikasi Go, database, Redis) di latar belakang.
    ```bash
    make up ENV=dev
    ```

2.  **Menjalankan Migrasi & Seeder Database:**
    Setelah container berjalan, buka terminal baru dan jalankan perintah berikut untuk menyiapkan skema database dan data awal.
    ```bash
    # Menjalankan migrasi
    make docker-migrate ENV=dev

    # Menjalankan seeder (opsional)
    make docker-seeder ENV=dev

    # Menjalankan keduanya
    make docker-both ENV=dev
    ```
    
3.  **Selesai!**
    Aplikasi sekarang berjalan dan dapat diakses di `http://<APP_HOST>:<APP_PORT>` sesuai konfigurasi `.env.dev` Anda.

4.  **Menghentikan Aplikasi:**
    Untuk menghentikan semua container, jalankan:
    ```bash
    make down ENV=dev
    ```

### Instalasi & Menjalankan (Lokal)

Gunakan metode ini jika Anda tidak ingin menggunakan Docker dan memilih untuk menjalankan semua layanan (Go, Postgres, Redis) secara manual di mesin lokal Anda.

1.  **Instalasi Dependensi Go:**
    ```bash
    make tidy
    ```

2.  **Menjalankan Migrasi & Seeder:**
    Pastikan database PostgreSQL dan Redis Anda sedang berjalan dan konfigurasinya sudah benar di `.env.local`.
    ```bash
    make both
    ```

3.  **Menjalankan Aplikasi:**
    Proyek ini menggunakan [Air](https://github.com/cosmtrek/air) untuk *live-reloading*. Instal Air terlebih dahulu:
    ```bash
    go install github.com/cosmtrek/air@latest
    ```
    Kemudian jalankan aplikasi:
    ```bash
    air
    ```
    Aplikasi akan otomatis di-restart setiap kali ada perubahan pada file `.go`.

### Perintah Penting

Gunakan `make help` untuk melihat semua perintah yang tersedia. Beberapa yang paling umum:

*   `make tidy`: Merapikan dependensi Go.
*   `make up ENV=dev`: Menjalankan lingkungan pengembangan Docker.
*   `make down ENV=dev`: Menghentikan lingkungan pengembangan Docker.
*   `make reset ENV=dev`: Menghentikan dan menghapus data (volume) dari lingkungan Docker.
*   `make build-docker ENV=dev`: Mem-build ulang image Docker jika ada perubahan pada Dockerfile.
*   `make docker-migrate ENV=dev`: Menjalankan migrasi database di dalam container Docker yang sedang berjalan.

---

## 2. Arsitektur Perangkat Lunak

### Filosofi Desain

Arsitektur proyek ini dirancang dengan beberapa prinsip utama:

*   **Modular**: Kode diorganisir ke dalam modul-modul fitur yang independen (`internal/modules`). Setiap modul bertanggung jawab atas satu domain bisnis (misalnya `auth`, `kenali_diri`). Ini membuat kode lebih mudah ditemukan, dikelola, dan diuji.
*   **Layered Architecture**: Setiap modul mengikuti arsitektur berlapis (Controller, Service, Repository) untuk memisahkan tanggung jawab dengan jelas.
*   **Dependency Injection**: Ketergantungan (seperti koneksi database atau service lain) di-pass sebagai parameter (diinjeksi), bukan dibuat di dalam komponen. Ini membuat komponen lebih fleksibel dan mudah di-test.
*   **Shared Packages**: Kode yang dapat digunakan kembali oleh banyak modul (seperti logger, koneksi cache, utilitas JWT) ditempatkan di `internal/pkg`.

### Struktur Direktori

Berikut adalah gambaran umum dari direktori-direktori penting:

```
rextra-backend/
├── cmd/                # Titik masuk aplikasi jika memiliki beberapa binary
├── db/                 # Migrasi database dan seeder
│   ├── migrations/
│   └── seeder/
├── internal/           # Kode utama aplikasi yang tidak untuk diimpor oleh proyek lain
│   ├── config/         # Konfigurasi dan inisialisasi aplikasi (wiring)
│   ├── dto/            # Data Transfer Objects (struct untuk request/response JSON)
│   ├── entity/         # Struct Entitas GORM (merepresentasikan tabel database)
│   ├── middleware/     # Middleware untuk Gin (misalnya, autentikasi)
│   ├── modules/        # **Direktori paling penting**: Berisi modul-modul fitur
│   │   ├── auth/       # Contoh: Modul Autentikasi
│   │   │   ├── controller/
│   │   │   ├── service/
│   │   │   ├── repository/
│   │   │   ├── routes/
│   │   │   └── module.go # Titik masuk/inisialisasi modul
│   │   └── kenali_diri/ # Contoh: Modul Kenali Diri
│   └── pkg/            # Paket-paket yang bisa digunakan di seluruh aplikasi (shared)
│       ├── cache/
│       ├── jwt/
│       └── ...
├── main.go             # Titik masuk utama aplikasi
├── Makefile            # Kumpulan shortcut perintah untuk development
├── go.mod              # Definisi modul dan dependensi Go
└── docker-compose.yml  # Definisi layanan untuk Docker
```

### Diagram Arsitektur

Secara konseptual, alur permintaan (request) dalam satu modul terlihat seperti ini:

```
Request Flow:
[Internet] -> [Gin Engine] -> [Middleware] -> [Controller] -> [Service] -> [Repository] -> [Database/Cache]
```

Dan struktur antar modul adalah sebagai berikut:

```
Application Structure:
+-------------------------------------------------------------+
|                          main.go                            |
|        (Inisialisasi Config, DB, Cache, dan Modul)          |
+-------------------------------------------------------------+
|                    internal/config                          |
|             (Memanggil InitModule() dari setiap modul)      |
+-------------------------------------------------------------+
|    internal/modules/                                        |
| +------------------+ +------------------+ +------------------+
| |      Auth        | |   Kenali Diri    | |      ...         |
| |------------------| |------------------| |------------------|
| | - Controller     | | - Controller     | | - Controller     |
| | - Service        | | - Service        | | - Service        |
| | - Repository     | | - Repository     | | - Repository     |
| +------------------+ +------------------+ +------------------+
+-------------------------------------------------------------+
|    internal/pkg/ (Cache, JWT, Logger, dll)                  |
|          (Digunakan oleh semua modul)                       |
+-------------------------------------------------------------+
```

---

## 3. Menambahkan Fitur Baru

Panduan ini menjelaskan alur kerja untuk menambahkan modul fitur baru (misalnya, modul "artikel").

### Langkah-langkah

1.  **Buat Direktori Modul Baru:**
    Buat direktori baru di dalam `internal/modules/`.
    ```bash
    mkdir internal/modules/artikel
    ```

2.  **Buat Struktur Direktori Internal:**
    Di dalam `internal/modules/artikel`, buat direktori untuk setiap lapisan.
    ```bash
    mkdir internal/modules/artikel/controller
    mkdir internal/modules/artikel/service
    mkdir internal/modules/artikel/repository
    mkdir internal/modules/artikel/routes
    ```

3.  **Definisikan Entity dan DTO:**
    *   Jika fitur ini membutuhkan tabel database baru, buat struct GORM di `internal/entity/artikel_entity.go`.
    *   Buat struct untuk request (input) dan response (output) JSON di `internal/dto/request/artikel_dto_request.go` dan `internal/dto/response/artikel_dto_response.go`.

4.  **Implementasikan Lapisan (dari bawah ke atas):**
    *   **Repository**: Buat file `internal/modules/artikel/repository/artikel_repository.go`. Implementasikan fungsi untuk mengakses database (CRUD - Create, Read, Update, Delete) terkait artikel.
    *   **Service**: Buat file `internal/modules/artikel/service/artikel_service.go`. Implementasikan logika bisnis di sini (misalnya, validasi sebelum menyimpan). Service ini akan memanggil fungsi dari repository.
    *   **Controller**: Buat file `internal/modules/artikel/controller/artikel_controller.go`. Buat fungsi-fungsi (handler) yang akan dipanggil oleh Gin. Fungsi ini menerima request, memanggil service, dan mengembalikan response.

5.  **Definisikan Rute API:**
    Buat file `internal/modules/artikel/routes/artikel_route.go`. Buat sebuah fungsi (misalnya `ServeArtikel`) yang mendaftarkan semua endpoint API untuk artikel (`/artikels`, `/artikels/:id`, dll) ke `*gin.Engine`.

6.  **Buat Titik Masuk Modul (`module.go`):**
    Buat file `internal/modules/artikel/module.go`. Di dalamnya, buat fungsi `InitModule` yang melakukan "wiring" atau penyambungan semua komponen yang telah dibuat (repository, service, controller, routes).

    ```go
    // internal/modules/artikel/module.go
    package artikel

    import (
        // ... imports
    )

    func InitModule(server *gin.Engine, db *gorm.DB, middleware middleware.Middleware) {
        // 1. Init Repository
        artikelRepo := repository.NewArtikel(db)
        
        // 2. Init Service
        artikelSvc := service.NewArtikel(artikelRepo)

        // 3. Init Controller
        artikelCtrl := controller.NewArtikel(artikelSvc)

        // 4. Init Routes
        routes.ServeArtikel(server, artikelCtrl, middleware)
    }
    ```

7.  **Daftarkan Modul ke Aplikasi:**
    Langkah terakhir adalah "menyalakan" modul Anda. Buka file `internal/config/rest_config.go` dan panggil `InitModule` dari modul artikel di dalam fungsi `NewRest`.

    ```go
    // internal/config/rest_config.go
    import "rextra-backend/internal/modules/artikel" // Jangan lupa import

    func NewRest(...) RestConfig {
        // ...
        
        // Module
        auth.InitModule(server, db, middleware)
        persona.InitModule(server, db, middleware)
        // ... modul lainnya
        artikel.InitModule(server, db, middleware) // <-- Daftarkan modul baru Anda di sini

        // ...
    }
    ```

8.  **Buat Migrasi Database (jika perlu):**
    Jika Anda menambahkan tabel baru, buat file migrasi di `db/migrations/` dan jalankan `make docker-migrate ENV=dev`.

Dengan mengikuti langkah-langkah ini, fitur baru Anda akan terintegrasi dengan rapi ke dalam arsitektur proyek yang sudah ada.
