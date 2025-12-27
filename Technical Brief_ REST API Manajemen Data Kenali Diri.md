---
title: 'Technical Brief: REST API Manajemen Data Kenali Diri'

---

# Brief Pengembangan REST API - Manajemen Data Kenali Diri
## Backend Golang - Admin Dashboard

---

## Ringkasan Eksekutif

Dokumen ini merupakan technical brief untuk pengembangan REST API backend Golang yang mendukung dashboard admin untuk manajemen data fitur Kenali Diri. Sistem ini dirancang untuk mengelola data hasil tes kepribadian dan karier (RIASEC + Ikigai), umpan balik pengguna (mahasiswa dan expert validator), serta master data RIASEC codes.

**Catatan Penting:**
- Entity database (12 tables) sudah tersedia dan tidak perlu dibuat ulang
- Fokus pengembangan: DTO, Repository, Service, Controller, Routes, dan Config
- Arsitektur: Clean Architecture dengan pemisahan Admin dan User
- Framework: Gin (HTTP), GORM (ORM), PostgreSQL (Database)

---

## Struktur Folder Target

```
internal/api/kenali_diri/
├── controller/
│   ├── kenalidiri_admin_controller.go
│   └── kenalidiri_user_controller.go
├── service/
│   ├── kenalidiri_admin_service.go
│   └── kenalidiri_user_service.go
├── repository/
│   ├── kenalidiri_history_repository.go
│   ├── kenalidiri_category_repository.go
│   ├── riasec_code_repository.go
│   ├── test_session_repository.go
│   ├── riasec_repository.go
│   ├── ikigai_repository.go
│   ├── feedback_repository.go
│   └── recommendation_repository.go
└── routes/
    ├── kenalidiri_admin_route.go
    └── kenalidiri_user_route.go

internal/dto/
├── request/
│   ├── kenalidiri_admin_dto_request.go
│   └── feedback_dto_request.go
└── response/
    ├── kenalidiri_admin_dto_response.go
    ├── feedback_dto_response.go
    └── pagination_response.go

internal/pkg/
├── export/
│   ├── csv_exporter.go
│   ├── excel_exporter.go
│   └── pdf_exporter.go
└── cache/
    ├── redis_cache.go
    └── cache_service.go
```

---

## FASE 1: Komponen Foundation (Tidak Terikat UI Spesifik)

### 1.1. Data Transfer Objects (DTO)

#### File: `internal/dto/response/pagination_response.go`
**Tujuan:** Struktur response pagination yang digunakan oleh semua list endpoint

**Komponen:**
- Struct `PaginationMeta` dengan field:
  - `CurrentPage` (int) - halaman saat ini
  - `TotalPages` (int) - total halaman
  - `TotalRecords` (int64) - total data
  - `PerPage` (int) - jumlah data per halaman

**Penggunaan:** Semua endpoint list/tabel yang memerlukan pagination

---

#### File: `internal/dto/request/kenalidiri_admin_dto_request.go`
**Tujuan:** Request validation untuk semua admin endpoints

**Structs yang Diperlukan:**

1. **GetTestHistoryRequest**
   - Purpose: Filter dan pagination untuk list riwayat tes
   - Fields:
     - `Page` (int, optional, min=1)
     - `Limit` (int, optional, min=1, max=100)
     - `CategoryID` (*int64, optional)
     - `Status` (string, optional, oneof: "ongoing|completed|abandoned")
     - `UserName` (string, optional)
     - `StartDate` (string, optional, datetime format)
     - `EndDate` (string, optional, datetime format)
     - `SortBy` (string, optional, oneof: "name_asc|name_desc|date_asc|date_desc")

2. **ExportTestHistoryRequest**
   - Purpose: Request untuk ekspor data dengan filter
   - Fields:
     - `CategoryID` (*int64, optional)
     - `Status` (string, optional)
     - `StartDate` (string, optional)
     - `EndDate` (string, optional)
     - `Format` (string, required, oneof: "csv|excel|pdf")

3. **DeleteTestDataRequest**
   - Purpose: Bulk delete test data
   - Fields:
     - `TestIDs` ([]int64, required, min=1)

4. **UpdateRiasecCodeRequest**
   - Purpose: Update master data RIASEC code
   - Fields:
     - `RiasecTitle` (string, required)
     - `RiasecDescription` (string, required)
     - `Strengths` ([]string, required, min=1)
     - `Challenges` ([]string, required, min=1)
     - `Strategies` ([]string, required, min=1)
     - `WorkEnvironments` ([]string, required, min=1)
     - `InteractionStyles` ([]string, required, min=1)

5. **GetRiasecCodeListRequest**
   - Purpose: Filter untuk list RIASEC codes
   - Fields:
     - `CodeType` (string, optional, oneof: "single|dual|triple")
     - `Search` (string, optional)

---

#### File: `internal/dto/request/feedback_dto_request.go`
**Tujuan:** Request validation untuk feedback endpoints

**Structs yang Diperlukan:**

1. **GetFeedbackListRequest** (Mahasiswa)
   - Fields:
     - `Page` (int, optional, min=1)
     - `Limit` (int, optional, min=1, max=100)
     - `CategoryID` (*int64, optional)
     - `UserName` (string, optional)
     - `HasObstacles` (*bool, optional) - filter kendala
     - `SortBy` (string, optional)

2. **GetExpertFeedbackListRequest** (Expert)
   - Fields:
     - `Page` (int, optional)
     - `Limit` (int, optional)
     - `CategoryID` (*int64, optional)
     - `ExpertName` (string, optional)
     - `TopNStatus` (string, optional, oneof: "P1|P2|P3-5|not_found")
     - `SortBy` (string, optional)

---

#### File: `internal/dto/response/kenalidiri_admin_dto_response.go`
**Tujuan:** Response structures untuk admin endpoints

**Structs yang Diperlukan:**

1. **TestHistoryListResponse**
   - Fields:
     - `Data` ([]TestHistoryItem)
     - `Pagination` (PaginationMeta)

2. **TestHistoryItem**
   - Fields:
     - `TestID` (string) - format "PK{id}"
     - `UserName` (string)
     - `CategoryName` (string)
     - `Status` (string)
     - `ResultCode` (string) - RIASEC code jika completed
     - `StartedAt` (string) - formatted datetime
     - `CompletedAt` (*string) - optional

3. **TestDetailResponse**
   - Fields:
     - `TestID` (string)
     - `UserName` (string)
     - `CategoryName` (string)
     - `Status` (string)
     - `StartedAt` (string)
     - `CompletedAt` (*string)
     - `RiasecResult` (*RiasecResultDetail)
     - `IkigaiResult` (*IkigaiResultDetail)
     - `Recommendations` ([]RecommendationDetail)

4. **RiasecResultDetail**
   - Fields:
     - `ScoreR`, `ScoreI`, `ScoreA`, `ScoreS`, `ScoreE`, `ScoreC` (int)
     - `RiasecCode` (string)
     - `RiasecTitle` (string)
     - `ClassificationType` (string)
     - `IsInconsistent` (bool)

5. **IkigaiResultDetail**
   - Fields:
     - `LoveNarrative` (string)
     - `GoodAtNarrative` (string)
     - `WorldNeedsNarrative` (string)
     - `PaidForNarrative` (string)

