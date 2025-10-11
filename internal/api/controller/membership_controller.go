package controller

import "rextra-backend/internal/api/service"

type (
	MembershipController interface {
	}
	membershipController struct {
		membershipService service.MembershipService
	}
)

func NewMembership(membershipService service.MembershipService) MembershipController {
	return &membershipController{
		membershipService: membershipService,
	}
}
