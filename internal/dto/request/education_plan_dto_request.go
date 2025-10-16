package dto_request

type CreateEducationPlanRequest struct {
    UserID            string `json:"user_id" binding:"required"`
    InstitutionName   string `json:"institution_name" binding:"required"`
    Major             string `json:"major" binding:"required"`
    Faculty           string `json:"faculty" binding:"required"`
    ExpectedEntryYear int    `json:"expected_entry_year" binding:"required"`
    EducationProgram  string `json:"education_program" binding:"required"`
    EducationLevel    string `json:"education_level" binding:"required"`
}

type UpdateEducationPlanRequest struct {
    ID                string `json:"id" binding:"required"`
    InstitutionName   string `json:"institution_name" binding:"required"`
    Major             string `json:"major" binding:"required"`
    Faculty           string `json:"faculty" binding:"required"`
    ExpectedEntryYear int    `json:"expected_entry_year" binding:"required"`
    EducationProgram  string `json:"education_program" binding:"required"`
    EducationLevel    string `json:"education_level" binding:"required"`
}