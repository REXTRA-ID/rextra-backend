package seeds

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"rextra-backend/internal/entity"
	mylog "rextra-backend/internal/pkg/logger"
)

func SeederKenaliDiri(db *gorm.DB) error {
	mylog.Infof("[PROCESS] Seeding Kenali Diri fixtures...")
	return db.Transaction(func(tx *gorm.DB) error {
		tables := []string{
			"careerprofile_feedback_expert_obstacles",
			"careerprofile_feedback_obstacles",
			"careerprofile_feedback_expert",
			"careerprofile_feedback_student",
			"kenalidiri_feedback",
			"user_career_profiles",
			"career_recommendations",
			"ikigai_total_scores",
			"riasec_results",
			"kenalidiri_history",
			"careerprofile_test_sessions",
			"kenalidiri_categories",
		}
		for _, tbl := range tables {
			if err := tx.Exec("TRUNCATE TABLE " + tbl + " RESTART IDENTITY CASCADE;").Error; err != nil {
				return err
			}
		}

		var user entity.User
		if err := tx.Where("email = ?", "user@email.com").First(&user).Error; err != nil {
			return err
		}
		userID := user.ID

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
			{
				ID:              3,
				CategoryCode:    "CAREER_PROFILE",
				CategoryName:    "Tes Profil Karier",
				Description:     ptrString("Tes profil karier lengkap"),
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
			UserID:             userID,
			SessionToken:       "session-kenalidiri-1",
			PersonaType:        "pathfinder",
			TestGoal:           entity.TestGoalRecommendation,
			UsesIkigai:         true,
			Status:             "completed",
			StartedAt:          mustTime("2024-01-02T10:00:00Z"),
			CompletedAt:        ptrTime("2024-01-02T10:30:00Z"),
			RiasecCompletedAt:  ptrTime("2024-01-02T10:15:00Z"),
			IkigaiCompletedAt:  ptrTime("2024-01-02T10:30:00Z"),
			AlgorithmVersion:   ptrString("v1.0.0"),
			QuestionSetVersion: ptrString("2024-Q1"),
		}
		if err := tx.Create(&session).Error; err != nil {
			return err
		}

		userProfile := entity.UserCareerProfile{
			UserID:          userID,
			ActiveSessionID: session.ID,
			Pinned:          true,
			SetSource:       entity.SetSourceAutoFirstTime,
			SetAt:           mustTime("2024-01-02T10:30:00Z"),
		}
		if err := tx.Create(&userProfile).Error; err != nil {
			return err
		}

		history := entity.KenaliDiriHistory{
			UserID:          userID,
			TestCategoryID:  3,
			DetailSessionID: session.ID,
			Status:          "completed",
			StartedAt:       mustTime("2024-01-02T10:00:00Z"),
			CompletedAt:     ptrTime("2024-01-02T10:30:00Z"),
		}
		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		riasecResult := entity.RiasecResult{
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
			TestSessionID: session.ID,
			RecommendationsData: mustJSON([]map[string]interface{}{
				{"profession_id": 101, "profession_name": "Data Analyst", "match_percentage": 92, "reasoning": "Kuat di analisis dan insight bisnis"},
				{"profession_id": 102, "profession_name": "Business Intelligence", "match_percentage": 86, "reasoning": "Mampu menyampaikan data ke stakeholder"},
				{"profession_id": 103, "profession_name": "Data Scientist", "match_percentage": 84, "reasoning": "Kombinasi analisis dan machine learning"},
				{"profession_id": 104, "profession_name": "Data Engineer", "match_percentage": 80, "reasoning": "Keahlian teknis dalam infrastruktur data"},
				{"profession_id": 105, "profession_name": "Product Analyst", "match_percentage": 78, "reasoning": "Analisis untuk keputusan produk"},
			}),
			TopProfession1ID: ptrInt64(101),
			TopProfession2ID: ptrInt64(102),
			GeneratedAt:      mustTime("2024-01-02T10:26:00Z"),
			AIModelUsed:      "gemini-1.5-flash",
		}
		if err := tx.Create(&recommendation).Error; err != nil {
			return err
		}

		studentFeedbackHeader := entity.KenaliDiriFeedback{
			TestCategory:     "CAREER_PROFILE",
			TestSessionID:    session.ID,
			RespondentType:   entity.RespondentTypeStudent,
			RespondentUserID: userID,
			SubmittedAt:      mustTime("2024-01-05T09:00:00Z"),
		}
		if err := tx.Create(&studentFeedbackHeader).Error; err != nil {
			return err
		}

		studentFeedbackDetail := entity.CareerProfileFeedbackStudent{
			FeedbackID:        studentFeedbackHeader.ID,
			EaseScore:         6,
			RelevanceScore:    7,
			SatisfactionScore: 6,
			MessageToTeam:     ptrString("Pengalaman cukup baik secara keseluruhan"),
		}
		if err := tx.Create(&studentFeedbackDetail).Error; err != nil {
			return err
		}

		expertFeedbackHeader := entity.KenaliDiriFeedback{
			TestCategory:     "CAREER_PROFILE",
			TestSessionID:    session.ID,
			RespondentType:   entity.RespondentTypeExpert,
			RespondentUserID: userID,
			SubmittedAt:      mustTime("2024-01-06T10:00:00Z"),
		}
		if err := tx.Create(&expertFeedbackHeader).Error; err != nil {
			return err
		}

		expertFeedbackDetail := entity.CareerProfileFeedbackExpert{
			FeedbackID:      expertFeedbackHeader.ID,
			AccuracyScore:   6,
			LogicScore:      6,
			UsefulnessScore: 5,

			ExpertName:            "Dr. Hana Wijaya",
			ExpertProfession:      "Data Analyst",
			ExpertDegree:          "M.T.",
			ExpertExperienceYears: ptrInt32(5),
			ExpertEducationLevel:  "Magister",
			ExpertUniversity:      "ITB",
			ExpertStudyProgram:    "Teknik Industri",

			ExpertProfessionID: ptrInt64(101),
			Top5RecommendationsJSON: mustJSON([]map[string]interface{}{
				{"rank": 1, "profession_id": 101, "profession_name": "Data Analyst"},
				{"rank": 2, "profession_id": 102, "profession_name": "Business Intelligence"},
				{"rank": 3, "profession_id": 103, "profession_name": "Data Scientist"},
				{"rank": 4, "profession_id": 104, "profession_name": "Data Engineer"},
				{"rank": 5, "profession_id": 105, "profession_name": "Product Analyst"},
			}),
			SuggestionText: ptrString("Tambahkan studi kasus lokal untuk konteks Indonesia."),
		}
		if err := tx.Create(&expertFeedbackDetail).Error; err != nil {
			return err
		}

		return nil
	})
}

func ptrString(s string) *string {
	return &s
}

func ptrInt32(v int32) *int32 {
	return &v
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
