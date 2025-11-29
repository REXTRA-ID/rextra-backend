package service

import (
	"context"
	"math"
	"rextra-backend/internal/api/repository"
	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	myerror "rextra-backend/internal/pkg/error"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	missionLimitPathfinder = 4
	missionLimitBuilder    = 7
	missionLimitAchiever   = 9
)

var staticMissions = []dto_response.PersonaMission{
	{Key: "education_saved", Name: "Simpan Data Pendidikan", Description: "Lengkapi informasi pendidikan terkini dan terdahulu", Order: 1},
	{Key: "career_recommendation_tried", Name: "Coba Kenali Diri", Description: "Gunakan fitur Kenali Diri untuk rekomendasi karier", Order: 2},
	{Key: "career_dictionary_accessed", Name: "Buka Kamus Karier", Description: "Jelajahi dunia kerja digital melalui Kamus Karier", Order: 3},
	{Key: "career_plan_created", Name: "Buat Rencana Karier", Description: "Mulai membuat rencana karier yang jelas dan spesifik", Order: 4},
	{Key: "portfolio_recorded", Name: "Selesaikan Portofolio", Description: "Catat & rencanakan aktivitas portofolio setiap semester", Order: 5},
	{Key: "exploration_ai_used", Name: "Rekomendasi portofolio", Description: "Gunakan rekomendasi pengisian portofolio karier", Order: 6},
	{Key: "cv_created", Name: "Buat CV", Description: "Buat CV profesional dan relevan dengan mudahnya", Order: 7},
	{Key: "interview_simulated", Name: "Simulasi interview", Description: "Penuhi kebutuhan seleksi kerja dengan mentoring", Order: 8},
	{Key: "linkedin_optimize", Name: "Optimasi LinkedIn", Description: "Optimasi branding dengan LinkedIn Analyzer", Order: 9},
}

type (
	PersonaService interface {
		Create(ctx context.Context, req dto_request.CreatePersonaRequest) (dto_response.CreatePersonaResponse, error)
		Get(ctx context.Context, userID string) (*dto_response.GetPersonaResponse, error)
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
	if _, exists, _ := s.personaRepository.GetByUserID(ctx, nil, req.UserID); exists {
		return dto_response.CreatePersonaResponse{}, myerror.RecordAlreadyExist("persona")
	}

	userPersonaType := determinePersonaType(req)

	newPersona := entity.Persona{
		UserID:      uuid.MustParse(req.UserID),
		PersonaType: userPersonaType,
	}

	createdResult, err := s.personaRepository.Create(ctx, nil, newPersona)
	if err != nil {
		return dto_response.CreatePersonaResponse{}, err
	}

	userMissions := getMissionsWithStatus(createdResult)
	progress := calculateProgress(createdResult, len(userMissions))

	return dto_response.CreatePersonaResponse{
		ID:               createdResult.ID.String(),
		UserID:           createdResult.UserID.String(),
		PersonaType:      string(createdResult.PersonaType),
		PersonaStartedAt: createdResult.CreatedAt.String(),
		Progress:         progress,
		Missions:         userMissions,
	}, nil
}

func (s *personaService) Get(ctx context.Context, userID string) (*dto_response.GetPersonaResponse, error) {
	persona, found, err := s.personaRepository.GetByUserID(ctx, nil, userID)
	if err != nil {
		return &dto_response.GetPersonaResponse{}, err
	}
	if !found {
		return nil, nil
	}

	userMissions := getMissionsWithStatus(persona)
	progress := calculateProgress(persona, len(userMissions))

	return &dto_response.GetPersonaResponse{
		ID:               persona.ID.String(),
		UserID:           persona.UserID.String(),
		PersonaType:      string(persona.PersonaType),
		PersonaStartedAt: persona.CreatedAt.String(),
		Progress:         progress,
		Missions:         userMissions,
	}, nil
}

func determinePersonaType(req dto_request.CreatePersonaRequest) entity.PersonaType {
	if req.HasCareerGoal && req.BuildingPortofolio && req.InRecruitmentProcess {
		return entity.Achiever
	}
	if !req.HasCareerGoal && !req.BuildingPortofolio && !req.InRecruitmentProcess {
		return entity.Pathfinder
	}
	return entity.Builder
}

func getMissionsWithStatus(p entity.Persona) []dto_response.PersonaMission {
	var limit int
	switch p.PersonaType {
	case entity.Pathfinder:
		limit = missionLimitPathfinder
	case entity.Builder:
		limit = missionLimitBuilder
	case entity.Achiever:
		limit = missionLimitAchiever
	default:
		return []dto_response.PersonaMission{}
	}

	statusMap := map[string]bool{
		"education_saved":             p.EducationSaved,
		"career_recommendation_tried": p.CareerRecommendationTired,
		"career_dictionary_accessed":  p.CareerDictionaryAccessed,
		"career_plan_created":         p.CareerPlanCreated,
		"portfolio_recorded":          p.PorfolioRecorded,
		"exploration_ai_used":         p.ExplorationAIUsed,
		"cv_created":                  p.CVCreated,
		"interview_simulated":         p.InterviewSimulated,
		"linkedin_optimize":           p.LinkedinOptimaze,
	}

	templateSlice := staticMissions[:limit]

	result := make([]dto_response.PersonaMission, len(templateSlice))

	for i, m := range templateSlice {
		result[i] = m
		result[i].IsCompleted = statusMap[m.Key]
	}

	return result
}

func getNextPersona(userPersonaType entity.PersonaType) entity.PersonaType {
	switch userPersonaType {
	case entity.Pathfinder:
		return entity.Builder
	case entity.Builder:
		return entity.Achiever
	default:
		return entity.Achiever
	}
}

func calculateProgress(p entity.Persona, totalMissions int) dto_response.PersonaProgress {
	completedCount := countCompletedMissions(p)

	percentage := 0.0
	if totalMissions > 0 {
		percentage = (float64(completedCount) / float64(totalMissions)) * 100
		percentage = math.Round(percentage*100) / 100
	}

	return dto_response.PersonaProgress{
		CompletedMissions: completedCount,
		TotalMissions:     totalMissions,
		Percentage:        percentage,
		CurrentPhase:      string(p.PersonaType),
		NextPersona:       string(getNextPersona(p.PersonaType)),
	}
}

func countCompletedMissions(p entity.Persona) int {
	count := 0
	if p.EducationSaved {
		count++
	}
	if p.CareerRecommendationTired {
		count++
	}
	if p.CareerDictionaryAccessed {
		count++
	}
	if p.CareerPlanCreated {
		count++
	}
	if p.PorfolioRecorded {
		count++
	}
	if p.ExplorationAIUsed {
		count++
	}
	if p.CVCreated {
		count++
	}
	if p.InterviewSimulated {
		count++
	}
	if p.LinkedinOptimaze {
		count++
	}
	if p.IntershipPlanReported {
		count++
	}
	return count
}
