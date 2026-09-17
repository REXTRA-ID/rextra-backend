package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"rextra-backend/internal/api/repository"
	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	myerror "rextra-backend/internal/pkg/error"

	"gorm.io/gorm"
)

type (
	AssesmentService interface {
		ValidateHash(ctx context.Context, req dto_request.ValidateHashRequest, cookieHeader string) ([]string, error)
		GetRiasecQuestion(ctx context.Context) ([]dto_response.RiasecQuestionResponse, error)
		SubmitRiasecAnswer(ctx context.Context, req dto_request.RiasecQuestionSubmitRequest, userID string, cookieHeader string) (dto_response.RiasecQuestionSubmitResponse, []string, error)
		GetRiasecResult(ctx context.Context, userID string) ([]dto_response.RiasecResultResponse, error)
		GetIkigaiQuestion(ctx context.Context, cookieHeader string) (dto_response.IkigaiQuestionResponse, []string, error)
		SubmitIkigaiAnswer(ctx context.Context, userID string, cookieHeader string, req dto_request.IkigaiQuestionSubmitRequest) (dto_response.IkigaiQuestionSubmitResponse, []string, error)
		GetIkigaiResult(ctx context.Context, userID string) ([]dto_response.IkigaiResultResponse, error)
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

func (s *assesmentService) ValidateHash(ctx context.Context, req dto_request.ValidateHashRequest, cookieHeader string) ([]string, error) {
	if req.Hash == "" {
		return nil, myerror.InvalidRequest(myerror.New("Voucher code cannot be empty", myerror.Error_InvalidRequest))
	}

	var voucher entity.Voucher
	if err := s.db.WithContext(ctx).Where("code = ?", req.Hash).First(&voucher).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, myerror.InvalidRequest(myerror.New("Invalid voucher code", myerror.Error_InvalidRequest))
		}
		return nil, myerror.ProcessingError(err)
	}

	if voucher.IsUsed {
		return nil, myerror.InvalidRequest(myerror.New("Voucher code has already been used", myerror.Error_InvalidRequest))
	}

	// Wait, we don't mark it as used yet? Or do we? Let's just mark it as used.
	// Normally we would associate it with UserID, but we might not have UserID in ctx for ValidateHash yet?
	// The controller doesn't pass userId to ValidateHash. Let's just mark IsUsed=true for now.
	// Or maybe just leave it as validated, and let the submission process mark it used?
	// For simplicity, let's just mark it used here if it's meant to be a single-use token to start a session.
	voucher.IsUsed = true
	if err := s.db.WithContext(ctx).Save(&voucher).Error; err != nil {
		return nil, myerror.ProcessingError(err)
	}

	return nil, nil
}

func (s *assesmentService) GetRiasecQuestion(ctx context.Context) ([]dto_response.RiasecQuestionResponse, error) {
	// Not supported in AI backend; mocked in frontend.
	return nil, myerror.ProcessingError(myerror.New("RIASEC questions are not served by the backend anymore. Please use local frontend mock.", myerror.SystemError))
}

func (s *assesmentService) SubmitRiasecAnswer(ctx context.Context, req dto_request.RiasecQuestionSubmitRequest, userID string, cookieHeader string) (dto_response.RiasecQuestionSubmitResponse, []string, error) {
	base := os.Getenv("AI_BACKEND_URL")
	if base == "" {
		base = "http://localhost:8010"
	}

	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		base = "http://" + base
	}
	base = strings.TrimRight(base, "/")
	url := base + "/api/v1/career-profile/riasec/submit"

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return dto_response.RiasecQuestionSubmitResponse{}, nil,myerror.ProcessingError(err)
	}


	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return dto_response.RiasecQuestionSubmitResponse{}, nil,myerror.ProcessingError(err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if userID != "" {
		httpReq.Header.Set("X-User-Id", userID)
	}
	if cookieHeader != "" {
		httpReq.Header.Set("Cookie", cookieHeader)
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return dto_response.RiasecQuestionSubmitResponse{}, nil,myerror.ProcessingError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto_response.RiasecQuestionSubmitResponse{}, nil, myerror.ProcessingError(err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errorResp struct {
			Error string `json:"error"`
		}
		_ = json.Unmarshal(body, &errorResp)
		msg := strings.TrimSpace(errorResp.Error)
		if msg == "" {
			msg = string(body)
		}
		return dto_response.RiasecQuestionSubmitResponse{}, nil, myerror.InvalidRequest(myerror.New(msg, myerror.Error_InvalidRequest))
	}

	setCookies := resp.Header["Set-Cookie"]

	var result dto_response.RiasecQuestionSubmitResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return dto_response.RiasecQuestionSubmitResponse{}, nil, myerror.ProcessingError(err)
	}

	return result, setCookies, nil
}

