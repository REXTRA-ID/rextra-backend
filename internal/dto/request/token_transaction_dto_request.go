package dto_request

import (
	"encoding/json"
	"rextra-backend/internal/entity"
)

type AssessmentUsageMetaData struct {
	QuestionCount         int    `json:"question_count"`
	CompletionTimeSeconds int    `json:"completion_time_seconds"`
	ResultPersona         string `json:"builder"`
}

type CVGeneratorUsageMetaData struct {
	SectionsGenerated []string `json:"sections_generated"`
	TotalPages        int      `json:"total_pages"`
	Format            string   `json:"format"`
}

type AiInterviewerUsageMetaData struct {
	JobPosition   string `json:"job_position"`
	QuestionCount int    `json:"question_count"`
	AverageScore  int    `json:"average_score"`
}

type MetaData struct {
	AssessmentMetadata *AssessmentUsageMetaData    `json:"assessment"`
	CVGenerator        *CVGeneratorUsageMetaData   `json:"cv_generator"`
	AiInterviewer      *AiInterviewerUsageMetaData `json:"ai_interviewer"`
}

type UseTokenRequest struct {
	FeatureName   entity.EnumFeature `json:"feature_name"`
	TokenRequired int                `json:"token_required"`
	UsageMetaData MetaData           `json:"usage_metadata"`
}

func ConvertAssessmentToBytes(a *AssessmentUsageMetaData) ([]byte, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

func ConvertCVGeneratorToBytes(a *CVGeneratorUsageMetaData) ([]byte, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

func ConvertAiInterviewerToBytes(a *AiInterviewerUsageMetaData) ([]byte, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}
