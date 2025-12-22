package seeds

import (
	"encoding/json"
	"errors"
	"os"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"rextra-backend/internal/entity"
	mylog "rextra-backend/internal/pkg/logger"
)

const riasecCodesJSONPath = "./db/seeder/data/riasec_codes.json"

type riasecCodeSeed struct {
	ID                 int64    `json:"id"`
	RiasecCode         string   `json:"riasec_code"`
	RiasecTitle        string   `json:"riasec_title"`
	RiasecDescription  *string  `json:"riasec_description"`
	Strengths          []string `json:"strengths"`
	Challenges         []string `json:"challenges"`
	Strategies         []string `json:"strategies"`
	WorkEnvironments   []string `json:"work_environments"`
	InteractionStyles  []string `json:"interaction_styles"`
}

func SeederRiasecCodes(db *gorm.DB) error {
	mylog.Infof("[PROCESS] Seeding riasec_codes from JSON...")

	jsonBytes, err := os.ReadFile(riasecCodesJSONPath)
	if err != nil {
		return err
	}

	var seeds []riasecCodeSeed
	if err := json.Unmarshal(jsonBytes, &seeds); err != nil {
		return err
	}

	if len(seeds) == 0 {
		return errors.New("riasec_codes.json is empty; please provide seed data")
	}

	if err := db.Exec("TRUNCATE TABLE riasec_codes RESTART IDENTITY CASCADE;").Error; err != nil {
		return err
	}

	for _, seed := range seeds {
		strengthsJSON, err := json.Marshal(seed.Strengths)
		if err != nil {
			return err
		}
		challengesJSON, err := json.Marshal(seed.Challenges)
		if err != nil {
			return err
		}
		strategiesJSON, err := json.Marshal(seed.Strategies)
		if err != nil {
			return err
		}
		workEnvJSON, err := json.Marshal(seed.WorkEnvironments)
		if err != nil {
			return err
		}
		interactionJSON, err := json.Marshal(seed.InteractionStyles)
		if err != nil {
			return err
		}

		record := entity.RiasecCode{
			ID:                seed.ID,
			RiasecCode:        seed.RiasecCode,
			RiasecTitle:       seed.RiasecTitle,
			RiasecDescription: seed.RiasecDescription,
			Strengths:         datatypes.JSON(strengthsJSON),
			Challenges:        datatypes.JSON(challengesJSON),
			Strategies:        datatypes.JSON(strategiesJSON),
			WorkEnvironments:  datatypes.JSON(workEnvJSON),
			InteractionStyles: datatypes.JSON(interactionJSON),
		}

		if err := db.Create(&record).Error; err != nil {
			return err
		}
	}

	mylog.Infof("[COMPLETE] Seeding riasec_codes completed")
	return nil
}
