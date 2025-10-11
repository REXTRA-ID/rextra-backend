package service

import (
	"rextra-backend/internal/api/repository"

	"gorm.io/gorm"
)

type (
	MembershipService interface {
	}

	membershipService struct {
		membershipRepository repository.MembershipRepository
		db                   *gorm.DB
	}
)

func NewMembership(membershipRepository repository.MembershipRepository, db *gorm.DB) MembershipService {
	return &membershipService{
		membershipRepository: membershipRepository,
		db:                   db,
	}
}
