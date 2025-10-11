package repository

import "gorm.io/gorm"

type (
	MembershipRepository interface {
	}
	membershipRepository struct {
		db *gorm.DB
	}
)

func newMembeship(db *gorm.DB) MembershipRepository {
	return &membershipRepository{
		db: db,
	}
}
