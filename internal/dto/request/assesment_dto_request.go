package dto_request

import "gorm.io/datatypes"

type ValidateHashRequest struct {
	Hash string `json:"hash"`
}

type RiasecQuestionSubmitRequest struct {
	Answer []int `json:"answers"`
}

type IkigaiQuestionSubmitRequest struct {
	Answer datatypes.JSON `json:"answers"`
}

type IkigaiTest struct {
	ChartData     datatypes.JSON `json:"chart_data"`
	Hash          string       `json:"hash"`
	Profile       string       `json:"profile"`
	Results       datatypes.JSON `json:"results"`
	RiasecExplanations datatypes.JSON `json:"riasec_explanations"`
	RiasecMapFull  datatypes.JSON `json:"riasec_map_full"`
}