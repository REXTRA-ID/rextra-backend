package seeders

import (
	"fmt"
	"rextra-backend/db/seeder/seeds"
	mylog "rextra-backend/internal/pkg/logger"

	"gorm.io/gorm"
)

func Seeding(db *gorm.DB) error {
	seeders := []func(*gorm.DB) error{
		seeds.SeederUser,
		seeds.SeederMembershipPlan,
		seeds.SeedMembershipDurations,
		seeds.SeedPromos,
		seeds.SeedTokenBundles,
	}

	fmt.Println(mylog.ColorizeInfo("\n=========== Start Seeding ==========="))
	for _, seeder := range seeders {
		if err := seeder(db); err != nil {
			return err
		}
	}

	return nil
}
