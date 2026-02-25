package entity

import "time"

type ProfessionMainCategory struct {
	ID          int64     `gorm:"primaryKey;autoIncrement"               json:"id"`
	Code        string    `gorm:"type:varchar(50);uniqueIndex;not null"  json:"code"`
	Name        string    `gorm:"type:varchar(100);not null"             json:"name"`
	Description string    `gorm:"type:text;not null"                     json:"description"`
	CreatedAt   time.Time `gorm:"type:timestamp without time zone"       json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:timestamp without time zone"       json:"updated_at"`
}

func (ProfessionMainCategory) TableName() string {
	return "profession_main_categories"
}
