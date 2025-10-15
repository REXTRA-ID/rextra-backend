package dto_request

type CreateEducationRequest struct {
	UserID                 string `json:"user_id" binding:"required"`
	InstitutionName        string `json:"institution_name" binding:"required"`
	Major                  string `json:"major" binding:"required"`
	Faculty                string `json:"faculty" binding:"required"`
	EntryYear              int    `json:"entry_year" binding:"required"`
	ExpectedGraduationYear int    `json:"expected_graduation_year" binding:"required"`
	ActualGraduationYear   *int   `json:"actual_graduation_year,omitempty"`
	CurrentSemester        int    `json:"current_semester" binding:"required"`
	TotalSemester          int    `json:"total_semester" binding:"required"`
	EducationLevel         string `json:"education_level" binding:"required"`
	Status                 string `json:"status" binding:"required"`
}
