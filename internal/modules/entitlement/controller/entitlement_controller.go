package controller

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/entitlement/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

type (
	EntitlementController interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetById(ctx *gin.Context)
		UpdateRestriction(ctx *gin.Context)
		Delete(ctx *gin.Context)
	}

	entitlementController struct {
		entitlementService service.EntitlementService
	}
)

func NewEntitlementController(svc service.EntitlementService) EntitlementController {
	return &entitlementController{entitlementService: svc}
}

func (c *entitlementController) Create(ctx *gin.Context) {
	var req dto_request.CreateEntitlementRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("invalid request body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}
	result, err := c.entitlementService.Create(ctx, req)
	if err != nil {
		response.NewFailed("failed to create entitlement", err).Send(ctx)
		return
	}
	response.NewSuccess("entitlement created successfully", result).Send(ctx)
}

func (c *entitlementController) GetAll(ctx *gin.Context) {
	featureID := ctx.Query("feature_id")
	var featureIDPtr *string
	if featureID != "" {
		featureIDPtr = &featureID
	}

	result, err := c.entitlementService.GetAll(ctx, featureIDPtr)
	if err != nil {
		response.NewFailed("failed to retrieve entitlements", err).Send(ctx)
		return
	}
	response.NewSuccess("entitlements retrieved successfully", result).Send(ctx)
}

func (c *entitlementController) GetById(ctx *gin.Context) {
	id := ctx.Param("entitlementId")
	result, err := c.entitlementService.GetById(ctx, id)
	if err != nil {
		response.NewFailed("failed to retrieve entitlement", err).Send(ctx)
		return
	}
	response.NewSuccess("entitlement retrieved successfully", result).Send(ctx)
}

func (c *entitlementController) UpdateRestriction(ctx *gin.Context) {
	id := ctx.Param("entitlementId")

	var req dto_request.UpdateEntitlementRestrictionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("invalid request body", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	result, err := c.entitlementService.UpdateRestriction(ctx, id, req)
	if err != nil {
		response.NewFailed("failed to update entitlement restriction", err).Send(ctx)
		return
	}
	response.NewSuccess("entitlement restriction updated successfully", result).Send(ctx)
}

func (c *entitlementController) Delete(ctx *gin.Context) {
	id := ctx.Param("entitlementId")
	if err := c.entitlementService.Delete(ctx, id); err != nil {
		response.NewFailed("failed to delete entitlement", err).Send(ctx)
		return
	}
	response.NewSuccess("entitlement deleted successfully", nil).Send(ctx)
}
