**BRIEF CRUD SYSTEM JELAJAH PROFESI**

Dokumen use case untuk fitur Jelajah Profesi hingga Profesi Favorit bukan hanya deskripsi UI, melainkan representasi kebutuhan sistem yang harus diterjemahkan ke dalam lapisan arsitektur backend.

Dalam konteks implementasi menggunakan Golang dengan pola arsitektur yang clean dan scalable, setiap use case umumnya diterjemahkan menjadi tiga lapisan utama:

**Controller (Handler Layer)** → menerima request HTTP dari client, melakukan validasi awal, dan memanggil service.

**Service (Business Logic Layer)** → berisi logika bisnis utama sesuai use case.

**Repository (Data Access Layer)** → bertanggung jawab terhadap query database dan interaksi dengan entity.

**Artinya,** setiap baris use case di tabel bukan langsung menjadi controller atau repository, tetapi diterjemahkan menjadi fungsi di dalam service layer yang kemudian dipanggil oleh controller dan mengakses data melalui repository.

## **Karakteristik Sistem Ini**

Fitur Jelajah Profesi hingga Profesi Favorit merupakan sistem CRUD untuk USER (mahasiswa), namun secara karakteristik lebih dominan pada operasi READ.

**Sebagian besar fitur melibatkan:**

* Membaca data profesi  
* Membaca relasi kategori dan sub kategori  
* Membaca aktivitas, kompetensi, prospek  
* Membaca hasil kecocokan  
* Membaca data favorit

**Operasi tulis (WRITE) hanya terjadi pada:**

* Create favorit  
* Delete favorit  
* Update favorit (implicit toggle)  
* Create fit check result (di fitur lain)  
* Create career planning (future integration)

Dengan demikian, sistem ini secara desain merupakan **read-heavy system**, terutama pada entity `profession` dan relasinya.

## **Translasi Use Case ke Lapisan Backend**

Contoh sederhana:

Use case:  
 “Menampilkan detail profesi berdasarkan slug”

Translasi:

* Controller: `GET /professions/{slug}`  
* Service: `GetProfessionDetail(slug string, userID optional)`  
* Repository:  
  * GetBySlug  
  * GetActivitiesByProfessionID  
  * GetSkillsByProfessionID  
  * GetCareerPathsByProfessionID  
  * GetMarketInsightsByProfessionID  
  * GetFitCheckByUserAndProfessionID

Service akan mengorkestrasi beberapa repository call, lalu membentuk response DTO yang siap dikirim ke frontend.

# **Daftar Endpoint yang Diharapkan**

Berikut tabel endpoint yang perlu disediakan untuk fitur Jelajah Profesi hingga Favorit.

## **Jelajah Profesi**

| Method | Endpoint | Deskripsi |
| ----- | ----- | ----- |
| GET | /professions | Mengambil daftar profesi (dengan filter kategori, sub kategori, search, pagination) |
| GET | /professions/{slug} | Mengambil detail profesi lengkap |
| GET | /professions/recommended | Mengambil profesi rekomendasi berdasarkan RIASEC user |
| GET | /profession-categories | Mengambil daftar kategori utama |
| GET | /profession-categories/{id}/subcategories | Mengambil sub kategori berdasarkan kategori utama |

## **Fit Check (Personalized Section)**

| Method | Endpoint | Deskripsi |
| ----- | ----- | ----- |
| GET | /fit-check/{profession\_id} | Mengambil hasil kecocokan user terhadap profesi tertentu |

## **Profesi Favorit**

| Method | Endpoint | Deskripsi |
| ----- | ----- | ----- |
| GET | /user/favorites | Mengambil daftar profesi favorit user |
| POST | /user/favorites | Menambahkan profesi ke favorit |
| DELETE | /user/favorites/{profession\_id} | Menghapus profesi dari favorit |

## **Statistik Ringkasan**

| Method | Endpoint | Deskripsi |
| ----- | ----- | ----- |
| GET | /professions/{id}/metrics | Mengambil statistik agregat (target karier, jumlah rekomendasi, dll) |

# **Ringkasan Arsitektur**

* Controller hanya menangani HTTP dan response format.  
* Service mengimplementasikan seluruh use case.  
* Repository hanya berinteraksi dengan database.  
* Entity tetap bersih tanpa logika bisnis.  
* Response disusun dalam DTO agar tidak expose struktur internal database secara langsung.

