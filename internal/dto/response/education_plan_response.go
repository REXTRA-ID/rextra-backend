package dto_response

type EducationPlanResponse struct {
	ID                string `json:"id"`
	UserID            string `json:"user_id"`
	InstitutionName   string `json:"institution_name"`
	Major             string `json:"major"`
	Faculty           string `json:"faculty"`
	ExpectedEntryYear int    `json:"expected_entry_year"`
	EducationProgram  string `json:"education_program"`
	EducationLevel    string `json:"education_level"`
}