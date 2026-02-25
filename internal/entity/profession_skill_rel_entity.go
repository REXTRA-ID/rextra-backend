package entity

import "time"

type ProfessionSkill struct {
	ProfessionID int64     `gorm:"primaryKey"                       json:"profession_id"`
	SkillID      int64     `gorm:"primaryKey"                       json:"skill_id"`
	SkillType    string    `gorm:"type:varchar(20);not null"        json:"skill_type"` // "hard" | "soft"
	Priority     string    `gorm:"type:varchar(20);not null"        json:"priority"`   // "wajib" | "dianjurkan"
	CreatedAt    time.Time `gorm:"type:timestamp without time zone" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamp without time zone" json:"updated_at"`

	Profession Profession `gorm:"foreignKey:ProfessionID;constraint:OnDelete:CASCADE"  json:"profession,omitempty"`
	Skill      Skill      `gorm:"foreignKey:SkillID;constraint:OnDelete:RESTRICT"      json:"skill,omitempty"`
}

func (ProfessionSkill) TableName() string {
	return "profession_skill_rels"
}
