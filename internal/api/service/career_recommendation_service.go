package service

import (
	"context"
	"errors"
	"time"
	"rextra-backend/internal/api/repository"
	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	myerror "rextra-backend/internal/pkg/error"

	"gorm.io/gorm"
)

type (
	CareerRecommendationService interface {
		Create(ctx context.Context, req dto_request.CreateCareerRecommendationRequest) (dto_response.CreateCareerRecommendationResponse, error)
		GetByTestSessionID(ctx context.Context, testSessionID int64) (dto_response.CreateCareerRecommendationResponse, error)
		Update(ctx context.Context, testSessionID int64, req dto_request.CreateCareerRecommendationRequest) (dto_response.CreateCareerRecommendationResponse, error)
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
	if req.TestSessionID == 0 {
		return dto_response.CreateCareerRecommendationResponse{}, myerror.InvalidRequest(errors.New("test_session_id is required"))
	}

	existing, err := s.careerRecommendationRepository.GetByTestSessionID(ctx, nil, req.TestSessionID)
	if err == nil && existing.ID != 0 {
		return dto_response.CreateCareerRecommendationResponse{}, myerror.RecordAlreadyExist("career recommendation")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto_response.CreateCareerRecommendationResponse{}, err
	}

	if len(req.RecommendationsData) == 0 {
		return dto_response.CreateCareerRecommendationResponse{}, myerror.InvalidRequest(errors.New("recommendations_data is required"))
	}

	aiModelUsed := req.AIModelUsed
	if aiModelUsed == "" {
		aiModelUsed = "gemini-1.5-flash"
	}

	createResult, err := s.careerRecommendationRepository.Create(ctx, nil, entity.CareerRecommendation{
		TestSessionID:       req.TestSessionID,
		RecommendationsData: req.RecommendationsData,
		TopProfession1ID:    req.TopProfession1ID,
		TopProfession2ID:    req.TopProfession2ID,
		AIModelUsed:         aiModelUsed,
	})
	if err != nil {
		return dto_response.CreateCareerRecommendationResponse{}, err
	}

	return dto_response.CreateCareerRecommendationResponse{
		ID:                  createResult.ID,
		TestSessionID:       createResult.TestSessionID,
		RecommendationsData: createResult.RecommendationsData,
		TopProfession1ID:    createResult.TopProfession1ID,
		TopProfession2ID:    createResult.TopProfession2ID,
		GeneratedAt:         createResult.GeneratedAt,
		AIModelUsed:         createResult.AIModelUsed,
	}, nil
}

func (s *careerRecommendationService) GetByTestSessionID(ctx context.Context, testSessionID int64) (dto_response.CreateCareerRecommendationResponse, error) {
	careerRecommendation, err := s.careerRecommendationRepository.GetByTestSessionID(ctx, nil, testSessionID)
	if err != nil {
		return dto_response.CreateCareerRecommendationResponse{}, err
	}

	return dto_response.CreateCareerRecommendationResponse{
		ID:                  careerRecommendation.ID,
		TestSessionID:       careerRecommendation.TestSessionID,
		RecommendationsData: careerRecommendation.RecommendationsData,
		TopProfession1ID:    careerRecommendation.TopProfession1ID,
		TopProfession2ID:    careerRecommendation.TopProfession2ID,
		GeneratedAt:         careerRecommendation.GeneratedAt,
		AIModelUsed:         careerRecommendation.AIModelUsed,
	}, nil
}

func (s *careerRecommendationService) Update(ctx context.Context, testSessionID int64, req dto_request.CreateCareerRecommendationRequest) (dto_response.CreateCareerRecommendationResponse, error) {
	if testSessionID == 0 {
		return dto_response.CreateCareerRecommendationResponse{}, myerror.InvalidRequest(errors.New("test_session_id is required"))
	}

	careerRecommendation, err := s.careerRecommendationRepository.GetByTestSessionID(ctx, nil, testSessionID)
	if err != nil {
		return dto_response.CreateCareerRecommendationResponse{}, err
	}

	if len(req.RecommendationsData) == 0 {
		return dto_response.CreateCareerRecommendationResponse{}, myerror.InvalidRequest(errors.New("recommendations_data is required"))
	}

	aiModelUsed := req.AIModelUsed
	if aiModelUsed == "" {
		aiModelUsed = careerRecommendation.AIModelUsed
	}

	careerRecommendation.RecommendationsData = req.RecommendationsData
	careerRecommendation.TopProfession1ID = req.TopProfession1ID
	careerRecommendation.TopProfession2ID = req.TopProfession2ID
	careerRecommendation.AIModelUsed = aiModelUsed
	careerRecommendation.GeneratedAt = time.Now()

	updateResult, err := s.careerRecommendationRepository.Update(ctx, nil, careerRecommendation)
	if err != nil {
		return dto_response.CreateCareerRecommendationResponse{}, err
	}
	return dto_response.CreateCareerRecommendationResponse{
		ID:                  updateResult.ID,
		TestSessionID:       updateResult.TestSessionID,
		RecommendationsData: updateResult.RecommendationsData,
		TopProfession1ID:    updateResult.TopProfession1ID,
		TopProfession2ID:    updateResult.TopProfession2ID,
		GeneratedAt:         updateResult.GeneratedAt,
		AIModelUsed:         updateResult.AIModelUsed,
	}, nil
}
