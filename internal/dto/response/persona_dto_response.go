package dto_response

type (
	CreatePersonaResponse struct {
		ID             string `json:"id"`
		UserID         string `json:"user_id"`
		Institution    string `json:"institution"`
		Study          string `json:"study"`
		EducationLevel string `json:"education_level"`
		GraduationYear int    `json:"graduation_year"`
		CareerPlan     string `json:"career_plan"`
		CareerDreams   string `json:"career_dreams"`
		Portfolio      bool   `json:"portfolio"`
		Application    bool   `json:"application"`
		Status         string `json:"status"`
	}
)