6. **RecommendationDetail**
   - Fields:
     - `Rank` (int)
     - `ProfessionID` (int64)
     - `ProfessionName` (string)
     - `MatchPercentage` (int)
     - `MatchReasoning` (string)

7. **RiasecCodeListResponse**
   - Fields:
     - `Data` ([]RiasecCodeItem)
     - `TotalCodes` (int)

8. **RiasecCodeItem**
   - Fields:
     - `ID` (int64)
     - `Code` (string)
     - `Title` (string)
     - `CodeType` (string)

9. **RiasecCodeDetailResponse**
   - Fields:
     - `ID` (int64)
     - `Code` (string)
     - `Title` (string)
     - `Description` (string)
     - `Strengths` ([]string)
     - `Challenges` ([]string)
     - `Strategies` ([]string)
     - `WorkEnvironments` ([]string)
     - `InteractionStyles` ([]string)

10. **ExportFileResponse**
    - Fields:
      - `FileName` (string)
      - `FileURL` (string)
      - `ExportedAt` (string)
      - `TotalRecords` (int)

---

#### File: `internal/dto/response/feedback_dto_response.go`
**Tujuan:** Response structures untuk feedback endpoints

**Structs yang Diperlukan:**

1. **FeedbackListResponse** (Mahasiswa)
   - Fields:
     - `Data` ([]FeedbackItem)
     - `Pagination` (PaginationMeta)

2. **FeedbackItem** (Mahasiswa)
   - Fields:
     - `ID` (int64)
     - `UserName` (string)
     - `EaseOfUseScore` (int) - 1-7 Likert
     - `RelevanceScore` (int) - 1-7 Likert
     - `SatisfactionScore` (int) - 1-7 Likert
     - `Obstacles` ([]string) - list kendala
     - `SubmittedAt` (string)

3. **FeedbackStatsResponse** (Mahasiswa - untuk visualisasi)
   - Fields:
     - `TotalFeedback` (int)
     - `AvgEaseOfUse` (float64)
     - `AvgRelevance` (float64)
     - `AvgSatisfaction` (float64)
     - `ParticipationRate` (float64)
     - `TrendData` (map[string]interface{}) - untuk charts

4. **ExpertFeedbackListResponse**
   - Fields:
     - `Data` ([]ExpertFeedbackItem)
     - `Pagination` (PaginationMeta)

5. **ExpertFeedbackItem**
   - Fields:
     - `ID` (int64)
     - `ExpertName` (string)
     - `Profession` (string)
     - `TopNStatus` (string) - "P1"|"P2"|"P3-5"|"Not Found"
     - `AccuracyScore` (int) - 1-7 Likert
     - `LogicScore` (int) - 1-7 Likert
     - `BenefitScore` (int) - 1-7 Likert
     - `Obstacles` ([]string)
     - `SubmittedAt` (string)

6. **ExpertFeedbackDetailResponse**
   - Fields:
     - `ID` (int64)
     - `ExpertName` (string)
     - `Profession` (string)
     - `Degree` (string)
     - `Experience` (string)
     - `Education` (string)
     - `University` (string)
     - `StudyProgram` (string)
     - `CategoryTest` (string)
     - `TopFiveProfessions` ([]string)
     - `TopNStatus` (string)
     - `AccuracyScore` (int)
     - `LogicScore` (int)
     - `BenefitScore` (int)
     - `Obstacles` ([]string)
     - `Suggestions` (string)

---

### 1.2. Repository Layer

#### File: `internal/api/kenali_diri/repository/kenalidiri_history_repository.go`
**Tujuan:** Data access untuk tabel `kenalidiri_history`

**Interface Methods:**
1. `Create(ctx, tx, history) (entity.KenalidiriHistory, error)`
   - Insert new test history record
   
2. `GetByID(ctx, tx, id) (entity.KenalidiriHistory, error)`
   - Retrieve single history by ID dengan preload User dan TestCategory
   
3. `UpdateStatus(ctx, tx, id, status, completedAt) error`
   - Update status tes (ongoing/completed/abandoned)
   
4. `ListWithFilters(ctx, tx, filters) ([]entity.KenalidiriHistory, int64, error)`
   - List dengan filter multi-variabel dan pagination
   - Support: category_id, status, user_name, date range, sort
   - Return: data array + total count
   - Preload: User, TestCategory
   
5. `GetUserHistory(ctx, tx, userID, categoryID) ([]entity.KenalidiriHistory, error)`
   - Ambil riwayat tes user tertentu (untuk user dashboard)
   
6. `BulkDelete(ctx, tx, ids) error`
   - Hapus multiple test records sekaligus
   
7. `CountByStatus(ctx, tx, categoryID) (map[string]int64, error)`
   - Aggregate count per status untuk dashboard statistics

**Query Optimization Notes:**
- Index pada: user_id, test_category_id, status, started_at
- Join dengan users table untuk filter by name
- Use ILIKE untuk case-insensitive search

---

#### File: `internal/api/kenali_diri/repository/kenalidiri_category_repository.go`
**Tujuan:** Data access untuk tabel `kenalidiri_categories`

**Interface Methods:**
1. `GetAll(ctx, tx) ([]entity.KenalidiriCategory, error)`
   - Retrieve all active categories
   - Filter: is_active = true
   - Cacheable (data jarang berubah)
   
2. `GetByID(ctx, tx, id) (entity.KenalidiriCategory, error)`
   - Retrieve single category by ID
   
3. `GetByCode(ctx, tx, code) (entity.KenalidiriCategory, error)`
   - Retrieve category by category_code

**Caching Strategy:**
- Cache key: "kenalidiri_categories:all"
- TTL: 24 hours
- Invalidate: saat ada update (jarang terjadi)

---

#### File: `internal/api/kenali_diri/repository/riasec_code_repository.go`
**Tujuan:** Data access untuk tabel `riasec_codes` (156 codes)

**Interface Methods:**
1. `GetAll(ctx, tx) ([]entity.RiasecCode, error)`
   - Retrieve all 156 RIASEC codes
   - **MUST be cached** (data statis)
   
2. `GetByID(ctx, tx, id) (entity.RiasecCode, error)`
   - Retrieve single code by ID
   - Cacheable
   
3. `GetByCode(ctx, tx, code) (entity.RiasecCode, error)`
   - Retrieve by riasec_code string (e.g., "RIA")
   - Cacheable
   
4. `ListByType(ctx, tx, codeType) ([]entity.RiasecCode, error)`
   - Filter by: "single" (6 codes), "dual" (30 codes), "triple" (120 codes)
   
5. `Search(ctx, tx, keyword) ([]entity.RiasecCode, error)`
   - Search by title or code
   - Use ILIKE untuk case-insensitive
   
6. `Update(ctx, tx, code) error`
   - Update RIASEC code master data
   - **Must invalidate cache** after update

**Caching Strategy - CRITICAL:**
- Cache key patterns:
  - All codes: "riasec_codes:all"
  - By ID: "riasec_code:{id}"
  - By code: "riasec_code_by_name:{code}"
- TTL: 24 hours (atau permanent dengan manual invalidation)
- Total cached data: ~156 records (ringan)
- Invalidate on: Update operation

**Performance Impact:**
- Without cache: 156 queries untuk setiap candidate generation
- With cache: 0 database queries (100% cache hit)
- Improvement: ~500ms → <10ms per request

---

