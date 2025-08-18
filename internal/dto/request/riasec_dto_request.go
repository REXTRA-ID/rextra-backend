package dto_request

type (
	CreateRiasecRequest struct {
		UserID       string                       `json:"user_id"`
		RiasecCode   string                       `json:"riasec_code" binding:"required"`
		Explanations RiasecCodeExplanationRequest `json:"explanations"`
	}

	RiasecCodeExplanationRequest struct {
		R *string `json:"R,omitempty"`
		I *string `json:"I,omitempty"`
		A *string `json:"A,omitempty"`
		S *string `json:"S,omitempty"`
		E *string `json:"E,omitempty"`
		C *string `json:"C,omitempty"`
	}
)
