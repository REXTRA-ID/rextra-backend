package controller

import (
	"github.com/gin-gonic/gin"

	"rextra-backend/internal/api/service"
)

type (
	AssesmentController interface {
		ValidateHash(ctx *gin.Context)
		// GetRiasecQuestion(ctx *gin.Context)
		// SubmitRiasecAnswer(ctx *gin.Context)
		// GetRiasecResult(ctx *gin.Context)
		// GetIkigaiQuestion(ctx *gin.Context)
		// SubmitIkigaiAnswer(ctx *gin.Context)
		// GetIkigaiResult(ctx *gin.Context)
	}

	assesmentcontroller struct {
		assesmentService service.AssesmentService
	}
)

func NewAssesment(assesmentService service.AssesmentService) AssesmentController {
	return &assesmentcontroller{
		assesmentService: assesmentService,
	}
}

func (c *assesmentcontroller) ValidateHash(ctx *gin.Context) {
	
}