#### File: `internal/api/kenali_diri/repository/test_session_repository.go`
**Tujuan:** Data access untuk `careerprofile_test_sessions`

**Interface Methods:**
1. `Create(ctx, tx, session) (entity.CareerprofileTestSession, error)`
   - Insert new test session dengan UUID token
   
2. `GetByID(ctx, tx, id) (entity.CareerprofileTestSession, error)`
   - Retrieve session by ID
   
3. `GetByToken(ctx, tx, token) (entity.CareerprofileTestSession, error)`
   - Retrieve session by session_token (untuk validasi)
   
4. `UpdateStatus(ctx, tx, id, status, timestamps) error`
   - Update status: riasec_ongoing → riasec_completed → ikigai_ongoing → completed
   
5. `GetUserSessions(ctx, tx, userID) ([]entity.CareerprofileTestSession, error)`
   - Ambil semua session milik user

---

#### File: `internal/api/kenali_diri/repository/riasec_repository.go`
**Tujuan:** Data access untuk RIASEC-related tables (question_sets, responses, results)

**Interface Methods:**

1. `SaveQuestionSet(ctx, tx, questionSet) error`
   - Insert riasec_question_sets (12 question IDs)
   
2. `GetQuestionSet(ctx, tx, sessionID) (entity.RiasecQuestionSets, error)`
   - Retrieve question set by test_session_id
   
3. `SaveResponses(ctx, tx, responses) error`
   - Insert riasec_responses (all 12 answers in JSONB)
   
4. `GetResponses(ctx, tx, sessionID) (entity.RiasecResponses, error)`
   - Retrieve user responses
   
5. `SaveResult(ctx, tx, result) error`
   - Insert riasec_results (scores + classified code)
   
6. `GetResultBySessionID(ctx, tx, sessionID) (entity.RiasecResult, error)`
   - Retrieve RIASEC result dengan preload RiasecCode
   
7. `GetResultByID(ctx, tx, id) (entity.RiasecResult, error)`
   - Retrieve result by primary key

**Join Strategy:**
- Preload RiasecCode relation untuk mendapatkan title dan description
- Eager loading untuk menghindari N+1 queries

---

#### File: `internal/api/kenali_diri/repository/ikigai_repository.go`
**Tujuan:** Data access untuk Ikigai-related tables (candidates, responses, scores)

**Interface Methods:**

1. `SaveCandidates(ctx, tx, candidates) error`
   - Insert ikigai_candidate_professions (5-30 professions)
   
2. `GetCandidates(ctx, tx, sessionID) (entity.IkigaiCandidateProfessions, error)`
   - Retrieve candidate professions
   
3. `SaveResponses(ctx, tx, responses) error`
   - Insert ikigai_responses (4 dimensions)
   
4. `GetResponses(ctx, tx, sessionID) (entity.IkigaiResponses, error)`
   - Retrieve user responses untuk 4 dimensi
   
5. `SaveDimensionScores(ctx, tx, scores) error`
   - Insert ikigai_dimension_scores (AI scoring results)
   
6. `SaveTotalScores(ctx, tx, scores) error`
   - Insert ikigai_total_scores (aggregated scores)
   
7. `GetTotalScores(ctx, tx, sessionID) (entity.IkigaiTotalScores, error)`
   - Retrieve final scores dengan ranking

---

#### File: `internal/api/kenali_diri/repository/recommendation_repository.go`
**Tujuan:** Data access untuk `career_recommendations`

**Interface Methods:**

1. `Save(ctx, tx, recommendation) error`
   - Insert career recommendation (AI-generated narratives)
   
2. `GetBySessionID(ctx, tx, sessionID) (entity.CareerRecommendation, error)`
   - Retrieve recommendation result
   
3. `GetByID(ctx, tx, id) (entity.CareerRecommendation, error)`
   - Retrieve by primary key

---

#### File: `internal/api/kenali_diri/repository/feedback_repository.go`
**Tujuan:** Data access untuk feedback tables (mahasiswa & expert)

**Interface Methods:**

1. `ListStudentFeedback(ctx, tx, filters) ([]entity.StudentFeedback, int64, error)`
   - List feedback mahasiswa dengan filter dan pagination
   
2. `GetStudentFeedbackByID(ctx, tx, id) (entity.StudentFeedback, error)`
   - Retrieve single student feedback
   
3. `GetStudentFeedbackStats(ctx, tx, filters) (map[string]interface{}, error)`
   - Aggregate statistics untuk visualisasi (Chart 1-8)
   - Return: total feedback, avg scores, trend data, obstacle distribution
   
4. `ListExpertFeedback(ctx, tx, filters) ([]entity.ExpertFeedback, int64, error)`
   - List feedback expert dengan filter dan pagination
   
5. `GetExpertFeedbackByID(ctx, tx, id) (entity.ExpertFeedback, error)`
   - Retrieve single expert feedback dengan detail lengkap

**Aggregation Queries:**
- AVG(ease_of_use_score), AVG(relevance_score), AVG(satisfaction_score)
- COUNT by date range untuk trend charts
- GROUP BY obstacles untuk distribution chart
- Complex JSON field queries untuk JSONB columns

---

### 1.3. Export Services

#### File: `internal/pkg/export/csv_exporter.go`
**Tujuan:** Generate CSV files

**Methods:**
- `Generate(data []map[string]interface{}, filename string) (string, error)`
- Input: array of key-value maps
- Output: file path atau S3 URL
- Library: `encoding/csv` (standard)

---

#### File: `internal/pkg/export/excel_exporter.go`
**Tujuan:** Generate Excel files dengan formatting

**Methods:**
- `Generate(data []map[string]interface{}, filename, sheetName string) (string, error)`
- `GenerateMultiSheet(sheets map[string][]map[string]interface{}, filename) (string, error)`
- Features: Headers, auto-width columns, cell styling
- Library: `github.com/xuri/excelize/v2`

---

#### File: `internal/pkg/export/pdf_exporter.go`
**Tujuan:** Generate PDF documents

**Methods:**
- `Generate(data []map[string]interface{}, filename, title string) (string, error)`
- Features: Table layout, headers, pagination
- Library: `github.com/jung-kurt/gofpdf`

---

#### File: `internal/pkg/export/export_service.go`
**Tujuan:** Unified export service interface

**Methods:**
- `GenerateCSV(...) (string, error)`
- `GenerateExcel(...) (string, error)`
- `GeneratePDF(...) (string, error)`

**Workflow:**
1. Generate file locally di `/tmp`
2. Upload ke AWS S3
3. Return public URL
4. Cleanup local file

---

### 1.4. Cache Service

#### File: `internal/pkg/cache/cache_service.go`
**Tujuan:** Abstraction layer untuk caching

**Interface Methods:**
- `Get(ctx, key) (interface{}, error)`
- `Set(ctx, key, value, expiration) error`
- `Delete(ctx, key) error`
- `Clear(ctx, pattern) error`

---

#### File: `internal/pkg/cache/redis_cache.go`
**Tujuan:** Redis implementation

**Features:**
- Connection pooling
- JSON serialization
- TTL support
- Pattern-based deletion

**Library:** `github.com/go-redis/redis/v8`

---

## FASE 2: Service Layer Per UI Screend
Semua 13 UI screen ditangani dalam 1 file service: kenalidiri_admin_service.go dengan 11 methods.
Alasan:

