package entity

import (
	"time"
)

type KenaliDiriCategory struct {
	ID              int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	CategoryCode    string    `json:"category_code" gorm:"size:50;uniqueIndex;not null"`
	CategoryName    string    `json:"category_name" gorm:"size:255;not null"`
	Description     *string   `json:"description" gorm:"type:text"`
	DetailTableName string    `json:"detail_table_name" gorm:"size:100;not null"`
	IsActive        bool      `json:"is_active" gorm:"default:true;not null"`
	CreatedAt       time.Time `json:"created_at" gorm:"type:timestamptz;default:now();autoCreateTime"`
}

func (KenaliDiriCategory) TableName() string {
	return "kenalidiri_categories"
}
