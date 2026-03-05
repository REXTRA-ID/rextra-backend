package controller

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/jelajah_profesi/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	ProfessionController interface {
		ListProfessions(ctx *gin.Context)
		GetProfessionDetail(ctx *gin.Context)
		ListMainCategories(ctx *gin.Context)
		ListSubCategories(ctx *gin.Context)
	}

	professionController struct {
		service service.ProfessionService
	}
)

func NewProfession(s service.ProfessionService) ProfessionController {
	return &professionController{service: s}
}

func (c *professionController) ListProfessions(ctx *gin.Context) {
	var req dto_request.ListProfessionsRequest
	if err := ctx.ShouldBind(&req); err != nil {
		response.NewFailed("invalid request", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	userID := getUserID(ctx)

	result, err := c.service.ListProfessions(ctx.Request.Context(), req, userID)
	if err != nil {
		response.NewFailed("failed to fetch professions", err).Send(ctx)
		return
	}

	response.NewSuccess("success get profession list", result).Send(ctx)
}

func (c *professionController) GetProfessionDetail(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if slug == "" {
		response.NewFailed("invalid slug", myerror.InvalidRequest(nil)).Send(ctx)
		return
	}

	userID := getUserID(ctx)

	result, err := c.service.GetProfessionDetail(ctx.Request.Context(), slug, userID)
	if err != nil {
		response.NewFailed("failed to fetch profession detail", err).Send(ctx)
		return
	}

	response.NewSuccess("success get profession detail", result).Send(ctx)
}

func (c *professionController) ListMainCategories(ctx *gin.Context) {
	result, err := c.service.ListMainCategories(ctx.Request.Context())
	if err != nil {
		response.NewFailed("failed to fetch categories", err).Send(ctx)
		return
	}

	response.NewSuccess("success get main categories", result).Send(ctx)
}

func (c *professionController) ListSubCategories(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.NewFailed("invalid category id", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	result, err := c.service.ListSubCategories(ctx.Request.Context(), id)
	if err != nil {
		response.NewFailed("failed to fetch sub categories", err).Send(ctx)
		return
	}

	response.NewSuccess("success get sub categories", result).Send(ctx)
}

func getUserID(ctx *gin.Context) uuid.UUID {
	uid, exists := ctx.Get("user_id")
	if !exists {
		return uuid.Nil
	}

	switch v := uid.(type) {
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return uuid.Nil
		}
		return parsed
	case uuid.UUID:
		return v
	default:
		return uuid.Nil
	}
}
