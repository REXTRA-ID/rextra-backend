package service

import (
	"context"

	"rextra-backend/internal/api/repository"
	dto_request "rextra-backend/internal/dto/request"

	"gorm.io/gorm"
)

type (
	AssesmentService interface {
		ValidateHash(ctx context.Context, req dto_request.ValidateHashRequest) error
		// GetRiasecQuestion(ctx *gin.Context)
		// SubmitRiasecAnswer(ctx *gin.Context)
		// GetRiasecResult(ctx *gin.Context)
		// GetIkigaiQuestion(ctx *gin.Context)
		// SubmitIkigaiAnswer(ctx *gin.Context)
		// GetIkigaiResult(ctx *gin.Context)
	}

	assesmentService struct {
		riasecRepository repository.RiasecRepository
		db                *gorm.DB
	}
)

func NewAssesment(riasecRepository repository.RiasecRepository, db *gorm.DB) AssesmentService {
	return &assesmentService{
		riasecRepository: riasecRepository,
		db:                db,
	}
}

func (s *assesmentService) ValidateHash(ctx context.Context, req dto_request.ValidateHashRequest) error {
	return nil
}