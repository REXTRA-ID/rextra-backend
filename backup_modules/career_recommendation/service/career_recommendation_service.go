package service

import (
	"context"
	"encoding/json"
	"errors"
	"rextra-backend/internal/modules/career_recommendation/repository"
	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	myerror "rextra-backend/internal/pkg/error"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type (
	CareerRecommendationService interface {
		Create(ctx context.Context, req dto_request.CreateCareerRecommendationRequest) (dto_response.CreateCareerRecommendationResponse, error)
		GetByUserID(ctx context.Context, userId string) (dto_response.CreateCareerRecommendationResponse, error)
		Update(ctx context.Context, id string, req dto_request.CreateCareerRecommendationRequest) (dto_response.CreateCareerRecommendationResponse, error)
	}

	careerRecommendationService struct {
		careerRecommendationRepository repository.CareerRecommendationRepository
		db                             *gorm.DB
	}
)

func NewCareerRecommendation(careerRecommendationRepository repository.CareerRecommendationRepository, db *gorm.DB) CareerRecommendationService {
	return &careerRecommendationService{
		careerRecommendationRepository: careerRecommendationRepository,
		db:                             db,
	}
}

func (s *careerRecommendationService) Create(ctx context.Context, req dto_request.CreateCareerRecommendationRequest) (dto_response.CreateCareerRecommendationResponse, error) {
	_, err := s.careerRecommendationRepository.GetByUserID(ctx, nil, req.UserID)
	if err == nil {
		return dto_response.CreateCareerRecommendationResponse{}, myerror.RecordAlreadyExist("career recommendation")
	}

	if len(req.Analysis) != 3 {
		return dto_response.CreateCareerRecommendationResponse{}, myerror.InvalidRequest(errors.New("analysis must contain exactly 3 items"))
	}

	if len(req.TopProfessions) != 2 {
		return dto_response.CreateCareerRecommendationResponse{}, myerror.InvalidRequest(errors.New("top professions must contain exactly 2 items"))
	}

	analysisJSON, err := json.Marshal(req.Analysis)
	if err != nil {
		return dto_response.CreateCareerRecommendationResponse{}, err
	}

	createResult, err := s.careerRecommendationRepository.Create(ctx, nil, entity.CareerRecommendation{
		UserID:         uuid.MustParse(req.UserID),
		Analysis:       analysisJSON,
		TopProfessions: pq.StringArray(req.TopProfessions),
	})
	if err != nil {
		return dto_response.CreateCareerRecommendationResponse{}, err
	}

	var analysis []dto_response.CareerRecommendationAnalysisResponse
	if err := json.Unmarshal(createResult.Analysis, &analysis); err != nil {
		return dto_response.CreateCareerRecommendationResponse{}, err
	}

	return dto_response.CreateCareerRecommendationResponse{
		ID:             createResult.ID.String(),
		UserID:         createResult.UserID.String(),
		Analysis:       analysis,
		TopProfessions: createResult.TopProfessions,
	}, nil
}

func (s *careerRecommendationService) GetByUserID(ctx context.Context, userId string) (dto_response.CreateCareerRecommendationResponse, error) {
	careerRecommendation, err := s.careerRecommendationRepository.GetByUserID(ctx, nil, userId)
	if err != nil {
		return dto_response.CreateCareerRecommendationResponse{}, err
	}

	var analysis []dto_response.CareerRecommendationAnalysisResponse
	if err := json.Unmarshal(careerRecommendation.Analysis, &analysis); err != nil {
		return dto_response.CreateCareerRecommendationResponse{}, err
	}

	return dto_response.CreateCareerRecommendationResponse{
		ID:             careerRecommendation.ID.String(),
		UserID:         careerRecommendation.UserID.String(),
		Analysis:       analysis,
		TopProfessions: careerRecommendation.TopProfessions,
	}, nil
}

func (s *careerRecommendationService) Update(ctx context.Context, id string, req dto_request.CreateCareerRecommendationRequest) (dto_response.CreateCareerRecommendationResponse, error) {
	careerRecommendation, err := s.careerRecommendationRepository.GetByUserID(ctx, nil, id)
	if err != nil {
		return dto_response.CreateCareerRecommendationResponse{}, err
	}

	analysisJSON, err := json.Marshal(req.Analysis)
	if err != nil {
		return dto_response.CreateCareerRecommendationResponse{}, err
	}

	careerRecommendation.Analysis = analysisJSON
	careerRecommendation.TopProfessions = pq.StringArray(req.TopProfessions)

	updateResult, err := s.careerRecommendationRepository.Update(ctx, nil, careerRecommendation)
	if err != nil {
		return dto_response.CreateCareerRecommendationResponse{}, err
	}
	var analysis []dto_response.CareerRecommendationAnalysisResponse
	if err := json.Unmarshal(updateResult.Analysis, &analysis); err != nil {
		return dto_response.CreateCareerRecommendationResponse{}, err
	}

	return dto_response.CreateCareerRecommendationResponse{
		ID:             updateResult.ID.String(),
		UserID:         updateResult.UserID.String(),
		Analysis:       analysis,
		TopProfessions: updateResult.TopProfessions,
	}, nil
}
