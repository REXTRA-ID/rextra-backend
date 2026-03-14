package seeds

import (
	"rextra-backend/internal/entity"
	"rextra-backend/internal/utils"
	mylog "rextra-backend/internal/pkg/logger"
	"gorm.io/gorm"
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
		// Hapus dulu biar ga conflict, baru insert
		db.Unscoped().Where("email = ?", user.Email).Delete(&entity.User{})
		if err := db.Create(&user).Error; err != nil {
			return err
		}
	}

	mylog.Infof("[COMPLETE] Seeding users completed")
	return nil
}
