package dto_response

type RiasecQuestionResponse struct {
	ID             string   `json:"id"`
	Pertanyaan     string   `json:"pertanyaan"`
	Skala          string   `json:"skala"`
	TargetKategori []string `json:"target_kategori"`
	Type           string   `json:"type"`
}

type normalizedScores struct {
	Artistic      float64 `json:"Artistic"`
	Conventional  float64 `json:"Conventional"`
	Enterprising  float64 `json:"Enterprising"`
	Investigative float64 `json:"Investigative"`
	Realistic     float64 `json:"Realistic"`
	Social        float64 `json:"Social"`
}

type RiasecQuestionSubmitResponse struct {
	NormalizedScores normalizedScores `json:"normalized_scores"`
	Profile          string           `json:"profile"`
}
