package dto_request

type (
	MakeHakAksesRequest struct {
		MembershipPlanId string `json:"membership_plan_id"`
		FeatureID        string `json:"feature_id"`
		Action           string `json:"action"`
	}

	CheckAccessRequest struct {
		MembershipPlan string
		Feature        string
		Action         string
	}
)
