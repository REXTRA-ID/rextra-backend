package dto_request

type (
	CreatePersonaRequest struct {
		UserID               string `json:"user_id"`
		HasCareerGoal        bool   `json:"has_career_goal"`
		BuildingPortofolio   bool   `json:"building_portofolio"`
		InRecruitmentProcess bool   `json:"in_recruitment_process"`
	}
)
