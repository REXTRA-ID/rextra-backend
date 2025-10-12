package routes

import (
	"rextra-backend/internal/api/controller"
	"rextra-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ServeAssesment(app *gin.Engine, assesmentcontroller controller.AssesmentController, middleware middleware.Middleware) {
	routes := app.Group("/api/v1/assesment")
	{

		routes.POST("/validate_hash", middleware.Authenticate(), assesmentcontroller.ValidateHash)

		routes.GET("/test/riasec/question", middleware.Authenticate(), assesmentcontroller.GetRiasecQuestion)
		routes.POST("/test/riasec/submit", middleware.Authenticate(), assesmentcontroller.SubmitRiasecAnswer)
		routes.GET("/test/riasec/result", middleware.Authenticate(), assesmentcontroller.GetRiasecResult)

		routes.GET("/test/ikigai/question", middleware.Authenticate(), assesmentcontroller.GetIkigaiQuestion)
		// routes.POST("/test/ikigai/submit", middleware.Authenticate(), assesmentcontroller.SubmitIkigaiAnswer)
		// routes.GET("/test/ikigai/result", middleware.Authenticate(), assesmentcontroller.GetIkigaiResult)
	}
}
