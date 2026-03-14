package dto_response

type GetUserMembershipResponse struct {
	MembershipID string `json:"membership_id"`
	UserID       string `json:"user_id"`
	PlanName     string `json:"plan_name"`
	PlanID       string `json:"plan_id,omitempty"`
	DurationID     string `json:"duration_id,omitempty"`
	DurationMonths *int   `json:"duration_months,omitempty"`
	StartedAt *string `json:"started_at,omitempty"`
	ExpiredAt *string `json:"expired_at,omitempty"`
	RemainingDays int `json:"remaining_days"`
	IsActive  bool `json:"is_active"`
	AutoRenew bool `json:"auto_renew"`
	CurrentTokenBalance int `json:"current_token_balance"`
	CurrentPoinBalance  int `json:"current_poin_balance"`
	PaidCycleCount   int `json:"paid_cycle_count"`
	EntitlementCount int `json:"entitlement_count"`
	UpdatedAt string `json:"updated_at"`
}

type GetUserMembershipListResponse struct {
	MembershipID  string  `json:"membership_id"`
	UserID        string  `json:"user_id"`
	UserName      string  `json:"user_name"`
	UserEmail     string  `json:"user_email"`
	PlanName      string  `json:"plan_name"`
	DurationMonths *int   `json:"duration_months,omitempty"`
	IsActive      bool    `json:"is_active"`
	ExpiredAt     *string `json:"expired_at,omitempty"`
	RemainingDays int     `json:"remaining_days"`
	PaidCycleCount int    `json:"paid_cycle_count"`
	UpdatedAt     string  `json:"updated_at"`
}

type GetSubscriptionCycleResponse struct {
	ID             string  `json:"id"`
	MembershipID   string  `json:"membership_id"`
	CycleNumber    int     `json:"cycle_number"`
	PlanName       string  `json:"plan_name"`
	PlanCategory   string  `json:"plan_category"`
	DurationMonths int     `json:"duration_months"`
	PlanDurationID string  `json:"plan_duration_id,omitempty"`
	StartDate string  `json:"start_date"`
	EndDate   *string `json:"end_date,omitempty"`
	AmountPaid     int64   `json:"amount_paid"`
	PaymentChannel string  `json:"payment_channel"`
	TransactionID  string  `json:"transaction_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type PaginatedUserMembershipResponse struct {
	Data       []GetUserMembershipListResponse `json:"data"`
	Total      int                             `json:"total"`
	Page       int                             `json:"page"`
	PageSize   int                             `json:"page_size"`
	TotalPages int                             `json:"total_pages"`
}

type PaginatedSubscriptionCycleResponse struct {
	Data       []GetSubscriptionCycleResponse `json:"data"`
	Total      int                            `json:"total"`
	Page       int                            `json:"page"`
	PageSize   int                            `json:"page_size"`
	TotalPages int                            `json:"total_pages"`
}