Semua screen dalam 1 domain: Manajemen Data Kenali Diri
Lebih mudah maintain (tidak scattered)
Share dependencies yang sama (repositories)
Follow single responsibility: admin service untuk admin dashboard


### UI SCREEN 1: Halaman Utama Hasil Tes Kenali Diri - Tabel Data

**Screenshot Reference:** Tabel riwayat tes dengan filter dan pagination

#### File: `internal/api/kenali_diri/service/kenalidiri_admin_service.go`

**Method yang Diperlukan:**

1. **GetTestHistory(ctx, req GetTestHistoryRequest) (TestHistoryListResponse, error)**
   
   **Use Cases:**
   - Menampilkan tabel riwayat hasil tes dengan pagination (default 10 data)
   - Filter multi-variabel: kategori tes, status tes, date picker, nama pengguna
   - Search bar: cari berdasarkan nama pengguna
   - Sort: nama A-Z/Z-A, tanggal mulai/selesai
   
   **Business Logic:**
   - Set default pagination (page=1, limit=10)
   - Parse date filters dari string ke time.Time
   - Build repository filters dari request params
   - Call `historyRepo.ListWithFilters()` dengan filters
   - Transform entity → DTO response
   - Untuk status "completed": load RIASEC code dari `riasec_results`
   - Calculate pagination meta (total pages, current page)
   
   **Data Transformation:**
   - ID Tes: format "PK{id}"
   - Status: render sebagai label dengan warna (ongoing=orange, completed=green, abandoned=red)
   - Hasil Tes: tampilkan RIASEC code jika completed, "-" jika belum
   - Waktu: format "02 Jan 2006 15:04"
   
   **Dependencies:**
   - `kenalidiriHistoryRepository`
   - `riasecResultRepository`

---

2. **DeleteTestData(ctx, req DeleteTestDataRequest) error**
   
   **Use Cases:**
   - Bulk delete: hapus multiple test data sekaligus
   - Single delete: hapus 1 test data
   - Validasi: cek apakah test IDs exist
   
   **Business Logic:**
   - Start database transaction
   - Validate test IDs exist
   - Delete dari `kenalidiri_history` (cascade delete ke child tables)
   - Commit transaction
   
   **Transaction Flow:**
   ```
   BEGIN TRANSACTION
   - Validate test IDs exist
   - Delete related riasec_question_sets
   - Delete related riasec_responses
   - Delete related riasec_results
   - Delete related ikigai_candidate_professions
   - Delete related ikigai_responses
   - Delete related ikigai_dimension_scores (optional - bisa dipertahankan untuk audit)
   - Delete related ikigai_total_scores
   - Delete related career_recommendations
   - Delete kenalidiri_history records
   COMMIT TRANSACTION
   ```
   
   **Dependencies:**
   - `kenalidiriHistoryRepository`
   - Database transaction support

---

3. **ExportTestHistory(ctx, req ExportTestHistoryRequest) (ExportFileResponse, error)**
   
   **Use Cases:**
   - Export data riwayat tes ke format CSV/Excel/PDF
   - Filter data sebelum export (kategori, status, date range)
   - Multiple sheets untuk Excel (per kategori tes)
   
   **Business Logic:**
   - Parse filters dari request
   - Call `historyRepo.ListWithFilters()` tanpa pagination limit (ambil semua data yang match filter)
   - Transform data ke format map[string]interface{} untuk exporter
   - Call export service sesuai format request
   - Generate filename dengan timestamp: "riwayat_tes_20251223_150405"
   - Upload to S3 dan return URL
   
   **Data Format untuk Export:**
   ```
   Columns:
   - ID Tes (PK{id})
   - Nama Pengguna
   - Email
   - Kategori Tes
   - Status
   - Waktu Mulai (format: DD/MM/YYYY HH:MM)
   - Waktu Selesai (format: DD/MM/YYYY HH:MM)
   - Kode RIASEC (jika completed)
   - Durasi (dalam menit)
   ```
   
   **Dependencies:**
   - `kenalidiriHistoryRepository`
   - `riasecResultRepository`
   - `exportService`

---

### UI SCREEN 2-4: Detail Hasil Tes - Tab RIASEC, IKIGAI, REKOMENDASI

**Screenshot Reference:** Modal fullscreen dengan 3 tab menu

#### File: `internal/api/kenali_diri/service/kenalidiri_admin_service.go`

**Method yang Diperlukan:**

4. **GetTestDetail(ctx, sessionID int64) (TestDetailResponse, error)**
   
   **Use Cases:**
   - Menampilkan detail lengkap hasil tes dalam modal
   - Tab 1 (RIASEC): Tabel skor, ranking, classification type
   - Tab 2 (IKIGAI): 4 card dimensi dengan narrative
   - Tab 3 (REKOMENDASI): 2 card top professions dengan match percentage
   
   **Business Logic:**
   - Get kenalidiri_history by detail_session_id
   - Load RIASEC result dengan preload RiasecCode entity
   - Load ikigai_total_scores
   - Load career_recommendations
   - Transform semua data ke nested DTO structure
   
   **RIASEC Result Processing:**
   - Ambil 6 scores (R, I, A, S, E, C)
   - Rank scores descending (Rank 1 = gold, Rank 2 = silver, Rank 3 = bronze)
   - Get classification type: "single", "dual", "triple"
   - Get classification narrative berdasarkan type:
     - Single: "Tipe {X} sangat dominan. Anda adalah seorang '{archetype}'."
     - Dual: "Anda memiliki profil Seimbang (Dual) antara {X} dan {Y}."
     - Triple: "Anda memiliki profil Kompleks (Triple), fleksibilitas tinggi dalam peran teknis, analitis, dan kreatif."
   
   **IKIGAI Result Processing:**
   - Parse recommendations_data JSON
   - Extract 4 dimension narratives:
     - love_narrative (❤️ Love - Minat)
     - good_at_narrative (⭐ Good At - Keahlian)
     - world_needs_narrative (🌍 Needs - Kebutuhan Pasar)
     - paid_for_narrative (💰 Paid - Nilai Ekonomi)
   
   **Recommendation Processing:**
   - Get top 2 professions dari ikigai_total_scores
   - Get recommendation narratives dari career_recommendations
   - Format match_percentage (0-100)
   - Determine card styling:
     - Rank 1: Background putih, stroke hijau tebal, badge "Best Match" hijau
     - Rank 2: Background putih, stroke abu-abu, badge "Alternative" abu
   
   **Dependencies:**
   - `kenalidiriHistoryRepository`
   - `testSessionRepository`
   - `riasecResultRepository`
   - `ikigaiTotalScoresRepository`
   - `recommendationRepository`

---

### UI SCREEN 5: Popup Ekspor Data

**Screenshot Reference:** Popup dengan dropdown filter dan tombol unduh

**Note:** Use case sudah covered di method `ExportTestHistory()` (lihat di atas)

**Additional UI Logic:**
- Dropdown kategori tes (multi-select atau single)
- Dropdown status tes (multi-select)
- Month picker untuk rentang waktu
- Toggle "Ekspor seluruh data" (bypass date filter)
- Tombol "Unduh Data" trigger `ExportTestHistory()`

---

### UI SCREEN 6-7: Popup Hapus Data (Satuan & Massal)

**Screenshot Reference:** 
- Satuan: Konfirmasi hapus 1 data dengan ID dan nama
- Massal: Konfirmasi hapus N data dengan checkbox verifikasi

