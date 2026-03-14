package dto_response

type GetDiscountResponse struct {
	ID           string  `json:"id"`
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	DiscountType string  `json:"discount_type"`
	Value        float64 `json:"value"`
	AppliesTo    string  `json:"applies_to"`
	MembershipPlanTargets []string `json:"membership_plan_targets"`
	MaxDiscountAmount *int64 `json:"max_discount_amount,omitempty"`
	MinPurchaseAmount *int64 `json:"min_purchase_amount,omitempty"`
	MaxTotalRedemptions   *int `json:"max_total_redemptions,omitempty"`
	MaxRedemptionsPerUser *int `json:"max_redemptions_per_user,omitempty"`
	CurrentRedemptions    int  `json:"current_redemptions"`
	Priority  int  `json:"priority"`
	Stackable bool `json:"stackable"`
	StartsAt *string `json:"starts_at,omitempty"`
	EndsAt   *string `json:"ends_at,omitempty"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type PaginatedDiscountResponse struct {
	Data       []GetDiscountResponse `json:"data"`
	Total      int                   `json:"total"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"page_size"`
	TotalPages int                   `json:"total_pages"`
}

type GetDiscountRedemptionResponse struct {
	ID            string `json:"id"`
	DiscountID    string `json:"discount_id"`
	UserID        string `json:"user_id"`
	TransactionID string `json:"transaction_id"`
	CodeSnapshot  string  `json:"code_snapshot"`
	PlanSnapshot  *string `json:"plan_snapshot,omitempty"`
	UserName      string  `json:"user_name"`
	AppliesToType string  `json:"applies_to_type"`
	SubtotalAmount int64 `json:"subtotal_amount"`
	DiscountAmount int64 `json:"discount_amount"`
	FinalAmount    int64 `json:"final_amount"`
	Status    string  `json:"status"`
	AppliedAt string  `json:"applied_at"`
	ReversedAt    *string `json:"reversed_at,omitempty"`
	ReverseReason *string `json:"reverse_reason,omitempty"`
}

type PaginatedDiscountRedemptionResponse struct {
	Data       []GetDiscountRedemptionResponse `json:"data"`
	Total      int                             `json:"total"`
	Page       int                             `json:"page"`
	PageSize   int                             `json:"page_size"`
	TotalPages int                             `json:"total_pages"`
}

type ValidateDiscountResponse struct {
	DiscountID   string  `json:"discount_id"`
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	DiscountType string  `json:"discount_type"`
	Value        float64 `json:"value"`
	SubtotalAmount int64 `json:"subtotal_amount"`
	DiscountAmount int64 `json:"discount_amount"`
	FinalAmount    int64 `json:"final_amount"`
	IsStackable bool   `json:"is_stackable"`
	Message     string `json:"message"`
}
