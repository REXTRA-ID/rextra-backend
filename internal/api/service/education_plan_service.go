package service

import (
	"context"
	"rextra-backend/internal/api/repository"
	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	myerror "rextra-backend/internal/pkg/error"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	EducationPlanService interface {
		Create(ctx context.Context, tx *gorm.DB, req dto_request.CreateEducationPlanRequest) (dto_response.EducationPlanResponse, error)
		GetAll(ctx context.Context, userID string) ([]dto_response.EducationPlanResponse, error)
		GetById(ctx context.Context, userID string, id string) (dto_response.EducationPlanResponse, error)
		Update(ctx context.Context, userID string, req dto_request.UpdateEducationPlanRequest) (dto_response.EducationPlanResponse, error)
		Delete(ctx context.Context, userID string, id string) error
	}

	educationPlanService struct {
		repo repository.EducationPlanRepository
		db  *gorm.DB
	}
)

func NewEducationPlan(repo repository.EducationPlanRepository, db *gorm.DB) EducationPlanService {
	return &educationPlanService{repo: repo, db: db}
}

func (s *educationPlanService) Create(ctx context.Context, tx *gorm.DB, req dto_request.CreateEducationPlanRequest) (dto_response.EducationPlanResponse, error) {
	userId, err := uuid.Parse(req.UserID)
	if err != nil {
		return dto_response.EducationPlanResponse{}, err
	}

	createRequest := entity.EducationPlan{
		UserID:                 userId,
		InstitutionName:        req.InstitutionName,
		Major:                  req.Major,
		Faculty:                req.Faculty,
		ExpectedEntryYear:      req.ExpectedEntryYear,
		EducationProgram:       req.EducationProgram,
		EducationLevel:         entity.EducationLevel(req.EducationLevel),
	}

	createdEducationPlan, err := s.repo.Create(ctx, nil, createRequest)
	if err != nil {
		return dto_response.EducationPlanResponse{}, err
	}

	return dto_response.EducationPlanResponse{
		ID:                     createdEducationPlan.ID.String(),
		UserID:                 createdEducationPlan.UserID.String(),
		InstitutionName:        createdEducationPlan.InstitutionName,
		Major:                  createdEducationPlan.Major,
		Faculty:                createdEducationPlan.Faculty,
		ExpectedEntryYear:      createdEducationPlan.ExpectedEntryYear,
		EducationProgram:       createdEducationPlan.EducationProgram,
		EducationLevel:         string(createdEducationPlan.EducationLevel),
	}, nil
}

func (s *educationPlanService) GetAll(ctx context.Context, userID string) ([]dto_response.EducationPlanResponse, error) {
	plans, err := s.repo.GetAllByUserId(ctx, nil, userID)
	if err != nil {
		return nil, err
	}

	res := make([]dto_response.EducationPlanResponse, 0, len(plans))
	for _, p := range plans {
		res = append(res, dto_response.EducationPlanResponse{
			ID:               p.ID.String(),
			UserID:           p.UserID.String(),
			InstitutionName:  p.InstitutionName,
			Major:            p.Major,
			Faculty:          p.Faculty,
			ExpectedEntryYear: p.ExpectedEntryYear,
			EducationProgram: p.EducationProgram,
			EducationLevel:   string(p.EducationLevel),
		})
	}
	return res, nil
}

func (s *educationPlanService) GetById(ctx context.Context, userID string, id string) (dto_response.EducationPlanResponse, error) {
	plan, found, err := s.repo.GetByUserIdAndEducationPlanById(ctx, nil, userID, id)
	if err != nil {
		return dto_response.EducationPlanResponse{}, err
	}
	if !found {
		return dto_response.EducationPlanResponse{}, myerror.RecordNotFound("education plan")
	}
	return dto_response.EducationPlanResponse{
		ID:               plan.ID.String(),
		UserID:           plan.UserID.String(),
		InstitutionName:  plan.InstitutionName,
		Major:            plan.Major,
		Faculty:          plan.Faculty,
		ExpectedEntryYear: plan.ExpectedEntryYear,
		EducationProgram: plan.EducationProgram,
		EducationLevel:   string(plan.EducationLevel),
	}, nil
}

func (s *educationPlanService) Update(ctx context.Context, userID string, req dto_request.UpdateEducationPlanRequest) (dto_response.EducationPlanResponse, error) {
	existing, found, err := s.repo.GetByUserIdAndEducationPlanById(ctx, nil, userID, req.ID)
	if err != nil {
		return dto_response.EducationPlanResponse{}, err
	}
	if !found {
		return dto_response.EducationPlanResponse{}, myerror.RecordNotFound("education plan")
	}

	// apply updates
	existing.InstitutionName = req.InstitutionName
	existing.Major = req.Major
	existing.Faculty = req.Faculty
	existing.ExpectedEntryYear = req.ExpectedEntryYear
	existing.EducationProgram = req.EducationProgram
	existing.EducationLevel = entity.EducationLevel(req.EducationLevel)

	updated, err := s.repo.Update(ctx, nil, existing)
	if err != nil {
		return dto_response.EducationPlanResponse{}, err
	}

	return dto_response.EducationPlanResponse{
		ID:               updated.ID.String(),
		UserID:           updated.UserID.String(),
		InstitutionName:  updated.InstitutionName,
		Major:            updated.Major,
		Faculty:          updated.Faculty,
		ExpectedEntryYear: updated.ExpectedEntryYear,
		EducationProgram: updated.EducationProgram,
		EducationLevel:   string(updated.EducationLevel),
	}, nil
}

func (s *educationPlanService) Delete(ctx context.Context, userID string, id string) error {
	_, found, err := s.repo.GetByUserIdAndEducationPlanById(ctx, nil, userID, id)
	if err != nil {
		return err
	}
	if !found {
		return myerror.RecordNotFound("education plan")
	}
	
	return s.repo.Delete(ctx, nil, userID, id)
}
