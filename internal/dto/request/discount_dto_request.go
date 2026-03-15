package dto_request

import "time"

type CreateDiscountRequest struct {
	Code string `json:"code" binding:"required,min=3,max=50"`
	Name string `json:"name" binding:"required,min=3,max=100"`
	DiscountType string `json:"discount_type" binding:"required,oneof=PERCENTAGE FIXED"`
	Value float64 `json:"value" binding:"required,gt=0"`
	AppliesTo string `json:"applies_to" binding:"required,oneof=MEMBERSHIP TOKEN_TOPUP GLOBAL"`
	MembershipPlanTargets []string `json:"membership_plan_targets,omitempty"`
	MaxDiscountAmount *int64 `json:"max_discount_amount,omitempty"`
	MinPurchaseAmount *int64 `json:"min_purchase_amount,omitempty"`
	MaxTotalRedemptions *int `json:"max_total_redemptions,omitempty"`
	MaxRedemptionsPerUser *int `json:"max_redemptions_per_user,omitempty"`
	Priority int `json:"priority" binding:"min=0"`
	Stackable bool `json:"stackable"`
	IsPublic bool `json:"is_public"`
	StartsAt *time.Time `json:"starts_at,omitempty"`
	EndsAt *time.Time `json:"ends_at,omitempty"`
	Description string `json:"description" binding:"max=500"`
}

type UpdateDiscountRequest struct {
	Name string `json:"name" binding:"required,min=3,max=100"`
	DiscountType string  `json:"discount_type" binding:"required,oneof=PERCENTAGE FIXED"`
	Value        float64 `json:"value" binding:"required,gt=0"`
	AppliesTo    string  `json:"applies_to" binding:"required,oneof=MEMBERSHIP TOKEN_TOPUP GLOBAL"`
	MembershipPlanTargets []string `json:"membership_plan_targets,omitempty"`
	MaxDiscountAmount *int64 `json:"max_discount_amount,omitempty"`
	MinPurchaseAmount *int64 `json:"min_purchase_amount,omitempty"`
	MaxTotalRedemptions   *int `json:"max_total_redemptions,omitempty"`
	MaxRedemptionsPerUser *int `json:"max_redemptions_per_user,omitempty"`
	Priority  int  `json:"priority" binding:"min=0"`
	Stackable bool `json:"stackable"`
	StartsAt *time.Time `json:"starts_at,omitempty"`
	EndsAt   *time.Time `json:"ends_at,omitempty"`
	Description string `json:"description" binding:"max=500"`
	Status string `json:"status" binding:"required,oneof=ACTIVE INACTIVE"`
}

type DiscountFilterRequest struct {
	Page     int `form:"page" binding:"omitempty,min=0"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
	Status    string `form:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
	AppliesTo string `form:"applies_to" binding:"omitempty,oneof=MEMBERSHIP TOKEN_TOPUP GLOBAL"`
	Search    string `form:"search"`
}

type DiscountRedemptionFilterRequest struct {
	Page     int `form:"page" binding:"omitempty,min=0"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
	Status string `form:"status" binding:"omitempty,oneof=APPLIED REVERSED"`
}

type ValidateDiscountRequest struct {
	Code       string `json:"code" binding:"required"`
	PlanID     string `json:"plan_id" binding:"required"`
	DurationID string `json:"duration_id" binding:"required"`
}
