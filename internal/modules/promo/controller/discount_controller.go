package controller

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/promo/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"
	"rextra-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type (
	DiscountController interface {
		Create(ctx *gin.Context)
		GetAll(ctx *gin.Context)
		GetByID(ctx *gin.Context)
		Update(ctx *gin.Context)
		Delete(ctx *gin.Context)
		GetRedemptions(ctx *gin.Context)
		ValidateCode(ctx *gin.Context)
	}

	discountController struct {
		service service.DiscountService
	}
)

func NewDiscountController(svc service.DiscountService) DiscountController {
	return &discountController{service: svc}
}

func (c *discountController) Create(ctx *gin.Context) {
	adminNameVal, _ := ctx.Get("user_name")
	adminName, _ := adminNameVal.(string)
	if adminName == "" { adminName = "Admin" }
	var req dto_request.CreateDiscountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil { response.NewFailed("invalid body", myerror.InvalidRequest(err)).Send(ctx); return }
	result, err := c.service.Create(ctx, req, adminName)
	if err != nil { response.NewFailed("failed create", err).Send(ctx); return }
	response.NewSuccess("created", result).ChangeStatusCode(201).Send(ctx)
}

func (c *discountController) GetAll(ctx *gin.Context) {
	var req dto_request.DiscountFilterRequest
	ctx.ShouldBindQuery(&req)
	result, err := c.service.GetAll(ctx, req)
	if err != nil { response.NewFailed("failed get all", err).Send(ctx); return }
	response.NewSuccess("retrieved", result).Send(ctx)
}

func (c *discountController) GetByID(ctx *gin.Context) {
	result, err := c.service.GetByID(ctx, ctx.Param("discountId"))
	if err != nil { response.NewFailed("failed get", err).Send(ctx); return }
	response.NewSuccess("retrieved", result).Send(ctx)
}

func (c *discountController) Update(ctx *gin.Context) {
	adminNameVal, _ := ctx.Get("user_name")
	adminName, _ := adminNameVal.(string)
	var req dto_request.UpdateDiscountRequest
	ctx.ShouldBindJSON(&req)
	result, err := c.service.Update(ctx, ctx.Param("discountId"), req, adminName)
	if err != nil { response.NewFailed("failed update", err).Send(ctx); return }
	response.NewSuccess("updated", result).Send(ctx)
}

func (c *discountController) Delete(ctx *gin.Context) {
	if err := c.service.Delete(ctx, ctx.Param("discountId")); err != nil { response.NewFailed("failed delete", err).Send(ctx); return }
	response.NewSuccess("deactivated", nil).Send(ctx)
}

func (c *discountController) GetRedemptions(ctx *gin.Context) {
	var req dto_request.DiscountRedemptionFilterRequest
	ctx.ShouldBindQuery(&req)
	result, err := c.service.GetRedemptions(ctx, ctx.Param("discountId"), req)
	if err != nil { response.NewFailed("failed get redemptions", err).Send(ctx); return }
	response.NewSuccess("retrieved", result).Send(ctx)
}

func (c *discountController) ValidateCode(ctx *gin.Context) {
	userID, _ := utils.GetUserIdFromCtx(ctx)
	var req dto_request.ValidateDiscountRequest
	ctx.ShouldBindJSON(&req)
	result, err := c.service.ValidateCode(ctx, req, userID)
	if err != nil { response.NewFailed("invalid code", err).Send(ctx); return }
	response.NewSuccess("valid", result).Send(ctx)
}
