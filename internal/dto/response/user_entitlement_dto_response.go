package dto_response

import "time"

type MyEntitlementQuotaResponse struct {
	EntitlementKey  string `json:"entitlement_key"`
	EntitlementName string `json:"entitlement_name"`
	RestrictionType string `json:"restriction_type"`
	TokenCost       int `json:"token_cost,omitempty"`
	QuotaGranted   *int       `json:"quota_granted,omitempty"`
	QuotaUsed      *int       `json:"quota_used,omitempty"`
	QuotaRemaining *int       `json:"quota_remaining,omitempty"`
	CycleExpiredAt *time.Time `json:"cycle_expired_at,omitempty"`
}

type MyQuotaResponse struct {
	PlanName    string                       `json:"plan_name"`
	Entitlements []MyEntitlementQuotaResponse `json:"entitlements"`
}
