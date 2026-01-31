package seeds

import (
	"encoding/json"
	"os"
	"rextra-backend/internal/entity"
	mylog "rextra-backend/internal/pkg/logger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedTokenBundlePackages(db *gorm.DB) error {
	jsonFile, err := os.Open("db/seeder/data/token_bundle_package_data.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	var packages []entity.TokenBundlePackage
	if err := json.NewDecoder(jsonFile).Decode(&packages); err != nil {
		return err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}}, // Sesuai unique index di struct
			DoNothing: true,                            // Jika sudah ada, jangan timpa
		}).Create(&packages).Error
	})

	if err != nil {
		return err
	}

	mylog.Infof("[COMPLETE] Seeding token bundle packages completed")
	return nil
}
