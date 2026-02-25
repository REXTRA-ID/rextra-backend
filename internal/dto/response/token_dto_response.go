package dto_response

import (
	"time"

	"github.com/google/uuid"
)

type TokenWalletDTOResponse struct {
	UserId  string `json:"user_id"`
	Balance int64  `json:"balance"`
}

type TokenLedgerDTOResponse struct {
	ID         string    `json:"id"`
	OccurredAt time.Time `json:"occurred_at"`
	Amount     int64     `json:"amount"`
	Direction  string    `json:"direction"`
	SourceType string    `json:"source_type"`
}

type TokenLedgerDetailDTOResponse struct {
	ID            string     `json:"id"`
	OccurredAt    time.Time  `json:"occurred_at"`
	Direction     string     `json:"direction"`
	Amount        int64      `json:"amount"`
	BalanceBefore int64      `json:"balance_before"`
	BalanceAfter  int64      `json:"balance_after"`
	SourceType    string     `json:"source_type"`
	SourceID      *uint64    `json:"source_id,omitempty"`
	ReferenceID   *string    `json:"reference_id,omitempty"`
	Description   string     `json:"description"`
	Metadata      any        `json:"metadata"`
	OperatorID    *uuid.UUID `json:"operator_id,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

type TokenBundleDTOResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	TokenAmount  int64  `json:"token_amount"`
	PriceRp      int64  `json:"price_rp"`
	Label        string `json:"label"`
	DisplayOrder int    `json:"display_order"`
	IsActive     bool   `json:"is_available"`
}

type CustomPricingDTOResponse struct {
	ID                       string              `json:"id"`
	IsEnabled                bool                `json:"is_enabled"`
	Mintoken                 int64               `json:"mintoken"`
	Maxtoken                 int64               `json:"maxtoken"`
	RecommendedPricePerToken int64               `json:"recommended_price_per_token"`
	Tiers                    []CustomPricingTier `json:"tiers"`
}

type CustomPricingTier struct {
	ID          string  `json:"id"`
	FromToken   int64   `json:"from_token"`
	ToToken     int64   `json:"to_token"`
	DiscountPct float64 `json:"discount_pct"`
}

type TopupTransactionDTOResponse struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	Type            string    `json:"type"`
	BundlePackageID string    `json:"bundle_package_id,omitempty"`
	TokenAmount     int64     `json:"token_amount"`
	TotalPriceRp    int64     `json:"total_price_rp"`
	Status          string    `json:"status"`
	InvoiceID       string    `json:"invoice_id,omitempty"`
	Provider        string    `json:"provider,omitempty"`
	PaidAt          time.Time `json:"paid_at"`
	ExpiredAt       time.Time `json:"expired_at"`
	LedgerID        string    `json:"ledger_id,omitempty"`
}

type TokenLedgerActivityDTOResponse struct {
	ID            string    `json:"id"`
	Username      string    `json:"username"`
	SourceType    string    `json:"source_type"`
	BalanceBefore int64     `json:"balance_before"`
	BalanceAfter  int64     `json:"balance_after"`
	CreatedAt     time.Time `json:"created_at"`
}

type KPI struct {
	Value      int64   `json:"value"`
	Label      string  `json:"label"`
	Percentage float64 `json:"percentage"`
}

type TokenKPIOverviewDTOResponse struct {
	TokenIn              KPI `json:"token_in"`
	TokenOut             KPI `json:"token_out"`
	Netflow              KPI `json:"netflow"`
	TopupSuccess         KPI `json:"topup_success"`
	MembershipAllocation KPI `json:"membership_allocation"`
	TokenUsage           KPI `json:"token_usage"`
}

type TokenTrendByDirectionDTOResponse struct {
	Date     time.Time `json:"date"`
	TokenIn  int       `json:"token_in"`
	TokenOut int       `json:"token_out"`
	Amount   int       `json:"amount"`
}

type TokenTrendSummaryDTOResponse struct {
	Date     string `json:"date"`
	TokenIn  int    `json:"token_in"`
	TokenOut int    `json:"token_out"`
	Net      int    `json:"net"`
}

type TokenSourceTrendDTOResponse struct {
	Date  string `json:"date"`
	Value int    `json:"value"`
}

type CreateTransactionDTOResponse struct {
	ID        string                 `json:"id"`
	Amount    int                    `json:"amount"`
	Invoice   string                 `json:"invoice"`
	ExpiredAt time.Time              `json:"expired_at"`
	Status    string                 `json:"status"`
	Metadata  map[string]interface{} `json:"metadata"`
}

type CountPaymentPriceDTOResponse struct {
	TotalPrice int `json:"total_price"`
	Fee        int `json:"fee"`
	BasePrice  int `json:"base_price"`
}
