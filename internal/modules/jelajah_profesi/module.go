package jelajahprofesi

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/jelajah_profesi/controller"
	"rextra-backend/internal/modules/jelajah_profesi/repository"
	"rextra-backend/internal/modules/jelajah_profesi/routes"
	"rextra-backend/internal/modules/jelajah_profesi/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitModule(server *gin.Engine, db *gorm.DB, middleware middleware.Middleware) {
	var (
		professionRepo repository.ProfessionRepository         = repository.NewProfession(db)
		categoryRepo   repository.ProfessionCategoryRepository = repository.NewProfessionCategory(db)
		detailRepo     repository.ProfessionDetailRepository   = repository.NewProfessionDetail(db)
		favoriteRepo   repository.FavoriteProfessionRepository = repository.NewFavoriteProfession(db)
	)

	professionService := service.NewProfession(professionRepo, categoryRepo, detailRepo, favoriteRepo)
	favoriteService := service.NewFavorite(favoriteRepo, professionRepo)

	professionController := controller.NewProfession(professionService)
	favoriteController := controller.NewFavorite(favoriteService)

	routes.ServeJelajahProfesi(server, professionController, favoriteController, middleware)
}
