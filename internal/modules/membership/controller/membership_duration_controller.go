package controller

import (
	"rextra-backend/internal/modules/membership/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type (
	MembershipDurationController interface {
		GetAllMembershipDuration(ctx *gin.Context)
	}

	membershipDurationController struct {
		membershipDurationService service.MembershipDurationService
	}
)

func NewMembershipDurationController(membershipDurationService service.MembershipDurationService) MembershipDurationController {
	return &membershipDurationController{
		membershipDurationService: membershipDurationService,
	}
}

func (c *membershipDurationController) GetAllMembershipDuration(ctx *gin.Context) {
	membershipDuration, err := c.membershipDurationService.GetAllMembershipDuration(ctx)
	if err != nil {
		response.NewFailed("failed to fetch membership duration", myerror.ProcessingError(err)).Send(ctx)
		return
	}
	response.NewSuccess("", membershipDuration).Send(ctx)
}
