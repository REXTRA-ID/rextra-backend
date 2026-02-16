package controller

import (
	"rextra-backend/internal/modules/membership/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type (
	MembershipPlanController interface {
		GetAllMembershipPlan(ctx *gin.Context)
	}
	membershipPlanController struct {
		membershipPlanService service.MembershipPlanService
	}
)

func NewMembership(membershipPlanService service.MembershipPlanService) MembershipPlanController {
	return &membershipPlanController{
		membershipPlanService: membershipPlanService,
	}
}

func (c *membershipPlanController) GetAllMembershipPlan(ctx *gin.Context) {
	plans, err := c.membershipPlanService.GetAllMembershipPlan(ctx)
	if err != nil {
		response.NewFailed("failed to fetch membership plans", myerror.ProcessingError(err)).Send(ctx)
		return
	}

	response.NewSuccess("", plans).Send(ctx)
}
