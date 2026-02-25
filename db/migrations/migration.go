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
		&entity.UserCareerProfile{}, 
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
		&entity.TokenBundlePackage{},
		&entity.TokenLedger{},
		&entity.TokenWallet{},
		&entity.TopupTransaction{},
		&entity.CustomPricingTier{},
		&entity.CustomPricingConfig{},

		// Jelajah Profesi Level 1
		&entity.ProfessionMainCategory{},
		&entity.Skill{},
		&entity.Tool{},
		&entity.StudyProgram{},

		// Jelajah Profesi Level 2
		&entity.ProfessionSubCategory{},

		// Jelajah Profesi Level 3
		&entity.Profession{},

		// Jelajah Profesi Level 4
		&entity.ProfessionAlias{},
		&entity.ProfessionActivity{},
		&entity.ProfessionMarketInsight{},
		&entity.ProfessionCareerPath{},

		// Jelajah Profesi Level 5
		&entity.ProfessionSkill{},
		&entity.ProfessionTool{},
		&entity.ProfessionStudyProgram{},
	); err != nil {
		return err
	}

	mylog.Infof("Migration completed successfully")

	mylog.Infof("Adding feedback indexes...")
	if err := AddFeedbackIndexes(db); err != nil {
		return err
	}

	mylog.Infof("Adding feedback constraints...")
	if err := AddFeedbackConstraints(db); err != nil {
		return err
	}

	mylog.Infof("Adding jelajah profesi constraints...")
	if err := AddJelajahProfesiConstraints(db); err != nil {
		return err
	}

	return nil
}
