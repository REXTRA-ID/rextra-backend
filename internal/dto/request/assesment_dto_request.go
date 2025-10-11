package dto_request

type ValidateHashRequest struct {
	Hash string `json:"hash"`
}

type RiasecQuestionSubmitRequest struct {
	Answer []int `json:"answers"`
}