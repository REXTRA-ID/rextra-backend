package seeds

import (
	"encoding/json"
	"os"
	"path/filepath"

	"rextra-backend/internal/entity"
	mylog "rextra-backend/internal/pkg/logger"

	"gorm.io/gorm"
)

func SeedExpertObstacleOptions(db *gorm.DB) error {
	mylog.Infof("[PROCESS] Seeding Expert Obstacle Options...")

	var count int64
	if err := db.Model(&entity.CareerProfileExpertObstacleOption{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		mylog.Infof("[SKIP] Expert obstacle options already seeded (%d records)", count)
		return nil
	}

	dataPath := filepath.Join("db", "seeder", "data", "expert_obstacles.json")
	fileData, err := os.ReadFile(dataPath)
	if err != nil {
		mylog.Errorf("[ERROR] Failed to read expert obstacles JSON: %v", err)
		return err
	}

	var options []entity.CareerProfileExpertObstacleOption
	if err := json.Unmarshal(fileData, &options); err != nil {
		mylog.Errorf("[ERROR] Failed to parse expert obstacles JSON: %v", err)
		return err
	}

	for _, opt := range options {
		if err := db.Create(&opt).Error; err != nil {
			mylog.Errorf("[ERROR] Failed to seed expert obstacle option: %v", err)
			return err
		}
	}

	mylog.Infof("[SUCCESS] Seeded %d expert obstacle options", len(options))
	return nil
}
