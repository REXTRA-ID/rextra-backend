package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
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
	AssesmentService interface {
		ValidateHash(ctx context.Context, req dto_request.ValidateHashRequest) error
		GetRiasecQuestion(ctx context.Context) ([]dto_response.RiasecQuestionResponse, error)
		SubmitRiasecAnswer(ctx context.Context, req dto_request.RiasecQuestionSubmitRequest, userID string) (dto_response.RiasecQuestionSubmitResponse, []string, error)
		GetRiasecResult(ctx context.Context, userID string) ([]dto_response.RiasecResultResponse, error)
		GetIkigaiQuestion(ctx context.Context, cookieHeader string) (dto_response.IkigaiQuestionResponse, error)
		// SubmitIkigaiAnswer(ctx *gin.Context)
		// GetIkigaiResult(ctx *gin.Context)
	}

	assesmentService struct {
		assesmentRepository repository.AssesmentRepository
		db                *gorm.DB
	}
)

func NewAssesment(assesmentRepository repository.AssesmentRepository, db *gorm.DB) AssesmentService {
	return &assesmentService{
		assesmentRepository: assesmentRepository,
		db:                db,
	}
}

func (s *assesmentService) ValidateHash(ctx context.Context, req dto_request.ValidateHashRequest) error {
	base := os.Getenv("MONGODB_BACKEND")
	if base == "" {
		return myerror.ProcessingError(myerror.New("MONGODB_BACKEND not set", myerror.SystemError))
	}

	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		base = "http://" + base
	}
	base = strings.TrimRight(base, "/")
	url := base + "/api/assessment/validate_hash"

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return myerror.ProcessingError(err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return myerror.ProcessingError(err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return myerror.ProcessingError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return myerror.ProcessingError(myerror.New("external api returned non-2xx status", myerror.SystemError))
	}

	return nil
}

func (s *assesmentService) GetRiasecQuestion(ctx context.Context) ([]dto_response.RiasecQuestionResponse, error) {
	base := os.Getenv("MONGODB_BACKEND")
	if base == "" {
		return nil, myerror.ProcessingError(myerror.New("MONGODB_BACKEND not set", myerror.SystemError))
	}

	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		base = "http://" + base
	}
	base = strings.TrimRight(base, "/")
	url := base + "/api/riasec/questions"

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, myerror.ProcessingError(err)
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, myerror.ProcessingError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, myerror.ProcessingError(myerror.New("external api returned non-2xx status", myerror.SystemError))
	}

	var result []dto_response.RiasecQuestionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, myerror.ProcessingError(err)
	}

	return result, nil
}

func (s *assesmentService) SubmitRiasecAnswer(ctx context.Context, req dto_request.RiasecQuestionSubmitRequest, userID string) (dto_response.RiasecQuestionSubmitResponse, []string, error) {
	base := os.Getenv("MONGODB_BACKEND")
	if base == "" {
		return dto_response.RiasecQuestionSubmitResponse{}, nil, myerror.ProcessingError(myerror.New("MONGODB_BACKEND not set", myerror.SystemError))
	}

	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		base = "http://" + base
	}
	base = strings.TrimRight(base, "/")
	url := base + "/api/riasec/submit"

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return dto_response.RiasecQuestionSubmitResponse{}, nil,myerror.ProcessingError(err)
	}


	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return dto_response.RiasecQuestionSubmitResponse{}, nil,myerror.ProcessingError(err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return dto_response.RiasecQuestionSubmitResponse{}, nil,myerror.ProcessingError(err)
	}
	defer resp.Body.Close()
	

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return dto_response.RiasecQuestionSubmitResponse{}, nil, myerror.ProcessingError(myerror.New("external api returned non-2xx status", myerror.SystemError))
	}

	setCookies := resp.Header["Set-Cookie"]

	var result dto_response.RiasecQuestionSubmitResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return dto_response.RiasecQuestionSubmitResponse{}, nil, myerror.ProcessingError(err)
	}

	if _, err := s.assesmentRepository.CreateUserRiasec(ctx, s.db, entity.UserRiasec{
		UserID: uuid.MustParse(userID),
		Profile: result.Profile,
		NormalizedScores: result.NormalizedScores,
	}); err != nil {
		return dto_response.RiasecQuestionSubmitResponse{}, nil, myerror.ProcessingError(err)
	}

	return result, setCookies, nil
}

func (s *assesmentService) GetRiasecResult(ctx context.Context, userID string) ([]dto_response.RiasecResultResponse, error) {
	userRiasec, err := s.assesmentRepository.GetUserRiasec(ctx, s.db, userID)
	if err != nil {
		return nil, myerror.ProcessingError(err)
	}

	var result []dto_response.RiasecResultResponse
	for _, v := range userRiasec {
		result = append(result, dto_response.RiasecResultResponse{
			ID:             v.ID.String(),
			Profile:          v.Profile,
			NormalizedScores: v.NormalizedScores,
			CreatedAt:        v.CreatedAt,
		})
	}

	return result, nil	
}

func (s *assesmentService) GetIkigaiQuestion(ctx context.Context, cookieHeader string) (dto_response.IkigaiQuestionResponse, error) {
	base := os.Getenv("MONGODB_BACKEND")
	if base == "" {
		return dto_response.IkigaiQuestionResponse{}, myerror.ProcessingError(myerror.New("MONGODB_BACKEND not set", myerror.SystemError))
	}

	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		base = "http://" + base
	}
	base = strings.TrimRight(base, "/")
	url := base + "/api/assessment/start"

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return dto_response.IkigaiQuestionResponse{}, myerror.ProcessingError(err)
	}
	
	httpReq.Header.Set("Content-Type", "application/json")

	if cookieHeader != "" {
		httpReq.Header.Set("Cookie", cookieHeader)
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return dto_response.IkigaiQuestionResponse{}, myerror.ProcessingError(err)
	}
	
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return dto_response.IkigaiQuestionResponse{}, myerror.ProcessingError(myerror.New("external api returned non-2xx status", myerror.SystemError))
	}

	var result dto_response.IkigaiQuestionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return dto_response.IkigaiQuestionResponse{}, myerror.ProcessingError(err)
	}

	return result, nil
}