**Note:** Use case sudah covered di method `DeleteTestData()` (lihat di atas)

**Additional UI Logic:**
- Satuan: Tampilkan "Apakah Anda yakin akan menghapus data tes dengan ID #{ID} milik {nama}?"
- Massal: 
  - Tampilkan "Apakah Anda yakin akan menghapus sebanyak {N} data tes Kenali Diri?"
  - Checkbox: "Saya mengerti bahwa data ini akan dihapus permanen dan tidak dapat dikembalikan."
  - Disable tombol "Iya, Hapus Data" jika checkbox belum dicentang

---

### UI SCREEN 8: Umpan Balik Mahasiswa - Tabel Data

**Screenshot Reference:** Tabel feedback dengan filter dan search

#### File: `internal/api/kenali_diri/service/kenalidiri_admin_service.go`

**Method yang Diperlukan:**

5. **GetStudentFeedbackList(ctx, req GetFeedbackListRequest) (FeedbackListResponse, error)**
   
   **Use Cases:**
   - Menampilkan tabel feedback mahasiswa dengan pagination
   - Filter: kategori tes, ada/tidak ada kendala
   - Search: nama pengguna
   - Sort: nama A-Z/Z-A, feedback terbaru/terlama
   
   **Business Logic:**
   - Set default pagination
   - Build repository filters
   - Call `feedbackRepo.ListStudentFeedback()` dengan filters
   - Transform entity → DTO
   - Format Likert scores (X/7)
   - Parse obstacles dari JSONB → array of strings
   
   **Data Columns:**
   - ID Feedback
   - Nama Pengguna
   - Kemudahan Tes (1-7, display: "X/7" dengan badge)
   - Relevansi Rekomendasi (1-7, display: "X/7" dengan badge)
   - Kepuasan Fitur (1-7, display: "X/7" dengan badge)
   - Daftar Kendala (display: chips/labels, max 2-3 visible + "+N" tooltip)
   - Tanggal Submit
   
   **Obstacles Display Logic:**
   - Jika <= 3 kendala: tampilkan semua sebagai chips
   - Jika > 3 kendala: tampilkan 3 pertama + "+{N}" chip dengan tooltip all obstacles
   
   **Dependencies:**
   - `feedbackRepository`

---

### UI SCREEN 9: Umpan Balik Mahasiswa - Visualisasi (8 Charts)

**Screenshot Reference:** Dashboard dengan overview cards + 8 charts

#### File: `internal/api/kenali_diri/service/kenalidiri_admin_service.go`

**Method yang Diperlukan:**

6. **GetStudentFeedbackStats(ctx, categoryID *int64, timeRange string) (FeedbackStatsResponse, error)**
   
   **Use Cases:**
   - Overview Cards (4 cards):
     - Total Feedback
     - Rata-rata Kemudahan Tes (1-7)
     - Rata-rata Relevansi Rekomendasi (1-7)
     - Rata-rata Kepuasan Keseluruhan (1-7)
   - Chart 1: Tren Peserta Tes vs Pengisi Feedback (Line Chart)
   - Chart 2: Rasio Pengisi Feedback vs Tidak Mengisi (Pie Chart)
   - Chart 3: Tingkat Kemudahan Tes (Bar Chart Vertical, 1-7 scale)
   - Chart 4: Tingkat Relevansi Rekomendasi (Bar Chart Vertical, 1-7 scale)
   - Chart 5: Tingkat Kepuasan Keseluruhan (Bar Chart Vertical, 1-7 scale)
   - Chart 6: Tingkat Kendala (Pie Chart: Ada vs Tidak Ada)
   - Chart 7: Ragam Kendala (Horizontal Bar Chart)
   - Chart 8: Komposisi Positif-Netral-Negatif (100% Stacked Bar)
   
   **Business Logic:**
   
   **Time Range Aggregation:**
   - "mingguan": per hari (Sen-Min)
   - "bulanan": per minggu (Week 1-4)
   - "sepanjang_waktu": per bulan (Jan-Des)
   
   **Overview Cards Calculation:**
   - Total Feedback: `COUNT(*)`
   - Avg Kemudahan: `AVG(ease_of_use_score)`
   - Avg Relevansi: `AVG(relevance_score)`
   - Avg Kepuasan: `AVG(satisfaction_score)`
   - Trend comparison: calculate Δ% vs previous period
   
   **Chart 1 Data:**
   ```sql
   -- Jumlah peserta tes per bucket waktu
   SELECT date_bucket, COUNT(DISTINCT test_session_id) as test_count
   FROM kenalidiri_history
   GROUP BY date_bucket
   
   -- Jumlah pengisi feedback per bucket waktu
   SELECT date_bucket, COUNT(*) as feedback_count
   FROM student_feedback
   GROUP BY date_bucket
   ```
   Return: `map[string]interface{}` dengan keys: "dates", "test_counts", "feedback_counts"
   
   **Chart 2 Data:**
   - Pengisi Feedback = jumlah feedback records
   - Tidak Mengisi = jumlah tes completed - jumlah feedback
   - Response Rate % = (feedback / completed_tests) * 100
   
   **Chart 3, 4, 5 Data:**
   ```sql
   -- Distribution per score (1-7)
   SELECT score_value, COUNT(*) as count
   FROM (SELECT ease_of_use_score as score_value FROM student_feedback) scores
   GROUP BY score_value
   ORDER BY score_value
   ```
   Return: array of {score: 1-7, count: N, percentage: X%}
   
   **Chart 6 Data:**
   - Ada Kendala = COUNT where obstacles array is not empty
   - Tidak Ada Kendala = COUNT where obstacles array is empty OR contains "Tidak ada kendala"
   
   **Chart 7 Data:**
   ```sql
   -- Flatten JSONB array and count occurrences
   SELECT obstacle, COUNT(*) as count
   FROM student_feedback,
   jsonb_array_elements_text(obstacles) as obstacle
   WHERE obstacle != 'Tidak ada kendala'
   GROUP BY obstacle
   ORDER BY count DESC
   ```
   Return: array of {obstacle_name: string, count: N, percentage: X%}
   
   **Chart 8 Data:**
   - Classify each metric score into: Negatif (1-3), Netral (4), Positif (5-7)
   - Calculate percentage distribution for Kemudahan, Relevansi, Kepuasan
   - Return: 3 bars (one per metric) dengan 3 segments each (neg/neu/pos)
   
   **Wawasan Grafik (Automatic Insights):**
   - Based on data patterns, generate textual insights
   - Templates provided in Product Knowledge Doc
   - Logic rules:
     - Dominan positif: pos% >= 70%
     - Dominan negatif: neg% >= 55%
     - Terpolarisasi: pos% >= 35% AND neg% >= 35%
     - etc.
   
   **Dependencies:**
   - `feedbackRepository` with aggregation methods
   - `kenalidiriHistoryRepository` for participant counts

---

### UI SCREEN 10: Umpan Balik Expert - Tabel Data

**Screenshot Reference:** Tabel feedback expert dengan kolom spesifik

#### File: `internal/api/kenali_diri/service/kenalidiri_admin_service.go`

**Method yang Diperlukan:**

