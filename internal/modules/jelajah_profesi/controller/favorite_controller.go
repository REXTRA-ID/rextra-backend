package controller

import (
	dto_request "rextra-backend/internal/dto/request"
	"rextra-backend/internal/modules/jelajah_profesi/service"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type (
	FavoriteController interface {
		ListFavorites(ctx *gin.Context)
		AddFavorite(ctx *gin.Context)
		RemoveFavorite(ctx *gin.Context)
	}

	favoriteController struct {
		service service.FavoriteService
	}
)

func NewFavorite(s service.FavoriteService) FavoriteController {
	return &favoriteController{service: s}
}

func (c *favoriteController) ListFavorites(ctx *gin.Context) {
	userID := getUserID(ctx)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "25"))

	result, err := c.service.ListFavorites(ctx.Request.Context(), userID, page, limit)
	if err != nil {
		response.NewFailed("failed to fetch favorites", err).Send(ctx)
		return
	}

	response.NewSuccess("success get favorite professions", result).Send(ctx)
}

func (c *favoriteController) AddFavorite(ctx *gin.Context) {
	userID := getUserID(ctx)

	var req dto_request.AddFavoriteProfessionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.NewFailed("invalid request", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	if err := c.service.AddFavorite(ctx.Request.Context(), userID, req.ProfessionID); err != nil {
		response.NewFailed("failed to add favorite", err).Send(ctx)
		return
	}

	response.NewSuccess("success add favorite profession", nil).Send(ctx)
}

func (c *favoriteController) RemoveFavorite(ctx *gin.Context) {
	userID := getUserID(ctx)

	professionID, err := strconv.ParseInt(ctx.Param("profession_id"), 10, 64)
	if err != nil {
		response.NewFailed("invalid profession id", myerror.InvalidRequest(err)).Send(ctx)
		return
	}

	if err := c.service.RemoveFavorite(ctx.Request.Context(), userID, professionID); err != nil {
		response.NewFailed("failed to remove favorite", err).Send(ctx)
		return
	}

	response.NewSuccess("success remove favorite profession", nil).Send(ctx)
}
