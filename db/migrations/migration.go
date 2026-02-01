package migrations

import (
	"fmt"
	"rextra-backend/internal/entity"
	mylog "rextra-backend/internal/pkg/logger"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	fmt.Println(mylog.ColorizeInfo("\n=========== Start Migrate ==========="))
	mylog.Infof("Migrating Tables...")

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		return err
	}

	//migrate table
	if err := db.AutoMigrate(
		&entity.User{},
		&entity.KenaliDiriCategory{},
		&entity.KenaliDiriHistory{},
		&entity.SessionToken{},
		&entity.Persona{},
		// &entity.Riasec{},
		&entity.RiasecCode{},
		&entity.CareerProfileTestSession{},
		&entity.UserCareerProfile{}, // NEW: User active profile pointer
		&entity.RiasecQuestionSet{},
		&entity.RiasecResponse{},
		&entity.RiasecResult{},
		&entity.IkigaiCandidateProfession{},
		&entity.IkigaiResponse{},
		&entity.IkigaiDimensionScore{},
		&entity.IkigaiTotalScore{},
		
		&entity.StudentFeedback{},
		&entity.ExpertFeedback{},
		
		&entity.KenaliDiriFeedback{},                         
		
		&entity.CareerProfileFeedbackStudent{},               
		&entity.CareerProfileObstacleOption{},                
		&entity.CareerProfileFeedbackObstacle{},              
		
		&entity.CareerProfileFeedbackExpert{},                
		&entity.CareerProfileExpertObstacleOption{},          
		&entity.CareerProfileFeedbackExpertObstacle{},        
		
		&entity.CareerRecommendation{},
		// &entity.UserIkigai{},
		// &entity.UserRiasec{},
	); err != nil {
		return err
	}

	mylog.Infof("Migration completed successfully")

	mylog.Infof("Adding feedback indexes...")
	if err := AddFeedbackIndexes(db); err != nil {
		return err
	}

	return nil
}