**Halaman Utama Jelajah Profesi (Semua State)**

halaman ini muncul dengan banyak cara. Bisa dari klik icon fitur pada homepage, bisa dari hasil tes profil karier, dan lainnya

| Komponen / State | Use Case | Kebutuhan Fungsional Sistem (Teknis & Data) |
| :---- | :---- | :---- |
| **Load Halaman (Initial State)** | User membuka fitur Jelajah Profesi | Sistem memanggil **2 API paralel**: 1️⃣ API Membership (untuk akses fitur) 2️⃣ API Profesi (untuk kategori & katalog) Sistem menampilkan skeleton loader pada section membership, kategori, dan katalog profesi hingga response diterima. |
| **Overview Membership – Aktif & Mendukung** | Membership aktif dan memiliki akses Jelajah Profesi | Sistem menampilkan nama paket \+ status aktif. Tidak ada pembatasan pada akses profesi. Tidak menampilkan CTA upgrade. |
| **Overview Membership – Tidak Mendukung** | User Standard / tidak memiliki akses | Sistem menampilkan badge pembatasan akses. Section rekomendasi atau fitur premium dapat dikunci. Tampilkan CTA: “Upgrade Membership”. |
| **Overview Membership – Expired** | Membership sudah habis | Sistem menampilkan status expired \+ CTA “Perpanjang Membership”. Beberapa fitur bisa dikunci sesuai business rule. |
| **Hero Section \+ Search** | User melihat konteks fitur | Sistem menampilkan banner \+ search bar aktif. Saat user submit keyword → redirect ke halaman hasil pencarian dengan parameter `keyword`. Search akan mencocokkan `profession.name` dan `profession_alias.alias_name`. |
| **State: Sudah Memiliki Profil Karier (RIASEC)** | User telah menyelesaikan Tes Profil Karier | Sistem memanggil API rekomendasi berbasis `profession.riasec_code_id`. Sistem hanya menampilkan profesi yang memiliki `riasec_code_id` sesuai hasil user. Label: “Rekomendasi Terbaik” atau “Direkomendasikan” ditentukan oleh skor prioritas matching. |
| **State: Belum Memiliki Profil Karier** | User belum mengikuti tes | Sistem tidak memanggil API rekomendasi. Section rekomendasi disembunyikan dan diganti CTA “Ikuti Tes Profil Karier”. |
| **Kategori Utama Profesi** | User melihat daftar kategori utama | Sistem memanggil `profession_main_category` dan menampilkan maksimal 10 kategori dalam horizontal scroll. Data: `id`, `code`, `name`. Klik kategori → navigasi ke halaman hasil kategori dengan parameter `main_category_id`. |
| **Katalog Profesi Umum** | User melihat daftar profesi tanpa filter RIASEC | Sistem memanggil entity `profession` dengan join `profession_sub_category`. Data yang ditampilkan di card: \- `profession.slug` \- `profession.name` \- `profession.image_url` \- `profession_sub_category.name` |
| **Klik Detail Profesi** | User memilih salah satu profesi | Sistem navigasi menggunakan `profession.slug` sebagai public identifier → `/professions/{slug}` |
| **Simpan Favorit** | User klik icon favorit | Sistem mengirim request POST/DELETE ke endpoint favorit (butuh entity tambahan misal: `user_profession_favorite`). UI update secara optimistic tanpa reload. |
| **State: Data Kosong (Edge Case)** | API profesi berhasil tetapi tidak ada data | Sistem menampilkan empty state ilustrasi \+ CTA kembali ke kategori utama atau refresh. |
| **State: Loading** | API masih memuat data | Skeleton loader muncul pada: \- Overview membership \- Kategori utama \- Card profesi |
| **State: Error API** | API gagal (timeout / 500\) | Sistem menampilkan error ringan \+ tombol “Coba Lagi”. Tidak memblokir seluruh halaman jika hanya 1 API yang gagal (membership dan profesi dipisah). |

Ya, struktur ini sudah benar secara arsitektur:

* API Membership → domain akses  
* API Profesi → domain katalog  
* API Rekomendasi (opsional) → hanya dipanggil jika user punya profil RIASEC

Artinya: Halaman utama bisa memiliki **2–3 API paralel tergantung state user**, dan itu clean secara separation of concern.

**Halaman: Hasil Pencarian – Jelajah Profesi**

