package dto_request

type (
	MakeHakAksesRequest struct {
		MembershipPlanId string `json:"membership_plan_id"`
		Feature          string `json:"feature"`
		Action           string `json:"action"`
	}

	CheckAccessRequest struct {
		MembershipPlan string
		Feature        string
		Action         string
	}
)
