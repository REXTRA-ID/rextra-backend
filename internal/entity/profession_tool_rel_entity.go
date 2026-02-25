package entity

import "time"

type ProfessionTool struct {
	ProfessionID int64     `gorm:"primaryKey"                       json:"profession_id"`
	ToolID       int64     `gorm:"primaryKey"                       json:"tool_id"`
	UsageType    string    `gorm:"type:varchar(20);not null"        json:"usage_type"` // "wajib" | "umum"
	CreatedAt    time.Time `gorm:"type:timestamp without time zone" json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamp without time zone" json:"updated_at"`

	Profession Profession `gorm:"foreignKey:ProfessionID;constraint:OnDelete:CASCADE" json:"profession,omitempty"`
	Tool       Tool       `gorm:"foreignKey:ToolID;constraint:OnDelete:RESTRICT"      json:"tool,omitempty"`
}

func (ProfessionTool) TableName() string {
	return "profession_tool_rels"
}
