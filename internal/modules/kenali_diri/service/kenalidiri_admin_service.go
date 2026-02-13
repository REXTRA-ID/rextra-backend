package service

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	dto_request "rextra-backend/internal/dto/request"
	dto_response "rextra-backend/internal/dto/response"
	"rextra-backend/internal/entity"
	"rextra-backend/internal/modules/kenali_diri/repository"
	"rextra-backend/internal/pkg/cache"
	myerror "rextra-backend/internal/pkg/error"
	"rextra-backend/internal/pkg/export"
	"rextra-backend/internal/utils"

	"errors"

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

	filters := repository.TestSessionFilters{
		TestGoal:    req.TestGoal,
		PersonaType: req.PersonaType,
		Status:      req.Status,
		UserName:    req.UserName,
		SortBy:      req.SortBy,
		Limit:       limit,
		Offset:      (page - 1) * limit,
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

	sessions, total, err := s.testSessionRepo.ListWithFilters(ctx, nil, filters)
	if err != nil {
		return dto_response.TestHistoryListResponse{}, wrapNotFound(err, "careerprofile test session")
	}

	items := make([]dto_response.TestHistoryItem, 0, len(sessions))
	for _, sess := range sessions {
		var resultCode string
		var codeType string
		if strings.EqualFold(sess.Status, "completed") || strings.EqualFold(sess.Status, "riasec_completed") {
			if r, err := s.riasecRepo.GetResultBySessionID(ctx, nil, sess.ID); err == nil {
				resultCode = r.RiasecCode.RiasecCode
				codeType = r.RiasecCodeType
			}
		}

		item := dto_response.TestHistoryItem{
			TestID:         sess.ID,
			UserName:       sess.User.Fullname,
			TestGoal:       string(sess.TestGoal),
			PersonaType:    sess.PersonaType,
			Status:         sess.Status,
			ResultCode:     resultCode,
			StartedAt:      sess.StartedAt.Format(timeFormatDisplay),
			RiasecCodeType: codeType,
		}
		if sess.CompletedAt != nil {
			completed := sess.CompletedAt.Format(timeFormatDisplay)
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
		// Manually delete related data to ensure data integrity, even if DB cascade is enabled.
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
			{table: &entity.UserCareerProfile{}},
		}

		for _, d := range relatedDeletes {
			if err := tx.WithContext(ctx).Where("test_session_id = ?", id).Delete(d.table).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if err := s.testSessionRepo.BulkDelete(ctx, tx, req.TestIDs); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *kenalidiriAdminService) ExportTestHistory(ctx context.Context, req dto_request.ExportTestHistoryRequest) (dto_response.ExportFileResponse, error) {
	filters := repository.TestSessionFilters{
		TestGoal:    req.TestGoal,
		PersonaType: req.PersonaType,
		Status:      req.Status,
		Limit:       0,
		Offset:      0,
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

	sessions, _, err := s.testSessionRepo.ListWithFilters(ctx, nil, filters)
	if err != nil {
		return dto_response.ExportFileResponse{}, wrapNotFound(err, "careerprofile test session")
	}

	if len(sessions) == 0 {
		return dto_response.ExportFileResponse{}, myerror.RecordNotFound("careerprofile test session")
	}

	data := make([]map[string]interface{}, 0, len(sessions))
	sheets := make(map[string][]map[string]interface{})
	sheetUsage := make(map[string]int)
	for _, sess := range sessions {
		var resultCode string
		if strings.EqualFold(sess.Status, "completed") || strings.EqualFold(sess.Status, "riasec_completed") {
			if r, err := s.riasecRepo.GetResultBySessionID(ctx, nil, sess.ID); err == nil {
				resultCode = r.RiasecCode.RiasecCode
			}
		}

		// Group by TestGoal or Persona
		categoryName := string(sess.TestGoal)
		if sess.PersonaType != "" {
			categoryName = fmt.Sprintf("%s - %s", sess.TestGoal, sess.PersonaType)
		}

		record := map[string]interface{}{
			"ID Tes":         sess.ID,
			"Nama Pengguna":  sess.User.Fullname,
			"Email":          sess.User.Email,
			"Tujuan Tes":     string(sess.TestGoal),
			"Tipe Persona":   sess.PersonaType,
			"Status":         sess.Status,
			"Waktu Mulai":    sess.StartedAt.Format(dateFormatExport),
			"Waktu Selesai":  "",
			"Kode RIASEC":    resultCode,
			"Durasi (menit)": "",
		}

		if sess.CompletedAt != nil {
			record["Waktu Selesai"] = sess.CompletedAt.Format(dateFormatExport)
			duration := sess.CompletedAt.Sub(sess.StartedAt).Minutes()
			record["Durasi (menit)"] = fmt.Sprintf("%.0f", duration)
		}

		data = append(data, record)

		baseName := categoryName
		if strings.TrimSpace(baseName) == "" {
			baseName = "Uncategorized"
		}
		// Generate safe sheet name (max 31 chars)
		sheetName := buildSheetName(baseName, 0, sheetUsage)
		sheets[sheetName] = append(sheets[sheetName], record)
	}

	filename := fmt.Sprintf("riwayat_tes_%s", time.Now().Format("20060102_150405"))
	var fileURL string
	switch strings.ToLower(req.Format) {
	case "excel":
		if len(sheets) == 0 {
			sheets["Riwayat Tes"] = data
		}
		fileURL, err = s.exportService.GenerateExcelMultiSheet(sheets, filename)
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

func (s *kenalidiriAdminService) GetTestDetail(ctx context.Context, sessionID int64) (dto_response.TestDetailResponse, error) {
	sess, err := s.testSessionRepo.GetByID(ctx, nil, sessionID)
	if err != nil {
		return dto_response.TestDetailResponse{}, wrapNotFound(err, "careerprofile test session")
	}

	resp := dto_response.TestDetailResponse{
		TestID:      sess.ID,
		UserName:    sess.User.Fullname,
		TestGoal:    string(sess.TestGoal),
		PersonaType: sess.PersonaType,
		Status:      sess.Status,
		StartedAt:   sess.StartedAt.Format(timeFormatDisplay),
	}
	if sess.CompletedAt != nil {
		completed := sess.CompletedAt.Format(timeFormatDisplay)
		resp.CompletedAt = &completed
	}

	if riasec, err := s.riasecRepo.GetResultBySessionID(ctx, nil, sess.ID); err == nil {
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

	if ikigai, err := s.ikigaiRepo.GetTotalScores(ctx, nil, sess.ID); err == nil {
		resp.IkigaiResult = buildIkigaiResult(ikigai)
		resp.Recommendations = append(resp.Recommendations, buildRecommendations(ikigai)...)
	}

	if rec, err := s.recommendationRepo.GetBySessionID(ctx, nil, sess.ID); err == nil {
		resp.Recommendations = mergeRecommendationNarratives(resp.Recommendations, rec)
	}

	return resp, nil
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
		return dto_response.RiasecCodeListResponse{}, wrapNotFound(err, "riasec code")
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
		return dto_response.RiasecCodeDetailResponse{}, wrapNotFound(err, "riasec code")
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

func wrapNotFound(err error, item string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return myerror.RecordNotFound(item)
	}
	return err
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
		if idx == 1 { // Limit to top 2 recommendations
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

func buildSheetName(raw string, categoryID int64, usage map[string]int) string {
	fallback := fmt.Sprintf("Kategori %d", categoryID)
	base := sanitizeSheetBaseName(raw, fallback)
	count := usage[base]
	if count == 0 {
		usage[base] = 1
		return base
	}

	count++
	usage[base] = count
	suffix := fmt.Sprintf("_%d", count)
	maxLen := 31 - utf8.RuneCountInString(suffix)
	if maxLen < 1 {
		maxLen = 1
	}

	baseRunes := []rune(base)
	if len(baseRunes) > maxLen {
		baseRunes = baseRunes[:maxLen]
	}

	return string(baseRunes) + suffix
}

func sanitizeSheetBaseName(name, fallback string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		trimmed = fallback
	}

	for _, invalid := range []string{":", "\\", "/", "?", "*", "[", "]"} {
		trimmed = strings.ReplaceAll(trimmed, invalid, "")
	}

	trimmed = strings.TrimSpace(trimmed)
	if trimmed == "" {
		trimmed = fallback
	}

	if utf8.RuneCountInString(trimmed) > 31 {
		trimmed = string([]rune(trimmed)[:31])
	}

	return trimmed
}
