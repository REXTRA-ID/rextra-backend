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

	MissionPersonaCompleteResponse struct {
		UserID      string             `json:"user_id"`
		MissionKey  string             `json:"mission_key"`
		IsCompleted bool               `json:"is_completed"`
		Progress    PersonaProgress    `json:"progress"`
		Transition  *PersonaTransition `json:"transition"`
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

type PersonaTransition struct {
	From            string `json:"from"`
	To              string `json:"to"`
	Message         string `json:"message"`
	NewMissionCount int    `json:"new_mission_count"`
}
