package seeds

import (
	"encoding/json"
	"rextra-backend/internal/entity"
	mylog "rextra-backend/internal/pkg/logger"

	"os"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		// Use OnConflict to avoid duplicate key error
		err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "plan_name"}},
			DoUpdates: clause.AssignmentColumns([]string{"category", "tier_label", "description", "status", "pricing_mode", "duration_mode", "base_price1_m", "base_token1_m"}),
		}).Create(&item).Error
		if err != nil {
			return err
		}
	}

	mylog.Infof("[COMPLETE] Seeding membership plans completed")
	return nil
}
