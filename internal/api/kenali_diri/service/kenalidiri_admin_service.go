package service

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"rextra-backend/internal/api/kenali_diri/repository"
	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	"rextra-backend/internal/pkg/cache"
	"rextra-backend/internal/pkg/export"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/utils"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const (
	timeFormatDisplay = "02 Jan 2006 15:04"
	dateFormatExport  = "02/01/2006 15:04"
)

type (
	KenalidiriAdminService interface {
		GetTestHistory(ctx context.Context, req dto_request.GetTestHistoryRequest) (dto_response.TestHistoryListResponse, error)
		DeleteTestData(ctx context.Context, req dto_request.DeleteTestDataRequest) error
		ExportTestHistory(ctx context.Context, req dto_request.ExportTestHistoryRequest) (dto_response.ExportFileResponse, error)
		GetTestDetail(ctx context.Context, historyID int64) (dto_response.TestDetailResponse, error)
		GetStudentFeedbackList(ctx context.Context, req dto_request.GetFeedbackListRequest) (dto_response.FeedbackListResponse, error)
		GetStudentFeedbackStats(ctx context.Context, categoryID *int64, timeRange string) (dto_response.FeedbackStatsResponse, error)
		GetExpertFeedbackList(ctx context.Context, req dto_request.GetExpertFeedbackListRequest) (dto_response.ExpertFeedbackListResponse, error)
		GetExpertFeedbackDetail(ctx context.Context, feedbackID int64) (dto_response.ExpertFeedbackDetailResponse, error)
		GetRiasecCodeList(ctx context.Context, req dto_request.GetRiasecCodeListRequest) (dto_response.RiasecCodeListResponse, error)
		GetRiasecCodeDetail(ctx context.Context, codeID int64) (dto_response.RiasecCodeDetailResponse, error)
		UpdateRiasecCode(ctx context.Context, codeID int64, req dto_request.UpdateRiasecCodeRequest) error
	}

	kenalidiriAdminService struct {
		historyRepo        repository.KenalidiriHistoryRepository
		categoryRepo       repository.KenalidiriCategoryRepository
		riasecCodeRepo     repository.RiasecCodeRepository
		testSessionRepo    repository.TestSessionRepository
		riasecRepo         repository.RiasecRepository
		ikigaiRepo         repository.IkigaiRepository
		recommendationRepo repository.RecommendationRepository
		feedbackRepo       repository.FeedbackRepository
		exportService      export.ExportService
		cache              cache.CacheService
		db                 *gorm.DB
	}
)

func NewKenalidiriAdmin(
	historyRepo repository.KenalidiriHistoryRepository,
	categoryRepo repository.KenalidiriCategoryRepository,
	riasecCodeRepo repository.RiasecCodeRepository,
	testSessionRepo repository.TestSessionRepository,
	riasecRepo repository.RiasecRepository,
	ikigaiRepo repository.IkigaiRepository,
	recommendationRepo repository.RecommendationRepository,
	feedbackRepo repository.FeedbackRepository,
	exportService export.ExportService,
	cacheSvc cache.CacheService,
	db *gorm.DB,
) KenalidiriAdminService {
	return &kenalidiriAdminService{
		historyRepo:        historyRepo,
		categoryRepo:       categoryRepo,
		riasecCodeRepo:     riasecCodeRepo,
		testSessionRepo:    testSessionRepo,
		riasecRepo:         riasecRepo,
		ikigaiRepo:         ikigaiRepo,
		recommendationRepo: recommendationRepo,
		feedbackRepo:       feedbackRepo,
		exportService:      exportService,
		cache:              cacheSvc,
		db:                 db,
	}
}