| Use Case | Deskripsi Fungsional | Data / Entity yang Dibutuhkan |
| ----- | ----- | ----- |
| Load Daftar Profesi Berdasarkan Kategori Utama | Sistem menampilkan daftar profesi berdasarkan `main_category_id` yang dipilih user (misal: ENGINEERING). Menampilkan total hasil dan list card profesi. | **profession**: `id`, `slug`, `name`, `image_url`, `main_category_id`, `sub_category_id` **profession\_main\_category**: `id`, `name` **profession\_sub\_category**: `id`, `name` |
| Load Daftar Profesi Berdasarkan Sub Kategori | Sistem memfilter profesi berdasarkan `sub_category_id` ketika user memilih chip sub kategori (misal: Backend, Mobile, dll). | **profession**: `id`, `slug`, `name`, `image_url`, `sub_category_id` **profession\_sub\_category**: `id`, `name`, `main_category_id` |
| Search Profesi (by Name & Alias) | Sistem memfilter daftar profesi berdasarkan input search bar. Pencarian dilakukan terhadap `profession.name` dan `profession_alias.alias_name`. | **profession**: `id`, `slug`, `name`, `image_url` **profession\_alias**: `alias_name`, `profession_id` |
| Tampilkan Total Hasil (“Temukan 99 Profesi”) | Sistem menampilkan jumlah total hasil pencarian berdasarkan filter aktif (kategori \+ sub kategori \+ keyword). | COUNT dari **profession.id** sesuai query filter |
| Tampilkan Card Profesi | Setiap card menampilkan gambar, nama profesi, nama sub kategori (badge), tombol “Lihat Profesi”, dan icon favorite. | **profession**: `id`, `slug`, `name`, `image_url`, `sub_category_id` **profession\_sub\_category**: `name` |
| Navigasi ke Detail Profesi | Saat tombol “Lihat Profesi” ditekan, sistem redirect menggunakan `slug` sebagai public identifier (`/professions/{slug}`). | **profession.slug** |
| Toggle Favorite Profesi | User dapat menandai profesi sebagai favorit (heart icon). | (Butuh entity tambahan user\_favorite\_profession — belum ada di brief, perlu dibuat jika fitur aktif) |
| State: Empty Result (Pencarian Tidak Ditemukan) | Jika hasil query kosong, sistem menampilkan ilustrasi dan pesan “Pencarian Tidak Ditemukan”. | Query **profession** menghasilkan 0 row |
| State: Loading | Skeleton loader saat request API berjalan. | Tidak perlu entity tambahan |
| State: Error | Jika API gagal, tampilkan error state. | Tidak perlu entity tambahan |

**Use Case – Halaman Katalog Profesi \- Kategori Utama**  
halaman ini muncul hanya ketika user klik card kategori utama pada halaman utama Jelajah Profesi

