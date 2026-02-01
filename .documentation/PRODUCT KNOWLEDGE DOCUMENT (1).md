## 

| Nama Tugas | Pembuatan Entity Kenali Diri (Postgree SQL) \- Feedback  |
| :---- | :---- |
| **Penanggungjawab** | Dedy |
| **Status pengerjaan** | Selesai |
| **Status pengecekan** | Belum Diperiksa |
| **Tanggal Mulai** | 15 Desember  2025  |
| **Tanggal Selesai** | 20 November 2025 |

## **Brief Penugasan Entity Feedback Kenali Diri**

Pastikan tim memahami konteks fitur **Kenali Diri** sebagai rangkaian asesmen yang menghasilkan output **profil** dan **rekomendasi** (contoh awal: **Tes Profil Karier**). Pada tahap ini yang dibangun adalah **entity database (PostgreSQL)** untuk menampung **umpan balik pengguna** terkait pengalaman penggunaan dan kualitas hasil asesmen. Implementasi REST API **tidak termasuk** ruang lingkup penugasan ini (endpoint ditangani tim AI via FastAPI), namun seluruh data tetap tersimpan pada **database yang sama (PostgreSQL)**.  
Dokumen materi spesifikasi: [https://docs.google.com/document/d/1g5XwoupcnXB3e\_UUxSTwVBW0Gxoood4JSQ5mfzqpxWs/edit?tab=t.g6c6u4ziczte](https://docs.google.com/document/d/1g5XwoupcnXB3e_UUxSTwVBW0Gxoood4JSQ5mfzqpxWs/edit?tab=t.g6c6u4ziczte)

### **1\. Tujuan Entity Feedback**

1. Mengukur pengalaman pengguna setelah tes: **kemudahan**, **relevansi**, **kepuasan**, dan **kendala**.  
2. Mengumpulkan evaluasi **validitas output** dari Expert/Validator: **akurasi**, **logika**, dan **manfaat**.  
3. Menjadi sumber data dashboard admin (tabel mentah \+ visualisasi tren/partisipasi).

### **2\. Target Peran dan Kegunaan**

1. **Mahasiswa**: mengisi feedback setelah menyelesaikan tes (fokus experience).  
2. **Expert/Validator**: mengisi feedback untuk menilai kualitas hasil (fokus validitas output).  
3. **Admin**: memantau partisipasi, gap pengisian feedback, serta prioritas perbaikan.

### **3\. Mekanisme Flow Pengumpulan Feedback**

1. Feedback hanya dapat dikirim **setelah sesi tes selesai**.  
2. Setiap feedback wajib terikat ke **1 test session** agar metrik “Peserta Tes vs Pengisi Feedback” konsisten.  
3. Aturan anti-duplikasi: **maksimal 1 feedback per role untuk sesi yang sama** (Mahasiswa 1x, Expert 1x).  
4. Seluruh agregasi dashboard mengikuti filter: **kategori tes \+ rentang waktu**.  
   

### **4\. Struktur Entity yang Dibutuhkan**

Struktur dibagi menjadi:

1. **Tabel induk (universal)** untuk menyatukan feedback lintas kategori tes dan lintas role.  
2. **Tabel detail per kategori tes**, karena pertanyaan feedback dapat berbeda antar tes (agar schema tidak melebar).

#### **4.1 kenalidiri\_feedback (tabel induk / header)**

* Menyimpan **1 record per pengiriman feedback** dan menjadi entry point list feedback admin.  
* Versi awal mengikat sesi tes **Tes Profil Karier** ke careerprofile\_test\_sessions (langsung/atau pola kategori+session\_id sesuai desain final).  
* Fungsi: menyimpan konteks dasar (kategori tes, sesi tes, role pengisi, waktu submit) sebagai basis query total feedback, response rate, tren.

#### **4.2 careerprofile\_feedback\_student (detail Mahasiswa — Tes Profil Karier)**

* Menyimpan skor Likert 1–7: **kemudahan**, **relevansi**, **kepuasan** \+ pesan opsional untuk tim.  
* Kendala **tidak disimpan JSON**; menggunakan desain relasional:  
  * careerprofile\_obstacle\_options (master kendala mahasiswa)  
  * careerprofile\_feedback\_obstacles (pivot multi-select kendala mahasiswa)

#### **4.3 careerprofile\_feedback\_expert (detail Expert — Tes Profil Karier)**

* Menyimpan skor Likert 1–7: **akurasi**, **logika**, **manfaat** \+ saran perbaikan (opsional).  
* Menyimpan **snapshot Top 5 rekomendasi** dalam top5\_recommendations\_json.  
* Menyimpan expert\_profession\_id untuk matching, dan top5\_status dihitung sebagai **GENERATED COLUMN** (P1/P2/P3–5/NOT\_PRESENT).  
* Menyimpan **snapshot identitas expert** (untuk audit trail dan stabilitas tampilan admin).  
* Kendala expert memakai tabel terpisah (daftar berbeda dari mahasiswa):  
  * careerprofile\_expert\_obstacle\_options (master kendala expert)  
  * careerprofile\_feedback\_expert\_obstacles (pivot multi-select kendala expert)

### **5\. Prinsip Skalabilitas untuk Tes Baru**

Jika kategori tes baru ditambahkan dan pertanyaan feedback berbeda:

1. kenalidiri\_feedback tetap menjadi tabel induk.  
2. Ditambah tabel detail baru per kategori dan per role (mis. disc\_feedback\_student, disc\_feedback\_expert, dst).  
3. Kendala mengikuti pola yang sama: master \+ pivot khusus bila daftar kendala berbeda.

## **5\. Rincian Brief Pembuatan Entity**

## **5.1 Tabel `kenalidiri_feedback` (Revisi Final)**

kenalidiri\_feedback adalah **tabel induk (header)** yang mencatat **1 kali pengiriman feedback** untuk **1 sesi tes** pada **1 kategori tes** oleh **1 responden** (Mahasiswa atau Expert). Tabel ini menjadi **entry point** untuk seluruh rekap admin, sedangkan isi feedback detail disimpan pada tabel turunan yang spesifik per kategori tes.

### **2\. Aturan Bisnis Utama**

* **Satu sesi tes hanya boleh punya 1 feedback per responden (per role).**  
   Artinya:  
  * Mahasiswa: maksimal 1 feedback untuk sesi tersebut.  
  * Expert: maksimal 1 feedback untuk sesi tersebut (per expert).

* **Feedback bersifat final saat submit (immutable).**  
   Setelah tersimpan, feedback **tidak diedit** dan user **tidak diminta isi ulang** untuk sesi yang sama.

* **Record dibuat hanya saat submit.**  
  Tidak ada konsep draft di tabel ini (kalau user belum submit → tidak ada record).

### **3\. Desain Referensi Sesi (Polymorphic)**

Karena kategori tes akan bertambah, tabel ini memakai: test\_category \+ test\_session\_id  
Validitas test\_session\_id terhadap tabel sesi yang benar **wajib divalidasi di layer aplikasi** sebelum INSERT.

### **4\. Struktur Kolom**

| Kolom | Tipe | Constraint | Keterangan |
| :---- | :---- | :---- | :---- |
| id | BIGSERIAL | PK | ID unik feedback. |
| test\_category | VARCHAR/ENUM | NOT NULL | Kategori tes (mis. CAREER\_PROFILE). |
| test\_session\_id | BIGINT | NOT NULL | ID sesi tes sesuai kategori (polymorphic). |
| respondent\_type | VARCHAR/ENUM | NOT NULL | STUDENT / EXPERT. |
| respondent\_user\_id | BIGINT | NOT NULL, FK users.id | Pengirim feedback. |
| submitted\_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() | Waktu submit (disarankan simpan UTC di backend). |
| deleted\_at | TIMESTAMPTZ | NULL *(opsional)* | Soft delete untuk governance admin. |
| deleted\_by | BIGINT | NULL *(opsional)* | Admin yang menghapus (FK users.id jika ada role admin). |

Karena feedback **tidak diedit**, kolom updated\_at tidak dibutuhkan.

### **5\. ENUM yang Dipakai**

* respondent\_type: STUDENT, EXPERT  
* test\_category (v1): CAREER\_PROFILE (kategori lain menyusul)

### **6\. Constraint Anti-Duplikasi (Wajib)**

Karena soft delete opsional, bentuk paling aman adalah **partial unique index** (PostgreSQL):

* Jika **tanpa soft delete**:

UNIQUE (test\_category, test\_session\_id, respondent\_type, respondent\_user\_id)

* Jika **pakai soft delete**:

CREATE UNIQUE INDEX uq\_feedback\_once

ON kenalidiri\_feedback (test\_category, test\_session\_id, respondent\_type, respondent\_user\_id)

WHERE deleted\_at IS NULL;

Makna: **1 user hanya boleh submit 1 kali** untuk kombinasi sesi+kategori+role yang sama (kecuali record itu sudah dihapus admin via soft delete).

### **7\. Index yang Disarankan (Admin Dashboard)**

Wajib untuk kebutuhan filter \+ agregasi waktu:

* (test\_category, submitted\_at)  
* (test\_category, test\_session\_id)  
* **Tambahan yang sangat berguna untuk chart/rekap role**:  
   (test\_category, respondent\_type, submitted\_at)

### **8\. Relasi**

* users (1) → (N) kenalidiri\_feedback via respondent\_user\_id  
* kenalidiri\_feedback (1) → (0..1) \<detail table\> sesuai kategori \+ respondent\_type, contoh untuk v1:  
  * careerprofile\_feedback\_student  
  * careerprofile\_feedback\_expert

### **9\. Catatan Validasi Polymorphic (Wajib didokumentasikan)**

Sebelum INSERT, backend harus memastikan:

* test\_category valid  
* test\_session\_id ada di tabel sesi kategori tersebut  
* sesi tersebut **eligible** untuk feedback (mis. status completed)

Itu saja perubahan yang diperlukan—setelah ini brief sudah selaras dengan aturan “**1 sesi cukup 1 kali feedback**” dan tetap scalable untuk kategori tes baru.

## **5.2 Entity Feedback Mahasiswa — Tes Profil Karier (Revisi)**

### **5.2.1 Tabel careerprofile\_feedback\_student**

Tabel detail untuk menyimpan **isi feedback Mahasiswa** pada **Tes Profil Karier**. Satu record mewakili 1 feedback yang sudah disubmit dan terhubung 1:1 dengan kenalidiri\_feedback.

#### **2\) Tujuan**

* Menyimpan **skor Likert 1–7** yang jadi sumber chart kemudahan, relevansi, dan kepuasan.  
* Menyimpan **pesan bebas** untuk masukan kualitatif.  
* Kendala **tidak disimpan di sini**, tetapi melalui tabel pivot agar agregasi chart efisien.

#### **3\) Struktur Kolom**

| Kolom | Tipe | Constraint | Deskripsi |
| ----- | ----- | ----- | ----- |
| feedback\_id | BIGINT | PK, FK → kenalidiri\_feedback.id ON DELETE CASCADE | ID header feedback. |
| ease\_score | SMALLINT | NOT NULL, CHECK 1–7 | Kemudahan menyelesaikan tes. |
| relevance\_score | SMALLINT | NOT NULL, CHECK 1–7 | Relevansi rekomendasi profesi. |
| satisfaction\_score | SMALLINT | NOT NULL, CHECK 1–7 | Kepuasan keseluruhan. |
| message\_to\_team | TEXT | NULL | Pesan bebas untuk tim REXTRA. |

#### **4\) Aturan Bisnis**

* Wajib valid untuk header: test\_category=CAREER\_PROFILE dan respondent\_type=STUDENT (divalidasi di layer aplikasi).  
* Skor wajib 1–7.  
* Kendala dicatat lewat pivot (lihat 5.2.3).

### **5.2.2 Tabel careerprofile\_obstacle\_options (Master Kendala)**

Tabel referensi untuk daftar kendala yang ditampilkan pada checkbox (fixed list), termasuk flag khusus “Tidak ada kendala” dan “Lainnya”.

#### **2\) Struktur Kolom**

| Kolom | Tipe | Constraint | Deskripsi |
| ----- | ----- | ----- | ----- |
| id | SERIAL | PK | ID opsi kendala. |
| key | VARCHAR(50) | UNIQUE, NOT NULL | Key stabil untuk agregasi (mis. BUG\_ERROR). |
| label | TEXT | NOT NULL | Label UI (mis. “Mengalami error/bug”). |
| is\_no\_issue | BOOLEAN | NOT NULL, DEFAULT FALSE | True hanya untuk “Tidak ada kendala”. |
| is\_other | BOOLEAN | NOT NULL, DEFAULT FALSE | True hanya untuk “Lainnya”. |
| sort\_order | INT | NOT NULL | Urutan tampil/urutan chart. |
| is\_active | BOOLEAN | NOT NULL, DEFAULT TRUE | Untuk menonaktifkan opsi tanpa menghapus histori. |

#### **3\) Seed Data (contoh v1)**

* NO\_ISSUE – Tidak ada kendala (is\_no\_issue=true)  
* DURATION\_TOO\_LONG  
* QUESTIONS\_CONFUSING  
* BUG\_ERROR  
* RESULT\_EXPLANATION\_TOO\_LONG  
* UI\_NAV\_CONFUSING  
* OTHER (is\_other=true)

### **5.2.3 Tabel careerprofile\_feedback\_obstacles (Pivot Multi-select Kendala)**

#### **1\) Definisi**

Tabel pivot untuk menyimpan pilihan kendala (multi-select). Satu feedback dapat memilih banyak kendala.

#### **2\) Struktur Kolom**

| Kolom | Tipe | Constraint | Deskripsi |
| ----- | ----- | ----- | ----- |
| feedback\_id | BIGINT | FK → careerprofile\_feedback\_student.feedback\_id ON DELETE CASCADE | Relasi ke feedback mahasiswa. |
| obstacle\_id | INT | FK → careerprofile\_obstacle\_options.id | Opsi kendala yang dipilih. |
| other\_text | TEXT | NULL | Diisi hanya jika obstacle\_id adalah opsi OTHER. |
| **PK** |  | PRIMARY KEY (feedback\_id,obstacle\_id) | Cegah duplikasi pilihan kendala pada feedback yang sama. |

#### **3\) Index yang disarankan**

* INDEX (obstacle\_id) → untuk Chart 7 “Ragam Kendala”.  
* (Opsional) INDEX (feedback\_id) kalau sering lookup detail.

#### **4\) Aturan Bisnis (penting untuk konsistensi)**

* Jika memilih NO\_ISSUE, maka **tidak boleh memilih kendala lain** (enforce di layer aplikasi).  
* Jika memilih OTHER, maka other\_text **disarankan** terisi.

## **5.2.4 Relasi Ringkas**

* kenalidiri\_feedback (1) → (1) careerprofile\_feedback\_student  
* careerprofile\_feedback\_student (1) → (N) careerprofile\_feedback\_obstacles  
* careerprofile\_obstacle\_options (1) → (N) careerprofile\_feedback\_obstacles

## **5.2.5 Contoh Data (Ilustratif)**

**Header** kenalidiri\_feedback: id=12031, test\_category=CAREER\_PROFILE, respondent\_type=STUDENT, test\_session\_id=88901

**Detail** careerprofile\_feedback\_student:

* feedback\_id=12031  
* ease\_score=6, relevance\_score=5, satisfaction\_score=6  
* message\_to\_team="Overall bagus, tapi durasi terasa panjang."

**Kendala** careerprofile\_feedback\_obstacles:  
(feedback\_id=12031, obstacle\_id= \[DURATION\_TOO\_LONG\])

Bedanya ada di **target responden dan daftar opsinya**.

1. **`careerprofile_obstacle_options` \+ `careerprofile_feedback_obstacles`**  
    → Ini khusus **kendala Mahasiswa** (opsinya: durasi lama, pertanyaan membingungkan, UI bingung, dll) dan pivot-nya nge-link ke **`careerprofile_feedback_student`**.

2. **`careerprofile_expert_obstacle_options` \+ `careerprofile_feedback_expert_obstacles`**  
    → Ini khusus **kendala Expert** (opsinya: rekomendasi tidak relevan, penjelasan terlalu umum, istilah membingungkan, dll) dan pivot-nya nge-link ke **`careerprofile_feedback_expert`**.

## **5.3 Entity Feedback Expert — Tes Profil Karier** 

### **1\. Definisi**

careerprofile\_feedback\_expert adalah tabel detail yang menyimpan **feedback Expert/Validator** untuk kategori **Tes Profil Karier**. Satu record mewakili 1 pengiriman feedback expert dan terhubung **1:1** dengan kenalidiri\_feedback.

Feedback ini berfokus pada **validitas hasil** (akurasi, logika, manfaat), serta mencatat **snapshot Top 5 rekomendasi** pada saat feedback dibuat agar tampilan admin tetap konsisten meskipun algoritma rekomendasi berubah di masa depan.

### **2\. Mekanisme “Sistem tahu profesi Expert” (Wajib)**

Agar “Top 5 status” bisa dihitung dengan benar, sistem membutuhkan **profesi expert yang dipilih saat mengisi feedback**:

1. Saat form dibuka, sistem menampilkan **Top 5 rekomendasi** (rank 1–5).  
2. Expert memilih: **“Profesi Anda yang mana?”** dari dropdown yang memetakan ke profession\_id, dengan opsi tambahan **“Tidak ada di Top 5”**.  
3. Nilai pilihan disimpan pada expert\_profession\_id (boleh NULL jika memilih “Tidak ada”).

Dengan ini, status P1/P2/P3–5/Tidak muncul dapat dihitung konsisten dari kombinasi:

* expert\_profession\_id  
* top5\_recommendations\_json

### **3\. Prinsip Desain (Anti-Inkonsistensi)**

* top5\_recommendations\_json adalah **single source of truth** snapshot Top 5\.  
* top5\_status **tidak diinput manual**, melainkan **GENERATED COLUMN** agar:  
  * tidak redundant,  
  * tidak bisa tidak sinkron,  
  * tetap cepat untuk filter admin.

## **5.3.1 Tabel careerprofile\_feedback\_expert**

### **4\. Struktur Kolom dan Atribut**

| Kolom | Tipe | Constraint | Deskripsi |
| :---- | :---- | :---- | :---- |
| feedback\_id | BIGINT | PK, FK → kenalidiri\_feedback.id ON DELETE CASCADE | Relasi 1:1 ke header feedback. |
| accuracy\_score | SMALLINT | NOT NULL, CHECK 1–7 | Akurasi hasil menurut expert. |
| logic\_score | SMALLINT | NOT NULL, CHECK 1–7 | Kelogisan penjelasan hasil. |
| usefulness\_score | SMALLINT | NOT NULL, CHECK 1–7 | Potensi manfaat untuk user. |
| expert\_name | VARCHAR(255) | NOT NULL | Snapshot nama expert saat submit. |
| expert\_profession | VARCHAR(120) | NOT NULL | Snapshot teks profesi untuk display. |
| expert\_profession\_id | BIGINT | NULL, FK → professions.id | Dipilih expert dari dropdown Top 5; NULL jika “Tidak ada di Top 5”. |
| expert\_degree | VARCHAR(120) | NULL | Snapshot gelar (opsional). |
| expert\_experience\_years | SMALLINT | NULL | Snapshot pengalaman (tahun). |
| expert\_education\_level | VARCHAR(60) | NULL | Snapshot pendidikan (mis. “Sarjana (S1)”). |
| expert\_university | VARCHAR(255) | NULL | Snapshot perguruan tinggi (opsional). |
| expert\_study\_program | VARCHAR(255) | NULL | Snapshot program studi (opsional). |
| top5\_recommendations\_json | JSONB | NOT NULL | Snapshot Top 5 rekomendasi rank 1–5. |
| top5\_status | VARCHAR(20) | GENERATED STORED | Status otomatis: P1/P2/P3\_5/NOT\_PRESENT. |
| suggestion\_text | TEXT | NULL | Masukan/saran perbaikan dari expert. |

### **5\. Nilai top5\_status (Generated)**

* P1, P2, P3\_5, NOT\_PRESENT

### **6\. Format top5\_recommendations\_json**

Array 5 object wajib (rank 1–5). Minimal menyimpan:

* rank  
* profession\_id  
* profession\_name (snapshot)

**Contoh:**

\[

{"rank":1,"profession\_id":55,"profession\_name":"Product Manager"},

 {"rank":2,"profession\_id":101,"profession\_name":"UI/UX Designer"},

{"rank":3,"profession\_id":88,"profession\_name":"Graphic Designer"},

  {"rank":4,"profession\_id":330,"profession\_name":"UX Researcher"},

{"rank":5,"profession\_id":412,"profession\_name":"Interaction Designer"}

\]

### **7\. Generated Column top5\_status (konsep)**

**Aturan:**

* Jika expert\_profession\_id NULL → NOT\_PRESENT  
* Jika match profession\_id di rank 1 → P1, rank 2 → P2, rank 3–5 → P3\_5  
* Selain itu → NOT\_PRESENT

Index yang disarankan untuk tabel admin:

* INDEX (top5\_status) agar filter cepat.

## **5.3.2 Kendala Expert (Final: Tabel Terpisah)**

### **A) careerprofile\_expert\_obstacle\_options (Master)**

Opsi kendala expert berbeda dari mahasiswa, sehingga dibuat master sendiri:

* NO\_ISSUE (Tidak ada kendala/lancar)  
* IRRELEVANT\_RECOMMENDATIONS (Rekomendasi sangat tidak relevan)  
* CONFUSING\_TERMS (Istilah/kalimat membingungkan)  
* TOO\_GENERIC (Penjelasan terlalu umum)  
* TECHNICAL\_BUG (Bug teknis)  
* OTHER (Lainnya)

Struktur kolom sama seperti master mahasiswa: id, key, label, is\_no\_issue, is\_other, sort\_order, is\_active.

### **B) careerprofile\_feedback\_expert\_obstacles (Pivot)**

Pivot multi-select kendala expert:

* PK (feedback\_id, obstacle\_id)  
* other\_text untuk opsi OTHER  
* INDEX (obstacle\_id) untuk agregasi chart

**Aturan:**

* Jika pilih NO\_ISSUE → tidak boleh pilih opsi lain.  
* Jika pilih OTHER → other\_text disarankan wajib.

### **8\. Catatan Implementasi Penting**

* Snapshot identitas expert disimpan untuk menjaga **audit trail** (data tidak berubah walau profil user berubah).  
* top5\_status tidak boleh diinput manual; wajib generated agar tidak inkonsisten.  
* Kendala expert dipisah dari mahasiswa agar analisis tidak tercampur.

## **5.3.3 Tabel `careerprofile_expert_obstacle_options` (Master Kendala Expert)**

### **1\. Definisi**

Tabel ini menyimpan **master data** berupa daftar opsi **Kendala/Kejanggalan** yang dapat dipilih oleh **Expert/Validator** saat mengisi feedback untuk **Tes Profil Karier**. Master ini memastikan opsi kendala konsisten di UI, mudah dikelola, dan siap diagregasi untuk kebutuhan visualisasi admin.

### **2\. Tujuan**

1. Menstandarkan pilihan kendala expert (tidak bergantung pada teks bebas).  
2. Memudahkan agregasi statistik (COUNT/GROUP BY) untuk chart kendala.  
3. Memungkinkan penonaktifan opsi tanpa menghapus data historis.

### **3\. Struktur Kolom**

| Kolom | Tipe | Constraint | Deskripsi |
| ----- | ----- | ----- | ----- |
| `id` | SERIAL | PRIMARY KEY | ID unik opsi kendala. |
| `key` | VARCHAR(50) | UNIQUE, NOT NULL | Identifier stabil untuk sistem (mis. `TOO_GENERIC`). |
| `label` | TEXT | NOT NULL | Label yang ditampilkan pada UI. |
| `is_no_issue` | BOOLEAN | NOT NULL, DEFAULT FALSE | Penanda opsi “Tidak ada kendala (lancar)”. |
| `is_other` | BOOLEAN | NOT NULL, DEFAULT FALSE | Penanda opsi “Lainnya”. |
| `sort_order` | INT | NOT NULL | Urutan tampil di UI dan urutan default chart. |
| `is_active` | BOOLEAN | NOT NULL, DEFAULT TRUE | Jika false, opsi tidak boleh dipilih pada submission baru. |

### **4\. Seed Data (v1)**

1. `NO_ISSUE` — Tidak ada kendala (lancar) *(is\_no\_issue=true)*  
2. `IRRELEVANT_RECOMMENDATIONS` — Rekomendasi sangat tidak relevan  
3. `CONFUSING_TERMS` — Istilah/kalimat membingungkan  
4. `TOO_GENERIC` — Penjelasan terlalu umum  
5. `TECHNICAL_BUG` — Bug teknis (error/loading/tombol tidak berfungsi)  
6. `OTHER` — Lainnya *(is\_other=true)*

### **5\. Relasi**

Tabel ini direferensikan oleh tabel pivot `careerprofile_feedback_expert_obstacles` melalui kolom `obstacle_id`.

### **6\. Aturan Konsistensi**

1. Opsi dengan `is_active=false` tidak boleh digunakan untuk feedback baru (histori tetap ditampilkan apa adanya).  
2. Opsi `is_no_issue=true` diperlakukan eksklusif (detail aturan diterapkan pada tabel pivot).  
3. Opsi `is_other=true` membutuhkan detail teks yang disimpan pada tabel pivot.

## **5.3.4 Tabel `careerprofile_feedback_expert_obstacles` (Pivot Kendala Expert)**

### **1\. Definisi**

Tabel ini menyimpan **pilihan kendala** yang dipilih oleh expert pada suatu feedback. Karena kendala bersifat **multi-select**, tabel ini berperan sebagai pivot agar penyimpanan terstruktur dan agregasi untuk analitik tetap efisien.

### **2\. Tujuan**

1. Mencatat daftar kendala per feedback expert tanpa redundansi.  
2. Memudahkan perhitungan “Ragam Kendala” dan “Ada/Tidak Ada Kendala” secara cepat.  
3. Menyediakan tempat penyimpanan detail untuk opsi “Lainnya”.

### **3\. Struktur Kolom**

| Kolom | Tipe | Constraint | Deskripsi |
| ----- | ----- | ----- | ----- |
| `feedback_id` | BIGINT | NOT NULL, FK → `careerprofile_feedback_expert.feedback_id` ON DELETE CASCADE | Referensi ke feedback expert. |
| `obstacle_id` | INT | NOT NULL, FK → `careerprofile_expert_obstacle_options.id` | Opsi kendala yang dipilih. |
| `other_text` | TEXT | NULL | Diisi jika opsi yang dipilih adalah `OTHER`. |
| **PK** |  | PRIMARY KEY (`feedback_id`, `obstacle_id`) | Mencegah pilihan yang sama tersimpan ganda. |

### **4\. Index yang Disarankan**

1. `INDEX (obstacle_id)` untuk agregasi chart “Ragam Kendala Expert”.  
2. `INDEX (feedback_id)` untuk mempercepat load detail kendala pada drawer feedback.

### **5\. Relasi**

1. `careerprofile_feedback_expert` (1) → (N) `careerprofile_feedback_expert_obstacles`  
2. `careerprofile_expert_obstacle_options` (1) → (N) `careerprofile_feedback_expert_obstacles`

### **6\. Aturan Konsistensi (Validasi Wajib di Layer Aplikasi)**

1. **Eksklusif NO\_ISSUE**: jika memilih `NO_ISSUE`, maka tidak boleh memilih opsi kendala lain pada feedback yang sama.  
2. **OTHER wajib detail**: jika memilih `OTHER`, maka `other_text` wajib terisi dan bukan string kosong.  
3. **Opsi nonaktif**: sistem menolak pemilihan obstacle yang `is_active=false` untuk submission baru.

### **7\. Contoh Data (Ilustratif)**

* Feedback `feedback_id=12032` memilih `TOO_GENERIC` dan `TECHNICAL_BUG`:  
  * `(12032, obstacle_id=4, NULL)`  
  * `(12032, obstacle_id=5, NULL)`

* Jika memilih `OTHER`:  
  * `(12032, obstacle_id=6, 'Penjelasan belum menyebut alasan rekomendasi secara spesifik.')`

## **`🚀 IMPLEMENTATION SQL (Complete)`**

*`-- =============================================`*

*`-- COMPLETE MIGRATION - COPY-PASTE READY`*

*`-- =============================================`*

*`-- 1. Header Feedback (Polymorphic)`*

`CREATE TABLE kenalidiri_feedback (`  
  `id BIGSERIAL PRIMARY KEY,`  
  `test_category VARCHAR(50) NOT NULL,`  
  `test_session_id BIGINT NOT NULL,`  
  `respondent_type VARCHAR(20) NOT NULL,`  
  `respondent_user_id BIGINT NOT NULL REFERENCES users(id),`  
  `submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),`  
  `deleted_at TIMESTAMPTZ,`  
  `deleted_by BIGINT REFERENCES users(id)`  
`);`

`CREATE UNIQUE INDEX uq_feedback_once`  
`ON kenalidiri_feedback (test_category, test_session_id, respondent_type, respondent_user_id)`  
`WHERE deleted_at IS NULL;`

`CREATE INDEX idx_feedback_category_time ON kenalidiri_feedback(test_category, submitted_at);`  
`CREATE INDEX idx_feedback_session ON kenalidiri_feedback(test_category, test_session_id);`  
`CREATE INDEX idx_feedback_role_time ON kenalidiri_feedback(test_category, respondent_type, submitted_at);`

*`-- 2. Student Feedback Detail`*  
`CREATE TABLE careerprofile_feedback_student (`  
  `feedback_id BIGINT PRIMARY KEY REFERENCES kenalidiri_feedback(id) ON DELETE CASCADE,`  
  `ease_score SMALLINT NOT NULL CHECK (ease_score BETWEEN 1 AND 7),`  
  `relevance_score SMALLINT NOT NULL CHECK (relevance_score BETWEEN 1 AND 7),`  
  `satisfaction_score SMALLINT NOT NULL CHECK (satisfaction_score BETWEEN 1 AND 7),`  
  `message_to_team TEXT`  
`);`

*`-- 3. Student Obstacle Options (Master)`*  
`CREATE TABLE careerprofile_obstacle_options (`  
  `id SERIAL PRIMARY KEY,`  
  `key VARCHAR(50) UNIQUE NOT NULL,`  
  `label TEXT NOT NULL,`  
  `is_no_issue BOOLEAN NOT NULL DEFAULT FALSE,`  
  `is_other BOOLEAN NOT NULL DEFAULT FALSE,`  
  `sort_order INT NOT NULL,`  
  `is_active BOOLEAN NOT NULL DEFAULT TRUE`  
`);`

`INSERT INTO careerprofile_obstacle_options (key, label, is_no_issue, is_other, sort_order) VALUES`  
`('NO_ISSUE', 'Tidak ada kendala', TRUE, FALSE, 1),`  
`('DURATION_TOO_LONG', 'Durasi tes terasa terlalu lama', FALSE, FALSE, 2),`  
`('QUESTIONS_CONFUSING', 'Ada pertanyaan yang membingungkan', FALSE, FALSE, 3),`  
`('BUG_ERROR', 'Mengalami error/bug', FALSE, FALSE, 4),`  
`('RESULT_EXPLANATION_TOO_LONG', 'Penjelasan hasil tes terlalu panjang atau sulit dipahami', FALSE, FALSE, 5),`  
`('UI_NAV_CONFUSING', 'Tampilan atau navigasi membingungkan', FALSE, FALSE, 6),`  
`('OTHER', 'Lainnya', FALSE, TRUE, 7);`

*`-- 4. Student Obstacles (Pivot)`*  
`CREATE TABLE careerprofile_feedback_obstacles (`  
  `feedback_id BIGINT NOT NULL REFERENCES careerprofile_feedback_student(feedback_id) ON DELETE CASCADE,`  
  `obstacle_id INT NOT NULL REFERENCES careerprofile_obstacle_options(id),`  
  `other_text TEXT,`  
  `PRIMARY KEY (feedback_id, obstacle_id)`  
`);`

`CREATE INDEX idx_student_obstacles_agg ON careerprofile_feedback_obstacles(obstacle_id);`

*`-- 5. Expert Feedback Detail`*  
`CREATE TABLE careerprofile_feedback_expert (`  
  `feedback_id BIGINT PRIMARY KEY REFERENCES kenalidiri_feedback(id) ON DELETE CASCADE,`  
    
  `-- Scores`  
  `accuracy_score SMALLINT NOT NULL CHECK (accuracy_score BETWEEN 1 AND 7),`  
  `logic_score SMALLINT NOT NULL CHECK (logic_score BETWEEN 1 AND 7),`  
  `usefulness_score SMALLINT NOT NULL CHECK (usefulness_score BETWEEN 1 AND 7),`  
    
  `-- Expert Identity Snapshot`  
  `expert_name VARCHAR(255) NOT NULL,`  
  `expert_profession VARCHAR(120) NOT NULL,`  
  `expert_profession_id BIGINT REFERENCES professions(id),`  
  `expert_degree VARCHAR(120),`  
  `expert_experience_years SMALLINT,`  
  `expert_education_level VARCHAR(60),`  
  `expert_university VARCHAR(255),`  
  `expert_study_program VARCHAR(255),`  
    
  `-- Top 5 Snapshot`  
  `top5_recommendations_json JSONB NOT NULL,`  
    
  `-- GENERATED COLUMN (auto-calculated!)`  
  `top5_status VARCHAR(20) GENERATED ALWAYS AS (`  
    `CASE`   
      `WHEN expert_profession_id IS NULL THEN 'NOT_PRESENT'`  
      `WHEN jsonb_path_exists(`  
        `top5_recommendations_json,`   
        `'$[*] ? (@.rank == 1 && @.profession_id == $id)',`  
        `jsonb_build_object('id', expert_profession_id)`  
      `) THEN 'P1'`  
      `WHEN jsonb_path_exists(`  
        `top5_recommendations_json,`   
        `'$[*] ? (@.rank == 2 && @.profession_id == $id)',`  
        `jsonb_build_object('id', expert_profession_id)`  
      `) THEN 'P2'`  
      `WHEN jsonb_path_exists(`  
        `top5_recommendations_json,`   
        `'$[*] ? (@.rank >= 3 && @.rank <= 5 && @.profession_id == $id)',`  
        `jsonb_build_object('id', expert_profession_id)`  
      `) THEN 'P3_5'`  
      `ELSE 'NOT_PRESENT'`  
    `END`  
  `) STORED,`  
    
  `suggestion_text TEXT`  
`);`

`CREATE INDEX idx_expert_top5_status ON careerprofile_feedback_expert(top5_status);`  
`CREATE INDEX idx_expert_profession ON careerprofile_feedback_expert(expert_profession_id) WHERE expert_profession_id IS NOT NULL;`

*`-- 6. Expert Obstacle Options (Master)`*  
`CREATE TABLE careerprofile_expert_obstacle_options (`  
  `id SERIAL PRIMARY KEY,`  
  `key VARCHAR(50) UNIQUE NOT NULL,`  
  `label TEXT NOT NULL,`  
  `is_no_issue BOOLEAN NOT NULL DEFAULT FALSE,`  
  `is_other BOOLEAN NOT NULL DEFAULT FALSE,`  
  `sort_order INT NOT NULL,`  
  `is_active BOOLEAN NOT NULL DEFAULT TRUE`  
`);`

`INSERT INTO careerprofile_expert_obstacle_options (key, label, is_no_issue, is_other, sort_order) VALUES`  
`('NO_ISSUE', 'Tidak ada kendala (lancar)', TRUE, FALSE, 1),`  
`('IRRELEVANT_RECOMMENDATIONS', 'Rekomendasi sangat tidak relevan', FALSE, FALSE, 2),`  
`('CONFUSING_TERMS', 'Istilah/kalimat membingungkan', FALSE, FALSE, 3),`  
`('TOO_GENERIC', 'Penjelasan terlalu umum', FALSE, FALSE, 4),`  
`('TECHNICAL_BUG', 'Bug teknis (error/loading/tombol tidak berfungsi)', FALSE, FALSE, 5),`  
`('OTHER', 'Lainnya', FALSE, TRUE, 6);`

*`-- 7. Expert Obstacles (Pivot)`*  
`CREATE TABLE careerprofile_feedback_expert_obstacles (`  
  `feedback_id BIGINT NOT NULL REFERENCES careerprofile_feedback_expert(feedback_id) ON DELETE CASCADE,`  
  `obstacle_id INT NOT NULL REFERENCES careerprofile_expert_obstacle_options(id),`  
  `other_text TEXT,`  
  `PRIMARY KEY (feedback_id, obstacle_id)`  
`);`

`CREATE INDEX idx_expert_obstacles_agg ON careerprofile_feedback_expert_obstacles(obstacle_id);`  
`CREATE INDEX idx_expert_obstacles_feedback ON careerprofile_feedback_expert_obstacles(feedba`