func (s *kenalidiriAdminService) GetTestHistory(ctx context.Context, req dto_request.GetTestHistoryRequest) (dto_response.TestHistoryListResponse, error) {
	page := req.Page
	limit := req.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	filters := repository.KenalidiriHistoryFilters{
		CategoryID: req.CategoryID,
		Status:     req.Status,
		UserName:   req.UserName,
		SortBy:     req.SortBy,
		Limit:      limit,
		Offset:     (page - 1) * limit,
	}

	if req.StartDate != "" {
		if t, err := time.Parse(time.RFC3339, req.StartDate); err == nil {
			filters.StartDate = &t
		}
	}
	if req.EndDate != "" {
		if t, err := time.Parse(time.RFC3339, req.EndDate); err == nil {
			filters.EndDate = &t
		}
	}

	histories, total, err := s.historyRepo.ListWithFilters(ctx, nil, filters)
	if err != nil {
		return dto_response.TestHistoryListResponse{}, err
	}

	items := make([]dto_response.TestHistoryItem, 0, len(histories))
	for _, h := range histories {
		var resultCode string
		var codeType string
		if strings.EqualFold(h.Status, "completed") {
			if r, err := s.riasecRepo.GetResultBySessionID(ctx, nil, h.DetailSessionID); err == nil {
				resultCode = r.RiasecCode.RiasecCode
				codeType = r.RiasecCodeType
			}
		}

		item := dto_response.TestHistoryItem{
			TestID:         fmt.Sprintf("PK%d", h.ID),
			UserName:       h.User.Fullname,
			CategoryName:   h.TestCategory.CategoryName,
			Status:         h.Status,
			ResultCode:     resultCode,
			StartedAt:      h.StartedAt.Format(timeFormatDisplay),
			RiasecCodeType: codeType,
		}
		if h.CompletedAt != nil {
			completed := h.CompletedAt.Format(timeFormatDisplay)
			item.CompletedAt = &completed
		}
		items = append(items, item)
	}

	meta := utils.CalculatePaginationMeta(page, limit, total)
	return dto_response.TestHistoryListResponse{
		Data: items,
		Pagination: dto_response.PaginationMeta{
			CurrentPage:  meta.CurrentPage,
			TotalPages:   meta.TotalPages,
			TotalRecords: meta.TotalRecords,
			PerPage:      meta.PerPage,
		},
	}, nil
}

