package seeds

import (
	"encoding/json"
	"os"

	"rextra-backend/internal/entity"
	mylog "rextra-backend/internal/pkg/logger"

	"gorm.io/gorm"
)

func SeederMembershipDuration(db *gorm.DB) error {
	mylog.Infof("[PROCESS] Seeding membership_duration...")

	jsonFile, err := os.Open("./db/seeder/data/membership_duration_data.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	var listEntity []entity.MembershipDuration
	if err := json.NewDecoder(jsonFile).Decode(&listEntity); err != nil {
		return err
	}

	for _, item := range listEntity {
		if err := db.Create(&item).Error; err != nil {
			return err
		}
	}

	mylog.Infof("[COMPLETE] Seeding membership_duration completed")
	return nil
}
