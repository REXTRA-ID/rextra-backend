package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/jelajah_profesi/controller"

	"github.com/gin-gonic/gin"
)

func ServeJelajahProfesi(
	app *gin.Engine,
	professionController controller.ProfessionController,
	favoriteController controller.FavoriteController,
	mw middleware.Middleware,
) {
	routes := app.Group("/api/v1/jelajah-profesi")
	routes.Use(mw.Authenticate())
	{
		// Profession catalog
		routes.GET("/professions", professionController.ListProfessions)
		routes.GET("/professions/:slug", professionController.GetProfessionDetail)

		// Categories
		routes.GET("/categories", professionController.ListMainCategories)
		routes.GET("/categories/:id/subcategories", professionController.ListSubCategories)

		// Favorites
		routes.GET("/favorites", favoriteController.ListFavorites)
		routes.POST("/favorites", favoriteController.AddFavorite)
		routes.DELETE("/favorites/:profession_id", favoriteController.RemoveFavorite)
	}
}
