package dto_request

import "gorm.io/datatypes"

type ValidateHashRequest struct {
	Hash string `json:"hash"`
}

type RiasecQuestionSubmitRequest map[string]interface{}

type IkigaiQuestionSubmitRequest map[string]interface{}

type IkigaiTest struct {
	ChartData     datatypes.JSON `json:"chart_data"`
	Hash          string       `json:"hash"`
	Profile       string       `json:"profile"`
	Results       datatypes.JSON `json:"results"`
	RiasecExplanations datatypes.JSON `json:"riasec_explanations"`
	RiasecMapFull  datatypes.JSON `json:"riasec_map_full"`
}

type GenerateVouchersRequest struct {
	Amount int    `json:"amount" binding:"required,min=1,max=1000"`
	Prefix string `json:"prefix"` // optional, default "REXTRA"
}