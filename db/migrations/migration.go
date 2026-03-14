package migrations

import (
	"fmt"
	"rextra-backend/internal/entity"
	mylog "rextra-backend/internal/pkg/logger"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	fmt.Println(mylog.ColorizeInfo("\n=========== Start Migrate ==========="))
	mylog.Infof("Migrating Tables...")

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		return err
	}

	// Migrate tables in order to handle dependencies
	if err := db.AutoMigrate(
		// 1. Base Data
		&entity.User{},
		&entity.SessionToken{},
		
		// 2. Settings (Singleton)
		&entity.InvoiceSettings{},
		&entity.TransactionIdSettings{},
		&entity.MembershipNotificationSettings{},

		// 3. Master Data (Fitur & Akses)
		&entity.Feature{},
		&entity.SubFeature{},
		&entity.ActionCategory{},
		&entity.Entitlement{},

		// 4. Membership Plans & Config
		&entity.MembershipPlans{},
		&entity.PlanDuration{},
		&entity.DurationAccessMapping{},

		// 5. User Membership & Wallet
		&entity.Memberships{},
		&entity.TokenWallet{},
		&entity.TokenLedger{},
		&entity.PointsLedger{},

		// 6. Transactions & Promo
		&entity.Discounts{},
		&entity.DiscountRedemption{},
		&entity.PaymentTransactions{},
		&entity.TopupTransaction{},
		&entity.SubscriptionCycle{},
		&entity.PoinTransactions{},

		// 7. Logs & Quotas
		&entity.UsageLog{},
		&entity.TokenUsageHistory{},
		&entity.UserEntitlementQuota{},

		// 8. Assessment
		&entity.Persona{},
		&entity.Riasec{},
		&entity.CareerRecommendation{},
		&entity.FitCheckResults{},
		&entity.IkigaiTotalScores{},
		&entity.UserCareerProfile{},
	); err != nil {
		return err
	}

	mylog.Infof("Migration completed successfully")
	return nil
}