| Komponen / State | Use Case | Kebutuhan Fungsional Sistem (Teknis & Data) |
| ----- | ----- | ----- |
| **Load Halaman Berdasarkan Kategori Utama** | User klik salah satu kategori utama (misal: Engineering) | Sistem menerima parameter `main_category_id` atau `code`. Query: ambil semua `profession` dengan `main_category_id` tersebut. Join: `profession_sub_category` untuk badge sub kategori. |
| **Tampilkan Header Konteks Kategori** | User melihat konteks kategori aktif | Sistem menampilkan nama kategori dari `profession_main_category.name`. Data: `id`, `code`, `name`, `description`. |
| **Load Sub Kategori Berdasarkan Kategori Aktif** | User melihat chip sub kategori | Sistem memanggil `profession_sub_category` dengan filter `main_category_id = kategori_aktif`. Tidak boleh menampilkan sub kategori dari kategori lain (isolasi taksonomi wajib). |
| **Filter Berdasarkan Sub Kategori** | User klik salah satu chip sub kategori | Sistem memfilter `profession` berdasarkan kombinasi: `main_category_id` (tetap) \+ `sub_category_id` (dipilih). UI update tanpa reload penuh. |
| **Reset ke Semua Sub Kategori** | User klik chip “Semua” | Sistem menampilkan kembali seluruh `profession` dengan `main_category_id` aktif tanpa filter `sub_category_id`. |
| **Search dalam Kategori Aktif** | User mengetik di search bar | Sistem memfilter berdasarkan: `main_category_id` (tetap aktif) \+ keyword cocok di `profession.name` atau `profession_alias.alias_name`. |
| **Tampilkan Jumlah Profesi** | User melihat total hasil | Sistem menampilkan COUNT `profession.id` sesuai filter aktif (kategori \+ sub kategori \+ keyword). |
| **Tampilkan Card Profesi** | User melihat list profesi hasil filter | Data yang dibutuhkan per card: \- `profession.slug` \- `profession.name` \- `profession.image_url` \- `profession_sub_category.name` |
| **Navigasi ke Detail Profesi** | User klik “Lihat Profesi” | Sistem redirect menggunakan `slug` → `/professions/{slug}` |
| **Simpan Favorit** | User klik icon heart | Sistem kirim POST/DELETE ke endpoint favorit (entity user\_profession\_favorite jika tersedia). Update UI real-time (optimistic update). |
| **State: Sub Kategori Kosong (Edge Case)** | Kategori utama belum memiliki sub kategori | Sistem tetap menampilkan daftar profesi berdasarkan `main_category_id`. Section chip sub kategori disembunyikan. |
| **State: Data Kosong** | Tidak ada profesi dalam kategori tersebut | Sistem menampilkan empty state \+ CTA kembali ke kategori lain. |
| **State: Loading** | API masih memuat data | Skeleton loader muncul untuk header kategori, chip sub kategori, dan card profesi. |
| **State: Error API** | API gagal | Sistem menampilkan error ringan \+ tombol “Coba Lagi”. Jika sub kategori gagal tapi profesi berhasil, tetap tampilkan profesi (graceful degradation). |

**Use Case – Halaman Detail Profesi**

## **1\. Header dan Identitas Profesi**

| Komponen | Use Case | Sumber Data | State |
| ----- | ----- | ----- | ----- |
| Hero Image dan Nama Profesi | Menampilkan banner dan nama profesi yang dipilih berdasarkan slug | `profession.slug`, `profession.name`, `profession.image_url` | Jika slug tidak ditemukan → tampilkan 404 Profession Not Found |
| Kategori | Menampilkan kategori utama profesi | Join `profession.main_category_id` → `profession_main_category.name` | — |
| Sub Kategori | Menampilkan sub kategori profesi | Join `profession.sub_category_id` → `profession_sub_category.name` | — |

## **2\. Ringkasan Statistik**

| Komponen | Use Case | Sumber Data | Catatan |
| ----- | ----- | ----- | ----- |
| Jumlah Pengguna Menargetkan Karier | Menghitung jumlah user yang menjadikan profesi ini sebagai target karier | Tabel perencanaan karier user berdasarkan profession\_id (future integration) | Jika belum tersedia → tampilkan placeholder |
| Jumlah Pengguna Telah Berkarier | Menghitung jumlah user yang telah berkarier di bidang terkait | External metric / future integration | Ditandai sebagai metrik eksternal |
| Jumlah Direkomendasikan dari Tes Profil Karier | Menghitung berapa kali profesi ini muncul sebagai rekomendasi hasil tes | careerprofile\_test\_sessions berdasarkan profession\_id | Menggunakan agregasi COUNT |

**TAB PROFIL – Halaman Detail Profesi**

| Section | Komponen / Kondisi | Use Case | Sumber Data | State |
| ----- | ----- | ----- | ----- | ----- |
| Tentang Profesi | Deskripsi Profesi | Menampilkan penjelasan umum mengenai profesi | `profession.about_description` | Jika NULL → tampilkan empty state |
| Aktivitas Profesi | Daftar Aktivitas | Menampilkan aktivitas kerja utama profesi | `profession_activities.description`, `sort_order` | Jika tidak ada data → empty state |
| Kepribadian Profesi | Deskripsi RIASEC | Menampilkan kecocokan kepribadian profesi | `profession.riasec_description`, join `riasec_codes` | Jika `riasec_code_id` NULL → sembunyikan atau empty state |
| Hasil Kecocokan Profesi | User belum memiliki profil karier | Menampilkan banner ajakan mengikuti tes | — | Tidak menampilkan skor kecocokan |
| Hasil Kecocokan Profesi | User sudah memiliki profil karier | Menampilkan hasil kecocokan antara user dan profesi | `fit_check_results` berdasarkan `user_id` dan `profession_id` | Jika tidak ditemukan → tampilkan fallback state |
| Wawasan Profesi | Event / Edukasi | Menampilkan event atau kelas terkait profesi | Tabel event / explorations (entity terpisah) | Jika tidak ada data → empty state |

