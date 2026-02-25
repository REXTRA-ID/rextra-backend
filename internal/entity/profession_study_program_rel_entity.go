package entity

import "time"

type ProfessionStudyProgram struct {
	ProfessionID   int64     `gorm:"primaryKey"  json:"profession_id"`
	StudyProgramID int64     `gorm:"primaryKey"  json:"study_program_id"`
	CreatedAt      time.Time `gorm:"type:timestamp without time zone" json:"created_at"`
	UpdatedAt      time.Time `gorm:"type:timestamp without time zone" json:"updated_at"`

	Profession   Profession   `gorm:"foreignKey:ProfessionID;constraint:OnDelete:CASCADE"   json:"profession,omitempty"`
	StudyProgram StudyProgram `gorm:"foreignKey:StudyProgramID;constraint:OnDelete:RESTRICT" json:"study_program,omitempty"`
}

func (ProfessionStudyProgram) TableName() string {
	return "profession_study_program_rels"
}
