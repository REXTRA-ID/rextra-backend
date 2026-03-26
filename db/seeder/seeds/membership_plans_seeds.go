package seeds

import (
	"encoding/json"
	"rextra-backend/internal/entity"
	mylog "rextra-backend/internal/pkg/logger"

	"os"

	"gorm.io/gorm"
)

func SeederMembershipPlan(db *gorm.DB) error {
	mylog.Infof("[PROCESS] Seeding membership plans...")

	jsonFile, err := os.Open("./db/seeder/data/membership_plans_data.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	var listEntity []entity.MembershipPlans
	if err := json.NewDecoder(jsonFile).Decode(&listEntity); err != nil {
		return err
	}

	for _, item := range listEntity {
		if err := db.Create(&item).Error; err != nil {
			return err
		}
	}

	mylog.Infof("[COMPLETE] Seeding membership plans completed")
	return nil
}
