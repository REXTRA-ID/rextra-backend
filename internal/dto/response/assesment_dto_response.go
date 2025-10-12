package dto_response

import (
	"time"

	"gorm.io/datatypes"
)

type RiasecQuestionResponse struct {
	ID             string   `json:"id"`
	Pertanyaan     string   `json:"pertanyaan"`
	Skala          string   `json:"skala"`
	TargetKategori []string `json:"target_kategori"`
	Type           string   `json:"type"`
}

type RiasecQuestionSubmitResponse struct {
	NormalizedScores datatypes.JSON `json:"normalized_scores"`
	Profile          string           `json:"profile"`
}

type RiasecResultResponse struct {
	ID            	 string   `json:"id"`
	Profile          string           `json:"profile"`
	NormalizedScores datatypes.JSON `json:"normalized_scores"`
	CreatedAt    	time.Time `json:"created_at"`	
}
