package repository

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"gorm.io/gorm"

	dto_req "rextra-backend/internal/dto/request"
	dto_res "rextra-backend/internal/dto/response"
	"rextra-backend/internal/pkg/feedback"
)

type CareerProfileFeedbackRepository interface {
	GetStudentFeedbackList(ctx context.Context, req dto_req.GetStudentFeedbackRequest) ([]dto_res.CareerProfileStudentFeedbackItem, int64, error)
	GetExpertFeedbackList(ctx context.Context, req dto_req.GetExpertFeedbackRequest) ([]dto_res.CareerProfileExpertFeedbackItem, int64, error)
	GetExpertFeedbackDetail(ctx context.Context, id int64) (*dto_res.CareerProfileExpertFeedbackDetailResponse, error)
	GetFeedbackMetadata(ctx context.Context) (*dto_res.CareerProfileFeedbackMetadataResponse, error)
}

type careerProfileFeedbackRepository struct {
	db *gorm.DB
}

func NewCareerProfileFeedbackRepository(db *gorm.DB) CareerProfileFeedbackRepository {
	return &careerProfileFeedbackRepository{
		db: db,
	}
}

type studentFeedbackRow struct {
	ID                int64     `gorm:"column:id"`
	TestCategory      string    `gorm:"column:test_category"`
	TestCategoryLabel string    `gorm:"column:test_category_label"`
	RespondentUserID  string    `gorm:"column:respondent_user_id"`
	SubmittedAt       time.Time `gorm:"column:submitted_at"`
	StudentName       string    `gorm:"column:student_name"`
	StudentEmail      string    `gorm:"column:student_email"`
	EaseScore         int16     `gorm:"column:ease_score"`
	RelevanceScore    int16     `gorm:"column:relevance_score"`
	SatisfactionScore int16     `gorm:"column:satisfaction_score"`
	MessageToTeam     *string   `gorm:"column:message_to_team"`
	ObstaclesJSON     []byte    `gorm:"column:obstacles"`
}

