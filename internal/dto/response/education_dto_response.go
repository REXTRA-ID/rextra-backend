package dto_response

type EducationResponse struct {
	ID                     string `json:"id"`
	UserID                 string `json:"user_id"`
	InstitutionName        string `json:"institution_name"`
	Major                  string `json:"major"`
	Faculty                string `json:"faculty"`
	EntryYear              int    `json:"entry_year"`
	ExpectedGraduationYear int    `json:"expected_graduation_year"`
	ActualGraduationYear   *int   `json:"actual_graduation_year,omitempty"`
	CurrentSemester        int    `json:"current_semester"`
	TotalSemester          int    `json:"total_semester"`
	EducationLevel         string `json:"education_level"`
	Status                 string `json:"status"`
	IsActive               bool   `json:"is_active"`
}
