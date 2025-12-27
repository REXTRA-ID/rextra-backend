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

// SeederKenaliDiri seeds core Kenali Diri data to simplify endpoint testing.
func SeederKenaliDiri(db *gorm.DB) error {
	mylog.Infof("[PROCESS] Seeding Kenali Diri fixtures...")
	return db.Transaction(func(tx *gorm.DB) error {
		// Clean dependent tables to keep IDs stable for tests.
		tables := []string{
			"expert_feedback",
			"student_feedback",
			"career_recommendations",
			"ikigai_total_scores",
			"riasec_results",
			"kenalidiri_history",
			"careerprofile_test_sessions",
			"kenalidiri_categories",
		}
		for _, tbl := range tables {
			if err := tx.Exec("TRUNCATE TABLE "+tbl+" RESTART IDENTITY CASCADE;").Error; err != nil {
				return err
			}
		}

		userID := uuid.MustParse("ef9cf8e8-46b1-4e91-89d0-40f6c824319e")

		categories := []entity.KenaliDiriCategory{
			{
				ID:              1,
				CategoryCode:    "RIASEC",
				CategoryName:    "Tes RIASEC",
				Description:     ptrString("Tes minat RIASEC"),
				DetailTableName: "careerprofile_test_sessions",
				IsActive:        true,
			},
			{
				ID:              2,
				CategoryCode:    "IKIGAI",
				CategoryName:    "Tes IKIGAI",
				Description:     ptrString("Tes ikigai lanjutan"),
				DetailTableName: "careerprofile_test_sessions",
				IsActive:        true,
			},
		}
		for _, c := range categories {
			if err := tx.Create(&c).Error; err != nil {
				return err
			}
		}

		session := entity.CareerProfileTestSession{
			ID:                1,
			UserID:            userID,
			SessionToken:      "session-kenalidiri-1",
			Status:            "completed",
			StartedAt:         mustTime("2024-01-02T10:00:00Z"),
			CompletedAt:       ptrTime("2024-01-02T10:30:00Z"),
			RiasecCompletedAt: ptrTime("2024-01-02T10:15:00Z"),
			IkigaiCompletedAt: ptrTime("2024-01-02T10:30:00Z"),
		}
		if err := tx.Create(&session).Error; err != nil {
			return err
		}

		history := entity.KenaliDiriHistory{
			ID:              1,
			UserID:          userID,
			TestCategoryID:  1,
			DetailSessionID: session.ID,
			Status:          "completed",
			StartedAt:       mustTime("2024-01-02T10:00:00Z"),
			CompletedAt:     ptrTime("2024-01-02T10:30:00Z"),
		}
		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		riasecResult := entity.RiasecResult{
			ID:                    1,
			TestSessionID:         session.ID,
			ScoreR:                12,
			ScoreI:                14,
			ScoreA:                10,
			ScoreS:                9,
			ScoreE:                8,
			ScoreC:                11,
			RiasecCodeID:          1,
			RiasecCodeType:        "single",
			IsInconsistentProfile: false,
			CalculatedAt:          mustTime("2024-01-02T10:16:00Z"),
		}
		if err := tx.Create(&riasecResult).Error; err != nil {
			return err
		}

		ikigaiScores := entity.IkigaiTotalScore{
			ID:            1,
			TestSessionID: session.ID,
			ScoresData: mustJSON(map[string]interface{}{
				"love_narrative":        "Kamu menikmati pemecahan masalah teknis.",
				"good_at_narrative":     "Kemampuan analitis menjadi kekuatan utama.",
				"world_needs_narrative": "Pasar membutuhkan analis data untuk keputusan bisnis.",
				"paid_for_narrative":    "Peran analis data bernilai tinggi di industri.",
				"top_professions": []map[string]interface{}{
					{"profession_id": 101, "profession_name": "Data Analyst", "match_percentage": 92, "match_reasoning": "Sejalan dengan minat analitis"},
					{"profession_id": 102, "profession_name": "Business Intelligence", "match_percentage": 86, "match_reasoning": "Kombinasi analisis dan komunikasi"},
				},
			}),
			TopProfession1ID: 101,
			TopProfession2ID: 102,
			CalculatedAt:     mustTime("2024-01-02T10:25:00Z"),
		}
		if err := tx.Create(&ikigaiScores).Error; err != nil {
			return err
		}

		recommendation := entity.CareerRecommendation{
			ID:            1,
			TestSessionID: session.ID,
			RecommendationsData: mustJSON([]map[string]interface{}{
				{"profession_id": 101, "profession_name": "Data Analyst", "match_percentage": 92, "reasoning": "Kuat di analisis dan insight bisnis"},
				{"profession_id": 102, "profession_name": "Business Intelligence", "match_percentage": 86, "reasoning": "Mampu menyampaikan data ke stakeholder"},
			}),
			TopProfession1ID: ptrInt64(101),
			TopProfession2ID: ptrInt64(102),
			GeneratedAt:      mustTime("2024-01-02T10:26:00Z"),
			AIModelUsed:      "gemini-1.5-flash",
		}
		if err := tx.Create(&recommendation).Error; err != nil {
			return err
		}

		studentFeedback := entity.StudentFeedback{
			ID:                1,
			UserID:            userID,
			TestCategoryID:    1,
			EaseOfUseScore:    6,
			RelevanceScore:    7,
			SatisfactionScore: 6,
			Obstacles:         mustJSON([]string{"UI sedikit lambat"}),
			SubmittedAt:       mustTime("2024-01-05T09:00:00Z"),
		}
		if err := tx.Create(&studentFeedback).Error; err != nil {
			return err
		}

		expertFeedback := entity.ExpertFeedback{
			ID:                 1,
			TestSessionID:      session.ID,
			TestCategoryID:     1,
			ExpertName:         "Dr. Hana",
			Profession:         "Data Analyst",
			Degree:             "M.T.",
			Experience:         "5 tahun",
			Education:          "Magister Teknik Industri",
			University:         "ITB",
			StudyProgram:       "Teknik Industri",
			CategoryTest:       "RIASEC",
			TopFiveProfessions: mustJSON([]string{"Data Analyst", "BI Analyst", "Data Engineer", "Product Analyst", "Researcher"}),
			AccuracyScore:      6,
			LogicScore:         6,
			BenefitScore:       5,
			Obstacles:          mustJSON([]string{"Perlu data project nyata"}),
			Suggestions:        ptrString("Tambahkan studi kasus lokal."),
			SubmittedAt:        mustTime("2024-01-06T10:00:00Z"),
		}
		if err := tx.Create(&expertFeedback).Error; err != nil {
			return err
		}

		return nil
	})
}

func ptrString(s string) *string {
	return &s
}

func ptrInt64(v int64) *int64 {
	return &v
}

func ptrTime(val string) *time.Time {
	t := mustTime(val)
	return &t
}

func mustTime(val string) time.Time {
	t, err := time.Parse(time.RFC3339, val)
	if err != nil {
		panic(err)
	}
	return t
}

func mustJSON(v interface{}) datatypes.JSON {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return datatypes.JSON(b)
}
