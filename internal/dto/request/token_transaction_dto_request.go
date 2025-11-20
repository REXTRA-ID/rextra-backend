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

type UseTokenAssessmentRequest struct {
	FeatureName   entity.EnumFeature      `json:"feature_name"`
	TokenRequired int                     `json:"token_required"`
	UsageMetaData AssessmentUsageMetaData `json:"usage_metadata"`
}

type UseTokenCVGeneratorRequest struct {
	FeatureName   entity.EnumFeature       `json:"feature_name"`
	TokenRequired int                      `json:"token_required"`
	UsageMetaData CVGeneratorUsageMetaData `json:"usage_metadata"`
}

type UseTokenAiInteweviewerRequest struct {
	FeatureName   entity.EnumFeature         `json:"feature_name"`
	TokenRequired int                        `json:"token_required"`
	UsageMetaData AiInterviewerUsageMetaData `json:"usage_metadata"`
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
