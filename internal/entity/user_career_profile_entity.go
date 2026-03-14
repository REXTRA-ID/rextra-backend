package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserCareerProfile struct {
	ID               int64     json:"id" gorm:"primaryKey;autoIncrement"
	UserID           uuid.UUID json:"user_id" gorm:"type:uuid;not null;index:idx_user_career_profiles_user_id"
	TestSessionID    int64     json:"test_session_id" gorm:"not null;index:idx_user_career_profiles_session_id"
	TopProfession1ID *int64    json:"top_profession1_id" gorm:"type:bigint"
	TopProfession2ID *int64    json:"top_profession2_id" gorm:"type:bigint"
	RIASECCode       string    json:"riasec_code" gorm:"type:varchar(6)"
	IsActive         bool      json:"is_active" gorm:"default:true;not null;index:idx_user_career_profiles_user_active"
	CreatedAt        time.Time json:"created_at" gorm:"type:timestamptz;default:now();not null"
	ActivatedAt      time.Time json:"activated_at" gorm:"type:timestamptz;default:now()"

	User        User                     json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"
	TestSession CareerProfileTestSession json:"-" gorm:"foreignKey:TestSessionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"
}

func (UserCareerProfile) TableName() string {
	return "user_career_profiles"
}
