package entity

import "time"

type Skill struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"                json:"id"`
	Name      string    `gorm:"type:varchar(100);uniqueIndex;not null"  json:"name"`
	CreatedAt time.Time `gorm:"type:timestamp without time zone"        json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp without time zone"        json:"updated_at"`
}

func (Skill) TableName() string {
	return "skills"
}
