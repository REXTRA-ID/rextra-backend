package dto_response

import (
	"time"

	"github.com/google/uuid"
)

type (
	TokenWalletDTOResponse struct {
		UserId  string `json:"user_id"`
		Balance int64  `json:"balance"`
	}

	TokenLedgerDTOResponse struct {
		ID         string    `json:"id"`
		OccurredAt time.Time `json:"occurred_at"`
		Amount     int64     `json:"amount"`
		Direction  string    `json:"direction"`
		SourceType string    `json:"source_type"`
	}

	TokenLedgerDetailDTOResponse struct {
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

	TokenBundleDTOResponse struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		TokenAmount  int64  `json:"token_amount"`
		PriceRp      int64  `json:"price_rp"`
		Label        string `json:"label"`
		DisplayOrder int    `json:"display_order"`
		IsActive     bool   `json:"is_available"`
	}

	CustomPricingDTOResponse struct {
		ID                       string              `json:"id"`
		IsEnabled                bool                `json:"is_enabled"`
		Mintoken                 int64               `json:"mintoken"`
		Maxtoken                 int64               `json:"maxtoken"`
		RecommendedPricePerToken int64               `json:"recommended_price_per_token"`
		Tiers                    []CustomPricingTier `json:"tiers"`
	}

	CustomPricingTier struct {
		ID          string  `json:"id"`
		FromToken   int64   `json:"from_token"`
		ToToken     int64   `json:"to_token"`
		DiscountPct float64 `json:"discount_pct"`
	}

	TopupTransactionDTOResponse struct {
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

	TokenLedgerActivityDTOResponse struct {
		ID            string    `json:"id"`
		Username      string    `json:"username"`
		SourceType    string    `json:"source_type"`
		BalanceBefore int64     `json:"balance_before"`
		BalanceAfter  int64     `json:"balance_after"`
		CreatedAt     time.Time `json:"created_at"`
	}

	KPI struct {
		Value      int64   `json:"value"`
		Label      string  `json:"label"`
		Percentage float64 `json:"percentage"`
	}

	TokenKPIOverviewDTOResponse struct {
		TokenIn              KPI `json:"token_in"`
		TokenOut             KPI `json:"token_out"`
		Netflow              KPI `json:"netflow"`
		TopupSuccess         KPI `json:"topup_success"`
		MembershipAllocation KPI `json:"membership_allocation"`
		TokenUsage           KPI `json:"token_usage"`
	}

	TokenGraphActivityDTOResponse struct {
		Time  time.Time `json:"time"`
		Value int64     `json:"value"`
	}
)
