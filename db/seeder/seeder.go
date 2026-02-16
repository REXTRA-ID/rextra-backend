package seeders

import (
	"fmt"
	"rextra-backend/db/seeder/seeds"
	mylog "rextra-backend/internal/pkg/logger"

	"gorm.io/gorm"
)

func Seeding(db *gorm.DB) error {
	seeders := []func(*gorm.DB) error{
		// ── Urutan seeding untuk fresh database ──
		// 1. Master data tanpa dependency
		seeds.SeederUser,
		seeds.SeederRiasecCodes,
		seeds.SeedTokenBundlePackages,
		seeds.SeedStudentObstacleOptions,
		seeds.SeedExpertObstacleOptions,
		// 2. KenaliDiri (depends: User) — TRUNCATE + recreate
		seeds.SeederKenaliDiri,
		// 3. Remaining entities (depends: User, KenaliDiri sessions)
		seeds.SeedRemainingEntities,
		// 4. Additional career profile data (depends: User, RiasecCodes)
		seeds.SeedCareerProfileData,
		// 5. Feedback data (depends: User, sessions, obstacle options)
		seeds.SeedFeedbackData,
		// 6. Topup transactions (depends: User, TokenBundlePackage)
		seeds.SeedTopupTransactions,
	}

	fmt.Println(mylog.ColorizeInfo("\n=========== Start Seeding ==========="))
	for _, seeder := range seeders {
		if err := seeder(db); err != nil {
			return err
		}
	}

	return nil
}
