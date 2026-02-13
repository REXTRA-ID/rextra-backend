package seeds

import (
	"encoding/json"
	"os"
	"path/filepath"

	"rextra-backend/internal/entity"
	mylog "rextra-backend/internal/pkg/logger"

	"gorm.io/gorm"
)

func SeedStudentObstacleOptions(db *gorm.DB) error {
	mylog.Infof("[PROCESS] Seeding Student Obstacle Options...")

	var count int64
	if err := db.Model(&entity.CareerProfileObstacleOption{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		mylog.Infof("[SKIP] Student obstacle options already seeded (%d records)", count)
		return nil
	}

	dataPath := filepath.Join("db", "seeder", "data", "student_obstacles.json")
	fileData, err := os.ReadFile(dataPath)
	if err != nil {
		mylog.Errorf("[ERROR] Failed to read student obstacles JSON: %v", err)
		return err
	}

	var options []entity.CareerProfileObstacleOption
	if err := json.Unmarshal(fileData, &options); err != nil {
		mylog.Errorf("[ERROR] Failed to parse student obstacles JSON: %v", err)
		return err
	}

	for _, opt := range options {
		if err := db.Create(&opt).Error; err != nil {
			mylog.Errorf("[ERROR] Failed to seed student obstacle option: %v", err)
			return err
		}
	}

	mylog.Infof("[SUCCESS] Seeded %d student obstacle options", len(options))
	return nil
}
