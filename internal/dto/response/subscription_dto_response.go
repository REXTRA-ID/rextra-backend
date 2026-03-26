package dto_response

// GetUserMembershipResponse adalah response untuk satu user beserta status membership-nya.
// Dipakai untuk list admin (GetAllUsers) dan detail (GetUserById, GetMyMembership).
type GetUserMembershipResponse struct {
	MembershipID string `json:"membership_id"`
	UserID       string `json:"user_id"`

	// Snapshot plan aktif
	PlanName string `json:"plan_name"`
	PlanID   string `json:"plan_id,omitempty"`

	// Snapshot durasi aktif (nil jika plan tanpa durasi)
	DurationID     string `json:"duration_id,omitempty"`
	DurationMonths *int   `json:"duration_months,omitempty"`

	// Periode aktif
	StartedAt     *string `json:"started_at,omitempty"`
	ExpiredAt     *string `json:"expired_at,omitempty"`
	RemainingDays int     `json:"remaining_days"` // 0 jika tidak ada expired_at

	// Status
	IsActive  bool `json:"is_active"`
	AutoRenew bool `json:"auto_renew"`

	// Saldo snapshot
	CurrentTokenBalance int `json:"current_token_balance"`
	CurrentPoinBalance  int `json:"current_poin_balance"`

	// Statistik
	PaidCycleCount   int `json:"paid_cycle_count"`
	EntitlementCount int `json:"entitlement_count"`

	UpdatedAt string `json:"updated_at"`
}

// GetUserMembershipListResponse adalah satu item di list GetAllUsers admin.
// Lebih ringkas — tidak include saldo dan statistik detail.
type GetUserMembershipListResponse struct {
	MembershipID   string  `json:"membership_id"`
	UserID         string  `json:"user_id"`
	UserName       string  `json:"user_name"`
	UserEmail      string  `json:"user_email"`
	PlanName       string  `json:"plan_name"`
	DurationMonths *int    `json:"duration_months,omitempty"`
	IsActive       bool    `json:"is_active"`
	ExpiredAt      *string `json:"expired_at,omitempty"`
	RemainingDays  int     `json:"remaining_days"`
	PaidCycleCount int     `json:"paid_cycle_count"`
	UpdatedAt      string  `json:"updated_at"`
}

// GetSubscriptionCycleResponse adalah response untuk satu siklus berlangganan.
type GetSubscriptionCycleResponse struct {
	ID             string  `json:"id"`
	UserID         string  `json:"user_id"`
	MembershipID   string  `json:"membership_id"`
	PlanName       string  `json:"plan_name"`
	DurationMonths int     `json:"duration_months"`
	AmountPaid     float64 `json:"amount_paid"`
	PaymentChannel string  `json:"payment_channel"`
	StartDate      string  `json:"start_date"`
	EndDate        string  `json:"end_date"`
}

// PaginatedUserMembershipResponse membungkus list GetUserMembershipListResponse dengan metadata pagination.
type PaginatedUserMembershipResponse struct {
	Data       []GetUserMembershipListResponse `json:"data"`
	Total      int                             `json:"total"`
	Page       int                             `json:"page"`
	PageSize   int                             `json:"page_size"`
	TotalPages int                             `json:"total_pages"`
}

// PaginatedSubscriptionCycleResponse membungkus list GetSubscriptionCycleResponse dengan metadata pagination.
type PaginatedSubscriptionCycleResponse struct {
	Data       []GetSubscriptionCycleResponse `json:"data"`
	Total      int                            `json:"total"`
	Page       int                            `json:"page"`
	PageSize   int                            `json:"page_size"`
	TotalPages int                            `json:"total_pages"`
}