7. **GetExpertFeedbackList(ctx, req GetExpertFeedbackListRequest) (ExpertFeedbackListResponse, error)**
   
   **Use Cases:**
   - Menampilkan tabel feedback expert validator
   - Filter: kategori tes, nama expert, status Top-N
   - Search: nama expert
   - Sort: nama A-Z/Z-A, feedback terbaru/terlama
   
   **Business Logic:**
   - Set default pagination
   - Build repository filters
   - Call `feedbackRepo.ListExpertFeedback()` dengan filters
   - Transform entity → DTO
   - Calculate Top-N Status dari top_five_professions field
   
   **Top-N Status Calculation:**
   - Extract expert's actual profession dari form
   - Compare dengan top_five_professions array
   - Status:
     - "P1": profession di rank 1
     - "P2": profession di rank 2
     - "P3-5": profession di rank 3-5
     - "Tidak muncul": profession tidak ada di top 5
   
   **Data Columns:**
   - ID Feedback
   - Nama Expert
   - Profesi
   - Top 5 Status (badge: P1 gold, P2 silver, P3-5 bronze, Not Found gray)
   - Akurasi (1-7, display: "X/7")
   - Logika (1-7, display: "X/7")
   - Manfaat (1-7, display: "X/7")
   - Kendala (chips, max 2-3 + tooltip)
   - Tanggal Submit
   - Aksi: Tombol "Detail"
   
   **Dependencies:**
   - `feedbackRepository`

---

### UI SCREEN 11: Umpan Balik Expert - Modal Detail

**Screenshot Reference:** Drawer panel kanan dengan 5 sections

#### File: `internal/api/kenali_diri/service/kenalidiri_admin_service.go`

**Method yang Diperlukan:**

