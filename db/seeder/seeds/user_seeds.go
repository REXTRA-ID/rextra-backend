package seeds

import (
	"rextra-backend/internal/entity"
	"rextra-backend/internal/utils"
	mylog "rextra-backend/internal/pkg/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeederUser(db *gorm.DB) error {
	mylog.Infof("[PROCESS] Seeding users...")

	password, _ := utils.HashPassword("admin123")
	
	users := []entity.User{
		{
			Fullname:    "Admin Rextra",
			Email:       "admin@email.com",
			Password:    password,
			IsVerified:  true,
			PhoneNumber: "0888888888",
			Role:        "ADMIN",
		},
		{
			Fullname:    "User Tester",
			Email:       "user@email.com",
			Password:    password,
			IsVerified:  true,
			PhoneNumber: "0812345678",
			Role:        "USER",
		},
	}

	for _, user := range users {
		db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "email"}},
			DoUpdates: clause.AssignmentColumns([]string{"fullname", "password", "is_verified", "phone_number", "role"}),
		}).Create(&user)
	}

	mylog.Infof("[COMPLETE] Seeding users completed")
	return nil
}
