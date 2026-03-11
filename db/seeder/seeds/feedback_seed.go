package seeds

import (
	"encoding/json"
	"os"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"rextra-backend/internal/entity"
	"rextra-backend/internal/pkg/feedback"
	mylog "rextra-backend/internal/pkg/logger"
	"rextra-backend/internal/utils"
)

type FeedbackData struct {
	Students []StudentFeedbackSeed `json:"students"`
	Experts  []ExpertFeedbackSeed  `json:"experts"`
}

type StudentFeedbackSeed struct {
	Name             string `json:"name"`
	Email            string `json:"email"`
	Phone            string `json:"phone"`
	SubmittedDaysAgo int    `json:"submitted_days_ago"`
	Scores           struct {
		Ease         int16 `json:"ease"`
		Relevance    int16 `json:"relevance"`
		Satisfaction int16 `json:"satisfaction"`
	} `json:"scores"`
	Message    *string `json:"message"`
	ObstacleID int32   `json:"obstacle_id"`
}

type ExpertFeedbackSeed struct {
	Name             string `json:"name"`
	Email            string `json:"email"`
	Phone            string `json:"phone"`
	Profession       string `json:"profession"`
	ProfessionID     int64  `json:"profession_id"`
	SubmittedDaysAgo int    `json:"submitted_days_ago"`
	Top5Status       string `json:"top5_status"`
	Scores           struct {
		Accuracy   int16 `json:"accuracy"`
		Logic      int16 `json:"logic"`
		Usefulness int16 `json:"usefulness"`
	} `json:"scores"`
	Suggestion *string `json:"suggestion"`
	ExpertInfo struct {
		Degree string `json:"degree"`
		Uni    string `json:"uni"`
		Major  string `json:"major"`
		Exp    int32  `json:"exp"`
	} `json:"expert_info"`
	ObstacleID int32 `json:"obstacle_id"`
}

// SeedFeedbackData reads from JSON and seeds feedback data.
func SeedFeedbackData(db *gorm.DB) error {
	mylog.Infof("[PROCESS] Seeding Feedback Data from JSON...")

	// Read JSON Data
	jsonFile, err := os.Open("./db/seeder/data/feedback_data.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	var data FeedbackData
	if err := json.NewDecoder(jsonFile).Decode(&data); err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, s := range data.Students {
			// Create user if not exists
			var user entity.User
			if err := tx.Where("email = ?", s.Email).First(&user).Error; err != nil {
				hashedPwd, _ := utils.HashPassword("password")
				user = entity.User{
					Email:       s.Email,
					Password:    hashedPwd,
					IsVerified:  true,
					PhoneNumber: s.Phone,
					Role:        "USER",
					Fullname:    s.Name,
				}
				if err := tx.Create(&user).Error; err != nil {
					return err
				}
			}

			submittedAt := time.Now().AddDate(0, 0, -s.SubmittedDaysAgo)

			kp := entity.KenaliDiriFeedback{
				TestCategory:     feedback.TestCategoryCareerProfile,
				TestSessionID:    1, // assume session 1 exists
				RespondentType:   entity.RespondentTypeStudent,
				RespondentUserID: user.ID,
				SubmittedAt:      submittedAt,
			}
			if err := tx.Create(&kp).Error; err != nil {
				return err
			}

			detail := entity.CareerProfileFeedbackStudent{
				FeedbackID:        kp.ID,
				EaseScore:         s.Scores.Ease,
				RelevanceScore:    s.Scores.Relevance,
				SatisfactionScore: s.Scores.Satisfaction,
				MessageToTeam:     s.Message,
			}
			if err := tx.Create(&detail).Error; err != nil {
				return err
			}

			obs := entity.CareerProfileFeedbackObstacle{
				FeedbackID: kp.ID,
				ObstacleID: s.ObstacleID,
			}
			if err := tx.Create(&obs).Error; err != nil {
				return err
			}
		}

		// Seed Experts
		for _, e := range data.Experts {
			var user entity.User
			if err := tx.Where("email = ?", e.Email).First(&user).Error; err != nil {
				hashedPwd, _ := utils.HashPassword("password")
				user = entity.User{
					Email:       e.Email,
					Password:    hashedPwd,
					IsVerified:  true,
					PhoneNumber: e.Phone,
					Role:        "EXPERT",
					Fullname:    e.Name,
				}
				if err := tx.Create(&user).Error; err != nil {
					return err
				}
			}

			submittedAt := time.Now().AddDate(0, 0, -e.SubmittedDaysAgo)

			kp := entity.KenaliDiriFeedback{
				TestCategory:     feedback.TestCategoryCareerProfile,
				TestSessionID:    1,
				RespondentType:   entity.RespondentTypeExpert,
				RespondentUserID: user.ID,
				SubmittedAt:      submittedAt,
			}
			if err := tx.Create(&kp).Error; err != nil {
				return err
			}

			// Construct Top 5 JSON logic based on status
			baseList := []map[string]interface{}{
				{"rank": 1, "profession_id": 901, "profession_name": "Pro A"},
				{"rank": 2, "profession_id": 902, "profession_name": "Pro B"},
				{"rank": 3, "profession_id": 903, "profession_name": "Pro C"},
				{"rank": 4, "profession_id": 904, "profession_name": "Pro D"},
				{"rank": 5, "profession_id": 905, "profession_name": "Pro E"},
			}

			switch e.Top5Status {
			case "P1":
				baseList[0]["profession_id"] = e.ProfessionID
				baseList[0]["profession_name"] = e.Profession
			case "P2":
				baseList[1]["profession_id"] = e.ProfessionID
			case "P3_5", "P3":
				baseList[2]["profession_id"] = e.ProfessionID
			}

			jsonBytes, _ := json.Marshal(baseList)
			var profID *int64 = &e.ProfessionID

			detail := entity.CareerProfileFeedbackExpert{
				FeedbackID:              kp.ID,
				AccuracyScore:           e.Scores.Accuracy,
				LogicScore:              e.Scores.Logic,
				UsefulnessScore:         e.Scores.Usefulness,
				ExpertName:              e.Name,
				ExpertProfession:        e.Profession,
				ExpertProfessionID:      profID,
				Top5RecommendationsJSON: datatypes.JSON(jsonBytes),
				SuggestionText:          e.Suggestion,
				ExpertDegree:            e.ExpertInfo.Degree,
				ExpertUniversity:        e.ExpertInfo.Uni,
				ExpertStudyProgram:      e.ExpertInfo.Major,
				ExpertExperienceYears:   &e.ExpertInfo.Exp,
			}
			if err := tx.Create(&detail).Error; err != nil {
				return err
			}

			obs := entity.CareerProfileFeedbackExpertObstacle{
				FeedbackID: kp.ID,
				ObstacleID: e.ObstacleID,
			}
			if err := tx.Create(&obs).Error; err != nil {
				return err
			}
		}

		mylog.Infof("[COMPLETE] Seeding Feedback Data from JSON completed.")
		return nil
	})
}
