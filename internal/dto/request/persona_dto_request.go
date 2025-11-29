package dto_request

type (
	CreatePersonaRequest struct {
		UserID               string `json:"user_id" binding:"required"`
		HasCareerGoal        bool   `json:"has_career_goal" binding:"required"`
		BuildingPortofolio   bool   `json:"building_portofolio" binding:"required"`
		InRecruitmentProcess bool   `json:"in_recruitment_process" binding:"required"`
	}
	MissionPersonaCompleteRequest struct {
		UserID     string `json:"user_id" binding:"required"`
		MissionKey string `json:"mission_key" binding:"required"`
	}
)
