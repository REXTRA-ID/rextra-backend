package entity

type Persona struct {
	ID          string `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID      string `json:"user_id" gorm:"not null"`
	Institution string `json:"institution" gorm:"not null"`
	Study       string `json:"study" gorm:"not null"`

	// enumerated [0 = D3, 1= D4/S1, 2 = S2, 3 = S3]
	EducationLevel string `json:"education_level" gorm:"not null"`
	GraduationYear int    `json:"graduation_year" gorm:"not null"`
	CareerPlan     string `json:"career_plan" gorm:"not null"`
	CareerDreams   string `json:"career_dreams" gorm:"not null"`
	Portfolio      bool   `json:"portfolio" gorm:"default:false;not null"`
	Application    bool   `json:"application" gorm:"not null"`

	// enumerated [0 = mahasiswa aktif, 1= fresh graduate, 2 = professional]
	Status string `json:"status" gorm:"not null"`

	Timestamp
}
