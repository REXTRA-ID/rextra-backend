package dto_response

type (
	CreatePersonaResponse struct {
		ID                 string           `json:"id"`
		UserID             string           `json:"user_id"`
		PersonaType        string           `json:"persona_type"`
		PersonaDescription string           `json:"persona_description"`
		PersonaStartedAt   string           `json:"persona_started_at"`
		Progress           PersonaProgress  `json:"progress"`
		Missions           []PersonaMission `json:"missions"`
	}

	GetPersonaResponse struct {
		ID                 string           `json:"id"`
		UserID             string           `json:"user_id"`
		PersonaType        string           `json:"persona_type"`
		PersonaDescription string           `json:"persona_description"`
		PersonaStartedAt   string           `json:"persona_started_at"`
		Progress           PersonaProgress  `json:"progress"`
		Missions           []PersonaMission `json:"missions"`
	}
)

type PersonaProgress struct {
	CompletedMissions int     `json:"completed_missions"`
	TotalMissions     int     `json:"total_missions"`
	Percentage        float64 `json:"percentage"`
	CurrentPhase      string  `json:"current_phase"`
	NextPersona       string  `json:"next_persona"`
}

type PersonaMission struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
	Order       int    `json:"order"`
}
