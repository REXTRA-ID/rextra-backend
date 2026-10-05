package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/firestore"

	"rextra-backend/internal/api/repository"
	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	myerror "rextra-backend/internal/pkg/error"

	"gorm.io/gorm"
)

type (
	AssesmentService interface {
		ValidateHash(ctx context.Context, req dto_request.ValidateHashRequest, cookieHeader string, userID string) ([]string, error)
		GetRiasecQuestion(ctx context.Context) ([]dto_response.RiasecQuestionResponse, error)
		SubmitRiasecAnswer(ctx context.Context, req dto_request.RiasecQuestionSubmitRequest, userID string, cookieHeader string) (dto_response.RiasecQuestionSubmitResponse, []string, error)
		GetRiasecResult(ctx context.Context, userID string) ([]dto_response.RiasecResultResponse, error)
		GetIkigaiQuestion(ctx context.Context, cookieHeader string) (dto_response.IkigaiQuestionResponse, []string, error)
		SubmitIkigaiAnswer(ctx context.Context, userID string, cookieHeader string, req dto_request.IkigaiQuestionSubmitRequest) (dto_response.IkigaiQuestionSubmitResponse, []string, error)
		GetIkigaiResult(ctx context.Context, userID string) ([]dto_response.IkigaiResultResponse, error)
	}

	assesmentService struct {
		assesmentRepository repository.AssesmentRepository
		db *gorm.DB
		firestore *firestore.Client
	}
)

func NewAssesment(assesmentRepository repository.AssesmentRepository, db *gorm.DB, fs *firestore.Client) AssesmentService {
	return &assesmentService{
		assesmentRepository: assesmentRepository,
		db: db,
		firestore: fs,
	}
}

func (s *assesmentService) ValidateHash(ctx context.Context, req dto_request.ValidateHashRequest, cookieHeader string, userID string) ([]string, error) {
	if req.Hash == "" {
		return nil, myerror.InvalidRequest(myerror.New("Voucher code cannot be empty", myerror.Error_InvalidRequest))
	}

	docRef := s.firestore.Collection("vouchers").Doc(req.Hash)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		return nil, myerror.InvalidRequest(myerror.New("Invalid voucher code", myerror.Error_InvalidRequest))
	}

	data := docSnap.Data()
	if isUsed, ok := data["is_used"].(bool); ok && isUsed {
		return nil, myerror.InvalidRequest(myerror.New("Voucher code has already been used", myerror.Error_InvalidRequest))
	}

	updateData := []firestore.Update{
		{Path: "is_used", Value: true},
		{Path: "used_at", Value: time.Now()},
	}

	if userID != "" {
		updateData = append(updateData, firestore.Update{Path: "user_id", Value: userID})
	}

	if _, err := docRef.Update(ctx, updateData); err != nil {
		return nil, myerror.ProcessingError(err)
	}

	return nil, nil
}

func (s *assesmentService) GetRiasecQuestion(ctx context.Context) ([]dto_response.RiasecQuestionResponse, error) {
	base := os.Getenv("AI_BACKEND_URL")
	if base == "" {
		base = "http://localhost:8010"
	}
	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		base = "http://" + base
	}
	base = strings.TrimRight(base, "/")
	url := base + "/api/v1/career-profile/riasec/questions"

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Println("ERROR in GetRiasecQuestion: ", err); return nil, myerror.ProcessingError(err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		fmt.Println("ERROR in GetRiasecQuestion: ", err); return nil, myerror.ProcessingError(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR in GetRiasecQuestion: ", err); return nil, myerror.ProcessingError(err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("ERROR AI returned status: ", resp.StatusCode); return nil, myerror.ProcessingError(myerror.New("AI Backend returned error", myerror.SystemError))
	}

	// Parse JSON
	var aiResp struct {
		Data []dto_response.RiasecQuestionResponse `json:"data"`
	}
	if err := json.Unmarshal(bodyBytes, &aiResp); err != nil {
		fmt.Println("ERROR in GetRiasecQuestion: ", err); return nil, myerror.ProcessingError(err)
	}

	return aiResp.Data, nil
}


type aiSessionStartResponse struct {
	SessionToken string `json:"session_token"`
}

type aiRiasecSubmitPayload struct {
	SessionToken string        `json:"session_token"`
	Responses    []interface{} `json:"responses"`
}

func getRiasecType(qID int) string {
	if qID >= 1 && qID <= 12 { return "R" }
	if qID >= 13 && qID <= 24 { return "I" }
	if qID >= 25 && qID <= 36 { return "A" }
	if qID >= 37 && qID <= 48 { return "S" }
	if qID >= 49 && qID <= 60 { return "E" }
	if qID >= 61 && qID <= 72 { return "C" }
	return "R"
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

	// 1. Get Session Token
	startReqBody, _ := json.Marshal(map[string]string{"persona_type": "PATHFINDER"})
	startReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, base+"/api/v1/career-profile/recommendation/start", bytes.NewReader(startReqBody))
	startReq.Header.Set("Content-Type", "application/json")
	if userID != "" {
		startReq.Header.Set("X-User-Id", userID)
	}
	
	client := &http.Client{Timeout: 10 * time.Second}
	startResp, err := client.Do(startReq)
	if err != nil || startResp.StatusCode != 200 {
		return dto_response.RiasecQuestionSubmitResponse{}, nil, myerror.ProcessingError(fmt.Errorf("Failed to start session"))
	}
	defer startResp.Body.Close()
	
	var sessionData aiSessionStartResponse
	json.NewDecoder(startResp.Body).Decode(&sessionData)
	
	// 2. Format answers
	rawAnswers, ok := req["answers"].([]interface{})
	if !ok {
		return dto_response.RiasecQuestionSubmitResponse{}, nil, myerror.ProcessingError(fmt.Errorf("Invalid answers format"))
	}
	
	var aiResponses []interface{}
	for i, ans := range rawAnswers {
		val, _ := ans.(float64)
		qID := i + 1
		aiResponses = append(aiResponses, map[string]interface{}{
			"question_id": qID,
			"question_type": getRiasecType(qID),
			"answer_value": int(val),
			"answered_at": time.Now().Format(time.RFC3339),
		})
	}
	
	submitPayload := aiRiasecSubmitPayload{
		SessionToken: sessionData.SessionToken,
		Responses: aiResponses,
	}
	
	submitBody, _ := json.Marshal(submitPayload)
	submitReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, base+"/api/v1/career-profile/riasec/submit", bytes.NewReader(submitBody))
	submitReq.Header.Set("Content-Type", "application/json")
	if userID != "" {
		submitReq.Header.Set("X-User-Id", userID)
	}
	
	submitResp, err := client.Do(submitReq)
	if err != nil {
		return dto_response.RiasecQuestionSubmitResponse{}, nil, myerror.ProcessingError(err)
	}
	defer submitResp.Body.Close()
	
	submitBytes, _ := io.ReadAll(submitResp.Body)
	if submitResp.StatusCode != 200 {
		return dto_response.RiasecQuestionSubmitResponse{}, nil, myerror.ProcessingError(fmt.Errorf("AI Backend submit failed: %s", string(submitBytes)))
	}
	
	var aiSubmitResp dto_response.RiasecQuestionSubmitResponse
	json.Unmarshal(submitBytes, &aiSubmitResp)
	
	// 3. Return response with cookie
	cookie := fmt.Sprintf("rextra_session=%s; Path=/; HttpOnly", sessionData.SessionToken)
	
	return aiSubmitResp, []string{cookie}, nil
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






