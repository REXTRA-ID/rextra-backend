package entity

import "time"

type StudyProgram struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"                json:"id"`
	Name      string    `gorm:"type:varchar(150);uniqueIndex;not null"  json:"name"`
	CreatedAt time.Time `gorm:"type:timestamp without time zone"        json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp without time zone"        json:"updated_at"`
}

func (StudyProgram) TableName() string {
	return "study_programs"
}
