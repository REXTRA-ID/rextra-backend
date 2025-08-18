package dto_response

type (
	CreateRiasecResponse struct {
		ID           string                        `json:"id"`
		RiasecCode   string                        `json:"riasec_code"`
		Explanations RiasecCodeExplanationResponse `json:"explanations"`
	}

	RiasecCodeExplanationResponse struct {
		R *string `json:"R,omitempty"`
		I *string `json:"I,omitempty"`
		A *string `json:"A,omitempty"`
		S *string `json:"S,omitempty"`
		E *string `json:"E,omitempty"`
		C *string `json:"C,omitempty"`
	}
)
