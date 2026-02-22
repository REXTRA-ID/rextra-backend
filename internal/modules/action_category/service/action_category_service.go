package service

import (
	"context"
	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	"rextra-backend/internal/modules/action_category/repository"
)

type (
	ActionCategoryService interface {
		Create(ctx context.Context, req dto_request.CreateActionCategoryRequest) (dto_response.GetActionCategoryResponse, error)
		GetAll(ctx context.Context) ([]dto_response.GetActionCategoryResponse, error)
	}

	actionCategoryService struct {
		actionCategoryRepository repository.ActionCategoryRepository
	}
)

func NewActionCategoryService(repo repository.ActionCategoryRepository) ActionCategoryService {
	return &actionCategoryService{
		actionCategoryRepository: repo,
	}
}

func (s *actionCategoryService) Create(ctx context.Context, req dto_request.CreateActionCategoryRequest) (dto_response.GetActionCategoryResponse, error) {
	newActionCategory := entity.NewActionCategory(req.Name, req.Slug, req.Description, req.Status)
	result, err := s.actionCategoryRepository.Create(ctx, nil, newActionCategory)
	if err != nil {
		return dto_response.GetActionCategoryResponse{}, err
	}
	return dto_response.GetActionCategoryResponse{
		Id:          result.ID.String(),
		Name:        result.Name,
		Description: result.Description,
		Status:      string(result.Status),
		Slug:        result.Slug,
	}, nil
}

func (s *actionCategoryService) GetAll(ctx context.Context) ([]dto_response.GetActionCategoryResponse, error) {
	result, err := s.actionCategoryRepository.GetAll(ctx, nil)
	if err != nil {
		return nil, err
	}
	var res []dto_response.GetActionCategoryResponse
	for i := range result {
		res = append(res, dto_response.GetActionCategoryResponse{
			Id:          result[i].ID.String(),
			Name:        result[i].Name,
			Description: result[i].Description,
			Slug:        result[i].Slug,
			Status:      string(result[i].Status),
		})
	}
	return res, nil
}
