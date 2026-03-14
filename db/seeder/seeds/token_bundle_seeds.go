package seeds

import (
	"rextra-backend/internal/entity"
	mylog "rextra-backend/internal/pkg/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func strPtr(s string) *string {
	return &s
}

func SeedTokenBundles(db *gorm.DB) error {
	mylog.Infof("[PROCESS] Seeding token bundles...")

	bundles := []entity.TokenBundlePackage{
		{Name: "Starter Pack", TokenAmount: 10, PriceRp: 15000, Label: strPtr("Hemat")},
		{Name: "Basic Pack", TokenAmount: 25, PriceRp: 35000, Label: strPtr("Terjangkau")},
		{Name: "Standard Pack", TokenAmount: 50, PriceRp: 70000, Label: strPtr("Populer")},
		{Name: "Medium Pack", TokenAmount: 100, PriceRp: 130000, Label: strPtr("Terbaik")},
		{Name: "Large Pack", TokenAmount: 200, PriceRp: 250000, Label: strPtr("Sultan")},
		{Name: "Super Pack", TokenAmount: 500, PriceRp: 600000, Label: strPtr("Pro")},
		{Name: "Ultra Pack", TokenAmount: 1000, PriceRp: 1100000, Label: strPtr("Enterprise")},
	}

	for _, b := range bundles {
		db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoUpdates: clause.AssignmentColumns([]string{"token_amount", "price_rp", "label"}),
		}).Create(&b)
	}

	return nil
}

func SeedPromos(db *gorm.DB) error {
	mylog.Infof("[PROCESS] Seeding promo codes...")

	promos := []entity.Discounts{
		{
			Code: "REXTRACLUB", Name: "Diskon Member Club", DiscountType: entity.DiscountTypeFixed, 
			Value: 10000, Status: entity.DiscountStatusActive,
		},
		{
			Code: "PROMO50", Name: "Diskon 50 Persen", DiscountType: entity.DiscountTypePercentage, 
			Value: 50, Status: entity.DiscountStatusActive,
		},
	}

	for _, p := range promos {
		db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "code"}},
			DoUpdates: clause.AssignmentColumns([]string{"value", "status"}),
		}).Create(&p)
	}

	return nil
}
