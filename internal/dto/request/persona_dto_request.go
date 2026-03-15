package dto_request

type (
	CreatePersonaRequest struct {
		UserID      string `json:"user_id"`
		Institution string `json:"institution"`
		Study       string `json:"study"`

		// enumerated [0 = D3, 1= D4/S1, 2 = S2, 3 = S3]
		EducationLevel string `json:"education_level"`
		GraduationYear int    `json:"graduation_year"`
		CareerPlan     string `json:"career_plan"`
		CareerDreams   string `json:"career_dreams"`
		Portfolio      bool   `json:"portfolio"`
		Application    bool   `json:"application"`

		// enumerated [0 = mahasiswa aktif, 1= fresh graduate, 2 = professional]
		Status string `json:"status"`
	}
)
