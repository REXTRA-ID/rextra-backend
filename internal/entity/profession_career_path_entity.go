package entity

import "time"

type ProfessionCareerPath struct {
	ID              int64     `gorm:"primaryKey;autoIncrement"          json:"id"`
	ProfessionID    int64     `gorm:"not null;index"                    json:"profession_id"`
	Title           string    `gorm:"type:varchar(100);not null"        json:"title"`
	ExperienceRange string    `gorm:"type:varchar(50);not null"         json:"experience_range"`
	SalaryMin       *int      `gorm:"type:int"                          json:"salary_min"`
	SalaryMax       *int      `gorm:"type:int"                          json:"salary_max"`
	SortOrder       int       `gorm:"not null"                          json:"sort_order"`
	CreatedAt       time.Time `gorm:"type:timestamp without time zone"  json:"created_at"`
	UpdatedAt       time.Time `gorm:"type:timestamp without time zone"  json:"updated_at"`

	Profession Profession `gorm:"foreignKey:ProfessionID;constraint:OnDelete:CASCADE" json:"profession,omitempty"`
}

func (ProfessionCareerPath) TableName() string {
	return "profession_career_paths"
}
