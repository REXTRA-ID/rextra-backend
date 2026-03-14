package seeds

import (
	"rextra-backend/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedMembershipDurations(db *gorm.DB) error {
	var plans []entity.MembershipPlans
	db.Find(&plans)

	for _, plan := range plans {
		// Durations: 1, 3, 6, 12 months
		durations := []int{1, 3, 6, 12}
		for _, months := range durations {
			// Basic Price Logic
			basePrice := plan.BasePrice1M
			if basePrice == 0 {
				basePrice = 15000 // Default value if 0
			}
			
			subtotal := basePrice * int64(months)
			tokenAmount := int64(plan.BaseToken1M) * int64(months)

			// Give discount for longer durations
			discount := 0.0
			if months == 3 { discount = 5.0 }
			if months == 6 { discount = 10.0 }
			if months == 12 { discount = 20.0 }

			finalPrice := subtotal
			if discount > 0 {
				finalPrice = int64(float64(subtotal) * (1 - discount/100))
			}

			dur := entity.PlanDuration{
				PlanID:         plan.ID,
				DurationMonths: int(months),
				Price:          subtotal,
				DiscountPct:    discount,
				FinalPrice:     finalPrice,
				DurationPrice:  finalPrice,
				TokenAmount:    int(tokenAmount),
				IsActive:       true,
			}

			// Use OnConflict to avoid duplicate key error
			db.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "plan_id"}, {Name: "duration_months"}},
				DoUpdates: clause.AssignmentColumns([]string{"price", "discount_pct", "final_price", "duration_price", "token_amount"}),
			}).Create(&dur)
		}
	}
	return nil
}
