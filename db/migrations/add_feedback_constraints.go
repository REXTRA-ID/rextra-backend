package migrations

import (
	"fmt"
	mylog "rextra-backend/internal/pkg/logger"

	"gorm.io/gorm"
)

func AddFeedbackConstraints(db *gorm.DB) error {
	mylog.Infof("[PROCESS] Adding feedback constraints...")

	constraints := []struct {
		Table      string
		Constraint string
		Check      string
	}{
		// Student Feedback Constraints
		{
			Table:      "careerprofile_feedback_student",
			Constraint: "chk_student_ease_score",
			Check:      "ease_score BETWEEN 1 AND 7",
		},
		{
			Table:      "careerprofile_feedback_student",
			Constraint: "chk_student_relevance_score",
			Check:      "relevance_score BETWEEN 1 AND 7",
		},
		{
			Table:      "careerprofile_feedback_student",
			Constraint: "chk_student_satisfaction_score",
			Check:      "satisfaction_score BETWEEN 1 AND 7",
		},

		// Expert Feedback Constraints
		{
			Table:      "careerprofile_feedback_expert",
			Constraint: "chk_expert_accuracy_score",
			Check:      "accuracy_score BETWEEN 1 AND 7",
		},
		{
			Table:      "careerprofile_feedback_expert",
			Constraint: "chk_expert_logic_score",
			Check:      "logic_score BETWEEN 1 AND 7",
		},
		{
			Table:      "careerprofile_feedback_expert",
			Constraint: "chk_expert_usefulness_score",
			Check:      "usefulness_score BETWEEN 1 AND 7",
		},
	}

	for _, c := range constraints {
		// PostgreSQL-specific block to safely add constraint only if it doesn't exist
		query := fmt.Sprintf(`
			DO $$
			BEGIN
				IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = '%s') THEN
					ALTER TABLE %s ADD CONSTRAINT %s CHECK (%s);
				END IF;
			END
			$$;
		`, c.Constraint, c.Table, c.Constraint, c.Check)

		if err := db.Exec(query).Error; err != nil {
			mylog.Errorf("[ERROR] Failed to add constraint %s: %v", c.Constraint, err)
			return err
		}
	}

	mylog.Infof("[SUCCESS] Verified/Added %d constraints for feedback tables", len(constraints))
	return nil
}
