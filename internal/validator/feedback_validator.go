package validator

import (
	"errors"
	"fmt"
	"rextra-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func ValidateStudentFeedback(feedback *entity.CareerProfileFeedbackStudent) error {
	if feedback.EaseScore < 1 || feedback.EaseScore > 7 {
		return errors.New("ease_score must be between 1 and 7")
	}
	if feedback.RelevanceScore < 1 || feedback.RelevanceScore > 7 {
		return errors.New("relevance_score must be between 1 and 7")
	}
	if feedback.SatisfactionScore < 1 || feedback.SatisfactionScore > 7 {
		return errors.New("satisfaction_score must be between 1 and 7")
	}

	return nil
}

func ValidateExpertFeedback(feedback *entity.CareerProfileFeedbackExpert) error {
	if feedback.AccuracyScore < 1 || feedback.AccuracyScore > 7 {
		return errors.New("accuracy_score must be between 1 and 7")
	}
	if feedback.LogicScore < 1 || feedback.LogicScore > 7 {
		return errors.New("logic_score must be between 1 and 7")
	}
	if feedback.UsefulnessScore < 1 || feedback.UsefulnessScore > 7 {
		return errors.New("usefulness_score must be between 1 and 7")
	}

	// Validate expert identity snapshot fields are not empty
	if feedback.ExpertName == "" {
		return errors.New("expert_name is required")
	}
	if feedback.ExpertProfession == "" {
		return errors.New("expert_profession is required")
	}

	return nil
}

func ValidateStudentObstacles(db *gorm.DB, obstacles []entity.CareerProfileFeedbackObstacle) error {
	if len(obstacles) == 0 {
		return nil
	}
	var obstacleIDs []int32
	for _, obs := range obstacles {
		obstacleIDs = append(obstacleIDs, obs.ObstacleID)
	}

	var options []entity.CareerProfileObstacleOption
	if err := db.Where("id IN ?", obstacleIDs).Find(&options).Error; err != nil {
		return fmt.Errorf("failed to load obstacle options: %w", err)
	}

	optionMap := make(map[int32]*entity.CareerProfileObstacleOption)
	for i := range options {
		optionMap[options[i].ID] = &options[i]
	}

	hasNoIssue := false
	hasOtherObstacles := false

	for _, obs := range obstacles {
		opt, exists := optionMap[obs.ObstacleID]
		if !exists {
			return fmt.Errorf("invalid obstacle_id: %d", obs.ObstacleID)
		}

		if !opt.IsActive {
			return fmt.Errorf("obstacle option '%s' is no longer active", opt.Key)
		}

		if opt.IsNoIssue {
			hasNoIssue = true
		} else {
			hasOtherObstacles = true
		}
		if opt.IsOther && (obs.OtherText == nil || *obs.OtherText == "") {
			return errors.New("other_text is required when selecting 'OTHER' option")
		}
	}
	if hasNoIssue && hasOtherObstacles {
		return errors.New("cannot select 'Tidak ada kendala' with other obstacles")
	}

	return nil
}

func ValidateExpertObstacles(db *gorm.DB, obstacles []entity.CareerProfileFeedbackExpertObstacle) error {
	if len(obstacles) == 0 {
		return nil
	}

	var obstacleIDs []int32
	for _, obs := range obstacles {
		obstacleIDs = append(obstacleIDs, obs.ObstacleID)
	}

	var options []entity.CareerProfileExpertObstacleOption
	if err := db.Where("id IN ?", obstacleIDs).Find(&options).Error; err != nil {
		return fmt.Errorf("failed to load obstacle options: %w", err)
	}

	optionMap := make(map[int32]*entity.CareerProfileExpertObstacleOption)
	for i := range options {
		optionMap[options[i].ID] = &options[i]
	}

	hasNoIssue := false
	hasOtherObstacles := false

	for _, obs := range obstacles {
		opt, exists := optionMap[obs.ObstacleID]
		if !exists {
			return fmt.Errorf("invalid obstacle_id: %d", obs.ObstacleID)
		}

		// Check if option is active
		if !opt.IsActive {
			return fmt.Errorf("obstacle option '%s' is no longer active", opt.Key)
		}

		if opt.IsNoIssue {
			hasNoIssue = true
		} else {
			hasOtherObstacles = true
		}

		// Validate OTHER requires text
		if opt.IsOther && (obs.OtherText == nil || *obs.OtherText == "") {
			return errors.New("other_text is required when selecting 'OTHER' option")
		}
	}

	if hasNoIssue && hasOtherObstacles {
		return errors.New("cannot select 'Tidak ada kendala' with other obstacles")
	}

	return nil
}

func ValidateTestSession(db *gorm.DB, testCategory string, testSessionID int64) error {
	switch testCategory {
	case "CAREER_PROFILE":
		var session entity.CareerProfileTestSession
		if err := db.First(&session, testSessionID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("test session not found")
			}
			return fmt.Errorf("failed to load test session: %w", err)
		}

		if session.Status != "completed" {
			return fmt.Errorf("feedback can only be submitted for completed sessions (current status: %s)", session.Status)
		}

		return nil

	default:
		return fmt.Errorf("invalid test category: %s", testCategory)
	}
}

func CheckDuplicateFeedback(db *gorm.DB, testCategory string, testSessionID int64, respondentType entity.RespondentType, respondentUserID uuid.UUID) (bool, error) {
	var count int64
	err := db.Model(&entity.KenaliDiriFeedback{}).
		Where("test_category = ? AND test_session_id = ? AND respondent_type = ? AND respondent_user_id = ? AND deleted_at IS NULL",
			testCategory, testSessionID, respondentType, respondentUserID).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check duplicate feedback: %w", err)
	}

	return count > 0, nil
}