func (r *careerProfileFeedbackRepository) GetStudentFeedbackList(ctx context.Context, req dto_req.GetStudentFeedbackRequest) ([]dto_res.CareerProfileStudentFeedbackItem, int64, error) {
	var rows []studentFeedbackRow
	var total int64

	whereClauses := []string{"f.test_category = ?", "f.respondent_type = ?", "f.deleted_at IS NULL"}
	args := []interface{}{feedback.TestCategoryCareerProfile, feedback.RespondentTypeStudent}

	if req.Query != "" {
		q := "%" + req.Query + "%"
		whereClauses = append(whereClauses, "(u.fullname ILIKE ? OR u.email ILIKE ?)")
		args = append(args, q, q)
	}
	if req.StartDate != "" {
		whereClauses = append(whereClauses, "DATE(f.submitted_at AT TIME ZONE 'UTC' AT TIME ZONE 'Asia/Jakarta') >= ?")
		args = append(args, req.StartDate)
	}
	if req.EndDate != "" {
		whereClauses = append(whereClauses, "DATE(f.submitted_at AT TIME ZONE 'UTC' AT TIME ZONE 'Asia/Jakarta') <= ?")
		args = append(args, req.EndDate)
	}

	whereSQL := strings.Join(whereClauses, " AND ")

	countSQL := `
        SELECT COUNT(DISTINCT f.id)
        FROM kenalidiri_feedback f
        JOIN users u ON u.id = f.respondent_user_id
        JOIN careerprofile_feedback_student cfs ON cfs.feedback_id = f.id
        WHERE ` + whereSQL

	if err := r.db.WithContext(ctx).Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	sql := `
        SELECT 
            f.id,
            f.test_category,
            f.respondent_user_id,
            f.submitted_at,
            u.fullname AS student_name,
            u.email AS student_email,
            cfs.ease_score,
            cfs.relevance_score,
            cfs.satisfaction_score,
            cfs.message_to_team,
            COALESCE(
              JSON_AGG(
                JSON_BUILD_OBJECT(
                  'key', o.key,
                  'label', o.label,
                  'other_text', cfo.other_text
                ) ORDER BY o.sort_order
              ) FILTER (WHERE o.id IS NOT NULL),
              '[]'::json
            ) AS obstacles
        FROM kenalidiri_feedback f
        JOIN users u ON u.id = f.respondent_user_id
        JOIN careerprofile_feedback_student cfs ON cfs.feedback_id = f.id
        LEFT JOIN careerprofile_feedback_obstacles cfo ON cfo.feedback_id = f.id
        LEFT JOIN careerprofile_obstacle_options o ON o.id = cfo.obstacle_id
        WHERE ` + whereSQL + `
        GROUP BY 
            f.id, f.test_category, f.respondent_user_id, f.submitted_at,
            u.fullname, u.email,
            cfs.ease_score, cfs.relevance_score, cfs.satisfaction_score, cfs.message_to_team
    `

	// Sorting
	sortDir := "DESC"
	if strings.ToLower(req.SortDir) == "asc" {
		sortDir = "ASC"
	}

	orderClause := "f.submitted_at " + sortDir
	if req.SortBy == "name" {
		orderClause = "u.fullname " + sortDir
	} else if req.SortBy == "submitted_at" {
		orderClause = "f.submitted_at " + sortDir
	}

	sql += " ORDER BY " + orderClause
	sql += " LIMIT ? OFFSET ?"

	args = append(args, req.PageSize, (req.Page-1)*req.PageSize)

	if err := r.db.WithContext(ctx).Raw(sql, args...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	result := make([]dto_res.CareerProfileStudentFeedbackItem, len(rows))
	for i, row := range rows {
		var obstacles []dto_res.FeedbackObstacleItem
		if len(row.ObstaclesJSON) > 0 {
			_ = json.Unmarshal(row.ObstaclesJSON, &obstacles)
		}

		result[i] = dto_res.CareerProfileStudentFeedbackItem{
			ID: row.ID,
			Student: dto_res.CareerProfileStudentItem{
				UserID: row.RespondentUserID,
				Name:   row.StudentName,
				Email:  row.StudentEmail,
			},
			Scores: dto_res.CareerProfileStudentScores{
				Ease:         row.EaseScore,
				Relevance:    row.RelevanceScore,
				Satisfaction: row.SatisfactionScore,
			},
			Obstacles:         obstacles,
			MessageToTeam:     row.MessageToTeam,
			SubmittedAt:       row.SubmittedAt,
			TestCategory:      row.TestCategory,
			TestCategoryLabel: "Tes Profil Karier",
		}
	}

	return result, total, nil
}

type expertFeedbackRow struct {
	ID                 int64     `gorm:"column:id"`
	TestCategory       string    `gorm:"column:test_category"`
	SubmittedAt        time.Time `gorm:"column:submitted_at"`
	ExpertName         string    `gorm:"column:expert_name"`
	ExpertProfession   string    `gorm:"column:expert_profession"`
	ExpertProfessionID int64     `gorm:"column:expert_profession_id"`
	Top5Status         string    `gorm:"column:top5_status"`
	AccuracyScore      int16     `gorm:"column:accuracy_score"`
	LogicScore         int16     `gorm:"column:logic_score"`
	UsefulnessScore    int16     `gorm:"column:usefulness_score"`
	SuggestionText     *string   `gorm:"column:suggestion_text"`
	ObstaclesJSON      []byte    `gorm:"column:obstacles"`
}

func (r *careerProfileFeedbackRepository) GetExpertFeedbackList(ctx context.Context, req dto_req.GetExpertFeedbackRequest) ([]dto_res.CareerProfileExpertFeedbackItem, int64, error) {
	var rows []expertFeedbackRow
	var total int64

	whereClauses := []string{"f.test_category = ?", "f.respondent_type = ?", "f.deleted_at IS NULL"}
	args := []interface{}{feedback.TestCategoryCareerProfile, feedback.RespondentTypeExpert}

	if req.Query != "" {
		q := "%" + req.Query + "%"
		whereClauses = append(whereClauses, "(cfe.expert_name ILIKE ? OR cfe.expert_profession ILIKE ?)")
		args = append(args, q, q)
	}
	if req.Top5Status != "" {
		whereClauses = append(whereClauses, "cfe.top5_status = ?")
		args = append(args, req.Top5Status)
	}
	// Date filters same as student
	if req.StartDate != "" {
		whereClauses = append(whereClauses, "DATE(f.submitted_at AT TIME ZONE 'UTC' AT TIME ZONE 'Asia/Jakarta') >= ?")
		args = append(args, req.StartDate)
	}
	if req.EndDate != "" {
		whereClauses = append(whereClauses, "DATE(f.submitted_at AT TIME ZONE 'UTC' AT TIME ZONE 'Asia/Jakarta') <= ?")
		args = append(args, req.EndDate)
	}

	whereSQL := strings.Join(whereClauses, " AND ")

	countSQL := `
        SELECT COUNT(DISTINCT f.id)
        FROM kenalidiri_feedback f
        JOIN careerprofile_feedback_expert cfe ON cfe.feedback_id = f.id
        WHERE ` + whereSQL

	if err := r.db.WithContext(ctx).Raw(countSQL, args...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	sql := `
        SELECT 
            f.id,
            f.test_category,
            f.submitted_at,
            cfe.expert_name,
            cfe.expert_profession,
            cfe.expert_profession_id,
            cfe.top5_status,
            cfe.accuracy_score,
            cfe.logic_score,
            cfe.usefulness_score,
            cfe.suggestion_text,
            COALESCE(
              JSON_AGG(
                JSON_BUILD_OBJECT(
                  'key', ceo.key,
                  'label', ceo.label,
                  'other_text', cfeo.other_text
                ) ORDER BY ceo.sort_order
              ) FILTER (WHERE ceo.id IS NOT NULL),
              '[]'::json
            ) AS obstacles
        FROM kenalidiri_feedback f
        JOIN careerprofile_feedback_expert cfe ON cfe.feedback_id = f.id
        LEFT JOIN careerprofile_feedback_expert_obstacles cfeo ON cfeo.feedback_id = f.id
        LEFT JOIN careerprofile_expert_obstacle_options ceo ON ceo.id = cfeo.obstacle_id
        WHERE ` + whereSQL + `
        GROUP BY 
            f.id, f.test_category, f.submitted_at,
            cfe.expert_name, cfe.expert_profession, cfe.expert_profession_id,
            cfe.top5_status, cfe.accuracy_score, cfe.logic_score, cfe.usefulness_score,
            cfe.suggestion_text
    `

	sortDir := "DESC"
	if strings.ToLower(req.SortDir) == "asc" {
		sortDir = "ASC"
	}

	orderClause := "f.submitted_at " + sortDir
	if req.SortBy == "name" {
		orderClause = "cfe.expert_name " + sortDir
	} else if req.SortBy == "submitted_at" {
		orderClause = "f.submitted_at " + sortDir
	}
	sql += " ORDER BY " + orderClause
	sql += " LIMIT ? OFFSET ?"

	args = append(args, req.PageSize, (req.Page-1)*req.PageSize)

	if err := r.db.WithContext(ctx).Raw(sql, args...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	result := make([]dto_res.CareerProfileExpertFeedbackItem, len(rows))
	for i, row := range rows {
		var obstacles []dto_res.FeedbackObstacleItem
		if len(row.ObstaclesJSON) > 0 {
			_ = json.Unmarshal(row.ObstaclesJSON, &obstacles)
		}

		hasSuggestion := row.SuggestionText != nil && *row.SuggestionText != ""

		statusLabel := row.Top5Status
		switch row.Top5Status {
		case "P1":
			statusLabel = "Muncul di Peringkat 1"
		case "P2":
			statusLabel = "Muncul di Peringkat 2"
		case "P3_5", "P3-5":
			statusLabel = "Muncul di Peringkat 3-5"
		case "NOT_PRESENT":
			statusLabel = "Tidak Muncul"
		}

		result[i] = dto_res.CareerProfileExpertFeedbackItem{
			ID: row.ID,
			Expert: dto_res.CareerProfileExpertItem{
				Name:         row.ExpertName,
				Profession:   row.ExpertProfession,
				ProfessionID: row.ExpertProfessionID,
			},
			Top5Status:      row.Top5Status,
			Top5StatusLabel: statusLabel,
			Scores: dto_res.CareerProfileExpertScores{
				Accuracy:   row.AccuracyScore,
				Logic:      row.LogicScore,
				Usefulness: row.UsefulnessScore,
			},
			Obstacles:         obstacles,
			HasSuggestion:     hasSuggestion,
			SubmittedAt:       row.SubmittedAt,
			TestCategory:      row.TestCategory,
			TestCategoryLabel: "Tes Profil Karier",
		}
	}

	return result, total, nil
}

func (r *careerProfileFeedbackRepository) GetExpertFeedbackDetail(ctx context.Context, id int64) (*dto_res.CareerProfileExpertFeedbackDetailResponse, error) {

	sql := `
        SELECT 
          	f.id,
            f.test_category,
            f.submitted_at,
            cfe.expert_name,
            cfe.expert_profession,
            cfe.expert_profession_id,
            cfe.expert_degree,
            cfe.expert_experience_years,
            cfe.expert_education_level,
            cfe.expert_university,
            cfe.expert_study_program,
            cfe.top5_status,
            cfe.top5_recommendations_json,
            cfe.accuracy_score,
            cfe.logic_score,
            cfe.usefulness_score,
            cfe.suggestion_text,
            COALESCE(
              JSON_AGG(
                JSON_BUILD_OBJECT(
                  'key', ceo.key,
                  'label', ceo.label,
                  'other_text', cfeo.other_text
                ) ORDER BY ceo.sort_order
              ) FILTER (WHERE ceo.id IS NOT NULL),
              '[]'::json
            ) AS obstacles
        FROM kenalidiri_feedback f
        JOIN careerprofile_feedback_expert cfe ON cfe.feedback_id = f.id
        LEFT JOIN careerprofile_feedback_expert_obstacles cfeo ON cfeo.feedback_id = f.id
        LEFT JOIN careerprofile_expert_obstacle_options ceo ON ceo.id = cfeo.obstacle_id
        WHERE f.id = ? AND f.respondent_type = 'EXPERT' AND f.deleted_at IS NULL
        GROUP BY 
             f.id, f.test_category, f.submitted_at,
             cfe.expert_name, cfe.expert_profession, cfe.expert_profession_id,
             cfe.expert_degree, cfe.expert_experience_years, cfe.expert_education_level,
             cfe.expert_university, cfe.expert_study_program,
             cfe.top5_status, cfe.top5_recommendations_json,
             cfe.accuracy_score, cfe.logic_score, cfe.usefulness_score,
             cfe.suggestion_text
    `

	var row struct {
		ID                    int64     `gorm:"column:id"`
		TestCategory          string    `gorm:"column:test_category"`
		SubmittedAt           time.Time `gorm:"column:submitted_at"`
		ExpertName            string    `gorm:"column:expert_name"`
		ExpertProfession      string    `gorm:"column:expert_profession"`
		ExpertProfessionID    int64     `gorm:"column:expert_profession_id"`
		ExpertDegree          string    `gorm:"column:expert_degree"`
		ExpertExperienceYears int       `gorm:"column:expert_experience_years"`
		ExpertEducationLevel  string    `gorm:"column:expert_education_level"`
		ExpertUniversity      string    `gorm:"column:expert_university"`
		ExpertStudyProgram    string    `gorm:"column:expert_study_program"`
		Top5Status            string    `gorm:"column:top5_status"`
		TopFiveProfessions    []byte    `gorm:"column:top5_recommendations_json"`
		AccuracyScore         int16     `gorm:"column:accuracy_score"`
		LogicScore            int16     `gorm:"column:logic_score"`
		UsefulnessScore       int16     `gorm:"column:usefulness_score"`
		SuggestionText        *string   `gorm:"column:suggestion_text"`
		ObstaclesJSON         []byte    `gorm:"column:obstacles"`
	}

	if err := r.db.WithContext(ctx).Raw(sql, id).Scan(&row).Error; err != nil {
		return nil, err
	}

	if row.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var obstacles []dto_res.FeedbackObstacleItem
	if len(row.ObstaclesJSON) > 0 {
		_ = json.Unmarshal(row.ObstaclesJSON, &obstacles)
	}

	var top5List []dto_res.Top5RecommendationItem
	if len(row.TopFiveProfessions) > 0 {
		_ = json.Unmarshal(row.TopFiveProfessions, &top5List)
	}

	var statusLabel string
	switch row.Top5Status {
	case "P1":
		statusLabel = "Profesi expert muncul di Peringkat 1"
	case "P2":
		statusLabel = "Profesi expert muncul di Peringkat 2"
	case "P3_5", "P3-5":
		statusLabel = "Profesi expert muncul di Peringkat 3-5"
	case "NOT_PRESENT":
		statusLabel = "Profesi expert Tidak Muncul"
	default:
		statusLabel = row.Top5Status
	}

	suggText := ""
	if row.SuggestionText != nil {
		suggText = *row.SuggestionText
	}

	res := &dto_res.CareerProfileExpertFeedbackDetailResponse{
		FeedbackID:        row.ID,
		SubmittedAt:       row.SubmittedAt,
		TestCategory:      row.TestCategory,
		TestCategoryLabel: "Tes Profil Karier",
		Identity: dto_res.ExpertIdentity{
			ExpertName:            row.ExpertName,
			ExpertProfession:      row.ExpertProfession,
			ExpertProfessionID:    row.ExpertProfessionID,
			ExpertDegree:          row.ExpertDegree,
			ExpertExperienceYears: row.ExpertExperienceYears,
			ExpertEducationLevel:  row.ExpertEducationLevel,
			ExpertUniversity:      row.ExpertUniversity,
			ExpertStudyProgram:    row.ExpertStudyProgram,
		},
		Top5Recommendations: dto_res.Top5RecommendationData{
			Status:      row.Top5Status,
			StatusLabel: statusLabel,
			List:        top5List,
		},
		Scores: dto_res.ExpertScoreDetails{
			Accuracy:   dto_res.ScoreDetail{Value: row.AccuracyScore, Max: 7, Label: "Akurasi Profil"},
			Logic:      dto_res.ScoreDetail{Value: row.LogicScore, Max: 7, Label: "Logika Penjelasan"},
			Usefulness: dto_res.ScoreDetail{Value: row.UsefulnessScore, Max: 7, Label: "Potensi Manfaat"},
		},
		Obstacles: obstacles,
		Suggestion: dto_res.ExpertSuggestion{
			HasText: suggText != "",
			Text:    suggText,
		},
	}

	return res, nil
}
func (r *careerProfileFeedbackRepository) GetFeedbackMetadata(ctx context.Context) (*dto_res.CareerProfileFeedbackMetadataResponse, error) {
	res := &dto_res.CareerProfileFeedbackMetadataResponse{
		TestCategories: []dto_res.LabelValue{
			{Label: "Tes Profil Karier", Value: feedback.TestCategoryCareerProfile},
		},
		SortOptions: []dto_res.SortOption{
			{Label: "Nama A–Z", Value: "name_asc", SortBy: feedback.SortByName, SortDir: feedback.SortDirAsc},
			{Label: "Nama Z–A", Value: "name_desc", SortBy: feedback.SortByName, SortDir: feedback.SortDirDesc},
			{Label: "Feedback terbaru", Value: "submitted_latest", SortBy: feedback.SortBySubmittedAt, SortDir: feedback.SortDirDesc},
			{Label: "Feedback terlama", Value: "submitted_oldest", SortBy: feedback.SortBySubmittedAt, SortDir: feedback.SortDirAsc},
		},
		Top5Statuses: []dto_res.LabelValue{
			{Label: "Muncul di P1", Value: string(feedback.Top5StatusP1)},
			{Label: "Muncul di P2", Value: string(feedback.Top5StatusP2)},
			{Label: "Muncul di P3–5", Value: string(feedback.Top5StatusP3_5)},
			{Label: "Tidak muncul", Value: string(feedback.Top5StatusNotPresent)},
		},
	}

	return res, nil
}
