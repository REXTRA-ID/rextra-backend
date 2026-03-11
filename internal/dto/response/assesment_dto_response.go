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

type IkigaiQuestionResponse struct {
	IkigaiQuestions []datatypes.JSON `json:"ikigai_questions"`
}

type IkigaiQuestionSubmitResponse struct {
	ChartData datatypes.JSON `json:"chart_data"`
    Hash string `json:"hash"`
    Profile string `json:"profile"`
    Results datatypes.JSON `json:"results"`
    RiasecExplanations datatypes.JSON `json:"riasec_explanations"`
    RiasecMapFull datatypes.JSON `json:"riasec_map_full"`
}

type IkigaiResultResponse struct {
	ID             string   `json:"id"`
	Profile          string           `json:"profile"`
	ChartData datatypes.JSON `json:"chart_data"`
    Hash string `json:"hash"`
    Results datatypes.JSON `json:"results"`
    RiasecExplanations datatypes.JSON `json:"riasec_explanations"`
    RiasecMapFull datatypes.JSON `json:"riasec_map_full"`
	CreatedAt time.Time `json:"created_at"`
}