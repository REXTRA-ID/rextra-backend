package entity

import "time"

type ProfessionSubCategory struct {
	ID             int64     `gorm:"primaryKey;autoIncrement"          json:"id"`
	MainCategoryID int64     `gorm:"not null;index"                    json:"main_category_id"`
	Code           string    `gorm:"type:varchar(50);not null"         json:"code"`
	Name           string    `gorm:"type:varchar(100);not null"        json:"name"`
	Description    string    `gorm:"type:text;not null"                json:"description"`
	CreatedAt      time.Time `gorm:"type:timestamp without time zone"  json:"created_at"`
	UpdatedAt      time.Time `gorm:"type:timestamp without time zone"  json:"updated_at"`

	MainCategory ProfessionMainCategory `gorm:"foreignKey:MainCategoryID;constraint:OnDelete:RESTRICT" json:"main_category,omitempty"`
}

func (ProfessionSubCategory) TableName() string {
	return "profession_sub_categories"
}
