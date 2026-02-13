package seeds

import (
	"encoding/json"
	"rextra-backend/internal/entity"
	mylog "rextra-backend/internal/pkg/logger"
	"rextra-backend/internal/utils"

	"os"

	"gorm.io/gorm"
)

func SeederUser(db *gorm.DB) error {
	mylog.Infof("[PROCESS] Seeding users...")
	jsonFile, err := os.Open("./db/seeder/data/user_data.json")
	if err != nil {
		return err
	}

	defer jsonFile.Close()

	var listEntity []entity.User
	if err := json.NewDecoder(jsonFile).Decode(&listEntity); err != nil {
		return err
	}

	for _, entity := range listEntity {
		var count int64
		if err := db.Model(&entity).Where("email = ?", entity.Email).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			mylog.Infof("[SKIP] User %s already exists", entity.Email)
			continue
		}

		hashedPwd, _ := utils.HashPassword(entity.Password)
		entity.Password = hashedPwd
		if err := db.Save(&entity).Error; err != nil {
			return err
		}
	}

	mylog.Infof("[COMPLETE] Seeding users completed")
	return nil
}
