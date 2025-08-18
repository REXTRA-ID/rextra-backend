package entity

import "github.com/google/uuid"

type Riasec struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	RiasecCode   string    `json:"riasec_code" gorm:"not null"`
	RDescription *string   `json:"r_description" gorm:""`
	IDescription *string   `json:"i_description" gorm:""`
	ADescription *string   `json:"a_description" gorm:""`
	SDescription *string   `json:"s_description" gorm:""`
	EDescription *string   `json:"e_description" gorm:""`
	CDescription *string   `json:"c_description" gorm:""`

	Timestamp
}

func (r *Riasec) TableName() string {
	return "riasec"
}
