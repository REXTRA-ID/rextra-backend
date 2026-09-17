package dto_response



type RiasecQuestionResponse struct {
	ID             string   `json:"id"`
	Pertanyaan     string   `json:"pertanyaan"`
	Skala          string   `json:"skala"`
	TargetKategori []string `json:"target_kategori"`
	Type           string   `json:"type"`
}

type RiasecQuestionSubmitResponse map[string]interface{}

type RiasecResultResponse map[string]interface{}

type IkigaiQuestionResponse map[string]interface{}

type IkigaiQuestionSubmitResponse map[string]interface{}

type IkigaiResultResponse map[string]interface{}