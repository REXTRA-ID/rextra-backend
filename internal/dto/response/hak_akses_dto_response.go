package dto_response

type (
	MakeHakAksesResponse struct {
		Id               string `json:"id"`
		MembershipPlanId string `json:"membership_plan_id"`
		FeatureId        string `json:"feature"`
		Action           string `json:"action"`
	}
)
