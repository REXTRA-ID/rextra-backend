package seeders

import (
	"fmt"
	"rextra-backend/db/seeder/seeds"
	mylog "rextra-backend/internal/pkg/logger"

	"gorm.io/gorm"
)

func Seeding(db *gorm.DB) error {
	seeders := []func(*gorm.DB) error{
		// seeds.SeederUser,
		// seeds.SeederRiasecCodes,
		// seeds.SeedTokenBundlePackages,
		seeds.SeedTopupTransactions,
		// seeds.SeederKenaliDiri,
    // seeds.SeedCareerProfileData,
		// seeds.SeedFeedbackData,
    // seeds.SeedStudentObstacleOptions,
		// seeds.SeedExpertObstacleOptions,
	}

	fmt.Println(mylog.ColorizeInfo("\n=========== Start Seeding ==========="))
	for _, seeder := range seeders {
		if err := seeder(db); err != nil {
			return err
		}
	}

	return nil
}