8. **GetExpertFeedbackDetail(ctx, feedbackID int64) (ExpertFeedbackDetailResponse, error)**
   
   **Use Cases:**
   - Drawer panel tampilkan detail lengkap expert feedback
   - Section 1: Identitas Responden (grid 2 kolom)
   - Section 2: 5 Rekomendasi Profesi Teratas (list dengan status badge)
   - Section 3: Penilaian Likert 1-7 (3 cards: Akurasi, Logika, Manfaat)
   - Section 4: Kendala yang Dilaporkan (chips)
   - Section 5: Masukan/Saran Perbaikan (text read-only)
   
   **Business Logic:**
   - Get expert feedback by ID
   - Load related test session data untuk mendapatkan rekomendasi
   - Parse JSONB fields: top_five_professions, obstacles
   - Format all data untuk UI consumption
   
   **Section 1 Data:**
   - ID Feedback (format: #FBX-{id})
   - Tanggal (format: DD MMM YYYY)
   - Nama Expert
   - Profesi
   - Gelar (full: "Sarjana Desain (S.Ds.)")
   - Pengalaman (format: "3 tahun" atau "3-5 tahun")
   - Pendidikan Terakhir ("Sarjana (S1)")
   - Perguruan Tinggi
   - Program Studi (colspan 2)
   - Kategori Tes (badge/label, colspan 2)
   
   **Section 2 Data:**
   - Header: judul + Top-N Status badge (kanan atas)
   - List 5 professions dengan badge angka (1-5)
   - Styling: background ringan + border tipis per item
   
   **Section 3 Data:**
   - 3 mini cards (grid 3 columns):
     - Card 1: Akurasi Profil (badge X/7)
     - Card 2: Logika Penjelasan (badge X/7)
     - Card 3: Potensi Manfaat (badge X/7)
   - Optional: icon per card untuk visual enhancement
   
   **Section 4 Data:**
   - Jika empty: tampilkan "Tidak ada kendala dilaporkan" (italic/muted)
   - Jika ada: tampilkan chips per obstacle
   - Jika "Lainnya": tampilkan chip "Lainnya" + detail text di bawah atau tooltip
   
   **Section 5 Data:**
   - Teks suggestions dalam box read-only
   - Whitespace preserved (pre-wrap)
   - Jika empty: "Tidak ada masukan tertulis."
   
   **Dependencies:**
   - `feedbackRepository`

---

### UI SCREEN 12: Master RIASEC - Daftar Item (156 Codes)

**Screenshot Reference:** Card view dengan 156 kode RIASEC

#### File: `internal/api/kenali_diri/service/kenalidiri_admin_service.go`

**Method yang Diperlukan:**

9. **GetRiasecCodeList(ctx, req GetRiasecCodeListRequest) (RiasecCodeListResponse, error)**
   
   **Use Cases:**
   - Menampilkan semua 156 kode RIASEC dalam card view
   - Filter: button group (Semua, 1 Huruf, 2 Huruf, 3 Huruf)
   - Search bar: cari berdasarkan kode atau nama
   
   **Business Logic:**
   - Check cache first: "riasec_codes:all"
   - If cache miss: load dari database + cache for 24h
   - Filter by code_type if specified ("single", "dual", "triple")
   - If search keyword: filter by code ILIKE or title ILIKE
   - Transform entity → DTO for card display
   
   **Card Display Data:**
   - Kotak Huruf RIASEC: jumlah kotak = panjang code (R=1, RI=2, RIA=3)
   - Setiap huruf diberi warna sesuai tipe:
     - R (Realistic): biru
     - I (Investigative): hijau
     - A (Artistic): ungu
     - S (Social): kuning
     - E (Enterprising): orange
     - C (Conventional): merah
   - Nama Kode: riasec_title (e.g., "Realistic-Investigative")
   - Tombol Detail: navigasi ke edit page
   - Label ID: tampilkan di pojok kanan bawah (e.g., "#45")
   
   **Caching Strategy - CRITICAL:**
   - All 156 codes MUST be cached
   - Cache key: "riasec_codes:all"
   - TTL: 24 hours atau permanent
   - Reason: data statis, sering diakses (setiap test generation)
   
   **Dependencies:**
   - `riasecCodeRepository` with caching
   - `cacheService`

---

### UI SCREEN 13: Master RIASEC - Detail & Edit

**Screenshot Reference:** Form edit dengan 5 sections

#### File: `internal/api/kenali_diri/service/kenalidiri_admin_service.go`

**Method yang Diperlukan:**

10. **GetRiasecCodeDetail(ctx, codeID int64) (RiasecCodeDetailResponse, error)**
    
    **Use Cases:**
    - Load data RIASEC code untuk edit form
    - Display 5 sections dengan data lengkap
    
    **Business Logic:**
    - Try cache first: "riasec_code:{id}"
    - If cache miss: load from database + cache
    - Transform entity → DTO
    - Parse JSONB arrays (strengths, challenges, etc.)
    
    **Section Data:**
    - Section 1 (Informasi Dasar):
      - Kode: read-only
      - Judul: editable text input
      - Deskripsi: editable textarea (supports markdown)
    - Section 2 (Kekuatan Profil): array of strings
    - Section 3 (Tantangan Profil): array of strings
    - Section 4 (Strategi Pengembangan Diri): array of strings
    - Section 5 (Lingkungan Kerja Ideal): array of strings
    
    **Dependencies:**
    - `riasecCodeRepository` with caching

---

11. **UpdateRiasecCode(ctx, codeID int64, req UpdateRiasecCodeRequest) error**
    
    **Use Cases:**
    - Update master data RIASEC code
    - Validasi input dari form
    - Invalidate cache setelah update
    
    **Business Logic:**
    - Validate codeID exists
    - Update entity fields dari request
    - Save to database
    - **CRITICAL: Invalidate cache**
      - Delete key: "riasec_code:{id}"
      - Delete key: "riasec_code_by_name:{code}"
      - Delete key: "riasec_codes:all"
    
    **Validation:**
    - riasec_title: required, max 255 chars
    - riasec_description: required, markdown format
    - strengths, challenges, strategies, work_environments, interaction_styles: 
      - required arrays
      - min 1 item per array
      - each item max 500 chars
    
    **Cache Invalidation - CRITICAL:**
    - Must delete all related cache keys
    - If not invalidated: users akan melihat data lama
    - Impact: candidate generation akan gunakan data outdated
    
    **Dependencies:**
    - `riasecCodeRepository`
    - `cacheService`

---

## FASE 3: Controller Layer Per UI Screen

### File: `internal/api/kenali_diri/controller/kenalidiri_admin_controller.go`

**Purpose:** HTTP request handlers untuk semua admin endpoints

**Interface Definition:**
```go
type KenalidiriAdminController interface {
    // Test History Management (UI Screen 1)
    GetTestHistory(ctx *gin.Context)
    GetTestDetail(ctx *gin.Context)
    DeleteTestData(ctx *gin.Context)
    ExportTestHistory(ctx *gin.Context)
    
    // Feedback Management (UI Screen 8-11)
    GetStudentFeedbackList(ctx *gin.Context)
    GetStudentFeedbackStats(ctx *gin.Context)
    GetExpertFeedbackList(ctx *gin.Context)
    GetExpertFeedbackDetail(ctx *gin.Context)
    
    // RIASEC Master Data (UI Screen 12-13)
    GetRiasecCodeList(ctx *gin.Context)
    GetRiasecCodeDetail(ctx *gin.Context)
    UpdateRiasecCode(ctx *gin.Context)
}
```

**Tanggung Jawab Per Method:**

1. **GetTestHistory**
   - Bind query params ke `GetTestHistoryRequest`
   - Validate binding errors
   - Call service method
   - Send response (success/error)

2. **GetTestDetail**
   - Extract `:id` dari URL param
   - Parse string → int64
   - Call service method
   - Send response

3. **DeleteTestData**
   - Bind JSON body ke `DeleteTestDataRequest`
   - Validate: test_ids array not empty
   - Call service method
   - Send success response

4. **ExportTestHistory**
   - Bind JSON body ke `ExportTestHistoryRequest`
   - Validate: format must be csv|excel|pdf
   - Call service method
   - Send response dengan file URL

5. **GetStudentFeedbackList**
   - Bind query params ke `GetFeedbackListRequest`
   - Call service method
   - Send response

6. **GetStudentFeedbackStats**
   - Bind query params: category_id, time_range
   - Call service method
   - Send response dengan chart data

7. **GetExpertFeedbackList**
   - Bind query params ke `GetExpertFeedbackListRequest`
   - Call service method
   - Send response

8. **GetExpertFeedbackDetail**
   - Extract `:id` dari URL param
   - Call service method
   - Send response

9. **GetRiasecCodeList**
   - Bind query params ke `GetRiasecCodeListRequest`
   - Call service method
   - Send response

10. **GetRiasecCodeDetail**
    - Extract `:id` dari URL param
    - Call service method
    - Send response

11. **UpdateRiasecCode**
    - Extract `:id` dari URL param
    - Bind JSON body ke `UpdateRiasecCodeRequest`
    - Validate all arrays have min 1 item
    - Call service method
    - Send success response

**Error Handling Pattern:**
```
if err := ctx.ShouldBind(&req); err != nil {
    response.NewFailed("invalid request", myerror.InvalidRequest(err)).Send(ctx)
    return
}

result, err := c.adminService.Method(ctx, req)
if err != nil {
    response.NewFailed("operation failed", err).Send(ctx)
    return
}

response.NewSuccess("operation successful", result).Send(ctx)
```

---

## FASE 4: Routes Registration

### File: `internal/api/kenali_diri/routes/kenalidiri_admin_route.go`

**Purpose:** Register all admin endpoints dengan middleware

**Route Structure:**
```
/api/v1/admin/kenali-diri
├── /history                          GET     (UI Screen 1: List)
├── /history/:id                      GET     (UI Screen 2-4: Detail)
├── /history                          DELETE  (UI Screen 6-7: Bulk Delete)
├── /history/export                   POST    (UI Screen 5: Export)
│
├── /feedback/student                 GET     (UI Screen 8: List)
├── /feedback/student/stats           GET     (UI Screen 9: Stats)
├── /feedback/expert                  GET     (UI Screen 10: List)
├── /feedback/expert/:id              GET     (UI Screen 11: Detail)
│
├── /riasec-codes                     GET     (UI Screen 12: List)
├── /riasec-codes/:id                 GET     (UI Screen 13: Detail)
└── /riasec-codes/:id                 PUT     (UI Screen 13: Update)
```

**Middleware Chain:**
- `middleware.Authenticate()` - verify JWT token
- `middleware.RoleAdmin()` - ensure role = ADMIN

**Implementation Pattern:**
```go
func ServeKenalidiriAdmin(
    app *gin.Engine,
    adminController controller.KenalidiriAdminController,
    middleware middleware.Middleware,
) {
    adminRoutes := app.Group("/api/v1/admin/kenali-diri")
    adminRoutes.Use(middleware.Authenticate(), middleware.RoleAdmin())
    {
        // Test History Management
        adminRoutes.GET("/history", adminController.GetTestHistory)
        adminRoutes.GET("/history/:id", adminController.GetTestDetail)
        adminRoutes.DELETE("/history", adminController.DeleteTestData)
        adminRoutes.POST("/history/export", adminController.ExportTestHistory)
        
        // Feedback Management
        adminRoutes.GET("/feedback/student", adminController.GetStudentFeedbackList)
        adminRoutes.GET("/feedback/student/stats", adminController.GetStudentFeedbackStats)
        adminRoutes.GET("/feedback/expert", adminController.GetExpertFeedbackList)
        adminRoutes.GET("/feedback/expert/:id", adminController.GetExpertFeedbackDetail)
        
        // RIASEC Master Data
        adminRoutes.GET("/riasec-codes", adminController.GetRiasecCodeList)
        adminRoutes.GET("/riasec-codes/:id", adminController.GetRiasecCodeDetail)
        adminRoutes.PUT("/riasec-codes/:id", adminController.UpdateRiasecCode)
    }
}
```

---

## FASE 5: Dependency Injection (Config)

### File: `internal/config/rest_config.go`

**Purpose:** Wire all components together

**Initialization Order:**
```
1. Package/Utilities
   - exportService
   - cacheService

2. Repositories (Shared)
   - kenalidiriHistoryRepository
   - kenalidiriCategoryRepository
   - riasecCodeRepository (with cache)
   - testSessionRepository
   - riasecRepository
   - ikigaiRepository
   - recommendationRepository
   - feedbackRepository

3. Services (Admin)
   - kenalidiriAdminService (inject all repos + export + cache)

4. Controllers (Admin)
   - kenalidiriAdminController (inject admin service)

5. Routes Registration
   - ServeKenalidiriAdmin(app, controller, middleware)
```

**Config Code Pattern:**
```go
func NewRest() RestConfig {
    db := db.New()
    app := gin.Default()
    server := NewRouter(app)
    firebaseApp := myfirebase.New()
    middleware := middleware.New(db, firebaseApp.MustGetClient())
    
    var (
        // Packages
        exportService export.ExportService = export.New()
        cacheService  cache.CacheService   = cache.NewRedis()
        
        // Repositories
        kdHistoryRepo     = repository.NewKenalidiriHistory(db)
        kdCategoryRepo    = repository.NewKenalidiriCategory(db, cacheService)
        kdRiasecCodeRepo  = repository.NewRiasecCode(db, cacheService)
        kdTestSessionRepo = repository.NewTestSession(db)
        kdRiasecRepo      = repository.NewRiasec(db)
        kdIkigaiRepo      = repository.NewIkigai(db)
        kdRecommendationRepo = repository.NewRecommendation(db)
        kdFeedbackRepo    = repository.NewFeedback(db)
        
        // Service
        kdAdminService = service.NewKenalidiriAdmin(
            kdHistoryRepo,
            kdCategoryRepo,
            kdRiasecCodeRepo,
            kdTestSessionRepo,
            kdRiasecRepo,
            kdIkigaiRepo,
            kdRecommendationRepo,
            kdFeedbackRepo,
            exportService,
            cacheService,
            db,
        )
        
        // Controller
        kdAdminController = controller.NewKenalidiriAdmin(kdAdminService)
    )
    
    // Register routes
    routes.ServeKenalidiriAdmin(server, kdAdminController, middleware)
    
    return RestConfig{server: server}
}
```

---

## Summary Checklist

### Files yang Perlu Dibuat

#### DTO Layer (8 files)
- [ ] `pagination_response.go` - shared pagination
- [ ] `kenalidiri_admin_dto_request.go` - 5 request structs
- [ ] `feedback_dto_request.go` - 2 request structs
- [ ] `kenalidiri_admin_dto_response.go` - 10 response structs
- [ ] `feedback_dto_response.go` - 6 response structs

#### Repository Layer (8 files)
- [ ] `kenalidiri_history_repository.go` - 7 methods
- [ ] `kenalidiri_category_repository.go` - 3 methods (with cache)
- [ ] `riasec_code_repository.go` - 6 methods (with cache, CRITICAL)
- [ ] `test_session_repository.go` - 5 methods
- [ ] `riasec_repository.go` - 7 methods
- [ ] `ikigai_repository.go` - 7 methods
- [ ] `recommendation_repository.go` - 3 methods
- [ ] `feedback_repository.go` - 5 methods

#### Export Services (4 files)
- [ ] `csv_exporter.go`
- [ ] `excel_exporter.go`
- [ ] `pdf_exporter.go`
- [ ] `export_service.go` - unified interface

#### Cache Services (2 files)
- [ ] `cache_service.go` - interface
- [ ] `redis_cache.go` - implementation

#### Service Layer (1 file)
- [ ] `kenalidiri_admin_service.go` - 11 methods

#### Controller Layer (1 file)
- [ ] `kenalidiri_admin_controller.go` - 11 handlers

#### Routes Layer (1 file)
- [ ] `kenalidiri_admin_route.go` - 11 endpoints

#### Config (update existing)
- [ ] Update `rest_config.go` - add kenali_diri dependencies

---

## Total: 26 New Files

**Breakdown:**
- DTO: 5 files
- Repository: 8 files
- Export: 4 files
- Cache: 2 files
- Service: 1 file
- Controller: 1 file
- Routes: 1 file
- Config: 1 update

---

## Critical Implementation Notes

### 1. Caching Strategy (MANDATORY)
- RIASEC codes (156 records) MUST be cached
- Cache all 156 codes on first access
- TTL: 24 hours atau permanent
- Impact tanpa cache: 500ms+ per request → <10ms with cache
- Invalidate cache on update operations

### 2. Export Performance
- Use transaction for data consistency
- Limit max export: 10,000 records per request
- Generate file async if > 1,000 records (optional enhancement)
- Upload to S3 untuk persistence
- Cleanup temp files after upload

### 3. Pagination Best Practices
- Default: page=1, limit=25
- Max limit: 100 per request
- Return total count untuk UI pagination controls
- Use OFFSET-LIMIT atau cursor-based (prefer cursor untuk large datasets)

### 4. Query Optimization
- Create indexes on:
  - `kenalidiri_history`: (user_id, test_category_id, status, started_at)
  - `riasec_results`: (test_session_id)
  - `feedback`: (category_id, submitted_at)
- Use preload untuk relations (avoid N+1)
- Use batch queries untuk multiple IDs

### 5. Transaction Management
- Use transaction untuk bulk delete
- Rollback on error
- Cascade delete ke child tables
- Consider soft delete untuk audit trail

### 6. Error Handling
- Return custom error types (`myerror` package)
- Log errors untuk debugging
- Don't expose internal errors ke client
- Use appropriate HTTP status codes

### 7. Security
- Validate all inputs
- Sanitize user-provided data
- Check authorization (admin role required)
- Rate limit export endpoints
- Prevent SQL injection (GORM handles this)

---

## API Response Format Standard

### Success Response
```json
{
  "success": true,
  "message": "operation successful",
  "data": { ... }
}
```

### Error Response
```json
{
  "success": false,
  "message": "operation failed",
  "error": {
    "code": "INVALID_REQUEST",
    "details": "validation error message"
  }
}
```

### List Response with Pagination
```json
{
  "success": true,
  "message": "data retrieved successfully",
  "data": {
    "data": [ ... ],
    "pagination": {
      "current_page": 1,
      "total_pages": 10,
      "total_records": 250,
      "per_page": 25
    }
  }
}
```

---

## Testing Strategy

### Unit Tests
- Test repository methods dengan mock database
- Test service business logic dengan mock repositories
- Test DTO validation rules

### Integration Tests
- Test controller handlers dengan mock service
- Test route registration
- Test middleware chain

### API Tests
- Use Bruno/Postman collections
- Test all endpoints dengan different scenarios
- Test pagination, filtering, sorting
- Test error cases (invalid input, not found, unauthorized)

### Performance Tests
- Load test export endpoints (simulate 100+ concurrent requests)
- Monitor cache hit rate (should be >95% for RIASEC codes)
- Profile database queries (identify slow queries)

---

## Deployment Considerations

### Environment Variables
```
# Cache
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Export
AWS_REGION=ap-southeast-1
AWS_S3_BUCKET=rextra-exports
AWS_ACCESS_KEY=
AWS_SECRET_KEY=

# Database
DB_MAX_CONNECTIONS=100
DB_IDLE_CONNECTIONS=10
```

### Database Migrations
- Create indexes before deploying
- Test migration on staging first
- Backup data before production deployment

### Cache Warm-up
- Pre-load RIASEC codes on application start
- Cache warm-up script untuk production

### Monitoring
- Log cache miss rate
- Monitor export file sizes
- Track API response times
- Alert on error rate spikes

---

**Dokumen ini dibuat untuk:** Technical Brief Pengembangan REST API Backend Golang - Manajemen Data Kenali Diri  
**Target Developer:** Backend Golang Engineers  
**Estimasi Waktu:** 4-6 minggu (1 developer) atau 2-3 minggu (2 developers)  
**Prioritas Implementasi:** FASE 1 → FASE 2 → FASE 3 → FASE 4 → FASE 5

---

**END OF DOCUMENT**