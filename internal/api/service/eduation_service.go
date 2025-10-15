package service

import (
	"context"
	"errors"
	"strings"

	"rextra-backend/internal/api/repository"
	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	myerror "rextra-backend/internal/pkg/error"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	EducationService interface {
		CreateEducation(ctx context.Context, education dto_request.CreateEducationRequest) (dto_response.EducationResponse, error)
		GetAllEducation(ctx context.Context, userID string) ([]dto_response.EducationResponse, error)
		GetEducationById(ctx context.Context, userID string, educationID string) (dto_response.EducationResponse, error)
	}
	educationService struct {
		educationRepository repository.EducationRepository
		db *gorm.DB
	}
)

func NewEducation(educationRepository repository.EducationRepository, db *gorm.DB) EducationService {
	return &educationService{
		educationRepository: educationRepository,
		db: db,
	}
}

func (s *educationService) CreateEducation(ctx context.Context, education dto_request.CreateEducationRequest) (dto_response.EducationResponse, error) {
	userId, err := uuid.Parse(education.UserID)
	if err != nil {
		return dto_response.EducationResponse{}, err
	}

	if err := EducationLevelValidation(education.EducationLevel); err != nil {
		return dto_response.EducationResponse{}, err
	}

	if err := EducationStatusValidation(education.Status); err != nil {
		return dto_response.EducationResponse{}, err
	}

	if education.ActualGraduationYear != nil {
		if *education.ActualGraduationYear < education.EntryYear {
			return dto_response.EducationResponse{}, errors.New("actual graduation year cannot be less than entry year")
		}
	}

	var actualGraduationYear int
	if education.ActualGraduationYear != nil {
		actualGraduationYear = *education.ActualGraduationYear
	}
	
	var isActive bool
	if strings.ToUpper(education.Status) == "ACTIVE" {
		_, flag , err := s.educationRepository.GetActiveEducationByUserId(ctx, s.db, education.UserID)
		if err != nil {
			return dto_response.EducationResponse{}, err
		}
		if flag {
			return dto_response.EducationResponse{}, errors.New("user already has active education")
		}
		isActive = true
	}

	createRequest := entity.Education{
		UserID:                 userId,
		InstitutionName:        education.InstitutionName,
		Major:                  education.Major,
		Faculty:                education.Faculty,
		EntryYear:              education.EntryYear,
		ExpectedGraduationYear: education.ExpectedGraduationYear,
		ActualGraduationYear:   actualGraduationYear,
		CurrentSemester:        education.CurrentSemester,
		TotalSemester:          education.TotalSemester,
		EducationLevel:         entity.EducationLevel(education.EducationLevel),
		Status:                 entity.EducationStatus(education.Status),
		IsActive:               isActive,
	}

	CreateEducation, err := s.educationRepository.Create(ctx, s.db, createRequest)

	if err != nil {
		return dto_response.EducationResponse{}, err
	}

	return dto_response.EducationResponse{
		ID:                     CreateEducation.ID.String(),
		UserID:                 CreateEducation.UserID.String(),
		InstitutionName:        CreateEducation.InstitutionName,
		Major:                  CreateEducation.Major,
		Faculty:                CreateEducation.Faculty,
		EntryYear:              CreateEducation.EntryYear,
		ExpectedGraduationYear: CreateEducation.ExpectedGraduationYear,
		ActualGraduationYear:   &CreateEducation.ActualGraduationYear,
		CurrentSemester:        CreateEducation.CurrentSemester,
		TotalSemester:          CreateEducation.TotalSemester,
		EducationLevel:         string(CreateEducation.EducationLevel),
		Status:                 string(CreateEducation.Status),
		IsActive:               CreateEducation.IsActive,
	}, nil
}

func (s *educationService) GetAllEducation(ctx context.Context, userID string) ([]dto_response.EducationResponse, error) {
	allEducation, err := s.educationRepository.GetAllEducationByUserId(ctx, s.db, userID)
	if err != nil {
		return nil, err
	}

	var educationResponse []dto_response.EducationResponse = []dto_response.EducationResponse{}
	for _, education := range allEducation {
		educationResponse = append(educationResponse, dto_response.EducationResponse{
			ID:                     education.ID.String(),
			UserID:                 education.UserID.String(),
			InstitutionName:        education.InstitutionName,
			Major:                  education.Major,
			Faculty:                education.Faculty,
			EntryYear:              education.EntryYear,
			ExpectedGraduationYear: education.ExpectedGraduationYear,
			ActualGraduationYear:   &education.ActualGraduationYear,
			CurrentSemester:        education.CurrentSemester,
			TotalSemester:          education.TotalSemester,
			EducationLevel:         string(education.EducationLevel),
			Status:                 string(education.Status),
			IsActive:               education.IsActive,
		})
	}

	return educationResponse, nil
}

func (s *educationService) GetEducationById(ctx context.Context, userID string, educationID string) (dto_response.EducationResponse, error) {
	education, flag, err := s.educationRepository.GetByUserIdAndEducationById(ctx, s.db, userID, educationID)
	if err != nil {
		return dto_response.EducationResponse{}, err
	}

	if !flag {
		return dto_response.EducationResponse{}, myerror.RecordNotFound("education")
	}

	return dto_response.EducationResponse{
		ID:                     education.ID.String(),
		UserID:                 education.UserID.String(),
		InstitutionName:        education.InstitutionName,
		Major:                  education.Major,
		Faculty:                education.Faculty,
		EntryYear:              education.EntryYear,
		ExpectedGraduationYear: education.ExpectedGraduationYear,
		ActualGraduationYear:   &education.ActualGraduationYear,
		CurrentSemester:        education.CurrentSemester,
		TotalSemester:          education.TotalSemester,
		EducationLevel:         string(education.EducationLevel),
		Status:                 string(education.Status),
		IsActive:               education.IsActive,
	}, nil
}

func EducationLevelValidation(educationLevel string) error {
	educationLevel = strings.ToUpper(educationLevel)
	if educationLevel != "D1" && educationLevel != "D2" && educationLevel != "D3" && educationLevel != "D4" && educationLevel != "S1" && educationLevel != "S2" && educationLevel != "S3" {
		return errors.New("education level is not valid")
	}
	return nil
}

func EducationStatusValidation(educationStatus string) error {
	educationStatus = strings.ToUpper(educationStatus)
	if educationStatus != "ACTIVE" && educationStatus != "GRADUATED" && educationStatus != "DROPPED" && educationStatus != "DEFERRED" && educationStatus != "TRANSFERRED" && educationStatus != "DISMISSED" {
		return errors.New("education status is not valid")
	}
	return nil
}
