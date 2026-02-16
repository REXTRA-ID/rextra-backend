package seeds

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"rextra-backend/internal/entity"
	mylog "rextra-backend/internal/pkg/logger"
)

// SeedRemainingEntities seeds all entities that don't have their own seeder yet.
// Depends on: SeederUser, SeederKenaliDiri having been run first (for User + CareerProfileTestSession).
//
// NOTE: Riasec, UserRiasec, UserIkigai are commented out in migrations (tables don't exist).
// If they are enabled in migration.go, add seeding for them here.
func SeedRemainingEntities(db *gorm.DB) error {
	mylog.Infof("[PROCESS] Seeding remaining entities...")

	return db.Transaction(func(tx *gorm.DB) error {
		// ── Lookup existing user & session ──────────────────────────────
		var user entity.User
		if err := tx.Where("email = ?", "user@email.com").First(&user).Error; err != nil {
			return err
		}
		userID := user.ID

		var session entity.CareerProfileTestSession
		if err := tx.Where("user_id = ?", userID).Order("id ASC").First(&session).Error; err != nil {
			return err
		}
		sessionID := session.ID

		var category entity.KenaliDiriCategory
		if err := tx.Where("category_code = ?", "CAREER_PROFILE").First(&category).Error; err != nil {
			return err
		}

		// ── 1. Persona ─────────────────────────────────────────────────
		persona := entity.Persona{
			UserID:                    userID,
			PersonaType:               entity.Pathfinder,
			EducationSaved:            true,
			CareerRecommendationTired: true,
			CareerDictionaryAccessed:  false,
			CareerPlanCreated:         false,
			PorfolioRecorded:          false,
			ExplorationAIUsed:         false,
			CVCreated:                 false,
			InterviewSimulated:        false,
			LinkedinOptimaze:          false,
			IntershipPlanReported:     false,
		}
		if err := tx.Create(&persona).Error; err != nil {
			return err
		}
		mylog.Infof("[OK] Persona seeded")

		// ── 2. TokenWallet ─────────────────────────────────────────────
		wallet := entity.TokenWallet{
			UserID:    userID,
			Balance:   150,
			UpdatedAt: time.Now(),
		}
		if err := tx.Create(&wallet).Error; err != nil {
			return err
		}
		mylog.Infof("[OK] TokenWallet seeded")

		// ── 3. TokenLedger ─────────────────────────────────────────────
		// Use map-based creation to avoid GORM issues with Ref relation field and metadata type
		ledgerMaps := []map[string]interface{}{
			{
				"occurred_at":    mustTime("2024-01-01T08:00:00Z"),
				"wallet_id":      wallet.ID,
				"direction":      string(entity.DirectionIN),
				"amount":         int64(200),
				"balance_before": int64(0),
				"balance_after":  int64(200),
				"source_type":    string(entity.SourceTypeTopup),
				"description":    "Top-up awal 200 token",
				"metadata":       `{"source":"initial_topup","note":"seeder data"}`,
				"created_at":     mustTime("2024-01-01T08:00:00Z"),
			},
			{
				"occurred_at":    mustTime("2024-01-02T11:00:00Z"),
				"wallet_id":      wallet.ID,
				"direction":      string(entity.DirectionOUT),
				"amount":         int64(50),
				"balance_before": int64(200),
				"balance_after":  int64(150),
				"source_type":    string(entity.SourceTypeUsage),
				"description":    "Pemakaian tes profil karier",
				"metadata":       `{"feature":"career_profile_test","note":"seeder data"}`,
				"created_at":     mustTime("2024-01-02T11:00:00Z"),
			},
		}
		for _, m := range ledgerMaps {
			if err := tx.Model(&entity.TokenLedger{}).Create(m).Error; err != nil {
				return err
			}
		}
		mylog.Infof("[OK] TokenLedger seeded (2 records)")

		// ── 4. CustomPricingConfig + Tiers ─────────────────────────────
		// Use map-based creation to avoid GORM issues with CustomPricingMetadata type
		configID := uuid.New()
		configMap := map[string]interface{}{
			"id":                          configID,
			"is_enabled":                  true,
			"min_token":                   int64(10),
			"max_token":                   int64(1000),
			"recommended_price_per_token": int64(5000),
			"is_current":                  true,
			"effective_from":              mustTime("2024-01-01T00:00:00Z"),
			"metadata":                    `{"version":"1.0","note":"seeder"}`,
		}
		if err := tx.Model(&entity.CustomPricingConfig{}).Create(configMap).Error; err != nil {
			return err
		}

		tiers := []entity.CustomPricingTier{
			{
				ConfigID:    configID,
				FromToken:   10,
				ToToken:     99,
				DiscountPct: 0,
			},
			{
				ConfigID:    configID,
				FromToken:   100,
				ToToken:     499,
				DiscountPct: 5,
			},
			{
				ConfigID:    configID,
				FromToken:   500,
				ToToken:     1000,
				DiscountPct: 10,
			},
		}
		for _, tier := range tiers {
			if err := tx.Create(&tier).Error; err != nil {
				return err
			}
		}
		mylog.Infof("[OK] CustomPricingConfig + 3 Tiers seeded")

		// ── 5. RiasecQuestionSet ────────────────────────────────────────
		questionIDs, _ := json.Marshal([]int{1, 5, 12, 18, 24, 30, 35, 42, 48, 55, 60, 3, 8, 15, 21, 27, 33, 39, 45, 51})
		riasecQuestionSet := entity.RiasecQuestionSet{
			TestSessionID: sessionID,
			QuestionIDs:   datatypes.JSON(questionIDs),
		}
		if err := tx.Create(&riasecQuestionSet).Error; err != nil {
			return err
		}
		mylog.Infof("[OK] RiasecQuestionSet seeded")

		// ── 6. RiasecResponse ──────────────────────────────────────────
		responsesData := mustJSON([]map[string]interface{}{
			{"question_id": 1, "answer": "A", "riasec_type": "R"},
			{"question_id": 5, "answer": "B", "riasec_type": "I"},
			{"question_id": 12, "answer": "A", "riasec_type": "A"},
			{"question_id": 18, "answer": "B", "riasec_type": "S"},
			{"question_id": 24, "answer": "A", "riasec_type": "E"},
			{"question_id": 30, "answer": "B", "riasec_type": "C"},
			{"question_id": 35, "answer": "A", "riasec_type": "R"},
			{"question_id": 42, "answer": "A", "riasec_type": "I"},
			{"question_id": 48, "answer": "B", "riasec_type": "A"},
			{"question_id": 55, "answer": "A", "riasec_type": "I"},
		})
		riasecResponse := entity.RiasecResponse{
			TestSessionID: sessionID,
			ResponsesData: responsesData,
		}
		if err := tx.Create(&riasecResponse).Error; err != nil {
			return err
		}
		mylog.Infof("[OK] RiasecResponse seeded")

		// ── 7. IkigaiResponse ─────────────────────────────────────────
		ikigaiResponse := entity.IkigaiResponse{
			TestSessionID:        sessionID,
			Dimension1Love:       mustJSON([]string{"analisis data", "memecahkan masalah", "membuat visualisasi"}),
			Dimension2GoodAt:     mustJSON([]string{"berpikir analitis", "programming Python", "statistik"}),
			Dimension3WorldNeeds: mustJSON([]string{"keputusan berbasis data", "efisiensi bisnis", "insight pasar"}),
			Dimension4PaidFor:    mustJSON([]string{"data analyst", "business intelligence", "data scientist"}),
			Completed:            true,
			CreatedAt:            mustTime("2024-01-02T10:17:00Z"),
			CompletedAt:          ptrTime("2024-01-02T10:25:00Z"),
		}
		if err := tx.Create(&ikigaiResponse).Error; err != nil {
			return err
		}
		mylog.Infof("[OK] IkigaiResponse seeded")

		// ── 8. IkigaiDimensionScore ───────────────────────────────────
		ikigaiDimScore := entity.IkigaiDimensionScore{
			TestSessionID: sessionID,
			ScoresData: mustJSON(map[string]interface{}{
				"love":        map[string]interface{}{"score": 85, "confidence": 0.92},
				"good_at":     map[string]interface{}{"score": 78, "confidence": 0.88},
				"world_needs": map[string]interface{}{"score": 90, "confidence": 0.95},
				"paid_for":    map[string]interface{}{"score": 88, "confidence": 0.90},
			}),
			CalculatedAt:  mustTime("2024-01-02T10:24:00Z"),
			AIModelUsed:   "gemini-1.5-flash",
			TotalAPICalls: 4,
		}
		if err := tx.Create(&ikigaiDimScore).Error; err != nil {
			return err
		}
		mylog.Infof("[OK] IkigaiDimensionScore seeded")

		// ── 9. IkigaiCandidateProfession ──────────────────────────────
		ikigaiCandidate := entity.IkigaiCandidateProfession{
			TestSessionID: sessionID,
			CandidatesData: mustJSON([]map[string]interface{}{
				{"profession_id": 101, "profession_name": "Data Analyst", "score": 92},
				{"profession_id": 102, "profession_name": "Business Intelligence Analyst", "score": 86},
				{"profession_id": 103, "profession_name": "Data Scientist", "score": 84},
				{"profession_id": 104, "profession_name": "Data Engineer", "score": 80},
				{"profession_id": 105, "profession_name": "Product Analyst", "score": 78},
			}),
			TotalCandidates:    5,
			GenerationStrategy: "4_tier_expansion",
			MaxCandidatesLimit: 15,
			GeneratedAt:        mustTime("2024-01-02T10:25:00Z"),
		}
		if err := tx.Create(&ikigaiCandidate).Error; err != nil {
			return err
		}
		mylog.Infof("[OK] IkigaiCandidateProfession seeded")

		// ── 10. ExpertFeedback ──────────────────────────────────────────
		expertSuggestion := "Pertimbangkan untuk menambahkan pengalaman magang di bidang data."
		expertFeedback := entity.ExpertFeedback{
			TestSessionID:  sessionID,
			TestCategoryID: category.ID,
			ExpertName:     "Dr. Budi Santoso",
			Profession:     "Data Scientist",
			Degree:         "Ph.D.",
			Experience:     "10 tahun",
			Education:      "S3",
			University:     "Universitas Indonesia",
			StudyProgram:   "Ilmu Komputer",
			CategoryTest:   "CAREER_PROFILE",
			TopFiveProfessions: mustJSON([]map[string]interface{}{
				{"rank": 1, "profession_name": "Data Analyst"},
				{"rank": 2, "profession_name": "Business Intelligence"},
				{"rank": 3, "profession_name": "Data Scientist"},
				{"rank": 4, "profession_name": "Data Engineer"},
				{"rank": 5, "profession_name": "Product Analyst"},
			}),
			AccuracyScore: 6,
			LogicScore:    7,
			BenefitScore:  6,
			Obstacles:     mustJSON([]string{"Kurang pengalaman kerja", "Perlu sertifikasi tambahan"}),
			Suggestions:   &expertSuggestion,
			SubmittedAt:   mustTime("2024-01-07T14:00:00Z"),
		}
		if err := tx.Create(&expertFeedback).Error; err != nil {
			return err
		}
		mylog.Infof("[OK] ExpertFeedback seeded")

		// ── 11. StudentFeedback ─────────────────────────────────────────
		studentFeedback := entity.StudentFeedback{
			UserID:            userID,
			TestCategoryID:    category.ID,
			EaseOfUseScore:    6,
			RelevanceScore:    7,
			SatisfactionScore: 6,
			Obstacles:         mustJSON([]string{"Pertanyaan terlalu banyak", "Butuh waktu lama"}),
			SubmittedAt:       mustTime("2024-01-05T09:30:00Z"),
		}
		if err := tx.Create(&studentFeedback).Error; err != nil {
			return err
		}
		mylog.Infof("[OK] StudentFeedback seeded")

		mylog.Infof("[COMPLETE] Seeding remaining entities completed.")
		return nil
	})
}
