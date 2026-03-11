package seeds

import (
	"encoding/json"
	"os"
	"rextra-backend/internal/entity"
	mylog "rextra-backend/internal/pkg/logger"

	"gorm.io/gorm"
)

func SeedCareerProfileData(db *gorm.DB) error {
	mylog.Infof("[PROCESS] Seeding Career Profile Data (Admin Verification)...")

	// Read JSON Data
	jsonFile, err := os.Open("./db/seeder/data/career_profile_data.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	type SeederData struct {
		Sessions      []entity.CareerProfileTestSession `json:"sessions"`
		RiasecResults []entity.RiasecResult             `json:"riasec_results"`
		UserProfile   entity.UserCareerProfile          `json:"user_profile"`
	}

	var data SeederData
	if err := json.NewDecoder(jsonFile).Decode(&data); err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// Get User
		var user entity.User
		if err := tx.Where("email = ?", "user@email.com").First(&user).Error; err != nil {
			return err
		}
		userID := user.ID

		// Ensure Categories Exist (Safeguard)
		var cat3 entity.KenaliDiriCategory
		if err := tx.FirstOrCreate(&cat3, entity.KenaliDiriCategory{
			ID:              3,
			CategoryCode:    "CAREER_PROFILE",
			CategoryName:    "Tes Profil Karier",
			Description:     ptrString("Tes profil karier lengkap"),
			DetailTableName: "careerprofile_test_sessions",
			IsActive:        true,
		}).Error; err != nil {
			return err
		}

		// Create Sessions
		for _, session := range data.Sessions {
			session.UserID = userID // Assign dynamic UserID
			if err := tx.Create(&session).Error; err != nil {
				return err
			}
		}

		// Create Results
		for _, result := range data.RiasecResults {
			if err := tx.Create(&result).Error; err != nil {
				return err
			}
		}

		// Set User Career Profile
		userProfile := data.UserProfile
		userProfile.UserID = userID

		var existingProfile entity.UserCareerProfile
		err := tx.Where("user_id = ?", userID).First(&existingProfile).Error
		if err == nil {
			// Update existing
			existingProfile.ActiveSessionID = userProfile.ActiveSessionID
			existingProfile.SetSource = userProfile.SetSource
			existingProfile.SetAt = userProfile.SetAt
			existingProfile.Pinned = userProfile.Pinned
			if err := tx.Save(&existingProfile).Error; err != nil {
				return err
			}
		} else {
			// Create new
			if err := tx.Create(&userProfile).Error; err != nil {
				return err
			}
		}

		mylog.Infof("[COMPLETE] Seeding Career Profile Data completed.")
		return nil
	})
}
