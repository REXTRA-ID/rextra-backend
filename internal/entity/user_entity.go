package entity

import (
	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin Role = "ADMIN"
	RoleUser  Role = "USER"
	RoleClub  Role = "CLUB"
)

type User struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProfileImageUrl string    `json:"profile_image_url"`
	Fullname        string    `json:"fullname" gorm:"not null"`
	Email           string    `json:"email" gorm:"uniqueIndex;not null"`
	Password        string    `json:"password" gorm:"not null"`
	IsVerified      bool      `json:"is_verified" gorm:"default:false;not null"`
	PhoneNumber     string    `json:"phone_number" gorm:"not null"`
	Role            Role      `json:"role" gorm:"default:USER;not null"`

	Membership Memberships `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (u *User) TableName() string {
	return "users"
}
