package migrations

import (
	mylog "rextra-backend/internal/pkg/logger"

	"gorm.io/gorm"
)

func AddFeedbackIndexes(db *gorm.DB) error {
	mylog.Infof("[PROCESS] Adding feedback indexes...")

	indexes := []string{
		`CREATE UNIQUE INDEX IF NOT EXISTS uq_feedback_once
		 ON kenalidiri_feedback (test_category, test_session_id, respondent_type, respondent_user_id)
		 WHERE deleted_at IS NULL`,

		`CREATE INDEX IF NOT EXISTS idx_feedback_category_time 
		 ON kenalidiri_feedback(test_category, submitted_at)`,

		`CREATE INDEX IF NOT EXISTS idx_feedback_category_session 
		 ON kenalidiri_feedback(test_category, test_session_id)`,

		`CREATE INDEX IF NOT EXISTS idx_feedback_category_role_time 
		 ON kenalidiri_feedback(test_category, respondent_type, submitted_at)`,

		`CREATE INDEX IF NOT EXISTS idx_student_obstacles_agg 
		 ON careerprofile_feedback_obstacles(obstacle_id)`,

		`CREATE INDEX IF NOT EXISTS idx_expert_obstacles_agg 
		 ON careerprofile_feedback_expert_obstacles(obstacle_id)`,

		`CREATE INDEX IF NOT EXISTS idx_expert_obstacles_feedback 
		 ON careerprofile_feedback_expert_obstacles(feedback_id)`,

		`CREATE INDEX IF NOT EXISTS idx_expert_top5_status 
		 ON careerprofile_feedback_expert(top5_status)`,

		`CREATE INDEX IF NOT EXISTS idx_expert_profession 
		 ON careerprofile_feedback_expert(expert_profession_id) 
		 WHERE expert_profession_id IS NOT NULL`,
	}

	for i, idx := range indexes {
		mylog.Infof("[%d/%d] Creating index...", i+1, len(indexes))
		if err := db.Exec(idx).Error; err != nil {
			mylog.Errorf("[ERROR] Failed to create index: %v", err)
			return err
		}
	}

	mylog.Infof("[SUCCESS] Created %d indexes for feedback tables", len(indexes))
	return nil
}
