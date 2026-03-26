package dto_request

type CreateMembershipDurationRequest struct {
	DurationMonth        int     `json:"duration_month" binding:"required,oneof=1 3 6 12"`
	TokenBonusPercentage float64 `json:"token_bonus_percentage" binding:"min=0"`
	RextraPoinMultiplier float64 `json:"rextra_poin_multiplier" binding:"min=0"`
}

type UpdateMembershipDurationRequest struct {
	TokenBonusPercentage float64 `json:"token_bonus_percentage" binding:"min=0"`
	RextraPoinMultiplier float64 `json:"rextra_poin_multiplier" binding:"min=0"`
	IsActive             bool    `json:"is_active"`
}
