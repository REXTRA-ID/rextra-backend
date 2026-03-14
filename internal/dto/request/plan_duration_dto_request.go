package dto_request

type CreatePlanDurationRequest struct {
	DurationMonths int     `json:"duration_months" binding:"required,oneof=1 3 6 12"`
	Price          int64   `json:"price" binding:"required"`
	DiscountPct    float64 `json:"discount_pct"`
	FinalPrice     int64   `json:"final_price"`
	DurationPrice  int64   `json:"duration_price" binding:"required"`
	TokenAmount    int     `json:"token_amount"`
	BonusToken     int     `json:"bonus_token"`
	PointsActive   bool    `json:"points_active"`
	PointsValue    int     `json:"points_value"`
	BonusPoints    int     `json:"bonus_points"`
}

type UpdatePlanDurationRequest struct {
	Price         int64   `json:"price" binding:"required"`
	DiscountPct   float64 `json:"discount_pct"`
	FinalPrice    int64   `json:"final_price" binding:"required"`
	DurationPrice int64   `json:"duration_price" binding:"required"`
	TokenAmount   int     `json:"token_amount" binding:"required"`
	BonusToken    int     `json:"bonus_token"`
	PointsActive  bool    `json:"points_active"`
	PointsValue   int     `json:"points_value"`
	BonusPoints   int     `json:"bonus_points"`
	IsActive      bool    `json:"is_active"`
}
