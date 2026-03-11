package seeds

import (
	"encoding/json"
	"rextra-backend/internal/entity"
	mylog "rextra-backend/internal/pkg/logger"
	"rextra-backend/internal/utils"

	"os"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	for i := range listEntity {
		hashedPwd, _ := utils.HashPassword(listEntity[i].Password)
		listEntity[i].Password = hashedPwd
	}

	err = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{"fullname", "phone_number", "role", "is_verified", "updated_at"}),
	}).Create(&listEntity).Error

	if err != nil {
		return err
	}

	mylog.Infof("[COMPLETE] Seeding users completed")
	return nil
}
