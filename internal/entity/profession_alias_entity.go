package entity

import "time"

type ProfessionAlias struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"          json:"id"`
	ProfessionID int64     `gorm:"not null;index"                    json:"profession_id"`
	AliasName    string    `gorm:"type:varchar(100);not null"        json:"alias_name"`
	CreatedAt    time.Time `gorm:"type:timestamp without time zone"  json:"created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamp without time zone"  json:"updated_at"`

	Profession Profession `gorm:"foreignKey:ProfessionID;constraint:OnDelete:CASCADE" json:"profession,omitempty"`
}

func (ProfessionAlias) TableName() string {
	return "profession_aliases"
}