# **TAB KUALIFIKASI – Halaman Detail Profesi**

| Section | Komponen | Use Case | Sumber Data | State |
| ----- | ----- | ----- | ----- | ----- |
| Kompetensi Profesi | Hard Skills | Menampilkan kompetensi utama yang wajib atau dianjurkan | `profession_skill_rel` (filter `skill_type = 'hard'`) join `skills.name` | Jika tidak ada data → empty state |
| Kompetensi Profesi | Soft Skills | Menampilkan kompetensi pendukung profesi | `profession_skill_rel` (filter `skill_type = 'soft'`) join `skills.name` | Jika tidak ada data → empty state |
| Perangkat & Teknologi | Tools / Teknologi | Menampilkan perangkat atau teknologi yang digunakan dalam profesi | `profession_tool_rel` join `tools.name` | Jika tidak ada data → empty state |
| Pendidikan Profesi | Pendidikan Formal | Menampilkan jurusan atau program studi yang relevan | `profession_study_program_rel` join `study_program.name` | Jika tidak ada data → empty state |
| Pendidikan Profesi | Pendidikan Non-Formal | Menampilkan bootcamp, sertifikasi, atau pelatihan relevan | Future integration (explorations / program non-formal) | Jika tidak ada data → empty state |

**TAB PROSPEK – Halaman Detail Profesi**

| Section | Komponen | Use Case | Sumber Data | State |
| ----- | ----- | ----- | ----- | ----- |
| Jenjang Karier | Timeline Karier | Menampilkan tahapan karier berdasarkan level, pengalaman, dan estimasi gaji | `profession_career_paths.title`, `experience_range`, `salary_min`, `salary_max`, `sort_order` | Jika tidak ada data → empty state |
| Kondisi Pasar Kerja | Insight Pasar Kerja | Menampilkan bullet list kondisi dan tren pasar kerja untuk profesi tersebut | `profession_market_insights.description`, `sort_order` | Jika tidak ada data → empty state |

# **State Global Halaman**

| State | Perilaku Sistem |
| ----- | ----- |
| Loading | Menampilkan skeleton loader pada seluruh section |
| Error API | Menampilkan error ringan dengan tombol retry |
| Profession Tidak Ditemukan | Menampilkan halaman 404 jika slug tidak valid |

**Halaman Profesi Favorit**

| State / Section | Use Case | Sumber Data | Perilaku Sistem |
| ----- | ----- | ----- | ----- |
| Load Halaman (Initial) | Sistem memuat daftar profesi yang telah difavoritkan oleh user | user\_favorite\_professions berdasarkan user\_id, join profession | Menampilkan skeleton loader sebelum data tersedia |
| State: Kosong | User belum memiliki profesi favorit | Query user\_favorite\_professions menghasilkan 0 row | Menampilkan ilustrasi empty state dan teks “Belum Ada Profesi Yang Kamu Favoritkan” |
| State: Tidak Kosong | User memiliki minimal 1 profesi favorit | user\_favorite\_professions join profession.name, profession.image\_url, profession.main\_category\_id | Menampilkan daftar card profesi favorit |
| Card Profesi Favorit | Menampilkan nama, kategori, dan thumbnail profesi | profession.name, profession.image\_url, join profession\_main\_category.name | — |
| Klik “Lihat Profesi” | User membuka detail profesi | profession.slug | Redirect ke Halaman Detail Profesi |
| Hapus Favorit (Icon Hati) | User menekan icon hati untuk menghapus favorit | DELETE ke user\_favorite\_professions berdasarkan user\_id \+ profession\_id | Card langsung hilang tanpa reload halaman |
| State: Setelah Hapus Terakhir | Jika profesi terakhir dihapus | Hasil query menjadi 0 row | Otomatis berpindah ke Empty State |
| Error API | API gagal mengambil atau menghapus data | — | Menampilkan error ringan \+ opsi retry |
| Loading | Data sedang dimuat | — | Menampilkan skeleton card list |

Secara arsitektur:

Halaman ini minimal butuh:

* 1 API GET daftar favorit user  
* 1 API DELETE untuk menghapus favorit  
* Relasi ke tabel `profession`  
* Relasi ke `profession_main_category` untuk badge kategori

