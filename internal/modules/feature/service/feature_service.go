package service

import (
	"context"
	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	"rextra-backend/internal/modules/feature/repository"

	"github.com/google/uuid"
)

type (
	FeatureService interface {
		Create(ctx context.Context, req dto_request.CreateFeatureRequest) (dto_response.GetFeatureResponse, error)
		GetAll(ctx context.Context) ([]dto_response.GetFeatureResponse, error)
		GetById(ctx context.Context, id string) (dto_response.GetFeatureResponse, error)
	}

	featureService struct {
		featureRepository repository.FeatureRepository
	}
)

func NewFeatureService(repository repository.FeatureRepository) FeatureService {
	return &featureService{
		featureRepository: repository,
	}
}

func (s *featureService) Create(ctx context.Context, req dto_request.CreateFeatureRequest) (dto_response.GetFeatureResponse, error) {
	var parentUUID *uuid.UUID
	if req.ParentId != nil {
		parsed := uuid.MustParse(*req.ParentId)
		parentUUID = &parsed
	}
	newFeature := entity.NewFeature(req.Name, req.Slug, parentUUID, req.Status)
	result, err := s.featureRepository.Create(ctx, nil, newFeature)
	if err != nil {
		return dto_response.GetFeatureResponse{}, err
	}
	return dto_response.GetFeatureResponse{
		Id:     result.ID.String(),
		Name:   result.Name,
		Slug:   result.Slug,
		Status: string(result.Status),
		Type:   string(result.Type),
	}, nil
}
func (s *featureService) GetAll(ctx context.Context) ([]dto_response.GetFeatureResponse, error) {
	result, err := s.featureRepository.GetAll(ctx, nil)
	if err != nil {
		return nil, err
	}
	var res []dto_response.GetFeatureResponse
	for i := range result {
		res = append(res, dto_response.GetFeatureResponse{
			Id:     result[i].ID.String(),
			Name:   result[i].Name,
			Slug:   result[i].Slug,
			Status: string(result[i].Status),
			Type:   string(result[i].Type),
		})
	}
	return res, nil
}
func (s *featureService) GetById(ctx context.Context, id string) (dto_response.GetFeatureResponse, error) {
	result, err := s.featureRepository.GetById(ctx, nil, uuid.MustParse(id))
	if err != nil {
		return dto_response.GetFeatureResponse{}, err
	}
	return dto_response.GetFeatureResponse{
		Id:     result.ID.String(),
		Name:   result.Name,
		Slug:   result.Slug,
		Status: string(result.Status),
		Type:   string(result.Type),
	}, nil
}