func (s *kenalidiriAdminService) DeleteTestData(ctx context.Context, req dto_request.DeleteTestDataRequest) error {
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, id := range req.TestIDs {
		history, err := s.historyRepo.GetByID(ctx, tx, id)
		if err != nil {
			tx.Rollback()
			return err
		}

		sessionID := history.DetailSessionID
		relatedDeletes := []struct {
			table interface{}
		}{
			{table: &entity.RiasecQuestionSet{}},
			{table: &entity.RiasecResponse{}},
			{table: &entity.RiasecResult{}},
			{table: &entity.IkigaiCandidateProfession{}},
			{table: &entity.IkigaiResponse{}},
			{table: &entity.IkigaiDimensionScore{}},
			{table: &entity.IkigaiTotalScore{}},
			{table: &entity.CareerRecommendation{}},
		}

		for _, d := range relatedDeletes {
			if err := tx.WithContext(ctx).Where("test_session_id = ?", sessionID).Delete(d.table).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if err := s.historyRepo.BulkDelete(ctx, tx, req.TestIDs); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *kenalidiriAdminService) ExportTestHistory(ctx context.Context, req dto_request.ExportTestHistoryRequest) (dto_response.ExportFileResponse, error) {
	filters := repository.KenalidiriHistoryFilters{
		CategoryID: req.CategoryID,
		Status:     req.Status,
		Limit:      0,
		Offset:     0,
	}

	if req.StartDate != "" {
		if t, err := time.Parse(time.RFC3339, req.StartDate); err == nil {
			filters.StartDate = &t
		}
	}
	if req.EndDate != "" {
		if t, err := time.Parse(time.RFC3339, req.EndDate); err == nil {
			filters.EndDate = &t
		}
	}

	histories, _, err := s.historyRepo.ListWithFilters(ctx, nil, filters)
	if err != nil {
		return dto_response.ExportFileResponse{}, err
	}

	data := make([]map[string]interface{}, 0, len(histories))
	for _, h := range histories {
		var resultCode string
		if strings.EqualFold(h.Status, "completed") {
			if r, err := s.riasecRepo.GetResultBySessionID(ctx, nil, h.DetailSessionID); err == nil {
				resultCode = r.RiasecCode.RiasecCode
			}
		}

		record := map[string]interface{}{
			"ID Tes":         fmt.Sprintf("PK%d", h.ID),
			"Nama Pengguna":  h.User.Fullname,
			"Email":          h.User.Email,
			"Kategori Tes":   h.TestCategory.CategoryName,
			"Status":         h.Status,
			"Waktu Mulai":    h.StartedAt.Format(dateFormatExport),
			"Waktu Selesai":  "",
			"Kode RIASEC":    resultCode,
			"Durasi (menit)": "",
		}

		if h.CompletedAt != nil {
			record["Waktu Selesai"] = h.CompletedAt.Format(dateFormatExport)
			duration := h.CompletedAt.Sub(h.StartedAt).Minutes()
			record["Durasi (menit)"] = fmt.Sprintf("%.0f", duration)
		}

		data = append(data, record)
	}

	filename := fmt.Sprintf("riwayat_tes_%s", time.Now().Format("20060102_150405"))
	var fileURL string
	switch strings.ToLower(req.Format) {
	case "excel":
		fileURL, err = s.exportService.GenerateExcel(data, filename, "Riwayat Tes")
	case "pdf":
		fileURL, err = s.exportService.GeneratePDF(data, filename, "Riwayat Tes")
	default:
		fileURL, err = s.exportService.GenerateCSV(data, filename)
	}
	if err != nil {
		return dto_response.ExportFileResponse{}, err
	}

	return dto_response.ExportFileResponse{
		FileName:     filepath.Base(fileURL),
		FileURL:      fileURL,
		ExportedAt:   time.Now().Format(time.RFC3339),
		TotalRecords: len(data),
	}, nil
}

func (s *kenalidiriAdminService) GetTestDetail(ctx context.Context, historyID int64) (dto_response.TestDetailResponse, error) {
	history, err := s.historyRepo.GetByID(ctx, nil, historyID)
	if err != nil {
		return dto_response.TestDetailResponse{}, err
	}

	sessionID := history.DetailSessionID
	resp := dto_response.TestDetailResponse{
		TestID:       fmt.Sprintf("PK%d", history.ID),
		UserName:     history.User.Fullname,
		CategoryName: history.TestCategory.CategoryName,
		Status:       history.Status,
		StartedAt:    history.StartedAt.Format(timeFormatDisplay),
	}
	if history.CompletedAt != nil {
		completed := history.CompletedAt.Format(timeFormatDisplay)
		resp.CompletedAt = &completed
	}

	if riasec, err := s.riasecRepo.GetResultBySessionID(ctx, nil, sessionID); err == nil {
		resp.RiasecResult = &dto_response.RiasecResultDetail{
			ScoreR:             riasec.ScoreR,
			ScoreI:             riasec.ScoreI,
			ScoreA:             riasec.ScoreA,
			ScoreS:             riasec.ScoreS,
			ScoreE:             riasec.ScoreE,
			ScoreC:             riasec.ScoreC,
			RiasecCode:         riasec.RiasecCode.RiasecCode,
			RiasecTitle:        riasec.RiasecCode.RiasecTitle,
			ClassificationType: riasec.RiasecCodeType,
			IsInconsistent:     riasec.IsInconsistentProfile,
		}
	}

	if ikigai, err := s.ikigaiRepo.GetTotalScores(ctx, nil, sessionID); err == nil {
		resp.IkigaiResult = buildIkigaiResult(ikigai)
		resp.Recommendations = append(resp.Recommendations, buildRecommendations(ikigai)...)
	}

	if rec, err := s.recommendationRepo.GetBySessionID(ctx, nil, sessionID); err == nil {
		resp.Recommendations = mergeRecommendationNarratives(resp.Recommendations, rec)
	}

	return resp, nil
}

func (s *kenalidiriAdminService) GetStudentFeedbackList(ctx context.Context, req dto_request.GetFeedbackListRequest) (dto_response.FeedbackListResponse, error) {
	page := req.Page
	limit := req.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	filters := repository.StudentFeedbackFilters{
		CategoryID:   req.CategoryID,
		UserName:     req.UserName,
		HasObstacles: req.HasObstacles,
		SortBy:       req.SortBy,
		Limit:        limit,
		Offset:       (page - 1) * limit,
	}

	data, total, err := s.feedbackRepo.ListStudentFeedback(ctx, nil, filters)
	if err != nil {
		return dto_response.FeedbackListResponse{}, err
	}

	items := make([]dto_response.FeedbackItem, 0, len(data))
	for _, fb := range data {
		item := dto_response.FeedbackItem{
			ID:                fb.ID,
			UserName:          fb.User.Fullname,
			EaseOfUseScore:    fb.EaseOfUseScore,
			RelevanceScore:    fb.RelevanceScore,
			SatisfactionScore: fb.SatisfactionScore,
			SubmittedAt:       fb.SubmittedAt.Format(timeFormatDisplay),
			Obstacles:         decodeStringArray(fb.Obstacles),
		}
		items = append(items, item)
	}

	meta := utils.CalculatePaginationMeta(page, limit, total)
	return dto_response.FeedbackListResponse{
		Data: items,
		Pagination: dto_response.PaginationMeta{
			CurrentPage:  meta.CurrentPage,
			TotalPages:   meta.TotalPages,
			TotalRecords: meta.TotalRecords,
			PerPage:      meta.PerPage,
		},
	}, nil
}

func (s *kenalidiriAdminService) GetStudentFeedbackStats(ctx context.Context, categoryID *int64, _ string) (dto_response.FeedbackStatsResponse, error) {
	filters := repository.StudentFeedbackFilters{
		CategoryID: categoryID,
	}

	stats, err := s.feedbackRepo.GetStudentFeedbackStats(ctx, nil, filters)
	if err != nil {
		return dto_response.FeedbackStatsResponse{}, err
	}

	trendData := map[string]interface{}{}
	if td, ok := stats["trend_data"].(map[string]interface{}); ok {
		trendData = td
	}

	return dto_response.FeedbackStatsResponse{
		TotalFeedback:     int(getFloat(stats["total_feedback"])),
		AvgEaseOfUse:      getFloat(stats["avg_ease_of_use"]),
		AvgRelevance:      getFloat(stats["avg_relevance"]),
		AvgSatisfaction:   getFloat(stats["avg_satisfaction"]),
		ParticipationRate: getFloat(stats["participation_rate"]),
		TrendData:         trendData,
	}, nil
}

func (s *kenalidiriAdminService) GetExpertFeedbackList(ctx context.Context, req dto_request.GetExpertFeedbackListRequest) (dto_response.ExpertFeedbackListResponse, error) {
	page := req.Page
	limit := req.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	filters := repository.ExpertFeedbackFilters{
		CategoryID: req.CategoryID,
		ExpertName: req.ExpertName,
		TopNStatus: req.TopNStatus,
		SortBy:     req.SortBy,
		Limit:      limit,
		Offset:     (page - 1) * limit,
	}

	data, total, err := s.feedbackRepo.ListExpertFeedback(ctx, nil, filters)
	if err != nil {
		return dto_response.ExpertFeedbackListResponse{}, err
	}

	items := make([]dto_response.ExpertFeedbackItem, 0, len(data))
	for _, fb := range data {
		topN := deriveTopNStatusFromJSON(fb.TopFiveProfessions, fb.Profession)
		item := dto_response.ExpertFeedbackItem{
			ID:            fb.ID,
			ExpertName:    fb.ExpertName,
			Profession:    fb.Profession,
			TopNStatus:    topN,
			AccuracyScore: fb.AccuracyScore,
			LogicScore:    fb.LogicScore,
			BenefitScore:  fb.BenefitScore,
			SubmittedAt:   fb.SubmittedAt.Format(timeFormatDisplay),
			Obstacles:     decodeStringArray(fb.Obstacles),
		}
		items = append(items, item)
	}

	meta := utils.CalculatePaginationMeta(page, limit, total)
	return dto_response.ExpertFeedbackListResponse{
		Data: items,
		Pagination: dto_response.PaginationMeta{
			CurrentPage:  meta.CurrentPage,
			TotalPages:   meta.TotalPages,
			TotalRecords: meta.TotalRecords,
			PerPage:      meta.PerPage,
		},
	}, nil
}

func (s *kenalidiriAdminService) GetExpertFeedbackDetail(ctx context.Context, feedbackID int64) (dto_response.ExpertFeedbackDetailResponse, error) {
	fb, err := s.feedbackRepo.GetExpertFeedbackByID(ctx, nil, feedbackID)
	if err != nil {
		return dto_response.ExpertFeedbackDetailResponse{}, err
	}

	topN := deriveTopNStatusFromJSON(fb.TopFiveProfessions, fb.Profession)
	return dto_response.ExpertFeedbackDetailResponse{
		ID:                 fb.ID,
		ExpertName:         fb.ExpertName,
		Profession:         fb.Profession,
		Degree:             fb.Degree,
		Experience:         fb.Experience,
		Education:          fb.Education,
		University:         fb.University,
		StudyProgram:       fb.StudyProgram,
		CategoryTest:       fb.CategoryTest,
		TopFiveProfessions: decodeStringArray(fb.TopFiveProfessions),
		TopNStatus:         topN,
		AccuracyScore:      fb.AccuracyScore,
		LogicScore:         fb.LogicScore,
		BenefitScore:       fb.BenefitScore,
		Obstacles:          decodeStringArray(fb.Obstacles),
		Suggestions:        derefString(fb.Suggestions),
		SubmittedAt:        fb.SubmittedAt.Format(timeFormatDisplay),
	}, nil
}

func (s *kenalidiriAdminService) GetRiasecCodeList(ctx context.Context, req dto_request.GetRiasecCodeListRequest) (dto_response.RiasecCodeListResponse, error) {
	var codes []entity.RiasecCode
	var err error

	if req.CodeType != "" {
		codes, err = s.riasecCodeRepo.ListByType(ctx, nil, req.CodeType)
	} else {
		codes, err = s.riasecCodeRepo.GetAll(ctx, nil)
	}
	if err != nil {
		return dto_response.RiasecCodeListResponse{}, err
	}

	items := make([]dto_response.RiasecCodeItem, 0, len(codes))
	for _, c := range codes {
		if req.Search != "" && !strings.Contains(strings.ToLower(c.RiasecTitle), strings.ToLower(req.Search)) && !strings.Contains(strings.ToLower(c.RiasecCode), strings.ToLower(req.Search)) {
			continue
		}

		codeType := deriveCodeType(c.RiasecCode)
		if req.CodeType != "" && codeType != req.CodeType {
			continue
		}

		items = append(items, dto_response.RiasecCodeItem{
			ID:       c.ID,
			Code:     c.RiasecCode,
			Title:    c.RiasecTitle,
			CodeType: codeType,
		})
	}

	return dto_response.RiasecCodeListResponse{
		Data:       items,
		TotalCodes: len(items),
	}, nil
}

func (s *kenalidiriAdminService) GetRiasecCodeDetail(ctx context.Context, codeID int64) (dto_response.RiasecCodeDetailResponse, error) {
	code, err := s.riasecCodeRepo.GetByID(ctx, nil, codeID)
	if err != nil {
		return dto_response.RiasecCodeDetailResponse{}, err
	}

	return dto_response.RiasecCodeDetailResponse{
		ID:                code.ID,
		Code:              code.RiasecCode,
		Title:             code.RiasecTitle,
		Description:       derefString(code.RiasecDescription),
		Strengths:         decodeStringArray(code.Strengths),
		Challenges:        decodeStringArray(code.Challenges),
		Strategies:        decodeStringArray(code.Strategies),
		WorkEnvironments:  decodeStringArray(code.WorkEnvironments),
		InteractionStyles: decodeStringArray(code.InteractionStyles),
	}, nil
}

func (s *kenalidiriAdminService) UpdateRiasecCode(ctx context.Context, codeID int64, req dto_request.UpdateRiasecCodeRequest) error {
	code, err := s.riasecCodeRepo.GetByID(ctx, nil, codeID)
	if err != nil {
		return myerror.RecordNotFound("riasec code")
	}

	code.RiasecTitle = req.RiasecTitle
	code.RiasecDescription = &req.RiasecDescription
	code.Strengths = mustMarshalJSON(req.Strengths)
	code.Challenges = mustMarshalJSON(req.Challenges)
	code.Strategies = mustMarshalJSON(req.Strategies)
	code.WorkEnvironments = mustMarshalJSON(req.WorkEnvironments)
	code.InteractionStyles = mustMarshalJSON(req.InteractionStyles)

	return s.riasecCodeRepo.Update(ctx, nil, code)
}

func buildIkigaiResult(total entity.IkigaiTotalScore) *dto_response.IkigaiResultDetail {
	var data map[string]interface{}
	if len(total.ScoresData) > 0 {
		_ = json.Unmarshal(total.ScoresData, &data)
	}

	return &dto_response.IkigaiResultDetail{
		LoveNarrative:       getString(data, "love_narrative"),
		GoodAtNarrative:     getString(data, "good_at_narrative"),
		WorldNeedsNarrative: getString(data, "world_needs_narrative"),
		PaidForNarrative:    getString(data, "paid_for_narrative"),
	}
}

func buildRecommendations(total entity.IkigaiTotalScore) []dto_response.RecommendationDetail {
	var results []dto_response.RecommendationDetail
	if len(total.ScoresData) == 0 {
		return results
	}

	type profession struct {
		ID              int64  `json:"profession_id"`
		Name            string `json:"profession_name"`
		MatchPercentage int    `json:"match_percentage"`
		Reasoning       string `json:"match_reasoning"`
	}

	var data struct {
		TopProfessions []profession `json:"top_professions"`
	}
	_ = json.Unmarshal(total.ScoresData, &data)

	for idx, p := range data.TopProfessions {
		results = append(results, dto_response.RecommendationDetail{
			Rank:            idx + 1,
			ProfessionID:    p.ID,
			ProfessionName:  p.Name,
			MatchPercentage: p.MatchPercentage,
			MatchReasoning:  p.Reasoning,
		})
		if idx == 1 { // only top 2
			break
		}
	}

	return results
}

func mergeRecommendationNarratives(existing []dto_response.RecommendationDetail, rec entity.CareerRecommendation) []dto_response.RecommendationDetail {
	var narratives []struct {
		ProfessionID   int64  `json:"profession_id"`
		ProfessionName string `json:"profession_name"`
		Reasoning      string `json:"reasoning"`
		Match          int    `json:"match_percentage"`
	}
	_ = json.Unmarshal(rec.RecommendationsData, &narratives)

	if len(existing) == 0 {
		for idx, n := range narratives {
			existing = append(existing, dto_response.RecommendationDetail{
				Rank:            idx + 1,
				ProfessionID:    n.ProfessionID,
				ProfessionName:  n.ProfessionName,
				MatchPercentage: n.Match,
				MatchReasoning:  n.Reasoning,
			})
		}
		return existing
	}

	for i := range existing {
		for _, n := range narratives {
			if n.ProfessionID == existing[i].ProfessionID || strings.EqualFold(n.ProfessionName, existing[i].ProfessionName) {
				existing[i].MatchReasoning = n.Reasoning
				if n.Match > 0 {
					existing[i].MatchPercentage = n.Match
				}
			}
		}
	}

	return existing
}

func decodeStringArray(data datatypes.JSON) []string {
	var arr []string
	if len(data) == 0 {
		return arr
	}

	_ = json.Unmarshal(data, &arr)
	return arr
}

func deriveTopNStatusFromJSON(data datatypes.JSON, profession string) string {
	if len(data) == 0 {
		return "not_found"
	}

	var arr []string
	_ = json.Unmarshal(data, &arr)

	prof := strings.ToLower(strings.TrimSpace(profession))
	for idx, v := range arr {
		if strings.ToLower(strings.TrimSpace(v)) == prof {
			switch idx {
			case 0:
				return "P1"
			case 1:
				return "P2"
			default:
				return "P3-5"
			}
		}
	}

	return "not_found"
}

func getString(data map[string]interface{}, key string) string {
	if data == nil {
		return ""
	}

	if v, ok := data[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getFloat(val interface{}) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int64:
		return float64(v)
	default:
		return 0
	}
}

func deriveCodeType(code string) string {
	switch len(strings.TrimSpace(code)) {
	case 1:
		return "single"
	case 2:
		return "dual"
	default:
		return "triple"
	}
}

func derefString(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func mustMarshalJSON(arr []string) datatypes.JSON {
	bytes, _ := json.Marshal(arr)
	return datatypes.JSON(bytes)
}
