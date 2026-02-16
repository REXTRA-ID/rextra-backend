package service

import (
	"context"
	"rextra-backend/internal/modules/persona/repository"
	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	myerror "rextra-backend/internal/pkg/error"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	PersonaService interface {
		Create(ctx context.Context, req dto_request.CreatePersonaRequest) (dto_response.CreatePersonaResponse, error)
		GetByUserID(ctx context.Context, userId string) (dto_response.CreatePersonaResponse, error)
		Update(ctx context.Context, id string, req dto_request.CreatePersonaRequest) (dto_response.CreatePersonaResponse, error)
	}

	personaService struct {
		personaRepository repository.PersonaRepository
		db                *gorm.DB
	}
)

func NewPersona(personaRepository repository.PersonaRepository, db *gorm.DB) PersonaService {
	return &personaService{
		personaRepository: personaRepository,
		db:                db,
	}
}

func (s *personaService) Create(ctx context.Context, req dto_request.CreatePersonaRequest) (dto_response.CreatePersonaResponse, error) {
	_, err := s.personaRepository.GetByUserID(ctx, nil, req.UserID)
	if err == nil {
		return dto_response.CreatePersonaResponse{}, myerror.RecordAlreadyExist("persona")
	}

	createResult, err := s.personaRepository.Create(ctx, nil, entity.Persona{
		UserID:         uuid.MustParse(req.UserID),
		Institution:    req.Institution,
		Study:          req.Study,
		EducationLevel: req.EducationLevel,
		GraduationYear: req.GraduationYear,
		CareerPlan:     req.CareerPlan,
		CareerDreams:   req.CareerDreams,
		Portfolio:      req.Portfolio,
		Application:    req.Application,
		Status:         req.Status,
	})
	if err != nil {
		return dto_response.CreatePersonaResponse{}, err
	}

	return dto_response.CreatePersonaResponse{
		ID:             createResult.ID.String(),
		UserID:         createResult.UserID.String(),
		Institution:    createResult.Institution,
		Study:          createResult.Study,
		EducationLevel: createResult.EducationLevel,
		GraduationYear: createResult.GraduationYear,
		CareerPlan:     createResult.CareerPlan,
		CareerDreams:   createResult.CareerDreams,
		Portfolio:      createResult.Portfolio,
		Application:    createResult.Application,
		Status:         createResult.Status,
	}, nil
}

func (s *personaService) GetByUserID(ctx context.Context, userId string) (dto_response.CreatePersonaResponse, error) {
	persona, err := s.personaRepository.GetByUserID(ctx, nil, userId)
	if err != nil {
		return dto_response.CreatePersonaResponse{}, err
	}

	return dto_response.CreatePersonaResponse{
		ID:             persona.ID.String(),
		UserID:         persona.UserID.String(),
		Institution:    persona.Institution,
		Study:          persona.Study,
		EducationLevel: persona.EducationLevel,
		GraduationYear: persona.GraduationYear,
		CareerPlan:     persona.CareerPlan,
		CareerDreams:   persona.CareerDreams,
		Portfolio:      persona.Portfolio,
		Application:    persona.Application,
		Status:         persona.Status,
	}, nil
}

func (s *personaService) Update(ctx context.Context, id string, req dto_request.CreatePersonaRequest) (dto_response.CreatePersonaResponse, error) {
	persona, err := s.personaRepository.GetByUserID(ctx, nil, id)
	if err != nil {
		return dto_response.CreatePersonaResponse{}, err
	}

	persona.Institution = req.Institution
	persona.Study = req.Study
	persona.EducationLevel = req.EducationLevel
	persona.GraduationYear = req.GraduationYear
	persona.CareerPlan = req.CareerPlan
	persona.CareerDreams = req.CareerDreams
	persona.Portfolio = req.Portfolio
	persona.Application = req.Application
	persona.Status = req.Status

	updateResult, err := s.personaRepository.Update(ctx, nil, persona)
	if err != nil {
		return dto_response.CreatePersonaResponse{}, err
	}

	return dto_response.CreatePersonaResponse{
		ID:             updateResult.ID.String(),
		UserID:         updateResult.UserID.String(),
		Institution:    updateResult.Institution,
		Study:          updateResult.Study,
		EducationLevel: updateResult.EducationLevel,
		GraduationYear: updateResult.GraduationYear,
		CareerPlan:     updateResult.CareerPlan,
		CareerDreams:   updateResult.CareerDreams,
		Portfolio:      updateResult.Portfolio,
		Application:    updateResult.Application,
		Status:         updateResult.Status,
	}, nil
}
