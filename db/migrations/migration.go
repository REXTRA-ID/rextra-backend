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
		&entity.SessionToken{},
		&entity.Persona{},
		&entity.Riasec{},
		&entity.CareerRecommendation{},
		&entity.Memberships{},
		&entity.RedemptionCode{},
		&entity.MembershipPlans{},
		&entity.MembershipDuration{},
		&entity.PaymentTransactions{},
		&entity.PoinTransactions{},
		&entity.TokenTransaction{},
		&entity.TokenUsageHistory{},
		&entity.PromoCodes{},
		&entity.PromoCodeUsage{},
	); err != nil {
		return err
	}

	return nil
}
