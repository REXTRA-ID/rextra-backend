package seeds

import (
	"encoding/json"
	"os"
	"rextra-backend/internal/entity"
	mylog "rextra-backend/internal/pkg/logger"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type tempTopupTransaction struct {
	ID              uuid.UUID          `json:"id"`
	UserEmail       string             `json:"user_email"`
	UserID          uuid.UUID          `json:"user_id"`
	Type            entity.TopupType   `json:"type"`
	BundlePackageID *uuid.UUID         `json:"bundle_package_id"`
	BundleName      *string            `json:"bundle_name"`
	TokenAmount     int64              `json:"token_amount"`
	TotalPriceRp    int64              `json:"total_price_rp"`
	Status          entity.TopupStatus `json:"status"`
	InvoiceID       string             `json:"invoice_id"`
	Provider        *string            `json:"provider"`
	PaidAt          *time.Time         `json:"paid_at"`
	ExpiredAt       *time.Time         `json:"expired_at"`
	LedgerID        *uuid.UUID         `json:"ledger_id"`
	Metadata        json.RawMessage    `json:"metadata"`
}

func SeedTopupTransactions(db *gorm.DB) error {
	mylog.Infof("[PROCESS] Seeding topup transactions...")
	jsonFile, err := os.Open("./db/seeder/data/topup_transaction_data.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	var tempSeeds []tempTopupTransaction
	if err := json.NewDecoder(jsonFile).Decode(&tempSeeds); err != nil {
		mylog.Errorf("Failed to decode topup_transaction_data.json: %v", err)
		return err
	}

	var mapsToCreate []map[string]interface{}
	for _, seed := range tempSeeds {
		// Resolve user_id from email if provided
		userID := seed.UserID
		if seed.UserEmail != "" {
			var user entity.User
			if err := db.Where("email = ?", seed.UserEmail).First(&user).Error; err != nil {
				mylog.Errorf("User with email %s not found, skipping transaction %s", seed.UserEmail, seed.InvoiceID)
				continue
			}
			userID = user.ID
		}

		// Resolve bundle_package_id from name if provided
		var bundlePackageID *uuid.UUID
		if seed.BundleName != nil && *seed.BundleName != "" {
			var pkg entity.TokenBundlePackage
			if err := db.Where("name = ?", *seed.BundleName).First(&pkg).Error; err != nil {
				mylog.Errorf("Bundle package %s not found, setting to nil", *seed.BundleName)
			} else {
				bundlePackageID = &pkg.ID
			}
		} else {
			bundlePackageID = seed.BundlePackageID
		}

		mapsToCreate = append(mapsToCreate, map[string]interface{}{
			"id":                seed.ID,
			"user_id":           userID,
			"type":              seed.Type,
			"bundle_package_id": bundlePackageID,
			"token_amount":      seed.TokenAmount,
			"total_price_rp":    seed.TotalPriceRp,
			"status":            seed.Status,
			"invoice_id":        seed.InvoiceID,
			"provider":          seed.Provider,
			"paid_at":           seed.PaidAt,
			"expired_at":        seed.ExpiredAt,
			"ledger_id":         seed.LedgerID,
			"metadata":          string(seed.Metadata),
		})
	}

	if len(mapsToCreate) == 0 {
		mylog.Infof("No topup transactions to seed.")
		return nil
	}

	err = db.Model(&entity.TopupTransaction{}).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"user_id",
			"type",
			"bundle_package_id",
			"token_amount",
			"total_price_rp",
			"status",
			"invoice_id",
			"provider",
			"paid_at",
			"expired_at",
			"ledger_id",
			"metadata",
			"updated_at",
		}),
	}).Create(&mapsToCreate).Error

	if err != nil {
		return err
	}

	mylog.Infof("[COMPLETE] Seeding topup transactions completed")
	return nil
}
