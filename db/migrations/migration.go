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

	if err := db.AutoMigrate(
		&entity.User{},
		&entity.SessionToken{},
		&entity.Persona{},
		&entity.Riasec{},
		&entity.CareerRecommendation{},
		&entity.UserIkigai{},
		&entity.UserRiasec{},
		&entity.Voucher{},
	); err != nil {
		return err
	}

	SeedVouchers(db)

	return nil
}

func SeedVouchers(db *gorm.DB) {
	var count int64
	db.Model(&entity.Voucher{}).Count(&count)
	
	if count == 0 {
		mylog.Infof("Database empty for Vouchers! Running Auto-Seeder...")
		vouchers := []entity.Voucher{
			{Code: "REXTRA-TEST1", IsUsed: false},
			{Code: "REXTRA-TEST2", IsUsed: false},
			{Code: "REXTRA-TEST3", IsUsed: false},
			{Code: "REXTRA-TEST4", IsUsed: false},
			{Code: "REXTRA-TEST5", IsUsed: false},
		}
		db.Create(&vouchers)
		mylog.Infof("Successfully seeded 5 testing vouchers!")
	}
}
