package routes

import (
	"rextra-backend/internal/middleware"
	"rextra-backend/internal/modules/kenali_diri/controller"

	"github.com/gin-gonic/gin"
)

func ServeCareerProfileFeedback(
	app *gin.Engine,
	c controller.CareerProfileFeedbackController,
	middleware middleware.Middleware,
) {

	r := app.Group("/api/v1/admin/kenali-diri/feedback")
	r.Use(middleware.Authenticate(), middleware.OnlyAdmin())
	{
		r.GET("/student", c.GetStudentFeedbacks)
		r.GET("/student/stats", c.GetStudentFeedbackStats)
		r.GET("/expert", c.GetExpertFeedbacks)
		r.GET("/expert/:id", c.GetExpertFeedbackDetail)
		r.GET("/meta", c.GetFeedbackMetadata)
	}
}
