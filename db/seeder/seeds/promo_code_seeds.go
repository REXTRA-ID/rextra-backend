package seeds

import (
	"encoding/json"
	"os"
	"rextra-backend/internal/entity"
	mylog "rextra-backend/internal/pkg/logger"

	"gorm.io/gorm"
)

func SeedPromoCodes(db *gorm.DB) error {
	mylog.Infof("[PROCESS] Seeding promo code...")
	jsonFile, err := os.Open("./db/seeder/data/promo_code_data.json")
	if err != nil {
		return err
	}

	defer jsonFile.Close()

	var listEntity []entity.PromoCodes
	if err := json.NewDecoder(jsonFile).Decode(&listEntity); err != nil {
		return err
	}

	for _, entity := range listEntity {
		if err := db.Save(&entity).Error; err != nil {
			return err
		}
	}

	mylog.Infof("[COMPLETE] Seeding promo code completed")
	return nil
}