func (s *assesmentService) GetRiasecResult(ctx context.Context, userID string) ([]dto_response.RiasecResultResponse, error) {
	// Results are now fetched directly from AI Backend using session_token or user-profile endpoints.
	return nil, myerror.ProcessingError(myerror.New("Results are now stored in AI Backend. Use AI Backend GET /career-profile/user-profile or /result/{session_token}", myerror.SystemError))
}

func (s *assesmentService) GetIkigaiQuestion(ctx context.Context, cookieHeader string) (dto_response.IkigaiQuestionResponse, []string,error) {
	// Not supported natively in Go without session token via POST now.
	// Returning dummy error to force frontend to use the new AI proxy routes directly or adapt.
	return dto_response.IkigaiQuestionResponse{}, nil, myerror.ProcessingError(myerror.New("Ikigai questions are now initialized via AI Backend POST /career-profile/ikigai/start with session_token", myerror.SystemError))
}


func (s *assesmentService) SubmitIkigaiAnswer(ctx context.Context, userID string, cookieHeader string, req dto_request.IkigaiQuestionSubmitRequest) (dto_response.IkigaiQuestionSubmitResponse, []string, error) {
	base := os.Getenv("AI_BACKEND_URL")
	if base == "" {
		base = "http://localhost:8010"
	}

	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		base = "http://" + base
	}
	base = strings.TrimRight(base, "/")
	url := base + "/api/v1/career-profile/ikigai/submit-with-clicks"

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return dto_response.IkigaiQuestionSubmitResponse{},nil, myerror.ProcessingError(err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return dto_response.IkigaiQuestionSubmitResponse{},nil, myerror.ProcessingError(err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if userID != "" {
		httpReq.Header.Set("X-User-Id", userID)
	}

	if cookieHeader != "" {
		httpReq.Header.Set("Cookie", cookieHeader)
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return dto_response.IkigaiQuestionSubmitResponse{},nil, myerror.ProcessingError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	
	if err != nil {
		return dto_response.IkigaiQuestionSubmitResponse{},nil, myerror.ProcessingError(err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errorResp struct {
			Error string `json:"error"`
		}
		_ = json.Unmarshal(body, &errorResp)
		msg := strings.TrimSpace(errorResp.Error)
		if msg == "" {
			msg = string(body)
		}
		return dto_response.IkigaiQuestionSubmitResponse{}, nil, myerror.InvalidRequest(myerror.New(msg, myerror.Error_InvalidRequest))
	}

	setCookies := resp.Header["Set-Cookie"]

	var result dto_response.IkigaiQuestionSubmitResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return dto_response.IkigaiQuestionSubmitResponse{}, nil, myerror.ProcessingError(err)
	}

	return result, setCookies, nil
}

func (s *assesmentService) GetIkigaiResult(ctx context.Context, userID string) ([]dto_response.IkigaiResultResponse, error) {
	// Results are now fetched directly from AI Backend.
	return nil, myerror.ProcessingError(myerror.New("Results are now stored in AI Backend. Use AI Backend GET /career-profile/user-profile or /result/{session_token}", myerror.SystemError))
}