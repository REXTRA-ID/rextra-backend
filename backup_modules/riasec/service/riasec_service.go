package service

import (
	"context"
	"rextra-backend/internal/modules/riasec/repository"
	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	myerror "rextra-backend/internal/pkg/error"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	RiasecService interface {
		Create(ctx context.Context, req dto_request.CreateRiasecRequest) (dto_response.CreateRiasecResponse, error)
		GetByUserID(ctx context.Context, userId string) (dto_response.CreateRiasecResponse, error)
		Update(ctx context.Context, id string, req dto_request.CreateRiasecRequest) (dto_response.CreateRiasecResponse, error)
	}

	riasecService struct {
		riasecRepository repository.RiasecRepository
		db               *gorm.DB
	}
)

func NewRiasec(riasecRepository repository.RiasecRepository, db *gorm.DB) RiasecService {
	return &riasecService{
		riasecRepository: riasecRepository,
		db:               db,
	}
}

func (s *riasecService) Create(ctx context.Context, req dto_request.CreateRiasecRequest) (dto_response.CreateRiasecResponse, error) {
	_, err := s.riasecRepository.GetByUserID(ctx, nil, req.UserID)
	if err == nil {
		return dto_response.CreateRiasecResponse{}, myerror.RecordAlreadyExist("riasec")
	}

	createResult, err := s.riasecRepository.Create(ctx, nil, entity.Riasec{
		UserID:       uuid.MustParse(req.UserID),
		RiasecCode:   req.RiasecCode,
		RDescription: req.Explanations.R,
		IDescription: req.Explanations.I,
		ADescription: req.Explanations.A,
		SDescription: req.Explanations.S,
		EDescription: req.Explanations.E,
		CDescription: req.Explanations.C,
	})
	if err != nil {
		return dto_response.CreateRiasecResponse{}, err
	}

	return dto_response.CreateRiasecResponse{
		ID:         createResult.ID.String(),
		RiasecCode: createResult.RiasecCode,
		Explanations: dto_response.RiasecCodeExplanationResponse{
			R: createResult.RDescription,
			I: createResult.IDescription,
			A: createResult.ADescription,
			S: createResult.SDescription,
			E: createResult.EDescription,
			C: createResult.CDescription,
		},
	}, nil
}

func (s *riasecService) GetByUserID(ctx context.Context, userId string) (dto_response.CreateRiasecResponse, error) {
	riasec, err := s.riasecRepository.GetByUserID(ctx, nil, userId)
	if err != nil {
		return dto_response.CreateRiasecResponse{}, err
	}

	return dto_response.CreateRiasecResponse{
		ID:         riasec.ID.String(),
		RiasecCode: riasec.RiasecCode,
		Explanations: dto_response.RiasecCodeExplanationResponse{
			R: riasec.RDescription,
			I: riasec.IDescription,
			A: riasec.ADescription,
			S: riasec.SDescription,
			E: riasec.EDescription,
			C: riasec.CDescription,
		},
	}, nil
}

func (s *riasecService) Update(ctx context.Context, id string, req dto_request.CreateRiasecRequest) (dto_response.CreateRiasecResponse, error) {
	riasec, err := s.riasecRepository.GetByUserID(ctx, nil, id)
	if err != nil {
		return dto_response.CreateRiasecResponse{}, err
	}

	riasec.RiasecCode = req.RiasecCode
	riasec.RDescription = req.Explanations.R
	riasec.IDescription = req.Explanations.I
	riasec.ADescription = req.Explanations.A
	riasec.SDescription = req.Explanations.S
	riasec.EDescription = req.Explanations.E
	riasec.CDescription = req.Explanations.C

	updateResult, err := s.riasecRepository.Update(ctx, nil, riasec)
	if err != nil {
		return dto_response.CreateRiasecResponse{}, err
	}

	return dto_response.CreateRiasecResponse{
		ID:         updateResult.ID.String(),
		RiasecCode: updateResult.RiasecCode,
		Explanations: dto_response.RiasecCodeExplanationResponse{
			R: updateResult.RDescription,
			I: updateResult.IDescription,
			A: updateResult.ADescription,
			S: updateResult.SDescription,
			E: updateResult.EDescription,
			C: updateResult.CDescription,
		},
	}, nil
}
