package entity

import "time"

type ProfessionMarketInsight struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"          json:"id"`
	ProfessionID int64     `gorm:"not null;index"                    json:"profession_id"`
	Description  string    `gorm:"type:text;not null"                json:"description"`
	SortOrder    int       `gorm:"not null"                          json:"sort_order"`
	CreatedAt    time.Time `gorm:"type:timestamp without time zone"  json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamp without time zone"  json:"updated_at"`

	Profession Profession `gorm:"foreignKey:ProfessionID;constraint:OnDelete:CASCADE" json:"profession,omitempty"`
}

func (ProfessionMarketInsight) TableName() string {
	return "profession_market_insights"
}
